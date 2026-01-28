package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/geo"

type Grid struct {
	geo.Rectangle
	units []Unit
}

func NewGrid(coord *geo.Coord, width, height int32) *Grid {
	return &Grid{
		Rectangle: geo.Rectangle{
			Coord:  *coord,
			Width:  width,
			Height: height,
		},
		units: make([]Unit, 0),
	}
}

// 添加地图单位
func (g *Grid) AddUnit(unit Unit) {
	g.units = append(g.units, unit)
}

// 移除地图单位
func (g *Grid) RemoveUnit(unit Unit) {
	for i, u := range g.units {
		if u == unit {
			g.units = append(g.units[:i], g.units[i+1:]...)
			return
		}
	}
}

// 获取地图单位数组
func (g *Grid) GetUnits() []Unit {
	return g.units
}

// 判断地图单位是否存在
func (g *Grid) IsExistUnit(unit Unit) bool {
	for _, u := range g.units {
		if u == unit {
			return true
		}
	}
	return false
}

// 遍历grid上的单位
func (g *Grid) RangeUnits(f func(unit Unit) bool) {
	for _, u := range g.units {
		if !f(u) {
			break
		}
	}
}
