package worldmap

import (
	"testing"

	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

func TestNewGridManager(t *testing.T) {
	// 测试创建 GridManager
	gm := NewGridManager(1000, 1000, 10, 10)

	// 验证基本属性
	if gm.MapWidth != 1000 {
		t.Errorf("期望 MapWidth = 1000, 得到 %d", gm.MapWidth)
	}
	if gm.MapHeight != 1000 {
		t.Errorf("期望 MapHeight = 1000, 得到 %d", gm.MapHeight)
	}
	if gm.GridWidth != 10 {
		t.Errorf("期望 GridWidth = 10, 得到 %d", gm.GridWidth)
	}
	if gm.GridHeight != 10 {
		t.Errorf("期望 GridHeight = 10, 得到 %d", gm.GridHeight)
	}

	// 验证网格计算
	expectedCols := int32(100) // 1000/10 = 100
	expectedRows := int32(100)
	if gm.gridCols != expectedCols {
		t.Errorf("期望 gridCols = %d, 得到 %d", expectedCols, gm.gridCols)
	}
	if gm.gridRows != expectedRows {
		t.Errorf("期望 gridRows = %d, 得到 %d", expectedRows, gm.gridRows)
	}

	// 验证总网格数
	totalGrids := gm.GetTotalGrids()
	expectedTotal := expectedCols * expectedRows // 100 * 100 = 10000
	if totalGrids != expectedTotal {
		t.Errorf("期望总网格数 = %d, 得到 %d", expectedTotal, totalGrids)
	}

	// 验证初始时没有创建的网格
	createdGrids := gm.GetCreatedGrids()
	if createdGrids != 0 {
		t.Errorf("期望初始已创建网格数 = 0, 得到 %d", createdGrids)
	}
}

func TestGetGridByPos_LazyInitialization(t *testing.T) {
	gm := NewGridManager(1000, 1000, 10, 10)

	// 第一次获取网格，应该创建
	grid1 := gm.GetGridByPos(50, 50)
	if grid1 == nil {
		t.Error("期望获取网格(50, 50)成功，得到 nil")
	}

	// 验证已创建网格数
	if gm.GetCreatedGrids() != 1 {
		t.Errorf("期望已创建网格数 = 1, 得到 %d", gm.GetCreatedGrids())
	}

	// 获取同一个网格，应该返回同一个对象
	grid2 := gm.GetGridByPos(50, 50)
	if grid2 != grid1 {
		t.Error("期望返回同一个网格对象")
	}

	// 验证已创建网格数没有增加
	if gm.GetCreatedGrids() != 1 {
		t.Errorf("期望已创建网格数仍然为 1, 得到 %d", gm.GetCreatedGrids())
	}

	// 获取另一个网格
	grid3 := gm.GetGridByPos(150, 150)
	if grid3 == nil {
		t.Error("期望获取网格(150, 150)成功，得到 nil")
	}

	if gm.GetCreatedGrids() != 2 {
		t.Errorf("期望已创建网格数 = 2, 得到 %d", gm.GetCreatedGrids())
	}
}

func TestGetGridByCoord(t *testing.T) {
	gm := NewGridManager(1000, 1000, 10, 10)

	coord := geo.NewCoord(75, 85)
	grid := gm.GetGridByCoord(coord)

	if grid == nil {
		t.Error("期望通过坐标获取网格成功，得到 nil")
	}

	// 验证坐标转换正确
	// (75, 85) 应该在网格 (70, 80) 到 (80, 90) 范围内
	// 因为网格宽度/高度为10，所以网格起始坐标为 (70, 80)
	if grid.X != 70 || grid.Y != 80 {
		t.Errorf("期望网格坐标 = (70, 80), 得到 (%d, %d)", grid.X, grid.Y)
	}
}

func TestCalcGridIndex(t *testing.T) {
	gm := NewGridManager(1000, 1000, 10, 10)

	// 测试有效坐标
	testCases := []struct {
		x, y          int32
		expectedIndex int32
		expectedOk    bool
	}{
		{0, 0, 0, true},               // 第一个网格
		{9, 9, 0, true},               // 第一个网格内
		{10, 0, 1, true},              // 第二个网格
		{999, 999, 99*100 + 99, true}, // 最后一个网格
		{50, 50, 5*100 + 5, true},     // 中间网格
		{-1, 0, -1, false},            // 无效坐标
		{0, -1, -1, false},            // 无效坐标
		{1000, 0, -1, false},          // 超出边界
		{0, 1000, -1, false},          // 超出边界
	}

	for _, tc := range testCases {
		index, ok := gm.calcGridIndex(tc.x, tc.y)
		if ok != tc.expectedOk {
			t.Errorf("坐标(%d, %d): 期望 ok = %v, 得到 %v", tc.x, tc.y, tc.expectedOk, ok)
		}
		if ok && index != tc.expectedIndex {
			t.Errorf("坐标(%d, %d): 期望 index = %d, 得到 %d", tc.x, tc.y, tc.expectedIndex, index)
		}
	}
}

func TestGridManager_MemoryEfficiency(t *testing.T) {
	// 创建大型地图管理器
	gm := NewGridManager(10000, 10000, 10, 10)

	// 总网格数应该是 1000x1000 = 1,000,000
	totalGrids := gm.GetTotalGrids()
	if totalGrids != 1000*1000 {
		t.Errorf("期望总网格数 = 1,000,000, 得到 %d", totalGrids)
	}

	// 初始时应该没有创建任何网格
	if gm.GetCreatedGrids() != 0 {
		t.Errorf("期望初始已创建网格数 = 0, 得到 %d", gm.GetCreatedGrids())
	}

	// 访问少量网格
	accessPoints := [][2]int32{
		{123, 456},
		{789, 123},
		{456, 789},
		{9999, 9999}, // 边界点
	}

	for _, point := range accessPoints {
		gm.GetGridByPos(point[0], point[1])
	}

	// 应该只创建了访问过的网格
	expectedCreated := int32(len(accessPoints))
	if gm.GetCreatedGrids() != expectedCreated {
		t.Errorf("期望已创建网格数 = %d, 得到 %d", expectedCreated, gm.GetCreatedGrids())
	}

	// 验证内存效率：已创建网格数远小于总网格数
	creationRatio := float64(gm.GetCreatedGrids()) / float64(gm.GetTotalGrids())
	if creationRatio > 0.01 {
		t.Errorf("期望创建比例 < 1%%, 得到 %.4f%%", creationRatio*100)
	}
}
