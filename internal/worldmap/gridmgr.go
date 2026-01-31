package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/geo"

// grid管理器
type GridManager struct {
	grids      []*Grid // 网格数组（惰性初始化，nil表示未创建）
	MapWidth   int32   // 地图的宽度（世界单位）
	MapHeight  int32   // 地图的高度（世界单位）
	GridWidth  int32   // 网格宽度（世界单位）
	GridHeight int32   // 网格高度（世界单位）
	gridCols   int32   // 网格列数
	gridRows   int32   // 网格行数
}

func NewGridManager(width, height, gridWidth, gridHeight int32) *GridManager {
	// 计算网格数量
	gridCols := (width + gridWidth - 1) / gridWidth
	gridRows := (height + gridHeight - 1) / gridHeight
	totalGrids := gridCols * gridRows

	mgr := &GridManager{
		grids:      make([]*Grid, totalGrids), // 预分配指针数组，所有元素为nil
		MapWidth:   width,
		MapHeight:  height,
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
		gridCols:   gridCols,
		gridRows:   gridRows,
	}
	return mgr
}

// 计算网格索引
func (gm *GridManager) calcGridIndex(x, y int32) (int32, bool) {
	if x < 0 || x >= gm.MapWidth || y < 0 || y >= gm.MapHeight {
		return -1, false
	}
	gridX := x / gm.GridWidth
	gridY := y / gm.GridHeight
	index := gridY*gm.gridCols + gridX
	return index, true
}

// 通过坐标获取网格（惰性初始化）
func (gm *GridManager) GetGridByPos(x, y int32) *Grid {
	index, ok := gm.calcGridIndex(x, y)
	if !ok {
		return nil
	}

	// 惰性初始化：如果格子为nil则创建
	if gm.grids[index] == nil {
		// 计算格子的世界坐标
		gridX := (x / gm.GridWidth) * gm.GridWidth
		gridY := (y / gm.GridHeight) * gm.GridHeight

		gm.grids[index] = NewGrid(
			geo.NewCoord(gridX, gridY),
			gm.GridWidth,
			gm.GridHeight,
		)
	}

	return gm.grids[index]
}

// 通过坐标对象获取网格
func (gm *GridManager) GetGridByCoord(coord *geo.Coord) *Grid {
	return gm.GetGridByPos(coord.X, coord.Y)
}

// 获取网格总数
func (gm *GridManager) GetTotalGrids() int32 {
	return gm.gridCols * gm.gridRows
}

// 获取已创建的网格数量
func (gm *GridManager) GetCreatedGrids() int32 {
	count := int32(0)
	for _, grid := range gm.grids {
		if grid != nil {
			count++
		}
	}
	return count
}

// 清理空的网格（可选：当网格变为空时释放内存）
func (gm *GridManager) CleanupEmptyGrids() {
	for idx, grid := range gm.grids {
		if grid != nil && len(grid.GetUnits()) == 0 {
			// 可以释放网格对象，设为nil
			// 注意：这会导致下次访问时重新创建，根据性能需求决定是否启用
			// gm.grids[idx] = nil
			_ = idx // 避免未使用变量警告
		}
	}
}

func (mgr *GridManager) AddUnit(unit Unit) {
	grid := mgr.GetGridByCoord(unit.GetCoord())
	if grid != nil {
		grid.AddUnit(unit)
	}
}

func (mgr *GridManager) RemoveUnit(unit Unit) {
	grid := mgr.GetGridByCoord(unit.GetCoord())
	if grid != nil {
		grid.RemoveUnit(unit)
	}
}

// 更新地图单位的位置
func (mgr *GridManager) UpdateInitCoord(unit Unit, coord *geo.Coord) {
	oldGrid := mgr.GetGridByCoord(unit.GetCoord())
	newGrid := mgr.GetGridByCoord(coord)
	unit.SetCoord(coord)

	if oldGrid != nil && newGrid != nil && oldGrid == newGrid {
		return
	}

	if oldGrid != nil {
		oldGrid.RemoveUnit(unit)
	}
	if newGrid != nil {
		newGrid.AddUnit(unit)
	}
}
