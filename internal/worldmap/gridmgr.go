package worldmap

import "github.com/GooLuck/WorldMap/internal/worldmap/geo"

// grid管理器
type GridManager struct {
	grids      []*Grid // 网格数组
	MapWidth   int32   // 地图的宽度
	MapHeight  int32   // 地图的高度
	GridWidth  int32   // 网格宽度
	GridHeight int32   // 网格高度
}

func NewGridManager(width, height, gridWidth, gridHeight int32) *GridManager {
	mgr := &GridManager{
		grids:      make([]*Grid, 0),
		MapWidth:   width,
		MapHeight:  height,
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}
	return mgr
}

// 通过坐标获取网格
func (gm *GridManager) GetGridByPos(x, y int32) *Grid {
	return nil
}

func (gm *GridManager) GetGridByCoord(coord *geo.Coord) *Grid {
	return nil
}
