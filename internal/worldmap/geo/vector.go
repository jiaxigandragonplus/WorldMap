package geo

import "math"

// 二维向量
type Vector2 struct {
	X float64
	Y float64
}

func NewVector2(start, end *Coord) *Vector2 {
	return &Vector2{
		X: float64(end.X - start.X),
		Y: float64(end.Y - start.Y),
	}
}

// 获取向量长度
func (v *Vector2) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v *Vector2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2) Add(vec *Vector2) *Vector2 {
	return &Vector2{
		X: v.X + vec.X,
		Y: v.Y + vec.Y,
	}
}

func (v *Vector2) Sub(vec *Vector2) *Vector2 {
	return &Vector2{
		X: v.X - vec.X,
		Y: v.Y - vec.Y,
	}
}

// 点乘
func (v *Vector2) Dot(vec *Vector2) float64 {
	return v.X*vec.X + v.Y*vec.Y
}

// 叉乘
func (v *Vector2) Cross(vec *Vector2) float64 {
	return v.X*vec.Y - v.Y*vec.X
}
