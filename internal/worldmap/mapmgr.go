package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/config"

// 地图管理器
type MapManager struct {
	maps map[int64]*WorldMap
}

func NewMapManager() *MapManager {
	mm := &MapManager{
		maps: make(map[int64]*WorldMap),
	}
	return mm
}

func (mm *MapManager) GetMap(mapID int64) *WorldMap {
	if _, ok := mm.maps[mapID]; !ok {
		return nil
	}
	return mm.maps[mapID]
}

func (mm *MapManager) CreateMap(mapConfig *config.MapConfig) *WorldMap {
	mapID := GetIDGenerator().GenerateNewID()
	mm.maps[mapID] = NewWorldMap(mapConfig)
	return mm.maps[mapID]
}
