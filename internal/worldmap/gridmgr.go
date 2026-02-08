package worldmap

import (
	"github.com/GooLuck/WorldMap/internal/worldmap/config"
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// grid管理器
type GridManager struct {
	grids    []*Grid // 网格数组（惰性初始化，nil表示未创建）
	mapSize  *config.MapSize
	gridCols int32 // 网格列数
	gridRows int32 // 网格行数
}

func NewGridManager(mapSize *config.MapSize) *GridManager {
	// 计算网格数量
	gridCols := (mapSize.Width + mapSize.GridWidth - 1) / mapSize.GridWidth
	gridRows := (mapSize.Height + mapSize.GridHeight - 1) / mapSize.GridHeight
	totalGrids := gridCols * gridRows

	mgr := &GridManager{
		grids:    make([]*Grid, totalGrids), // 预分配指针数组，所有元素为nil
		mapSize:  mapSize,
		gridCols: gridCols,
		gridRows: gridRows,
	}
	return mgr
}

// 计算网格索引
func (gm *GridManager) calcGridIndex(x, y int32) (int32, bool) {
	if x < 0 || x >= gm.mapSize.Width || y < 0 || y >= gm.mapSize.Height {
		return -1, false
	}
	gridX := x / gm.mapSize.GridWidth
	gridY := y / gm.mapSize.GridHeight
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
		gridX := (x / gm.mapSize.GridWidth) * gm.mapSize.GridWidth
		gridY := (y / gm.mapSize.GridHeight) * gm.mapSize.GridHeight

		gm.grids[index] = NewGrid(
			geo.NewCoord(gridX, gridY),
			gm.mapSize.GridWidth,
			gm.mapSize.GridHeight,
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

// rect 是否与grid 对齐，对齐的话，就不用一个一个unit判断了，整个grid的unit都满足
func (mgr *GridManager) isAlignGrid(rect *geo.Rectangle) bool {
	return rect.X%mgr.mapSize.GridWidth == 0 &&
		rect.Y%mgr.mapSize.GridHeight == 0 &&
		rect.Width%mgr.mapSize.GridWidth == 0 &&
		rect.Height%mgr.mapSize.GridHeight == 0
}

// 获取矩形范围内的单位
func (mgr *GridManager) GetRectUnits(rect *geo.Rectangle, align bool) []Unit {
	leftX, rightX, leftY, rightY := RectToGrid(mgr.mapSize, rect)
	retUnits := make([]Unit, 0)

	if !align && mgr.isAlignGrid(rect) {
		align = true
	}

	for y := leftY; y <= rightY; y++ {
		for x := leftX; x <= rightX; x++ {
			grid := mgr.GetGridByPos(x, y)
			if grid == nil {
				continue
			}

			if align {
				retUnits = append(retUnits, grid.GetUnits()...)
				continue
			}

			gridUnits := grid.GetUnits()
			for _, u := range gridUnits {
				if rect.IsCoordInRect(u.GetCoord()) {
					retUnits = append(retUnits, u)
				}
			}
		}
	}
	return retUnits
}

// 遍历矩形范围内的单位
func (mgr *GridManager) RangeRectUnits(rect *geo.Rectangle, align bool, callback func(unit Unit) bool) {
	leftX, rightX, leftY, rightY := RectToGrid(mgr.mapSize, rect)

	for y := leftY; y <= rightY; y++ {
		for x := leftX; x <= rightX; x++ {
			grid := mgr.GetGridByPos(x, y)
			if grid != nil {
				grid.RangeUnits(callback)
			}
		}
	}
}
