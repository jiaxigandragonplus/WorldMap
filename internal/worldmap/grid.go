package worldmap

import (
	"math"

	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

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

func MaxGridXY(mapSize *config.MapSize) (maxGridX, maxGridY int32) {
	return (mapSize.Width - 1) / mapSize.GridWidth, (mapSize.Height - 1) / mapSize.GridHeight
}

// rect转grid
func RectToGrid(mapSize *config.MapSize, rect *geo.Rectangle) (minGridX, maxGridX, minGridY, maxGridY int32) {
	maxGridXTmp, maxGridYTmp := MaxGridXY(mapSize)
	minGridX = int32(math.Max(float64(rect.X)/float64(mapSize.Width), 0))
	maxGridX = int32(math.Min(float64(rect.X+rect.Width)/float64(mapSize.Width), float64(maxGridXTmp)))
	minGridY = int32(math.Max(float64(rect.Y)/float64(mapSize.Height), 0))
	maxGridY = int32(math.Min(float64(rect.Y+rect.Height)/float64(mapSize.Height), float64(maxGridYTmp)))
	return
}
