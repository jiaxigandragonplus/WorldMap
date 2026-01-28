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

// 对象之间的关系
type Relation int32

const (
	Relation_Self   Relation = iota + 1 // 自己
	Relation_Union                      // 同盟
	Relation_Enemy                      // 敌人
	Relation_System                     // 系统
	Relation_All                        // 任意关系
)

// 地图单位所有者类型
type OwnerType int32

const (
	OwnerType_None   OwnerType = iota // 无
	OwnerType_Player                  // 玩家
	OwnerType_Npc                     // npc
	OwnerType_Union                   // 联盟
)
