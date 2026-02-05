package worldmap

import (
	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// 世界地图
type WorldMap struct {
	id        int64             // 地图实例id
	gridMgr   *GridManager      // 网格管理器
	unitMgr   *UnitManager      // 单位管理器
	playerMgr *MapPlayerManager // 玩家管理器
	mapConfig *config.MapConfig // 地图配置
}

type CityZoneArea struct {
	ZoneID   int32 // 区域id
	CurCount int32 // 当前区域内的主城数量
	CurIndex int32 // 区域索引
}

func NewWorldMap(config *config.MapConfig) *WorldMap {
	return &WorldMap{
		gridMgr:   NewGridManager(config.Width, config.Height, config.GridWidth, config.GridHeight),
		unitMgr:   NewUnitManager(),
		playerMgr: NewMapPlayerManager(),
		mapConfig: config,
	}
}

// 创建一个城市坐标
func (wm *WorldMap) NewCityCoord() (*geo.Coord, bool) {
	return nil, false
}

// 在指定区域内创建城市坐标
func (wm *WorldMap) NewCityCoordInArea(radius int32, area *CityZoneArea) (*geo.Coord, bool) {
	return nil, false
}

// 随机生成一个城市坐标
func (wm *WorldMap) RandomCityCoord() (*geo.Coord, bool) {
	return nil, false
}

// 在指定位置创建一个Npc部队
func (wm *WorldMap) NewNpcTroop(confId int32, level int32, coord *geo.Coord) Unit {
	return nil
}

// 获取可见单位
func (wm *WorldMap) GetVisibleUnits(playerId int64, rect *geo.Rectangle) map[int64]Unit {
	retUnits := make(map[int64]Unit)

	return retUnits
}

func (wm *WorldMap) GetConfig() *config.MapConfig {
	return wm.mapConfig
}
