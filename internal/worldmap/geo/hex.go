package geo

import (
	"fmt"
	"math"
)

// HexCoord 六边形坐标（轴向坐标系统）
// 使用 q, r 两个坐标表示，第三个坐标 s = -q - r 恒成立
type HexCoord struct {
	Q int32 // q轴坐标
	R int32 // r轴坐标
}

// NewHexCoord 创建新的六边形坐标
func NewHexCoord(q, r int32) *HexCoord {
	return &HexCoord{
		Q: q,
		R: r,
	}
}

// S 获取s坐标，s = -q - r
func (h *HexCoord) S() int32 {
	return -h.Q - h.R
}

// Add 坐标相加
func (h *HexCoord) Add(other *HexCoord) *HexCoord {
	return NewHexCoord(h.Q+other.Q, h.R+other.R)
}

// Sub 坐标相减
func (h *HexCoord) Sub(other *HexCoord) *HexCoord {
	return NewHexCoord(h.Q-other.Q, h.R-other.R)
}

// Multiply 坐标乘以标量
func (h *HexCoord) Multiply(k int32) *HexCoord {
	return NewHexCoord(h.Q*k, h.R*k)
}

// 六个方向的坐标增量
var hexDirections = [6]*HexCoord{
	{1, 0},  // 右
	{1, -1}, // 右上
	{0, -1}, // 左上
	{-1, 0}, // 左
	{-1, 1}, // 左下
	{0, 1},  // 右下
}

// GetNeighbor 获取指定方向的邻居坐标
func (h *HexCoord) GetNeighbor(direction int) *HexCoord {
	return h.Add(hexDirections[direction])
}

// GetAllNeighbors 获取所有六个邻居坐标
func (h *HexCoord) GetAllNeighbors() []*HexCoord {
	neighbors := make([]*HexCoord, 6)
	for i := 0; i < 6; i++ {
		neighbors[i] = h.GetNeighbor(i)
	}
	return neighbors
}

// DistanceTo 计算两个六边形之间的距离（步数）
func (h *HexCoord) DistanceTo(other *HexCoord) int32 {
	dq := abs(h.Q - other.Q)
	dr := abs(h.R - other.R)
	ds := abs(h.S() - other.S())
	return (dq + dr + ds) / 2
}

// Equal 检查两个六边形坐标是否相等
func (h *HexCoord) Equal(other *HexCoord) bool {
	if other == nil {
		return false
	}
	return h.Q == other.Q && h.R == other.R
}

// Hash 计算六边形坐标的哈希值
func (h *HexCoord) Hash() uint64 {
	return hashHex(h.Q, h.R)
}

// Clone 克隆六边形坐标
func (h *HexCoord) Clone() *HexCoord {
	return NewHexCoord(h.Q, h.R)
}

// String 返回六边形坐标的字符串表示
func (h *HexCoord) String() string {
	return fmt.Sprintf("Hex(%d, %d, %d)", h.Q, h.R, h.S())
}

// abs 绝对值函数
func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// HexLayout 六边形布局配置
type HexLayout struct {
	Radius   float64 // 六边形半径（边长）
	OriginX  float64 // 原点X偏移
	OriginY  float64 // 原点Y偏移
	IsPointy bool    // 是否是 pointy 朝向（顶点朝上，否则平边朝上）
}

// NewHexLayout 创建六边形布局
func NewHexLayout(radius float64, originX, originY float64, isPointy bool) *HexLayout {
	return &HexLayout{
		Radius:   radius,
		OriginX:  originX,
		OriginY:  originY,
		IsPointy: isPointy,
	}
}

// HexToWorld 将六边形坐标转换为世界坐标（中心坐标）
func (l *HexLayout) HexToWorld(h *HexCoord) (float64, float64) {
	var x, y float64
	if l.IsPointy {
		// 顶点朝上
		x = l.Radius * (math.Sqrt(3)*float64(h.Q) + math.Sqrt(3)/2*float64(h.R))
		y = l.Radius * (3.0 / 2 * float64(h.R))
	} else {
		// 平边朝上
		x = l.Radius * (3.0 / 2 * float64(h.Q))
		y = l.Radius * (math.Sqrt(3)/2*float64(h.Q) + math.Sqrt(3)*float64(h.R))
	}
	return x + l.OriginX, y + l.OriginY
}

