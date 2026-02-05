package worldmap

import (
	"time"

	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// ResourceUnit 资源点单位实现
type ResourceUnit struct {
	id              int64
	configId        int32
	coord           geo.Coord
	owner           *Owner
	config          *config.EnhancedResourcePointConfig
	lastHarvestTime time.Time
	lastRefreshTime time.Time
	currentAmount   int32
	isActive        bool
}

// NewResourceUnit 创建新的资源点单位（指定ID）
func NewResourceUnit(id int64, configId int32, coord geo.Coord, config *config.EnhancedResourcePointConfig) *ResourceUnit {
	return &ResourceUnit{
		id:              id,
		configId:        configId,
		coord:           coord,
		owner:           NewOwner(0, OwnerType_System), // 系统所有
		config:          config,
		currentAmount:   config.CurrentAmount,
		isActive:        true,
		lastHarvestTime: time.Now(),
		lastRefreshTime: time.Now(),
	}
}

// NewResourceUnitWithGeneratedID 创建新的资源点单位（自动生成ID）
func NewResourceUnitWithGeneratedID(configId int32, coord geo.Coord, config *config.EnhancedResourcePointConfig) (*ResourceUnit, error) {
	id := GetIDGenerator().GenerateNewID()

	return &ResourceUnit{
		id:              id,
		configId:        configId,
		coord:           coord,
		owner:           NewOwner(0, OwnerType_System), // 系统所有
		config:          config,
		currentAmount:   config.CurrentAmount,
		isActive:        true,
		lastHarvestTime: time.Now(),
		lastRefreshTime: time.Now(),
	}, nil
}

// GetId 获取单位ID
func (r *ResourceUnit) GetId() int64 {
	return r.id
}

// GetConfigId 获取配置ID
func (r *ResourceUnit) GetConfigId() int32 {
	return r.configId
}

// GetCoord 获取坐标
func (r *ResourceUnit) GetCoord() *geo.Coord {
	return &r.coord
}

// SetCoord 设置坐标
func (r *ResourceUnit) SetCoord(coord *geo.Coord) {
	r.coord = *coord
}

// GetType 获取单位类型
func (r *ResourceUnit) GetType() MapUnitType {
	return MapUnitType_Resource
}

// GetOwner 获取所有者
func (r *ResourceUnit) GetOwner() *Owner {
	return r.owner
}

// GetResourceType 获取资源类型
func (r *ResourceUnit) GetResourceType() string {
	return r.config.ResourceType
}

// GetCurrentAmount 获取当前资源量
func (r *ResourceUnit) GetCurrentAmount() int32 {
	return r.currentAmount
}

// Harvest 采集资源
func (r *ResourceUnit) Harvest(amount int32) int32 {
	if amount <= 0 || !r.isActive {
		return 0
	}

	harvested := amount
	if harvested > r.currentAmount {
		harvested = r.currentAmount
	}

	r.currentAmount -= harvested
	r.lastHarvestTime = time.Now()

	// 如果资源被采空，根据配置处理
	if r.currentAmount <= 0 {
		r.currentAmount = 0
		if r.config.PointType == config.ResourcePointType_RandomSpawn {
			r.isActive = false
		}
	}

	return harvested
}

// Update 更新资源点状态
func (r *ResourceUnit) Update(now time.Time) {
	if !r.isActive {
		// 检查是否需要重新激活（对于随机刷新点）
		if r.config.PointType == config.ResourcePointType_RandomSpawn {
			elapsed := now.Sub(r.lastHarvestTime).Seconds()
			if int32(elapsed) >= r.config.SpawnInterval {
				// 随机决定是否刷新
				// 这里简化处理，实际应该使用随机数
				r.isActive = true
				r.currentAmount = r.config.MaxAmount
				r.lastRefreshTime = now
			}
		}
		return
	}

	// 检查时间限制
	if !r.isInActiveTime(now) {
		return
	}

	// 恢复资源
	if r.currentAmount < r.config.MaxAmount {
		elapsed := now.Sub(r.lastHarvestTime).Seconds()
		if int32(elapsed) >= r.config.RegenDelay {
			// 根据策略计算恢复量
			regenAmount := r.calculateRegenAmount(now)
			r.currentAmount += regenAmount
			if r.currentAmount > r.config.MaxAmount {
				r.currentAmount = r.config.MaxAmount
			}
		}
	}
}

// isInActiveTime 检查当前时间是否在活跃时间内
func (r *ResourceUnit) isInActiveTime(now time.Time) bool {
	// 检查月份
	if len(r.config.SeasonMonths) > 0 {
		currentMonth := int32(now.Month())
		found := false
		for _, month := range r.config.SeasonMonths {
			if month == currentMonth {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查小时
	if len(r.config.ActiveHours) >= 2 {
		currentHour := int32(now.Hour())
		if currentHour < r.config.ActiveHours[0] || currentHour > r.config.ActiveHours[1] {
			return false
		}
	}

	return true
}

// calculateRegenAmount 根据刷新策略计算恢复量
func (r *ResourceUnit) calculateRegenAmount(now time.Time) int32 {
	elapsed := now.Sub(r.lastHarvestTime).Seconds()

	switch r.config.RefreshStrategy {
	case config.RefreshStrategy_Linear:
		// 线性恢复：每秒恢复 RegenRate
		return int32(elapsed * float64(r.config.RegenRate))

	case config.RefreshStrategy_Exponential:
		// 指数恢复：恢复速率随时间增加
		baseRate := float64(r.config.RegenRate)
		exponentialFactor := 1.0 + (elapsed / 3600.0) // 每小时增加1倍
		return int32(elapsed * baseRate * exponentialFactor)

	case config.RefreshStrategy_Stepwise:
		// 阶梯式恢复：每间隔一段时间恢复固定量
		interval := float64(r.config.RegenDelay)
		if interval <= 0 {
			interval = 300 // 默认5分钟
		}
		steps := int32(elapsed / interval)
		return steps * int32(r.config.RegenRate*float32(interval))

	case config.RefreshStrategy_Random:
		// 随机恢复：在0到最大恢复量之间随机
		maxRegen := int32(elapsed * float64(r.config.RegenRate) * 2)
		// 这里简化处理，实际应该使用随机数
		return maxRegen / 2

	default:
		return int32(elapsed * float64(r.config.RegenRate))
	}
}

// CanBeHarvestedBy 检查指定玩家是否可以采集
func (r *ResourceUnit) CanBeHarvestedBy(playerLevel int32, faction string) bool {
	if !r.isActive {
		return false
	}

	// 检查等级限制
	if r.config.MinPlayerLevel > 0 && playerLevel < r.config.MinPlayerLevel {
		return false
	}
	if r.config.MaxPlayerLevel > 0 && playerLevel > r.config.MaxPlayerLevel {
		return false
	}

	// 检查阵营限制
	if r.config.FactionRestrict != "" && r.config.FactionRestrict != faction {
		return false
	}

	return true
}

// GetConfig 获取资源配置
func (r *ResourceUnit) GetConfig() *config.EnhancedResourcePointConfig {
	return r.config
}

// IsActive 检查资源点是否活跃
func (r *ResourceUnit) IsActive() bool {
	return r.isActive
}
