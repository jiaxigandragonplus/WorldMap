package geo

// Coord 坐标
type Coord struct {
	X int32
	Y int32
}

func NewCoord(x, y int32) *Coord {
	return &Coord{
		X: x,
		Y: y,
	}
}

// 相加
func (c *Coord) Add(add *Coord) *Coord {
	c.X += add.X
	c.Y += add.Y
	return c
}

func (c *Coord) Sub(sub *Coord) *Coord {
	c.X -= sub.X
	c.Y -= sub.Y
	return c
}

// 平移
func (c *Coord) Translate(vec *Vector2) *Coord {
	c.X += int32(vec.X)
	c.Y += int32(vec.Y)
	return c
}