// WorldToHex 将世界坐标转换为六边形坐标（浮点数，需要取整）
func (l *HexLayout) WorldToHex(worldX, worldY float64) (float64, float64) {
	vx := worldX - l.OriginX
	vy := worldY - l.OriginY
	var q, r float64
	if l.IsPointy {
		q = (math.Sqrt(3)/3*vx - 1.0/3*vy) / l.Radius
		r = (2.0 / 3 * vy) / l.Radius
	} else {
		q = (2.0 / 3 * vx) / l.Radius
		r = (-1.0/3*vx + math.Sqrt(3)/3*vy) / l.Radius
	}
	return q, r
}

// RoundToHex 将浮点坐标取整为六边形整数坐标
func RoundToHex(q, r float64) *HexCoord {
	s := -q - r
	qi := int32(math.Round(q))
	ri := int32(math.Round(r))
	si := int32(math.Round(s))

	qDiff := math.Abs(float64(qi) - q)
	rDiff := math.Abs(float64(ri) - r)
	sDiff := math.Abs(float64(si) - s)

	if qDiff > rDiff && qDiff > sDiff {
		qi = -ri - si
	} else if rDiff > sDiff {
		ri = -qi - si
	} else {
		si = -qi - ri
	}

	return NewHexCoord(qi, ri)
}

// GetHexCorners 获取六边形六个顶点的世界坐标
func (l *HexLayout) GetHexCorners(h *HexCoord) [][2]float64 {
	cx, cy := l.HexToWorld(h)
	corners := make([][2]float64, 6)
	for i := 0; i < 6; i++ {
		angle := 2 * math.Pi * (float64(i) + l.getAngleOffset()) / 6
		cornerX := cx + l.Radius*math.Cos(angle)
		cornerY := cy + l.Radius*math.Sin(angle)
		corners[i] = [2]float64{cornerX, cornerY}
	}
	return corners
}

// getAngleOffset 获取起始角度偏移
func (l *HexLayout) getAngleOffset() float64 {
	if l.IsPointy {
		return 0.0
	}
	return 0.5
}

// HexRectangle 六边形地图范围
type HexRectangle struct {
	MinQ int32
	MaxQ int32
	MinR int32
	MaxR int32
}

// NewHexRectangle 创建六边形范围
func NewHexRectangle(minQ, maxQ, minR, maxR int32) *HexRectangle {
	return &HexRectangle{
		MinQ: minQ,
		MaxQ: maxQ,
		MinR: minR,
		MaxR: maxR,
	}
}

// hashHex 计算六边形坐标的哈希值，将两个 int32 编码为 uint64
func hashHex(q, r int32) uint64 {
	return (uint64(q) << 32) | uint64(uint32(r))
}

// Contains 检查坐标是否在范围内
func (hr *HexRectangle) Contains(h *HexCoord) bool {
	return h.Q >= hr.MinQ && h.Q <= hr.MaxQ && h.R >= hr.MinR && h.R <= hr.MaxR
}

// CalculateGridSize 计算地图边界能容纳的六边形数量
// mapWidth, mapHeight: 地图大小，单位和半径相同
// 返回值: q方向数量, r方向数量
func CalculateGridSize(mapWidth, mapHeight, radius float64, isPointy bool) (qCount, rCount int32) {
	var qf, rf float64
	if isPointy {
		// pointy 朝向（顶点朝上）
		hexWidth := radius * math.Sqrt(3)
		hexHeight := radius * 1.5
		qf = mapWidth / hexWidth
		rf = mapHeight / hexHeight
	} else {
		// flat 朝向（平边朝上）
		hexWidth := radius * 1.5
		hexHeight := radius * math.Sqrt(3)
		qf = mapWidth / hexWidth
		rf = mapHeight / hexHeight
	}
	return int32(math.Ceil(qf)), int32(math.Ceil(rf))
}

// GetBoundingRect 获取六边形地图的外接矩形
func (l *HexLayout) GetBoundingRect(qCount, rCount int32) (minX, minY, maxX, maxY float64) {
	if l.IsPointy {
		// pointy
		w := l.Radius * math.Sqrt(3) * float64(qCount)
		h := l.Radius * (1.5*float64(rCount) + 0.5)
		return 0, 0, w, h
	}
	// flat
	w := l.Radius * (1.5*float64(qCount) + 0.5)
	h := l.Radius * math.Sqrt(3) * float64(rCount)
	return 0, 0, w, h
}
