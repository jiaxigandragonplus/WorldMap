package worldmap

import (
	"container/heap"
	"math"

	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// HexGrid 六边形网格
type HexGrid struct {
	coord  *geo.HexCoord  // 六边形坐标
	units  map[int64]Unit // 网格上的单位 (id -> Unit)
	layout *geo.HexLayout // 六边形布局信息
}

// NewHexGrid 创建新的六边形网格
func NewHexGrid(coord *geo.HexCoord, layout *geo.HexLayout) *HexGrid {
	return &HexGrid{
		coord:  coord,
		units:  make(map[int64]Unit),
		layout: layout,
	}
}

// GetCoord 获取六边形坐标
func (g *HexGrid) GetCoord() *geo.HexCoord {
	return g.coord
}

// AddUnit 添加单位
func (g *HexGrid) AddUnit(unit Unit) {
	if unit != nil {
		g.units[unit.GetId()] = unit
	}
}

// RemoveUnit 移除单位
func (g *HexGrid) RemoveUnit(unit Unit) {
	if unit != nil {
		delete(g.units, unit.GetId())
	}
}

// RemoveUnitById 根据 ID 移除单位
func (g *HexGrid) RemoveUnitById(id int64) {
	delete(g.units, id)
}

// IsExistUnit 检查单位是否存在
func (g *HexGrid) IsExistUnit(unit Unit) bool {
	if unit == nil {
		return false
	}
	_, exists := g.units[unit.GetId()]
	return exists
}

// IsExistUnitById 根据 ID 检查单位是否存在
func (g *HexGrid) IsExistUnitById(id int64) bool {
	_, exists := g.units[id]
	return exists
}

// GetUnits 获取所有单位
func (g *HexGrid) GetUnits() []Unit {
	units := make([]Unit, 0, len(g.units))
	for _, u := range g.units {
		units = append(units, u)
	}
	return units
}

// GetUnitById 根据 ID 获取单位
func (g *HexGrid) GetUnitById(id int64) Unit {
	return g.units[id]
}

// GetUnitCount 获取单位数量
func (g *HexGrid) GetUnitCount() int {
	return len(g.units)
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

// GetUnitsByType 根据类型获取单位
func (g *HexGrid) GetUnitsByType(unitType MapUnitType) []Unit {
	units := make([]Unit, 0)
	for _, u := range g.units {
		if u.GetType() == unitType {
			units = append(units, u)
		}
	}
	return units
}

// HexGridManager 六边形网格管理器
type HexGridManager struct {
	layout  *geo.HexLayout      // 六边形布局
	qCount  int32               // q 方向网格数量
	rCount  int32               // r 方向网格数量
	bounds  *geo.HexRectangle   // 边界范围
	grids   map[uint64]*HexGrid // 所有网格 (key = hash(q, r))
	mapSize *config.MapSize     // 地图大小配置
}

// NewHexGridManager 创建新的六边形网格管理器
// mapConfig: 地图配置，hexRadius: 六边形边长
func NewHexGridManager(mapSize *config.MapSize, hexRadius float64, isPointy bool) *HexGridManager {
	// 计算能容纳多少个六边形
	qCount, rCount := geo.CalculateGridSize(float64(mapSize.Width), float64(mapSize.Height), hexRadius, isPointy)

	// 创建布局，原点在 (0,0)
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

// hashHex 计算六边形坐标的哈希值，将两个 int32 编码为 uint64
func hashHex(q, r int32) uint64 {
	return (uint64(q) << 32) | uint64(uint32(r))
}

// absInt32 返回 int32 的绝对值
func absInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
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

// GetQCount 获取 q 方向网格数量
func (hgm *HexGridManager) GetQCount() int32 {
	return hgm.qCount
}

// GetRCount 获取 r 方向网格数量
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

// GetNeighborGridsByCoord 根据坐标获取相邻网格
func (hgm *HexGridManager) GetNeighborGridsByCoord(hex *geo.HexCoord) []*HexGrid {
	neighbors := make([]*HexGrid, 0, 6)
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
// f: 遍历回调函数，返回 false 停止遍历
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

// GetHexGridsInRadius 获取指定半径范围内的所有六边形网格
func (hgm *HexGridManager) GetHexGridsInRadius(center *geo.HexCoord, radius int32) []*HexGrid {
	result := make([]*HexGrid, 0)
	for q := -radius; q <= radius; q++ {
		for r := -radius; r <= radius; r++ {
			if absInt32(q+r) > radius {
				continue
			}
			hex := geo.NewHexCoord(center.Q+q, center.R+r)
			if grid := hgm.GetGrid(hex); grid != nil {
				result = append(result, grid)
			}
		}
	}
	return result
}

// GetUnitsInRadius 获取指定半径范围内的所有单位
func (hgm *HexGridManager) GetUnitsInRadius(center *geo.HexCoord, radius int32) []Unit {
	result := make([]Unit, 0)
	grids := hgm.GetHexGridsInRadius(center, radius)
	for _, grid := range grids {
		result = append(result, grid.GetUnits()...)
	}
	return result
}

// GetUnitsInRadiusByWorld 获取世界坐标指定半径范围内的所有单位
func (hgm *HexGridManager) GetUnitsInRadiusByWorld(worldX, worldY float64, radius int32) []Unit {
	q, r := hgm.layout.WorldToHex(worldX, worldY)
	center := geo.RoundToHex(q, r)
	return hgm.GetUnitsInRadius(center, radius)
}

// GetDistance 计算两个六边形坐标之间的距离
func (hgm *HexGridManager) GetDistance(hex1, hex2 *geo.HexCoord) int32 {
	return hex1.DistanceTo(hex2)
}

// GetWorldDistance 计算两个世界坐标之间的距离（六边形步数）
func (hgm *HexGridManager) GetWorldDistance(x1, y1, x2, y2 float64) int32 {
	q1, r1 := hgm.layout.WorldToHex(x1, y1)
	q2, r2 := hgm.layout.WorldToHex(x2, y2)
	hex1 := geo.RoundToHex(q1, r1)
	hex2 := geo.RoundToHex(q2, r2)
	return hex1.DistanceTo(hex2)
}

// pathNode A*路径查找节点
type pathNode struct {
	hex    *geo.HexCoord
	gCost  int32 // 从起点到当前节点的实际成本
	hCost  int32 // 从当前节点到终点的估计成本
	fCost  int32 // 总成本 (gCost + hCost)
	parent *pathNode
	index  int // heap 需要的索引
}

// pathNodeHeap 实现 heap.Interface 用于优先队列
type pathNodeHeap []*pathNode

func (h pathNodeHeap) Len() int           { return len(h) }
func (h pathNodeHeap) Less(i, j int) bool { return h[i].fCost < h[j].fCost }
func (h pathNodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *pathNodeHeap) Push(x interface{}) {
	n := x.(*pathNode)
	n.index = len(*h)
	*h = append(*h, n)
}
func (h *pathNodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	old[n-1] = nil // 避免内存泄漏
	*h = old[:n-1]
	return node
}

// TerrainCostFunc 地形成本函数类型
type TerrainCostFunc func(hex *geo.HexCoord) int32

// FindPath A*算法查找路径
// start, end: 起点和终点六边形坐标
// terrainCost: 地形成本函数（可选，nil 表示默认成本为 1）
// 返回：路径上的六边形坐标列表（包含起点和终点）
func (hgm *HexGridManager) FindPath(start, end *geo.HexCoord, terrainCost TerrainCostFunc) []*geo.HexCoord {
	if !hgm.bounds.Contains(start) || !hgm.bounds.Contains(end) {
		return nil
	}

	// 如果起点和终点相同，直接返回
	if start.Equal(end) {
		return []*geo.HexCoord{start}
	}

	// 初始化地形成本函数
	if terrainCost == nil {
		terrainCost = func(hex *geo.HexCoord) int32 { return 1 }
	}

	// openSet: 待处理节点集合（优先队列）
	openSet := &pathNodeHeap{}
	heap.Init(openSet)

	// closedSet: 已处理节点集合
	closedSet := make(map[uint64]bool)

	// nodes: 节点映射
	nodes := make(map[uint64]*pathNode)

	// 创建起点节点
	startNode := &pathNode{
		hex:   start,
		gCost: 0,
		hCost: start.DistanceTo(end),
	}
	startNode.fCost = startNode.gCost + startNode.hCost
	nodes[start.Hash()] = startNode
	heap.Push(openSet, startNode)

	for openSet.Len() > 0 {
		// 取出 fCost 最小的节点
		current := heap.Pop(openSet).(*pathNode)
		currentHash := current.hex.Hash()

		// 如果已处理过，跳过
		if closedSet[currentHash] {
			continue
		}
		closedSet[currentHash] = true

		// 如果到达终点，重建路径
		if current.hex.Equal(end) {
			return hgm.reconstructPath(current)
		}

		// 遍历邻居
		for _, neighborHex := range current.hex.GetAllNeighbors() {
			if !hgm.bounds.Contains(neighborHex) {
				continue
			}

			neighborHash := neighborHex.Hash()
			if closedSet[neighborHash] {
				continue
			}

			// 计算新的 gCost
			newGCost := current.gCost + terrainCost(neighborHex)

			// 如果找到更好的路径或这是新节点
			existingNode, exists := nodes[neighborHash]
			if !exists || newGCost < existingNode.gCost {
				if !exists {
					existingNode = &pathNode{
						hex: neighborHex,
					}
					nodes[neighborHash] = existingNode
				}
				existingNode.gCost = newGCost
				existingNode.hCost = neighborHex.DistanceTo(end)
				existingNode.fCost = existingNode.gCost + existingNode.hCost
				existingNode.parent = current

				if !closedSet[neighborHash] {
					heap.Push(openSet, existingNode)
				}
			}
		}
	}

	// 没有找到路径
	return nil
}

// reconstructPath 重建路径
func (hgm *HexGridManager) reconstructPath(endNode *pathNode) []*geo.HexCoord {
	path := make([]*geo.HexCoord, 0)
	for node := endNode; node != nil; node = node.parent {
		path = append(path, node.hex)
	}
	// 反转路径
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// GetHexesInLine 获取两个六边形坐标之间的直线路径（Bresenham 算法）
func (hgm *HexGridManager) GetHexesInLine(start, end *geo.HexCoord) []*geo.HexCoord {
	n := start.DistanceTo(end)
	if n == 0 {
		return []*geo.HexCoord{start}
	}

	result := make([]*geo.HexCoord, 0, n+1)
	for i := int32(0); i <= n; i++ {
		t := float64(i) / float64(n)
		q := float64(start.Q) + (float64(end.Q)-float64(start.Q))*t
		r := float64(start.R) + (float64(end.R)-float64(start.R))*t
		result = append(result, geo.RoundToHex(q, r))
	}
	return result
}

// GetVisionRange 获取视野范围内的六边形（指定距离）
func (hgm *HexGridManager) GetVisionRange(center *geo.HexCoord, visionRange int32) []*HexGrid {
	return hgm.GetHexGridsInRadius(center, visionRange)
}

// GetVisibleUnits 获取视野范围内的所有单位
func (hgm *HexGridManager) GetVisibleUnits(center *geo.HexCoord, visionRange int32) []Unit {
	return hgm.GetUnitsInRadius(center, visionRange)
}

// IsVisible 检查目标是否在视野范围内（无障碍）
func (hgm *HexGridManager) IsVisible(from, to *geo.HexCoord, visionRange int32) bool {
	// 检查距离
	if from.DistanceTo(to) > visionRange {
		return false
	}

	// 检查直线路径上是否有障碍物
	line := hgm.GetHexesInLine(from, to)
	for _, hex := range line[1 : len(line)-1] { // 跳过起点和终点
		if grid := hgm.GetGrid(hex); grid != nil {
			// 检查是否有障碍物
			obstacles := grid.GetUnitsByType(MapUnitType_Obstacle)
			if len(obstacles) > 0 {
				return false
			}
		}
	}
	return true
}

// MoveUnit 移动单位到新的六边形
func (hgm *HexGridManager) MoveUnit(unit Unit, from, to *geo.HexCoord) bool {
	if !hgm.bounds.Contains(to) {
		return false
	}

	// 从原位置移除
	fromGrid := hgm.GetGrid(from)
	if fromGrid != nil {
		fromGrid.RemoveUnit(unit)
	}

	// 添加到新位置
	toGrid := hgm.GetGrid(to)
	if toGrid != nil {
		toGrid.AddUnit(unit)
		return true
	}
	return false
}

// GetAllUnits 获取地图上所有单位
func (hgm *HexGridManager) GetAllUnits() []Unit {
	result := make([]Unit, 0)
	for _, grid := range hgm.grids {
		result = append(result, grid.GetUnits()...)
	}
	return result
}

// GetUnitsByType 根据类型获取所有单位
func (hgm *HexGridManager) GetUnitsByType(unitType MapUnitType) []Unit {
	result := make([]Unit, 0)
	for _, grid := range hgm.grids {
		result = append(result, grid.GetUnitsByType(unitType)...)
	}
	return result
}

// RangeAllGrids 遍历所有网格
func (hgm *HexGridManager) RangeAllGrids(f func(grid *HexGrid) bool) {
	for _, grid := range hgm.grids {
		if !f(grid) {
			break
		}
	}
}

// GetEmptyGrids 获取所有空网格
func (hgm *HexGridManager) GetEmptyGrids() []*HexGrid {
	result := make([]*HexGrid, 0)
	for _, grid := range hgm.grids {
		if grid.GetUnitCount() == 0 {
			result = append(result, grid)
		}
	}
	return result
}

// GetRandomEmptyGrid 获取随机空网格
func (hgm *HexGridManager) GetRandomEmptyGrid() *HexGrid {
	emptyGrids := hgm.GetEmptyGrids()
	if len(emptyGrids) == 0 {
		return nil
	}
	return emptyGrids[0] // 简化实现，实际应该随机选择
}
