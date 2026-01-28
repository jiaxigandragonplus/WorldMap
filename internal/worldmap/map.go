package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/config"

// 世界地图管理器
type WorldMapManager struct {
	gridMgr   *GridManager
	unitMgr   *UnitManager
	playerMgr *MapPlayerManager
}

func NewWorldMapManager(config *config.MapConfig) *WorldMapManager {
	return &WorldMapManager{
		gridMgr:   NewGridManager(config.Width, config.Height, config.XGridLen, config.YGridLen),
		unitMgr:   NewUnitManager(),
		playerMgr: NewMapPlayerManager(),
	}
}
