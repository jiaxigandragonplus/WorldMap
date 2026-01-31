package worldmap

import (
	"math"

	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// ObstacleUnit 障碍物单位实现
type ObstacleUnit struct {
	id       int64
	configId int32
	coord    geo.Coord
	width    int32
	height   int32
	owner    *Owner
	config   *config.ObstacleConfig
	rect     geo.Rectangle
}

// NewObstacleUnit 创建新的障碍物单位
func NewObstacleUnit(id int64, configId int32, coord geo.Coord, width, height int32, config *config.ObstacleConfig) *ObstacleUnit {
	rect := geo.Rectangle{
		Coord:  coord,
		Width:  width,
		Height: height,
	}

	return &ObstacleUnit{
		id:       id,
		configId: configId,
		coord:    coord,
		width:    width,
		height:   height,
		owner:    NewOwner(0, OwnerType_System), // 系统所有
		config:   config,
		rect:     rect,
	}
}

// GetId 获取单位ID
func (o *ObstacleUnit) GetId() int64 {
	return o.id
}

// GetConfigId 获取配置ID
func (o *ObstacleUnit) GetConfigId() int32 {
	return o.configId
}

// GetCoord 获取坐标（返回障碍物中心点）
func (o *ObstacleUnit) GetCoord() *geo.Coord {
	return &o.coord
}

// SetCoord 设置坐标
func (o *ObstacleUnit) SetCoord(coord *geo.Coord) {
	o.coord = *coord
	// 更新矩形区域
	o.rect.Coord = *coord
}

// GetType 获取单位类型
func (o *ObstacleUnit) GetType() MapUnitType {
	return MapUnitType_Obstacle
}

// GetOwner 获取所有者
func (o *ObstacleUnit) GetOwner() *Owner {
	return o.owner
}

// GetObstacleType 获取障碍物类型
func (o *ObstacleUnit) GetObstacleType() string {
	return o.config.ObstacleType
}

// GetWidth 获取宽度
func (o *ObstacleUnit) GetWidth() int32 {
	return o.width
}

// GetHeight 获取高度
func (o *ObstacleUnit) GetHeight() int32 {
	return o.height
}

// GetRect 获取矩形区域
func (o *ObstacleUnit) GetRect() *geo.Rectangle {
	return &o.rect
}

// CanBuildAt 检查指定位置是否可以建造建筑
func (o *ObstacleUnit) CanBuildAt(x, y int32, buildingRadius int32) bool {
	if !o.config.BlockBuilding {
		return true
	}

	// 计算建筑位置到障碍物的距离
	distance := o.calculateDistanceToPoint(x, y)
	blockRadius := o.config.BuildingRadius

	// 如果配置了阻挡半径，使用阻挡半径，否则使用障碍物本身区域
	if blockRadius > 0 {
		return distance > blockRadius
	}

	// 检查是否在障碍物矩形区域内
	return !o.isPointInRect(x, y)
}

// CanSpawnResourceAt 检查指定位置是否可以刷新资源
func (o *ObstacleUnit) CanSpawnResourceAt(x, y int32) bool {
	if !o.config.BlockResource {
		return true
	}

	distance := o.calculateDistanceToPoint(x, y)
	blockRadius := o.config.ResourceRadius

	if blockRadius > 0 {
		return distance > blockRadius
	}

	return !o.isPointInRect(x, y)
}

// CanSpawnMonsterAt 检查指定位置是否可以刷新怪物
func (o *ObstacleUnit) CanSpawnMonsterAt(x, y int32) bool {
	if !o.config.BlockMonster {
		return true
	}

	distance := o.calculateDistanceToPoint(x, y)
	blockRadius := o.config.MonsterRadius

	if blockRadius > 0 {
		return distance > blockRadius
	}

	return !o.isPointInRect(x, y)
}

// CanMarchThrough 检查是否可以行军通过
func (o *ObstacleUnit) CanMarchThrough() bool {
	return o.config.AllowMarch
}

// GetSpecialEffects 获取特殊效果
func (o *ObstacleUnit) GetSpecialEffects() []string {
	return o.config.SpecialEffects
}

// GetEffectStrength 获取效果强度
func (o *ObstacleUnit) GetEffectStrength() float32 {
	return o.config.EffectStrength
}

// calculateDistanceToPoint 计算点到障碍物边缘的最短距离
func (o *ObstacleUnit) calculateDistanceToPoint(x, y int32) int32 {
	// 计算点到矩形边缘的最短距离
	rect := o.rect

	// 如果点在矩形内，距离为0
	if o.isPointInRect(x, y) {
		return 0
	}

	// 计算到四条边的最短距离
	dx := max(rect.Coord.X-x, x-(rect.Coord.X+rect.Width), 0)
	dy := max(rect.Coord.Y-y, y-(rect.Coord.Y+rect.Height), 0)

	// 欧几里得距离
	distance := int32(math.Sqrt(float64(dx*dx + dy*dy)))
	return distance
}

// isPointInRect 检查点是否在矩形区域内
func (o *ObstacleUnit) isPointInRect(x, y int32) bool {
	rect := o.rect
	return x >= rect.Coord.X && x <= rect.Coord.X+rect.Width &&
		y >= rect.Coord.Y && y <= rect.Coord.Y+rect.Height
}

// GetTerrainEffect 获取地形效果（用于障碍物区域）
func (o *ObstacleUnit) GetTerrainEffect(effectName string) (float32, bool) {
	// 单个障碍物没有地形效果，只有障碍物区域有
	return 0, false
}

// GetConfig 获取障碍物配置
func (o *ObstacleUnit) GetConfig() *config.ObstacleConfig {
	return o.config
}

// max 辅助函数，返回最大值
func max(values ...int32) int32 {
	if len(values) == 0 {
		return 0
	}

	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m
}
