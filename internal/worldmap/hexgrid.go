package worldmap

import (
	"math"

	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// HexGrid 六边形网格
type HexGrid struct {
	coord  *geo.HexCoord  // 六边形坐标
	units  []Unit         // 网格上的单位
	layout *geo.HexLayout // 六边形布局信息
}

// NewHexGrid 创建新的六边形网格
func NewHexGrid(coord *geo.HexCoord, layout *geo.HexLayout) *HexGrid {
	return &HexGrid{
		coord:  coord,
		units:  make([]Unit, 0),
		layout: layout,
	}
}

// GetCoord 获取六边形坐标
func (g *HexGrid) GetCoord() *geo.HexCoord {
	return g.coord
}

// AddUnit 添加单位
func (g *HexGrid) AddUnit(unit Unit) {
	g.units = append(g.units, unit)
}

// RemoveUnit 移除单位
func (g *HexGrid) RemoveUnit(unit Unit) {
	for i, u := range g.units {
		if u == unit {
			g.units = append(g.units[:i], g.units[i+1:]...)
			return
		}
	}
}

// IsExistUnit 检查单位是否存在
func (g *HexGrid) IsExistUnit(unit Unit) bool {
	for _, u := range g.units {
		if u == unit {
			return true
		}
	}
	return false
}

// GetUnits 获取所有单位
func (g *HexGrid) GetUnits() []Unit {
	return g.units
}

// RangeUnits 遍历所有单位
func (g *HexGrid) RangeUnits(f func(unit Unit) bool) {
	for _, u := range g.units {
		if !f(u) {
			break
		}
	}
}

// GetCenterWorld 获取中心世界坐标
func (g *HexGrid) GetCenterWorld() (x, y float64) {
	return g.layout.HexToWorld(g.coord)
}

// HexGridManager 六边形网格管理器
type HexGridManager struct {
	layout  *geo.HexLayout      // 六边形布局
	qCount  int32               // q方向网格数量
	rCount  int32               // r方向网格数量
	bounds  *geo.HexRectangle   // 边界范围
	grids   map[uint64]*HexGrid // 所有网格 (key = hash(q, r))
	mapSize *config.MapSize     // 地图大小配置
}

// 计算哈希值，将两个int32编码为uint64
func hashHex(q, r int32) uint64 {
	return (uint64(q) << 32) | uint64(uint32(r))
}

// 从哈希值解码出q和r
func unhashHex(hash uint64) (q, r int32) {
	q = int32(hash >> 32)
	r = int32(uint32(hash & 0xFFFFFFFF))
	return
}

// NewHexGridManager 创建新的六边形网格管理器
// mapConfig: 地图配置，hexRadius: 六边形边长
func NewHexGridManager(mapSize *config.MapSize, hexRadius float64, isPointy bool) *HexGridManager {
	// 计算能容纳多少个六边形
	qCount, rCount := geo.CalculateGridSize(float64(mapSize.Width), float64(mapSize.Height), hexRadius, isPointy)

	// 创建布局，原点在(0,0)
	layout := geo.NewHexLayout(hexRadius, 0, 0, isPointy)

	// 计算边界范围
	minQ := int32(0)
	maxQ := qCount - 1
	minR := int32(0)
	maxR := rCount - 1

	hgm := &HexGridManager{
		layout:  layout,
		qCount:  qCount,
		rCount:  rCount,
		bounds:  geo.NewHexRectangle(minQ, maxQ, minR, maxR),
		grids:   make(map[uint64]*HexGrid),
		mapSize: mapSize,
	}

	// 预分配所有网格
	hgm.initializeAllGrids()

	return hgm
}

// initializeAllGrids 初始化所有网格
func (hgm *HexGridManager) initializeAllGrids() {
	for q := hgm.bounds.MinQ; q <= hgm.bounds.MaxQ; q++ {
		for r := hgm.bounds.MinR; r <= hgm.bounds.MaxR; r++ {
			hex := geo.NewHexCoord(q, r)
			grid := NewHexGrid(hex, hgm.layout)
			hash := hashHex(q, r)
			hgm.grids[hash] = grid
		}
	}
}

// GetGrid 获取指定六边形坐标的网格
func (hgm *HexGridManager) GetGrid(hex *geo.HexCoord) *HexGrid {
	if !hgm.bounds.Contains(hex) {
		return nil
	}
	hash := hashHex(hex.Q, hex.R)
	return hgm.grids[hash]
}

// GetGridByWorld 获取世界坐标所在的六边形网格
func (hgm *HexGridManager) GetGridByWorld(worldX, worldY float64) *HexGrid {
	q, r := hgm.layout.WorldToHex(worldX, worldY)
	hex := geo.RoundToHex(q, r)
	return hgm.GetGrid(hex)
}

// GetGridCount 获取网格总数
func (hgm *HexGridManager) GetGridCount() int32 {
	return hgm.qCount * hgm.rCount
}

// GetQCount 获取q方向网格数量
func (hgm *HexGridManager) GetQCount() int32 {
	return hgm.qCount
}

// GetRCount 获取r方向网格数量
func (hgm *HexGridManager) GetRCount() int32 {
	return hgm.rCount
}

// Contains 检查六边形坐标是否在地图范围内
func (hgm *HexGridManager) Contains(hex *geo.HexCoord) bool {
	return hgm.bounds.Contains(hex)
}

// AddUnitToGrid 将单位添加到指定六边形
func (hgm *HexGridManager) AddUnitToGrid(unit Unit, hex *geo.HexCoord) bool {
	grid := hgm.GetGrid(hex)
	if grid == nil {
		return false
	}
	grid.AddUnit(unit)
	return true
}

// AddUnitToWorld 将单位添加到世界坐标所在六边形
func (hgm *HexGridManager) AddUnitToWorld(unit Unit, worldX, worldY float64) bool {
	grid := hgm.GetGridByWorld(worldX, worldY)
	if grid == nil {
		return false
	}
	grid.AddUnit(unit)
	return true
}

// RemoveUnitFromGrid 从指定六边形移除单位
func (hgm *HexGridManager) RemoveUnitFromGrid(unit Unit, hex *geo.HexCoord) bool {
	grid := hgm.GetGrid(hex)
	if grid == nil {
		return false
	}
	grid.RemoveUnit(unit)
	return true
}

// RemoveUnitFromWorld 从世界坐标所在六边形移除单位
func (hgm *HexGridManager) RemoveUnitFromWorld(unit Unit, worldX, worldY float64) bool {
	grid := hgm.GetGridByWorld(worldX, worldY)
	if grid == nil {
		return false
	}
	grid.RemoveUnit(unit)
	return true
}

// GetNeighborGrids 获取相邻网格
func (hgm *HexGridManager) GetNeighborGrids(grid *HexGrid) []*HexGrid {
	neighbors := make([]*HexGrid, 0, 6)
	hex := grid.GetCoord()
	for _, neighborHex := range hex.GetAllNeighbors() {
		if neighborGrid := hgm.GetGrid(neighborHex); neighborGrid != nil {
			neighbors = append(neighbors, neighborGrid)
		}
	}
	return neighbors
}

// RangeInRect 遍历矩形范围内的所有六边形网格
// rect: 世界坐标矩形范围
// includeOffScreen: 是否包含不在屏幕内的网格
// f: 遍历回调函数，返回false停止遍历
func (hgm *HexGridManager) RangeInRect(minX, minY, maxX, maxY float64, f func(grid *HexGrid) bool) {
	// 计算矩形范围覆盖的六边形坐标范围
	qMin, rMin := hgm.layout.WorldToHex(minX, minY)
	qMax, rMax := hgm.layout.WorldToHex(maxX, maxY)

	// 取整得到范围边界
	var minQ, maxQ, minR, maxR int32
	if qMin < qMax {
		minQ = int32(math.Floor(qMin))
		maxQ = int32(math.Ceil(qMax))
	} else {
		minQ = int32(math.Floor(qMax))
		maxQ = int32(math.Ceil(qMin))
	}
	if rMin < rMax {
		minR = int32(math.Floor(rMin))
		maxR = int32(math.Ceil(rMax))
	} else {
		minR = int32(math.Floor(rMax))
		maxR = int32(math.Ceil(rMin))
	}

	// 限制在地图范围内
	if minQ < hgm.bounds.MinQ {
		minQ = hgm.bounds.MinQ
	}
	if maxQ > hgm.bounds.MaxQ {
		maxQ = hgm.bounds.MaxQ
	}
	if minR < hgm.bounds.MinR {
		minR = hgm.bounds.MinR
	}
	if maxR > hgm.bounds.MaxR {
		maxR = hgm.bounds.MaxR
	}

	// 遍历所有覆盖的六边形
	for q := minQ; q <= maxQ; q++ {
		for r := minR; r <= maxR; r++ {
			hex := geo.NewHexCoord(q, r)
			grid := hgm.GetGrid(hex)
			if grid != nil {
				// 检查是否至少一个顶点在矩形内
				cx, cy := grid.GetCenterWorld()
				// 使用半径近似判断，如果中心点都不在范围内就跳过
				if cx+hgm.layout.Radius < minX || cx-hgm.layout.Radius > maxX ||
					cy+hgm.layout.Radius < minY || cy-hgm.layout.Radius > maxY {
					continue
				}
				if !f(grid) {
					return
				}
			}
		}
	}
}

// RangeUnitsInRect 遍历矩形范围内的所有单位
func (hgm *HexGridManager) RangeUnitsInRect(minX, minY, maxX, maxY float64, f func(unit Unit) bool) {
	hgm.RangeInRect(minX, minY, maxX, maxY, func(grid *HexGrid) bool {
		for _, unit := range grid.GetUnits() {
			if !f(unit) {
				return false
			}
		}
		return true
	})
}

// GetLayout 获取布局信息
func (hgm *HexGridManager) GetLayout() *geo.HexLayout {
	return hgm.layout
}

// GetBounds 获取网格边界
func (hgm *HexGridManager) GetBounds() *geo.HexRectangle {
	return hgm.bounds
}
