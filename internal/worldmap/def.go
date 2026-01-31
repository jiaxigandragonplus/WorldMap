package worldmap

// MapUnitType 地图单位类型
type MapUnitType int32

const (
	MapUnitType_None         MapUnitType = iota
	MapUnitType_Obstacle                 // 障碍物
	MapUnitType_PlayerCity               // 玩家主城
	MapUnitType_PlayerTroop              // 玩家部队
	MapUnitType_Npc                      // NPC
	MapUnitType_Resource                 // 资源点
	MapUnitType_ResourceZone             // 资源区域
)

// 资源类型
type ResourceType int32

const (
	ResourceType_None    ResourceType = iota
	ResourceType_Gold                 // 金币
	ResourceType_Wood                 // 木材
	ResourceType_Stone                // 石头
	ResourceType_Food                 // 食物
	ResourceType_Ore                  // 矿石
	ResourceType_Crystal              // 水晶
	ResourceType_Gem                  // 宝石
	ResourceType_Herb                 // 草药
)

// 障碍物类型
type ObstacleType int32

const (
	ObstacleType_None     ObstacleType = iota
	ObstacleType_Mountain              // 高山
	ObstacleType_Lake                  // 湖泊
	ObstacleType_Forest                // 森林（密集）
	ObstacleType_Swamp                 // 沼泽
	ObstacleType_Cliff                 // 悬崖
	ObstacleType_River                 // 河流
	ObstacleType_Volcano               // 火山
	ObstacleType_Desert                // 沙漠（特殊地形）
	ObstacleType_Ruins                 // 废墟
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
	OwnerType_System                  // 系统（资源点等）
)
