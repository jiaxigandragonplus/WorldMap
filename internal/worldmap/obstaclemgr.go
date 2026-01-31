package worldmap

import (
	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// ObstacleManager 障碍物管理器
type ObstacleManager struct {
	obstacles      map[int64]*ObstacleUnit              // 障碍物ID -> 障碍物单位
	obstacleZones  map[int32]*config.ObstacleZoneConfig // 障碍物区域ID -> 配置
	nextObstacleId int64
	gridMgr        *GridManager
}

// NewObstacleManager 创建障碍物管理器
func NewObstacleManager(gridMgr *GridManager) *ObstacleManager {
	return &ObstacleManager{
		obstacles:      make(map[int64]*ObstacleUnit),
		obstacleZones:  make(map[int32]*config.ObstacleZoneConfig),
		nextObstacleId: 2000,
		gridMgr:        gridMgr,
	}
}

// LoadConfig 加载障碍物配置
func (om *ObstacleManager) LoadConfig(mapConfig *config.MapConfig) {
	// 加载单个障碍物
	for _, obstacleConfig := range mapConfig.Obstacles {
		coord := geo.Coord{X: obstacleConfig.X, Y: obstacleConfig.Y}
		obstacle := NewObstacleUnit(
			om.nextObstacleId,
			obstacleConfig.ObstacleID,
			coord,
			obstacleConfig.Width,
			obstacleConfig.Height,
			&obstacleConfig,
		)
		om.obstacles[om.nextObstacleId] = obstacle
		om.nextObstacleId++

		// 将障碍物单位添加到网格
		if grid := om.gridMgr.GetGridByCoord(&coord); grid != nil {
			grid.AddUnit(obstacle)
		}
	}

	// 加载障碍物区域
	for _, zoneConfig := range mapConfig.ObstacleZones {
		om.obstacleZones[zoneConfig.ZoneID] = &zoneConfig
		// 障碍物区域可以动态生成障碍物，这里先只存储配置
	}
}

// CanBuildAt 检查指定位置是否可以建造建筑
func (om *ObstacleManager) CanBuildAt(x, y int32, buildingRadius int32) bool {
	// 检查所有障碍物
	for _, obstacle := range om.obstacles {
		if !obstacle.CanBuildAt(x, y, buildingRadius) {
			return false
		}
	}

	// 检查障碍物区域
	for _, zoneConfig := range om.obstacleZones {
		if zoneConfig.BlockBuilding && om.isPointInZone(x, y, zoneConfig) {
			return false
		}
	}

	return true
}

// CanSpawnResourceAt 检查指定位置是否可以刷新资源
func (om *ObstacleManager) CanSpawnResourceAt(x, y int32) bool {
	// 检查所有障碍物
	for _, obstacle := range om.obstacles {
		if !obstacle.CanSpawnResourceAt(x, y) {
			return false
		}
	}

	// 检查障碍物区域
	for _, zoneConfig := range om.obstacleZones {
		if zoneConfig.BlockResource && om.isPointInZone(x, y, zoneConfig) {
			return false
		}
	}

	return true
}

// CanSpawnMonsterAt 检查指定位置是否可以刷新怪物
func (om *ObstacleManager) CanSpawnMonsterAt(x, y int32) bool {
	// 检查所有障碍物
	for _, obstacle := range om.obstacles {
		if !obstacle.CanSpawnMonsterAt(x, y) {
			return false
		}
	}

	// 检查障碍物区域
	for _, zoneConfig := range om.obstacleZones {
		if zoneConfig.BlockMonster && om.isPointInZone(x, y, zoneConfig) {
			return false
		}
	}

	return true
}

// CanMarchThrough 检查是否可以行军通过指定位置
func (om *ObstacleManager) CanMarchThrough(x, y int32) bool {
	// 检查所有障碍物（只检查是否在障碍物区域内）
	for _, obstacle := range om.obstacles {
		rect := obstacle.GetRect()
		if om.isPointInRect(x, y, rect) && !obstacle.CanMarchThrough() {
			return false
		}
	}

	// 检查障碍物区域
	for _, zoneConfig := range om.obstacleZones {
		if !zoneConfig.AllowMarch && om.isPointInZone(x, y, zoneConfig) {
			return false
		}
	}

	return true
}

// GetTerrainEffect 获取指定位置的地形效果
func (om *ObstacleManager) GetTerrainEffect(x, y int32, effectName string) (float32, bool) {
	// 检查障碍物区域
	for _, zoneConfig := range om.obstacleZones {
		if om.isPointInZone(x, y, zoneConfig) {
			if effectValue, exists := zoneConfig.TerrainEffects[effectName]; exists {
				return effectValue, true
			}
		}
	}

	return 0, false
}

// GetObstaclesInArea 获取区域内的障碍物
func (om *ObstacleManager) GetObstaclesInArea(minX, minY, maxX, maxY int32) []*ObstacleUnit {
	result := make([]*ObstacleUnit, 0)

	for _, obstacle := range om.obstacles {
		coord := obstacle.GetCoord()
		if coord.X >= minX && coord.X <= maxX && coord.Y >= minY && coord.Y <= maxY {
			result = append(result, obstacle)
		}
	}

	return result
}

// GetObstacle 获取障碍物
func (om *ObstacleManager) GetObstacle(obstacleId int64) *ObstacleUnit {
	return om.obstacles[obstacleId]
}

// RemoveObstacle 移除障碍物
func (om *ObstacleManager) RemoveObstacle(obstacleId int64) {
	obstacle, exists := om.obstacles[obstacleId]
	if !exists {
		return
	}

	// 从网格中移除
	coord := obstacle.GetCoord()
	if grid := om.gridMgr.GetGridByCoord(coord); grid != nil {
		grid.RemoveUnit(obstacle)
	}

	// 从管理器中移除
	delete(om.obstacles, obstacleId)
}

// isPointInZone 检查点是否在障碍物区域内
func (om *ObstacleManager) isPointInZone(x, y int32, zoneConfig *config.ObstacleZoneConfig) bool {
	return x >= zoneConfig.MinX && x <= zoneConfig.MaxX &&
		y >= zoneConfig.MinY && y <= zoneConfig.MaxY
}

// isPointInRect 检查点是否在矩形区域内
func (om *ObstacleManager) isPointInRect(x, y int32, rect *geo.Rectangle) bool {
	return x >= rect.Coord.X && x <= rect.Coord.X+rect.Width &&
		y >= rect.Coord.Y && y <= rect.Coord.Y+rect.Height
}

// GenerateZoneObstacles 生成障碍物区域内的障碍物（按需生成）
func (om *ObstacleManager) GenerateZoneObstacles(zoneId int32) {
	zoneConfig, exists := om.obstacleZones[zoneId]
	if !exists {
		return
	}

	// 计算需要生成的障碍物数量
	areaWidth := zoneConfig.MaxX - zoneConfig.MinX
	areaHeight := zoneConfig.MaxY - zoneConfig.MinY
	areaSize := areaWidth * areaHeight
	expectedCount := int32(float32(areaSize) * zoneConfig.Density / 10000.0) // 假设单位面积

	// 简单实现：在区域内随机生成障碍物
	// 实际项目中应该使用更复杂的算法
	for i := int32(0); i < expectedCount; i++ {
		// 随机位置和大小
		x := zoneConfig.MinX + (areaWidth/10)*(i%10)
		y := zoneConfig.MinY + (areaHeight/10)*(i/10)
		width := zoneConfig.MinSize + (zoneConfig.MaxSize-zoneConfig.MinSize)/2
		height := width

		// 创建障碍物配置
		obstacleConfig := config.ObstacleConfig{
			ObstacleID:     int32(om.nextObstacleId),
			X:              x,
			Y:              y,
			Width:          width,
			Height:         height,
			ObstacleType:   zoneConfig.ObstacleType,
			Name:           zoneConfig.ZoneName + "_" + string(rune('A'+i)),
			BlockBuilding:  zoneConfig.BlockBuilding,
			BlockResource:  zoneConfig.BlockResource,
			BlockMonster:   zoneConfig.BlockMonster,
			AllowMarch:     zoneConfig.AllowMarch,
			BuildingRadius: 0,
			ResourceRadius: 0,
			MonsterRadius:  0,
		}

		coord := geo.Coord{X: x, Y: y}
		obstacle := NewObstacleUnit(
			om.nextObstacleId,
			obstacleConfig.ObstacleID,
			coord,
			width,
			height,
			&obstacleConfig,
		)
		om.obstacles[om.nextObstacleId] = obstacle
		om.nextObstacleId++

		if grid := om.gridMgr.GetGridByCoord(&coord); grid != nil {
			grid.AddUnit(obstacle)
		}
	}
}
