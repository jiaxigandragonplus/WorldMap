package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/geo"

// 观察者
type Observer struct {
	Id         int64          // 观察者id
	ViewWindow *geo.Rectangle // 观察窗口
	Lod        int32          // 观察等级
}

func NewObserver(id int64, viewWindow *geo.Rectangle, lod int32) *Observer {
	return &Observer{
		Id:         id,
		ViewWindow: viewWindow,
		Lod:        lod,
	}
}

func (o *Observer) ChangeLod(lod int32) {
	o.Lod = lod
}

// 指定实体对玩家是否可见
func (o *Observer) IsVisible(unit Unit) bool {
	return false
}
