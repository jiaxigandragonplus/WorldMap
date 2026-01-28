package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/geo"

type Unit interface {
	GetId() int64              // 获取unit实例id
	GetConfigId() int32        // 获取unit配置id
	GetCoord() *geo.Coord      // 获取unit坐标
	SetCoord(coord *geo.Coord) // 设置unit坐标
	GetType() MapUnitType      // 获取unit类型
}
