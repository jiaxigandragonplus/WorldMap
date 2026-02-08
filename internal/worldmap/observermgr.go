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

func (om *ObserverManager) GetObserver(id int64) *Observer {
	return om.observers[id]
}

func (om *ObserverManager) AddObserver(playerId int64, viewWindow *geo.Rectangle, lod int32) *Observer {
	observer := NewObserver(playerId, viewWindow, lod)
	om.observers[observer.Id] = observer

	return observer
}

func (om *ObserverManager) GetObserverViewByIndex(x, y int32) (*ObserverView, bool) {
	maxX, maxY := om.GetMaxViewSize()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return nil, false
	}
	return om.views[x][y], true
}

// 地图划分后的最大单边视口数量
func (om *ObserverManager) GetMaxViewSize() (int32, int32) {
	mapConfig := om.worldMap.GetConfig()
	return (mapConfig.MapSize.Width - 1) / mapConfig.MapSize.GridWidth, (mapConfig.MapSize.Height - 1) / mapConfig.MapSize.GridHeight
}

// 坐标到视野索引
func (om *ObserverManager) CoordToViewIndex(coord *geo.Coord) (int32, int32) {
	mapConfig := om.worldMap.GetConfig()
	return coord.X / mapConfig.MapSize.GridWidth, coord.Y / mapConfig.MapSize.GridHeight
}

func (om *ObserverManager) GetMarchingUnits(playerId int64, rect *geo.Rectangle) map[int64]Unit {
	retUnits := make(map[int64]Unit)
	observer := om.GetObserver(playerId)
	if observer == nil {
		return retUnits
	}

	return retUnits
}

// 获取矩形区域覆盖的view
func (om *ObserverManager) GetCoverViews(rect *geo.Rectangle) []*ObserverView {
	retViews := make([]*ObserverView, 0)

	mapSize := om.worldMap.GetMapSize()
	if rect.X == 0 && rect.Y == 0 && rect.Width == mapSize.Width && rect.Height == mapSize.Height {
		// 最大的区域
		for i := range om.views {
			retViews = append(retViews, om.views[i]...)
		}
		return retViews
	}

	leftX := rect.X / mapSize.GridWidth
	rightX := (rect.X + rect.Width) / mapSize.GridWidth
	leftY := rect.Y / mapSize.GridHeight
	rightY := (rect.Y + rect.Height) / mapSize.GridHeight

	for x := leftX; x <= rightX; x++ {
		for y := leftY; y <= rightY; y++ {
			view, ok := om.GetObserverViewByIndex(x, y)
			if ok {
				retViews = append(retViews, view)
			}
		}
	}
	return retViews
}
