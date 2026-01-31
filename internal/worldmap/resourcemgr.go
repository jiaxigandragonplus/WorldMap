package worldmap

import (
	"math/rand"
	"time"

	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// ResourceManager 资源管理器
type ResourceManager struct {
	resources       map[int64]*ResourceUnit              // 资源点ID -> 资源单位
	resourceZones   map[int32]*config.ResourceZoneConfig // 资源区域ID -> 配置
	globalConfig    *config.GlobalRefreshConfig
	lastRefreshTime time.Time
	gridMgr         *GridManager
	obstacleMgr     *ObstacleManager // 障碍物管理器
}

// NewResourceManager 创建资源管理器
func NewResourceManager(gridMgr *GridManager, obstacleMgr *ObstacleManager, globalConfig *config.GlobalRefreshConfig) *ResourceManager {
	return &ResourceManager{
		resources:       make(map[int64]*ResourceUnit),
		resourceZones:   make(map[int32]*config.ResourceZoneConfig),
		globalConfig:    globalConfig,
		lastRefreshTime: time.Now(),
		gridMgr:         gridMgr,
		obstacleMgr:     obstacleMgr,
	}
}

// LoadConfig 加载资源配置
func (rm *ResourceManager) LoadConfig(mapConfig *config.MapConfig) {
	// 加载增强资源点
	for _, pointConfig := range mapConfig.EnhancedResourcePoints {
		coord := geo.Coord{X: pointConfig.X, Y: pointConfig.Y}
		resource, err := NewResourceUnitWithGeneratedID(pointConfig.PointID, coord, &pointConfig)
		if err != nil {
			// 如果ID生成失败，使用配置ID作为回退（不推荐，但保证功能）
			resource = NewResourceUnit(int64(pointConfig.PointID), pointConfig.PointID, coord, &pointConfig)
		}
		rm.resources[resource.GetId()] = resource

		// 将资源单位添加到网格
		if grid := rm.gridMgr.GetGridByCoord(&coord); grid != nil {
			grid.AddUnit(resource)
		}
	}

	// 加载资源区域
	for _, zoneConfig := range mapConfig.ResourceZones {
		rm.resourceZones[zoneConfig.ZoneID] = &zoneConfig
	}

	// 加载旧版资源点（转换为增强版）
	for _, oldPoint := range mapConfig.ResourcePoints {
		enhancedConfig := config.EnhancedResourcePointConfig{
			PointID:         oldPoint.PointID,
			X:               oldPoint.X,
			Y:               oldPoint.Y,
			ResourceType:    oldPoint.ResourceType,
			PointType:       config.ResourcePointType_Fixed,
			RefreshStrategy: config.RefreshStrategy_Linear,
			MaxAmount:       oldPoint.MaxAmount,
			CurrentAmount:   oldPoint.MaxAmount,
			RegenRate:       oldPoint.RegenRate,
			RegenDelay:      oldPoint.RegenDelay,
			SpawnRadius:     0,
			SpawnInterval:   0,
			SpawnChance:     1.0,
		}

		coord := geo.Coord{X: oldPoint.X, Y: oldPoint.Y}
		resource, err := NewResourceUnitWithGeneratedID(oldPoint.PointID, coord, &enhancedConfig)
		if err != nil {
			// 如果ID生成失败，使用配置ID作为回退
			resource = NewResourceUnit(int64(oldPoint.PointID), oldPoint.PointID, coord, &enhancedConfig)
		}
		rm.resources[resource.GetId()] = resource

		if grid := rm.gridMgr.GetGridByCoord(&coord); grid != nil {
			grid.AddUnit(resource)
		}
	}
}

// Update 更新所有资源点状态
func (rm *ResourceManager) Update(now time.Time) {
	// 更新现有资源点
	for _, resource := range rm.resources {
		resource.Update(now)
	}

	// 检查全局刷新
	rm.checkGlobalRefresh(now)

	// 检查区域刷新
	rm.checkZoneRefresh(now)
}

// checkGlobalRefresh 检查全局刷新
func (rm *ResourceManager) checkGlobalRefresh(now time.Time) {
	if rm.globalConfig == nil {
		return
	}

	// 检查每日刷新
	if rm.globalConfig.DailyRefreshTime != "" {
		refreshTime, err := time.Parse("15:04", rm.globalConfig.DailyRefreshTime)
		if err == nil {
			nowTime := time.Date(0, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
			refreshTimeUTC := time.Date(0, 1, 1, refreshTime.Hour(), refreshTime.Minute(), 0, 0, time.UTC)

			if nowTime.Sub(refreshTimeUTC).Abs() < time.Minute &&
				now.Sub(rm.lastRefreshTime) > time.Hour {
				rm.performGlobalRefresh()
				rm.lastRefreshTime = now
			}
		}
	}

	// 检查动态平衡
	if rm.globalConfig.EnableDynamicBalance {
		rm.adjustDynamicBalance()
	}
}

// checkZoneRefresh 检查区域刷新
func (rm *ResourceManager) checkZoneRefresh(now time.Time) {
	for zoneId, zoneConfig := range rm.resourceZones {
		if !zoneConfig.RefreshEnabled {
			continue
		}

		// 计算区域内当前资源点数量
		currentCount := rm.countResourcesInZone(zoneId)

		// 如果资源点数量不足，尝试刷新
		if currentCount < zoneConfig.MaxPoints {
			rm.refreshZoneResources(zoneId, zoneConfig, now)
		}
	}
}

// countResourcesInZone 计算区域内资源点数量
func (rm *ResourceManager) countResourcesInZone(zoneId int32) int32 {
	count := int32(0)
	zoneConfig, exists := rm.resourceZones[zoneId]
	if !exists {
		return 0
	}

	for _, resource := range rm.resources {
		coord := resource.GetCoord()
		if rm.isCoordInZone(coord, zoneConfig) {
			count++
		}
	}

	return count
}

// isCoordInZone 检查坐标是否在区域内
func (rm *ResourceManager) isCoordInZone(coord *geo.Coord, zoneConfig *config.ResourceZoneConfig) bool {
	return coord.X >= zoneConfig.MinX && coord.X <= zoneConfig.MaxX &&
		coord.Y >= zoneConfig.MinY && coord.Y <= zoneConfig.MaxY
}

// refreshZoneResources 刷新区域资源
func (rm *ResourceManager) refreshZoneResources(zoneId int32, zoneConfig *config.ResourceZoneConfig, now time.Time) {
	// 计算需要刷新的数量
	currentCount := rm.countResourcesInZone(zoneId)
	needed := zoneConfig.MaxPoints - currentCount

	if needed <= 0 {
		return
	}

	// 限制每次刷新的数量
	if rm.globalConfig != nil && needed > rm.globalConfig.MaxRefreshPerTick {
		needed = rm.globalConfig.MaxRefreshPerTick
	}

	// 刷新资源点
	for i := int32(0); i < needed; i++ {
		rm.spawnResourceInZone(zoneConfig, now)
	}
}

// spawnResourceInZone 在区域内生成资源点
func (rm *ResourceManager) spawnResourceInZone(zoneConfig *config.ResourceZoneConfig, now time.Time) {
	// 随机选择资源类型
	resourceType := ""
	if len(zoneConfig.ResourceTypes) > 0 {
		idx := rand.Intn(len(zoneConfig.ResourceTypes))
		resourceType = zoneConfig.ResourceTypes[idx]
	} else {
		resourceType = "wood" // 默认木材
	}

	// 随机生成坐标（确保不与其他资源点太近且不在障碍物区域内）
	maxAttempts := 20 // 增加尝试次数，因为要考虑障碍物
	for attempt := 0; attempt < maxAttempts; attempt++ {
		x := zoneConfig.MinX + rand.Int31n(zoneConfig.MaxX-zoneConfig.MinX+1)
		y := zoneConfig.MinY + rand.Int31n(zoneConfig.MaxY-zoneConfig.MinY+1)

		// 检查距离其他资源点是否足够远
		if !rm.isPositionValid(x, y, zoneConfig.MinDistance) {
			continue
		}

		// 检查是否在障碍物阻挡区域内
		if rm.obstacleMgr != nil && !rm.obstacleMgr.CanSpawnResourceAt(x, y) {
			continue
		}

		// 创建资源点配置
		config := config.EnhancedResourcePointConfig{
			PointID:         0, // 临时值，将在创建资源点时使用生成的ID
			X:               x,
			Y:               y,
			ResourceType:    resourceType,
			PointType:       config.ResourcePointType_RandomSpawn,
			RefreshStrategy: config.RefreshStrategy_Linear,
			MaxAmount:       1000,
			CurrentAmount:   1000,
			RegenRate:       1.0,
			RegenDelay:      300,
			SpawnRadius:     50,
			SpawnInterval:   3600,
			SpawnChance:     0.7,
		}

		coord := geo.Coord{X: x, Y: y}
		resource, err := NewResourceUnitWithGeneratedID(0, coord, &config)
		if err != nil {
			// 如果ID生成失败，跳过这个资源点
			continue
		}

		// 更新配置中的PointID为生成的ID（转换为int32）
		config.PointID = int32(resource.GetId())
		rm.resources[resource.GetId()] = resource

		// 添加到网格
		if grid := rm.gridMgr.GetGridByCoord(&coord); grid != nil {
			grid.AddUnit(resource)
		}

		break
	}
}

// isPositionValid 检查位置是否有效（距离其他资源点足够远）
func (rm *ResourceManager) isPositionValid(x, y, minDistance int32) bool {
	if minDistance <= 0 {
		return true
	}

	for _, resource := range rm.resources {
		coord := resource.GetCoord()
		dx := coord.X - x
		dy := coord.Y - y
		distanceSq := dx*dx + dy*dy

		if distanceSq < minDistance*minDistance {
			return false
		}
	}

	return true
}

// performGlobalRefresh 执行全局刷新
func (rm *ResourceManager) performGlobalRefresh() {
	// 刷新所有资源点
	for _, resource := range rm.resources {
		resConfig := resource.GetConfig()
		if resConfig.PointType == config.ResourcePointType_Fixed {
			// 固定资源点恢复到最大量
			// 这里需要调用资源点的方法，简化处理
		}
	}
}

// adjustDynamicBalance 调整动态平衡
func (rm *ResourceManager) adjustDynamicBalance() {
	if rm.globalConfig == nil || !rm.globalConfig.EnableDynamicBalance {
		return
	}

	// 计算总体资源比例
	totalResources := 0
	totalMaxResources := 0

	for _, resource := range rm.resources {
		if resource.IsActive() {
			totalResources += int(resource.GetCurrentAmount())
			config := resource.GetConfig()
			totalMaxResources += int(config.MaxAmount)
		}
	}

	if totalMaxResources == 0 {
		return
	}

	resourceRatio := float32(totalResources) / float32(totalMaxResources)

	// 调整刷新速率
	if resourceRatio < rm.globalConfig.MinResourceRatio {
		// 资源过少，加速刷新
		rm.accelerateRefresh()
	} else if resourceRatio > rm.globalConfig.MaxResourceRatio {
		// 资源过多，减速刷新
		rm.decelerateRefresh()
	}
}

// accelerateRefresh 加速刷新
func (rm *ResourceManager) accelerateRefresh() {
	// 这里可以调整全局刷新参数或资源点刷新速率
	// 简化实现
}

// decelerateRefresh 减速刷新
func (rm *ResourceManager) decelerateRefresh() {
	// 这里可以调整全局刷新参数或资源点刷新速率
	// 简化实现
}

// GetResource 获取资源点
func (rm *ResourceManager) GetResource(resourceId int64) *ResourceUnit {
	return rm.resources[resourceId]
}

// GetResourcesInArea 获取区域内的资源点
func (rm *ResourceManager) GetResourcesInArea(minX, minY, maxX, maxY int32) []*ResourceUnit {
	result := make([]*ResourceUnit, 0)

	for _, resource := range rm.resources {
		coord := resource.GetCoord()
		if coord.X >= minX && coord.X <= maxX && coord.Y >= minY && coord.Y <= maxY {
			result = append(result, resource)
		}
	}

	return result
}

// HarvestResource 采集资源
func (rm *ResourceManager) HarvestResource(resourceId int64, amount int32, playerLevel int32, faction string) int32 {
	resource := rm.resources[resourceId]
	if resource == nil {
		return 0
	}

	if !resource.CanBeHarvestedBy(playerLevel, faction) {
		return 0
	}

	return resource.Harvest(amount)
}

// RemoveResource 移除资源点
func (rm *ResourceManager) RemoveResource(resourceId int64) {
	resource, exists := rm.resources[resourceId]
	if !exists {
		return
	}

	// 从网格中移除
	coord := resource.GetCoord()
	if grid := rm.gridMgr.GetGridByCoord(coord); grid != nil {
		grid.RemoveUnit(resource)
	}

	// 从管理器中移除
	delete(rm.resources, resourceId)
}
