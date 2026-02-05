package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/geo"

// 观察者视野
type ObserverView struct {
	geo.Rectangle
	observers map[int64]bool // 观察者集合
	marching  UnitSet        // 当前view中移动中的单位
}

func NewObserverView(coord *geo.Coord, width, height int32) *ObserverView {
	return &ObserverView{
		Rectangle: *geo.NewRectangle(coord.X, coord.Y, width, height),
		observers: make(map[int64]bool),
	}
}

// 添加观察者
func (ov *ObserverView) AddObserver(playerId int64) {
	ov.observers[playerId] = true
}

func (ov *ObserverView) RemoveObserver(playerId int64) {
	delete(ov.observers, playerId)
}

// 添加行军
func (ov *ObserverView) AddMarching(unit Unit) {
	ov.marching.Insert(unit)
}

func (ov *ObserverView) RemoveMarching(unit Unit) {
	ov.marching.Delete(unit)
}
