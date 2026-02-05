package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/geo"

// 观察者管理器

type ObserverManager struct {
	worldMap  *WorldMap              // 所属地图指针
	observers map[int64]*Observer    // 观察者
	views     [][]*ObserverView      // 视野管理
	marching  map[int64][]*geo.Coord // 所有的行军
}

func NewObserverManager(worldMap *WorldMap) *ObserverManager {
	return &ObserverManager{
		worldMap:  worldMap,
		observers: make(map[int64]*Observer),
	}
}

func (om *ObserverManager) AddObserver(playerId int64, viewWindow *geo.Rectangle, lod int32) *Observer {
	observer := NewObserver(playerId, viewWindow, lod)
	om.observers[observer.Id] = observer

	return observer
}

// 地图划分后的最大单边视口数量
func (om *ObserverManager) GetMaxViewSize() (int32, int32) {
	mapConfig := om.worldMap.GetConfig()
	return (mapConfig.Width - 1) / mapConfig.GridWidth, (mapConfig.Height - 1) / mapConfig.GridHeight
}

// 坐标到视野索引
func (om *ObserverManager) CoordToViewIndex(coord *geo.Coord) (int32, int32) {
	mapConfig := om.worldMap.GetConfig()
	return coord.X / mapConfig.GridWidth, coord.Y / mapConfig.GridHeight
}
