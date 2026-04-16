package worldmap

import (
	"math"
	"testing"

	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// TestHexCoord 测试六边形坐标基本功能
func TestHexCoord(t *testing.T) {
	// 测试创建坐标
	hex := geo.NewHexCoord(3, 5)
	if hex.Q != 3 || hex.R != 5 {
		t.Errorf("HexCoord 创建失败：期望 (3, 5), 得到 (%d, %d)", hex.Q, hex.R)
	}

	// 测试 S 坐标
	s := hex.S()
	if s != -8 {
		t.Errorf("S 坐标计算错误：期望 -8, 得到 %d", s)
	}

	// 测试坐标相加
	hex2 := geo.NewHexCoord(1, -1)
	result := hex.Add(hex2)
	if result.Q != 4 || result.R != 4 {
		t.Errorf("坐标相加错误：期望 (4, 4), 得到 (%d, %d)", result.Q, result.R)
	}

	// 测试坐标相减
	result = hex.Sub(hex2)
	if result.Q != 2 || result.R != 6 {
		t.Errorf("坐标相减错误：期望 (2, 6), 得到 (%d, %d)", result.Q, result.R)
	}

	// 测试距离计算
	hex3 := geo.NewHexCoord(5, 5)
	distance := hex.DistanceTo(hex3)
	if distance != 2 {
		t.Errorf("距离计算错误：期望 2, 得到 %d", distance)
	}

	// 测试 Equal 方法
	hex4 := geo.NewHexCoord(3, 5)
	if !hex.Equal(hex4) {
		t.Error("Equal 方法错误：期望相等")
	}
	if hex.Equal(hex3) {
		t.Error("Equal 方法错误：期望不相等")
	}

	// 测试 Hash 方法
	hash1 := hex.Hash()
	hash2 := hex4.Hash()
	if hash1 != hash2 {
		t.Error("Hash 方法错误：相等的坐标应该有相同的哈希值")
	}

	// 测试 Clone 方法
	cloned := hex.Clone()
	if !hex.Equal(cloned) {
		t.Error("Clone 方法错误：克隆的坐标应该相等")
	}
	cloned.Q = 100
	if hex.Equal(cloned) {
		t.Error("Clone 方法错误：克隆应该是深拷贝")
	}
}

// TestHexNeighbors 测试六边形邻居坐标
func TestHexNeighbors(t *testing.T) {
	hex := geo.NewHexCoord(0, 0)
	neighbors := hex.GetAllNeighbors()

	if len(neighbors) != 6 {
		t.Errorf("邻居数量错误：期望 6, 得到 %d", len(neighbors))
	}

	// 验证邻居坐标
	expected := [][2]int32{
		{1, 0},  // 右
		{1, -1}, // 右上
		{0, -1}, // 左上
		{-1, 0}, // 左
		{-1, 1}, // 左下
		{0, 1},  // 右下
	}

	for i, neighbor := range neighbors {
		if neighbor.Q != expected[i][0] || neighbor.R != expected[i][1] {
			t.Errorf("邻居 %d 坐标错误：期望 (%d, %d), 得到 (%d, %d)",
				i, expected[i][0], expected[i][1], neighbor.Q, neighbor.R)
		}
	}
}

// TestHexLayout 测试六边形布局转换
func TestHexLayout(t *testing.T) {
	radius := 10.0
	layout := geo.NewHexLayout(radius, 0, 0, true) // pointy 朝向

	// 测试 HexToWorld
	hex := geo.NewHexCoord(1, 1)
	worldX, worldY := layout.HexToWorld(hex)

	// 验证转换结果（允许小误差）
	expectedX := radius * (1.732*1 + 0.866*1) // sqrt(3)*Q + sqrt(3)/2*R
	expectedY := radius * (1.5 * 1)           // 3/2 * R

	if absFloat(worldX-expectedX) > 0.1 {
		t.Errorf("HexToWorld X 错误：期望 %.2f, 得到 %.2f", expectedX, worldX)
	}
	if absFloat(worldY-expectedY) > 0.1 {
		t.Errorf("HexToWorld Y 错误：期望 %.2f, 得到 %.2f", expectedY, worldY)
	}

	// 测试 WorldToHex（往返转换）
	q, r := layout.WorldToHex(worldX, worldY)
	rounded := geo.RoundToHex(q, r)
	if rounded.Q != hex.Q || rounded.R != hex.R {
		t.Errorf("WorldToHex 往返转换错误：期望 (%d, %d), 得到 (%d, %d)",
			hex.Q, hex.R, rounded.Q, rounded.R)
	}
}

// TestHexCorners 测试六边形顶点计算
func TestHexCorners(t *testing.T) {
	radius := 10.0
	layout := geo.NewHexLayout(radius, 0, 0, true)
	hex := geo.NewHexCoord(0, 0)

	corners := layout.GetHexCorners(hex)

	if len(corners) != 6 {
		t.Errorf("顶点数量错误：期望 6, 得到 %d", len(corners))
	}

	// 验证所有顶点到中心的距离都等于半径
	for i, corner := range corners {
		distance := sqrtFloat(corner[0]*corner[0] + corner[1]*corner[1])
		if absFloat(distance-radius) > 0.1 {
			t.Errorf("顶点 %d 到中心距离错误：期望 %.2f, 得到 %.2f", i, radius, distance)
		}
	}
}

// TestHexGridManager 测试六边形网格管理器
func TestHexGridManager(t *testing.T) {
	mapSize := &config.MapSize{
		Width:  1000,
		Height: 1000,
	}
	hexRadius := 50.0
	isPointy := true

	hgm := NewHexGridManager(mapSize, hexRadius, isPointy)

	// 测试网格数量
	if hgm.GetGridCount() <= 0 {
		t.Error("网格数量应该大于 0")
	}

	// 测试获取网格
	hex := geo.NewHexCoord(0, 0)
	grid := hgm.GetGrid(hex)
	if grid == nil {
		t.Error("获取网格失败")
	}

	// 测试边界检查
	if !hgm.Contains(hex) {
		t.Error("原点应该在地图范围内")
	}

	outOfBounds := geo.NewHexCoord(1000, 1000)
	if hgm.Contains(outOfBounds) {
		t.Error("超出边界的坐标不应该在地图范围内")
	}

	// 测试获取相邻网格
	neighbors := hgm.GetNeighborGrids(grid)
	if len(neighbors) == 0 {
		t.Error("应该至少有一个相邻网格")
	}
}

// TestHexGridUnitManagement 测试六边形网格单位管理
func TestHexGridUnitManagement(t *testing.T) {
	mapSize := &config.MapSize{
		Width:  1000,
		Height: 1000,
	}
	hgm := NewHexGridManager(mapSize, 50.0, true)

	// 创建测试单位
	unit := &TestUnit{
		id:       1,
		configId: 100,
		coord:    geo.Coord{X: 0, Y: 0},
		hexCoord: geo.NewHexCoord(0, 0),
		unitType: MapUnitType_PlayerCity,
		owner:    NewOwner(1, OwnerType_Player),
	}

	// 测试添加单位到网格
	hex := geo.NewHexCoord(0, 0)
	success := hgm.AddUnitToGrid(unit, hex)
	if !success {
		t.Error("添加单位到网格失败")
	}

	// 验证单位存在
	grid := hgm.GetGrid(hex)
	if !grid.IsExistUnit(unit) {
		t.Error("单位应该存在于网格中")
	}

	// 测试获取单位
	retrievedUnit := grid.GetUnitById(unit.GetId())
	if retrievedUnit == nil {
		t.Error("获取单位失败")
	}

	// 测试移除单位
	success = hgm.RemoveUnitFromGrid(unit, hex)
	if !success {
		t.Error("移除单位失败")
	}

	if grid.IsExistUnit(unit) {
		t.Error("单位应该已被移除")
	}
}

// TestHexGridManagerUnitsInRadius 测试范围内单位查询
func TestHexGridManagerUnitsInRadius(t *testing.T) {
	mapSize := &config.MapSize{
		Width:  1000,
		Height: 1000,
	}
	hgm := NewHexGridManager(mapSize, 50.0, true)

	// 在中心添加单位
	centerHex := geo.NewHexCoord(5, 5)
	centerUnit := &TestUnit{
		id:       1,
		configId: 100,
		coord:    geo.Coord{X: 0, Y: 0},
		hexCoord: centerHex,
		unitType: MapUnitType_PlayerCity,
		owner:    NewOwner(1, OwnerType_Player),
	}
	hgm.AddUnitToGrid(centerUnit, centerHex)

	// 在附近添加单位
	nearbyHex := geo.NewHexCoord(6, 5)
	nearbyUnit := &TestUnit{
		id:       2,
		configId: 100,
		coord:    geo.Coord{X: 0, Y: 0},
		hexCoord: nearbyHex,
		unitType: MapUnitType_PlayerTroop,
		owner:    NewOwner(1, OwnerType_Player),
	}
	hgm.AddUnitToGrid(nearbyUnit, nearbyHex)

	// 测试获取范围内单位
	units := hgm.GetUnitsInRadius(centerHex, 2)
	if len(units) != 2 {
		t.Errorf("范围内单位数量错误：期望 2, 得到 %d", len(units))
	}

	// 测试获取范围内网格
	grids := hgm.GetHexGridsInRadius(centerHex, 1)
	// 半径为 1 应该包含中心 + 6 个邻居 = 7 个网格
	if len(grids) != 7 {
		t.Errorf("范围内网格数量错误：期望 7, 得到 %d", len(grids))
	}
}

// TestHexPathFinding 测试路径查找
func TestHexPathFinding(t *testing.T) {
	mapSize := &config.MapSize{
		Width:  1000,
		Height: 1000,
	}
	hgm := NewHexGridManager(mapSize, 50.0, true)

	start := geo.NewHexCoord(0, 0)
	end := geo.NewHexCoord(5, 5)

	// 测试基本路径查找
	path := hgm.FindPath(start, end, nil)
	if path == nil {
		t.Error("路径查找失败")
	}

	// 验证路径起点和终点
	if len(path) == 0 {
		t.Error("路径不应该为空")
	}
	if !path[0].Equal(start) {
		t.Error("路径起点错误")
	}
	if !path[len(path)-1].Equal(end) {
		t.Error("路径终点错误")
	}

	// 测试相同起点和终点
	samePath := hgm.FindPath(start, start, nil)
	if samePath == nil || len(samePath) != 1 {
		t.Error("相同起点和终点的路径应该只包含一个点")
	}

	// 测试超出边界的路径
	outOfBounds := geo.NewHexCoord(1000, 1000)
	invalidPath := hgm.FindPath(start, outOfBounds, nil)
	if invalidPath != nil {
		t.Error("超出边界的路径应该返回 nil")
	}
}

// TestHexLine 测试六边形直线算法
func TestHexLine(t *testing.T) {
	mapSize := &config.MapSize{
		Width:  1000,
		Height: 1000,
	}
	hgm := NewHexGridManager(mapSize, 50.0, true)

	start := geo.NewHexCoord(0, 0)
	end := geo.NewHexCoord(5, 0)

	line := hgm.GetHexesInLine(start, end)

	if len(line) < 2 {
		t.Error("直线应该至少包含起点和终点")
	}

	if !line[0].Equal(start) {
		t.Error("直线起点错误")
	}

	if !line[len(line)-1].Equal(end) {
		t.Error("直线终点错误")
	}
}

// TestHexVisionRange 测试视野范围
func TestHexVisionRange(t *testing.T) {
	mapSize := &config.MapSize{
		Width:  1000,
		Height: 1000,
	}
	hgm := NewHexGridManager(mapSize, 50.0, true)

	center := geo.NewHexCoord(5, 5)
	visionRange := int32(3)

	// 测试获取视野范围内的网格
	visionGrids := hgm.GetVisionRange(center, visionRange)
	if len(visionGrids) == 0 {
		t.Error("视野范围内应该有网格")
	}

	// 测试可见性检查
	nearbyHex := geo.NewHexCoord(6, 5)
	if !hgm.IsVisible(center, nearbyHex, visionRange) {
		t.Error("附近的六边形应该是可见的")
	}

	farHex := geo.NewHexCoord(20, 20)
	if hgm.IsVisible(center, farHex, visionRange) {
		t.Error("远处的六边形应该是不可见的")
	}
}

// TestHexMoveUnit 测试单位移动
func TestHexMoveUnit(t *testing.T) {
	mapSize := &config.MapSize{
		Width:  1000,
		Height: 1000,
	}
	hgm := NewHexGridManager(mapSize, 50.0, true)

	unit := &TestUnit{
		id:       1,
		configId: 100,
		coord:    geo.Coord{X: 0, Y: 0},
		hexCoord: geo.NewHexCoord(0, 0),
		unitType: MapUnitType_PlayerTroop,
		owner:    NewOwner(1, OwnerType_Player),
	}

	from := geo.NewHexCoord(0, 0)
	to := geo.NewHexCoord(1, 0)

	// 先添加单位
	hgm.AddUnitToGrid(unit, from)

	// 测试移动单位
	success := hgm.MoveUnit(unit, from, to)
	if !success {
		t.Error("移动单位失败")
	}

	// 验证原位置没有单位
	fromGrid := hgm.GetGrid(from)
	if fromGrid.IsExistUnit(unit) {
		t.Error("单位应该已从原位置移除")
	}

	// 验证新位置有单位
	toGrid := hgm.GetGrid(to)
	if !toGrid.IsExistUnit(unit) {
		t.Error("单位应该已添加到新位置")
	}

	// 测试移动到边界外
	outOfBounds := geo.NewHexCoord(1000, 1000)
	success = hgm.MoveUnit(unit, to, outOfBounds)
	if success {
		t.Error("移动到边界外应该失败")
	}
}

// TestHexGetAllUnits 测试获取所有单位
func TestHexGetAllUnits(t *testing.T) {
	mapSize := &config.MapSize{
		Width:  1000,
		Height: 1000,
	}
	hgm := NewHexGridManager(mapSize, 50.0, true)

	// 添加多个单位
	units := make([]Unit, 5)
	for i := 0; i < 5; i++ {
		units[i] = &TestUnit{
			id:       int64(i + 1),
			configId: 100,
			coord:    geo.Coord{X: int32(i), Y: int32(i)},
			hexCoord: geo.NewHexCoord(int32(i), int32(i)),
			unitType: MapUnitType_PlayerTroop,
			owner:    NewOwner(1, OwnerType_Player),
		}
		hgm.AddUnitToGrid(units[i], geo.NewHexCoord(int32(i), int32(i)))
	}

	// 测试获取所有单位
	allUnits := hgm.GetAllUnits()
	if len(allUnits) != 5 {
		t.Errorf("所有单位数量错误：期望 5, 得到 %d", len(allUnits))
	}

	// 测试按类型获取单位
	troopUnits := hgm.GetUnitsByType(MapUnitType_PlayerTroop)
	if len(troopUnits) != 5 {
		t.Errorf("部队单位数量错误：期望 5, 得到 %d", len(troopUnits))
	}
}

// TestHexRectangle 测试六边形范围
func TestHexRectangle(t *testing.T) {
	rect := geo.NewHexRectangle(0, 10, 0, 10)

	// 测试范围内
	hexInside := geo.NewHexCoord(5, 5)
	if !rect.Contains(hexInside) {
		t.Error("坐标应该在范围内")
	}

	// 测试范围外
	hexOutside := geo.NewHexCoord(15, 15)
	if rect.Contains(hexOutside) {
		t.Error("坐标应该在范围外")
	}

	// 测试边界
	hexOnEdge := geo.NewHexCoord(10, 10)
	if !rect.Contains(hexOnEdge) {
		t.Error("边界坐标应该在范围内")
	}
}

// TestCalculateGridSize 测试网格大小计算
func TestCalculateGridSize(t *testing.T) {
	mapWidth := 1000.0
	mapHeight := 1000.0
	radius := 50.0

	// 测试 pointy 朝向
	qCount, rCount := geo.CalculateGridSize(mapWidth, mapHeight, radius, true)
	if qCount <= 0 || rCount <= 0 {
		t.Error("网格数量应该大于 0")
	}

	// 测试 flat 朝向
	qCountFlat, rCountFlat := geo.CalculateGridSize(mapWidth, mapHeight, radius, false)
	if qCountFlat <= 0 || rCountFlat <= 0 {
		t.Error("网格数量应该大于 0")
	}
}

// 辅助函数

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func sqrtFloat(x float64) float64 {
	return math.Sqrt(x)
}

// TestUnit 测试用单位实现
type TestUnit struct {
	id       int64
	configId int32
	coord    geo.Coord
	hexCoord *geo.HexCoord
	unitType MapUnitType
	owner    *Owner
}

func (t *TestUnit) GetId() int64                    { return t.id }
func (t *TestUnit) GetConfigId() int32              { return t.configId }
func (t *TestUnit) GetCoord() *geo.Coord            { return &t.coord }
func (t *TestUnit) SetCoord(coord *geo.Coord)       { t.coord = *coord }
func (t *TestUnit) GetHexCoord() *geo.HexCoord      { return t.hexCoord }
func (t *TestUnit) SetHexCoord(coord *geo.HexCoord) { t.hexCoord = coord }
func (t *TestUnit) GetType() MapUnitType            { return t.unitType }
func (t *TestUnit) GetOwner() *Owner                { return t.owner }
