package worldmap

import (
	"github.com/GooLuck/WorldMap/internal/worldmap/geo"
)

// TerrainType 地形类型
type TerrainType int32

const (
	TerrainType_None     TerrainType = iota // 无地形
	TerrainType_Plain                       // 平原
	TerrainType_Forest                      // 森林
	TerrainType_Mountain                    // 山地
	TerrainType_Swamp                       // 沼泽
	TerrainType_Desert                      // 沙漠
	TerrainType_Snow                        // 雪地
	TerrainType_Water                       // 水域（不可通行）
	TerrainType_Lava                        // 熔岩（不可通行）
)

// TerrainConfig 地形配置
type TerrainConfig struct {
	Type         TerrainType // 地形类型
	MoveCost     float32     // 移动成本系数（1.0 为正常，越高越难通行）
	DefenseBonus float32     // 防御加成（0.0-1.0）
	Visible      bool        // 是否可见（用于战争迷雾）
	Passable     bool        // 是否可通行
}

// DefaultTerrainConfigs 默认地形配置
var DefaultTerrainConfigs = map[TerrainType]*TerrainConfig{
	TerrainType_None: {
		Type:         TerrainType_None,
		MoveCost:     1.0,
		DefenseBonus: 0.0,
		Visible:      true,
		Passable:     true,
	},
	TerrainType_Plain: {
		Type:         TerrainType_Plain,
		MoveCost:     1.0,
		DefenseBonus: 0.0,
		Visible:      true,
		Passable:     true,
	},
	TerrainType_Forest: {
		Type:         TerrainType_Forest,
		MoveCost:     1.5,
		DefenseBonus: 0.3,
		Visible:      true,
		Passable:     true,
	},
	TerrainType_Mountain: {
		Type:         TerrainType_Mountain,
		MoveCost:     3.0,
		DefenseBonus: 0.5,
		Visible:      true,
		Passable:     true,
	},
	TerrainType_Swamp: {
		Type:         TerrainType_Swamp,
		MoveCost:     2.5,
		DefenseBonus: 0.1,
		Visible:      true,
		Passable:     true,
	},
	TerrainType_Desert: {
		Type:         TerrainType_Desert,
		MoveCost:     1.3,
		DefenseBonus: 0.0,
		Visible:      true,
		Passable:     true,
	},
	TerrainType_Snow: {
		Type:         TerrainType_Snow,
		MoveCost:     1.8,
		DefenseBonus: 0.1,
		Visible:      true,
		Passable:     true,
	},
	TerrainType_Water: {
		Type:         TerrainType_Water,
		MoveCost:     0,
		DefenseBonus: 0.0,
		Visible:      true,
		Passable:     false,
	},
	TerrainType_Lava: {
		Type:         TerrainType_Lava,
		MoveCost:     0,
		DefenseBonus: 0.0,
		Visible:      true,
		Passable:     false,
	},
}

// GetTerrainConfig 获取地形配置
func GetTerrainConfig(terrainType TerrainType) *TerrainConfig {
	if config, exists := DefaultTerrainConfigs[terrainType]; exists {
		return config
	}
	return DefaultTerrainConfigs[TerrainType_None]
}

// TerrainMap 地形地图
type TerrainMap struct {
	terrains map[uint64]TerrainType // hex hash -> terrain type
	bounds   *geo.HexRectangle      // 边界范围
}

// NewTerrainMap 创建地形地图
func NewTerrainMap(bounds *geo.HexRectangle) *TerrainMap {
	return &TerrainMap{
		terrains: make(map[uint64]TerrainType),
		bounds:   bounds,
	}
}

// SetTerrain 设置六边形地形
func (tm *TerrainMap) SetTerrain(hex *geo.HexCoord, terrainType TerrainType) {
	if tm.bounds.Contains(hex) {
		tm.terrains[hashHex(hex.Q, hex.R)] = terrainType
	}
}

// GetTerrain 获取六边形地形
func (tm *TerrainMap) GetTerrain(hex *geo.HexCoord) TerrainType {
	if !tm.bounds.Contains(hex) {
		return TerrainType_None
	}
	if terrain, exists := tm.terrains[hashHex(hex.Q, hex.R)]; exists {
		return terrain
	}
	return TerrainType_Plain // 默认为平原
}

