package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/geo"

// Unit 地图单位接口
type Unit interface {
	GetId() int64                    // 获取 unit 实例 id
	GetConfigId() int32              // 获取 unit 配置 id
	GetCoord() *geo.Coord            // 获取 unit 坐标（矩形网格坐标）
	SetCoord(coord *geo.Coord)       // 设置 unit 坐标
	GetHexCoord() *geo.HexCoord      // 获取 unit 六边形坐标
	SetHexCoord(coord *geo.HexCoord) // 设置 unit 六边形坐标
	GetType() MapUnitType            // 获取 unit 类型
	GetOwner() *Owner                // 获取 unit 所属者
}

// BaseUnit 基础单位实现，提供通用的坐标管理
type BaseUnit struct {
	id       int64
	configId int32
	coord    geo.Coord
	hexCoord *geo.HexCoord
	unitType MapUnitType
	owner    *Owner
}

// NewBaseUnit 创建基础单位
func NewBaseUnit(id int64, configId int32, coord geo.Coord, hexCoord *geo.HexCoord, unitType MapUnitType, owner *Owner) *BaseUnit {
	return &BaseUnit{
		id:       id,
		configId: configId,
		coord:    coord,
		hexCoord: hexCoord,
		unitType: unitType,
		owner:    owner,
	}
}

// GetId 获取单位 ID
func (b *BaseUnit) GetId() int64 {
	return b.id
}

// GetConfigId 获取配置 ID
func (b *BaseUnit) GetConfigId() int32 {
	return b.configId
}

// GetCoord 获取坐标
func (b *BaseUnit) GetCoord() *geo.Coord {
	return &b.coord
}

// SetCoord 设置坐标
func (b *BaseUnit) SetCoord(coord *geo.Coord) {
	b.coord = *coord
}

// GetHexCoord 获取六边形坐标
func (b *BaseUnit) GetHexCoord() *geo.HexCoord {
	return b.hexCoord
}

// SetHexCoord 设置六边形坐标
func (b *BaseUnit) SetHexCoord(coord *geo.HexCoord) {
	b.hexCoord = coord
}

// GetType 获取单位类型
func (b *BaseUnit) GetType() MapUnitType {
	return b.unitType
}

// GetOwner 获取所有者
func (b *BaseUnit) GetOwner() *Owner {
	return b.owner
}
