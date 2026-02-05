package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/config"

// 世界地图
type WorldMap struct {
	id        int64             // 地图实例id
	gridMgr   *GridManager      // 网格管理器
	unitMgr   *UnitManager      // 单位管理器
	playerMgr *MapPlayerManager // 玩家管理器
}

func NewWorldMap(config *config.MapConfig) *WorldMap {
	return &WorldMap{
		gridMgr:   NewGridManager(config.Width, config.Height, config.GridWidth, config.GridHeight),
		unitMgr:   NewUnitManager(),
		playerMgr: NewMapPlayerManager(),
	}
}