// GetTerrainConfig 获取六边形地形配置
func (tm *TerrainMap) GetTerrainConfig(hex *geo.HexCoord) *TerrainConfig {
	terrainType := tm.GetTerrain(hex)
	return GetTerrainConfig(terrainType)
}

// GetMoveCost 获取六边形移动成本
func (tm *TerrainMap) GetMoveCost(hex *geo.HexCoord) float32 {
	return tm.GetTerrainConfig(hex).MoveCost
}

// IsPassable 检查六边形是否可通行
func (tm *TerrainMap) IsPassable(hex *geo.HexCoord) bool {
	return tm.GetTerrainConfig(hex).Passable
}

// GetDefenseBonus 获取六边形防御加成
func (tm *TerrainMap) GetDefenseBonus(hex *geo.HexCoord) float32 {
	return tm.GetTerrainConfig(hex).DefenseBonus
}

// TerrainCostFunc 创建地形成本函数（用于路径查找）
func (tm *TerrainMap) TerrainCostFunc() TerrainCostFunc {
	return func(hex *geo.HexCoord) int32 {
		config := tm.GetTerrainConfig(hex)
		if !config.Passable {
			return 9999 // 不可通行的地形成本设为极高
		}
		return int32(config.MoveCost * 10) // 转换为整数成本
	}
}

// FindPathWithTerrain 考虑地形的路径查找
func (tm *TerrainMap) FindPathWithTerrain(hgm *HexGridManager, start, end *geo.HexCoord) []*geo.HexCoord {
	return hgm.FindPath(start, end, tm.TerrainCostFunc())
}

// TerrainGenerator 地形生成器
type TerrainGenerator struct {
	seed int64
}

// NewTerrainGenerator 创建地形生成器
func NewTerrainGenerator(seed int64) *TerrainGenerator {
	return &TerrainGenerator{
		seed: seed,
	}
}

// GenerateSimpleTerrain 生成简单地形（基于距离和随机性）
// centerHex: 中心六边形（通常是地图中心）
// waterDistance: 水域距离阈值（超过此距离可能有水）
// mountainDistance: 山地距离阈值
func (tg *TerrainGenerator) GenerateSimpleTerrain(hgm *HexGridManager, centerHex *geo.HexCoord, waterDistance, mountainDistance int32) *TerrainMap {
	terrainMap := NewTerrainMap(hgm.GetBounds())

	hgm.RangeAllGrids(func(grid *HexGrid) bool {
		hex := grid.GetCoord()
		distance := centerHex.DistanceTo(hex)

		var terrain TerrainType
		if distance < waterDistance {
			// 靠近中心可能是平原
			terrain = TerrainType_Plain
		} else if distance > mountainDistance {
			// 边缘可能是山地或水域
			if distance%2 == 0 {
				terrain = TerrainType_Mountain
			} else {
				terrain = TerrainType_Forest
			}
		} else {
			// 中间区域混合地形
			switch distance % 5 {
			case 0:
				terrain = TerrainType_Plain
			case 1:
				terrain = TerrainType_Forest
			case 2:
				terrain = TerrainType_Mountain
			case 3:
				terrain = TerrainType_Swamp
			default:
				terrain = TerrainType_Desert
			}
		}

		terrainMap.SetTerrain(hex, terrain)
		return true
	})

	return terrainMap
}

// GenerateBiomeTerrain 生成生物群系地形
func (tg *TerrainGenerator) GenerateBiomeTerrain(hgm *HexGridManager, biomes map[int32]TerrainType) *TerrainMap {
	terrainMap := NewTerrainMap(hgm.GetBounds())

	// biomes: distance -> terrain type
	hgm.RangeAllGrids(func(grid *HexGrid) bool {
		hex := grid.GetCoord()

		// 根据距离选择地形
		for distance, terrain := range biomes {
			if hex.DistanceTo(geo.NewHexCoord(0, 0)) <= distance {
				terrainMap.SetTerrain(hex, terrain)
				break
			}
		}

		return true
	})

	return terrainMap
}
