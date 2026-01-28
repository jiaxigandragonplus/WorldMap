package worldmap

// MapUnitType 地图单位类型
type MapUnitType int32

const (
	MapUnitType_None        MapUnitType = iota
	MapUnitType_Obstacle                // 障碍物
	MapUnitType_PlayerCity              // 玩家主城
	MapUnitType_PlayerTroop             // 玩家部队
	MapUnitType_Npc                     // NPC
)
