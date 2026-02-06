package geo

import (
	"math/rand"

	"github.com/GooLuck/WorldMap/internal/worldmap/config"
)

// 矩形
type Rectangle struct {
	Coord
	Width  int32
	Height int32
}

func NewRectangle(x, y, width, height int32) *Rectangle {
	return &Rectangle{
		Coord: Coord{
			X: x,
			Y: y,
		},
		Width:  width,
		Height: height,
	}
}

// 随机一个矩形内的坐标
func (r *Rectangle) RandomCoord() *Coord {
	return NewCoord(r.X+rand.Int31n(r.Width), r.Y+rand.Int31n(r.Height))
}

// 判断一个坐标是否在矩形内
func (r *Rectangle) IsCoordInRect(coord *Coord) bool {
	return coord.X >= r.X && coord.X < r.X+r.Width && coord.Y >= r.Y && coord.Y < r.Y+r.Height
}

// 返回矩形的中心点
func (r *Rectangle) Center() *Coord {
	return NewCoord(r.X+r.Width/2, r.Y+r.Height/2)
}

// 判断两个矩形是否相交
func (r *Rectangle) Intersects(other *Rectangle) bool {
	return r.X < other.X+other.Width && r.X+r.Width > other.X && r.Y < other.Y+other.Height && r.Y+r.Height > other.Y
}

func RectToGrid(mapConfig *config.MapConfig, rect *Rectangle) (minGridX, maxGridX, minGridY, maxGridY int32) {
	return
}
