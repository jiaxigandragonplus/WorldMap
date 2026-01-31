package config

// 资源点类型
type ResourcePointType int32

const (
	ResourcePointType_Fixed       ResourcePointType = iota // 固定资源点（永久存在）
	ResourcePointType_RandomSpawn                          // 随机刷新点
	ResourcePointType_Seasonal                             // 季节性资源点
	ResourcePointType_EventBased                           // 事件触发资源点
)

// 刷新策略类型
type RefreshStrategy int32

const (
	RefreshStrategy_Linear      RefreshStrategy = iota // 线性恢复
	RefreshStrategy_Exponential                        // 指数恢复
	RefreshStrategy_Stepwise                           // 阶梯式恢复
	RefreshStrategy_Random                             // 随机恢复
)

// slg地图配置
type MapConfig struct {
	MapID              int32  // 地图唯一ID
	MapName            string // 地图名称
	Width              int32  // 地图宽度（世界单位）
	Height             int32  // 地图高度（世界单位）
	XGridLen           int32  // 横向网格数量
	YGridLen           int32  // 纵向网格数量
	DefaultVisionRange int32  // 默认视野范围（网格数）
	MaxPlayers         int32  // 最大玩家数量

	// 出生点配置
	SpawnPoints []SpawnPointConfig // 出生点列表

	// 资源点配置（兼容旧版）
	ResourcePoints []ResourcePointConfig // 资源点列表

	// 增强资源点配置
	EnhancedResourcePoints []EnhancedResourcePointConfig // 增强资源点列表

	// 资源区域配置
	ResourceZones []ResourceZoneConfig // 资源区域列表

	// 障碍物配置
	Obstacles []ObstacleConfig // 障碍物列表

	// 障碍物区域配置
	ObstacleZones []ObstacleZoneConfig // 障碍物区域列表

	// 全局刷新配置
	GlobalRefreshConfig GlobalRefreshConfig // 全局刷新配置
}

// 出生点配置
type SpawnPointConfig struct {
	PointID   int32 // 出生点ID
	X         int32 // X坐标（世界单位）
	Y         int32 // Y坐标（世界单位）
	Radius    int32 // 安全区半径
	ForPlayer bool  // 是否为玩家出生点（true为玩家，false为NPC）
}

// 资源点配置（基础版，保持向后兼容）
type ResourcePointConfig struct {
	PointID      int32   // 资源点ID
	X            int32   // X坐标（世界单位）
	Y            int32   // Y坐标（世界单位）
	ResourceType string  // 资源类型：gold金币, wood木材, stone石头, food食物
	MaxAmount    int32   // 最大资源量
	RegenRate    float32 // 资源恢复速率（单位/秒）
	RegenDelay   int32   // 资源恢复延迟（秒）
}

// 增强资源点配置
type EnhancedResourcePointConfig struct {
	PointID         int32             // 资源点ID
	X               int32             // X坐标（世界单位）
	Y               int32             // Y坐标（世界单位）
	ResourceType    string            // 资源类型
	PointType       ResourcePointType // 资源点类型
	RefreshStrategy RefreshStrategy   // 刷新策略

	// 基础属性
	MaxAmount     int32   // 最大资源量
	CurrentAmount int32   // 当前资源量
	RegenRate     float32 // 基础恢复速率（单位/秒）
	RegenDelay    int32   // 恢复延迟（秒）

	// 随机刷新相关
	SpawnRadius   int32   // 刷新半径（对于随机刷新点，0表示固定位置）
	SpawnInterval int32   // 刷新间隔（秒）
	SpawnChance   float32 // 刷新概率（0.0-1.0）

	// 时间相关配置
	ActiveHours  []int32 // 活跃时间段（小时，如[9,18]表示9点到18点）
	SeasonMonths []int32 // 活跃月份（1-12）

	// 事件触发
	TriggerEventID int32 // 触发事件ID
	DespawnAfter   int32 // 消失时间（秒，0表示不消失）

	// 高级配置
	MinPlayerLevel  int32  // 最小玩家等级可见
	MaxPlayerLevel  int32  // 最大玩家等级可见
	FactionRestrict string // 阵营限制（如"alliance", "horde", ""表示无限制）
}

// 资源区域配置
type ResourceZoneConfig struct {
	ZoneID   int32  // 区域ID
	ZoneName string // 区域名称
	ZoneType string // 区域类型：forest森林, mountain山脉, plain平原, lake湖泊, desert沙漠
	MinX     int32  // 区域最小X坐标
	MinY     int32  // 区域最小Y坐标
	MaxX     int32  // 区域最大X坐标
	MaxY     int32  // 区域最大Y坐标

	// 资源分布
	ResourceTypes []string // 该区域可能出现的资源类型
	Density       float32  // 资源密度（0.0-1.0）
	MinDistance   int32    // 资源点最小间距
	MaxPoints     int32    // 最大资源点数

	// 刷新规则
	RefreshEnabled bool    // 是否启用自动刷新
	RefreshRate    float32 // 刷新速率（点/小时）
	MaxConcurrent  int32   // 最大同时存在资源点数
}

// 全局刷新配置
type GlobalRefreshConfig struct {
	// 时间相关
	DailyRefreshTime  string // 每日刷新时间（格式："HH:MM"）
	WeeklyRefreshDay  int32  // 每周刷新日（0-6，0表示周日）
	MonthlyRefreshDay int32  // 每月刷新日（1-31）

	// 资源平衡
	EnableDynamicBalance bool    // 启用动态平衡
	MinResourceRatio     float32 // 最小资源比例（当资源低于此比例时加速刷新）
	MaxResourceRatio     float32 // 最大资源比例（当资源高于此比例时减速刷新）

	// 玩家相关
	PlayerBasedRefresh bool  // 基于玩家数量的刷新
	PlayersPerResource int32 // 每个资源点对应的玩家数
	MaxRefreshDistance int32 // 最大刷新距离（离玩家最近距离）

	// 性能限制
	MaxRefreshPerTick   int32 // 每Tick最大刷新数量
	RefreshTickInterval int32 // 刷新Tick间隔（秒）
}

// 障碍物配置
type ObstacleConfig struct {
	ObstacleID   int32  // 障碍物ID
	X            int32  // X坐标（世界单位）
	Y            int32  // Y坐标（世界单位）
	Width        int32  // 宽度（世界单位）
	Height       int32  // 高度（世界单位）
	ObstacleType string // 障碍物类型：mountain高山, lake湖泊, forest森林, swamp沼泽, cliff悬崖, river河流, volcano火山, desert沙漠, ruins废墟
	Name         string // 障碍物名称（可选）

	// 影响范围配置
	BlockBuilding bool // 是否阻挡建筑摆放
	BlockResource bool // 是否阻挡资源刷新
	BlockMonster  bool // 是否阻挡怪物刷新
	AllowMarch    bool // 是否允许行军通过（true为允许，false为阻挡）

	// 影响半径
	BuildingRadius int32 // 建筑阻挡半径（0表示仅障碍物本身区域）
	ResourceRadius int32 // 资源阻挡半径
	MonsterRadius  int32 // 怪物阻挡半径

	// 视觉效果
	VisualEffect string // 视觉效果标识
	MinimapIcon  string // 小地图图标
	Description  string // 描述文本

	// 特殊效果
	SpecialEffects []string // 特殊效果列表，如["fog", "damage_over_time", "slow_movement"]
	EffectStrength float32  // 效果强度
}

// 障碍物区域配置（用于大片连续障碍物）
type ObstacleZoneConfig struct {
	ZoneID       int32  // 区域ID
	ZoneName     string // 区域名称
	ObstacleType string // 主要障碍物类型
	MinX         int32  // 区域最小X坐标
	MinY         int32  // 区域最小Y坐标
	MaxX         int32  // 区域最大X坐标
	MaxY         int32  // 区域最大Y坐标

	// 密度和分布
	Density     float32 // 障碍物密度（0.0-1.0）
	MinSize     int32   // 最小障碍物尺寸
	MaxSize     int32   // 最大障碍物尺寸
	MinDistance int32   // 障碍物间最小距离

	// 影响配置
	BlockBuilding bool // 是否阻挡建筑摆放
	BlockResource bool // 是否阻挡资源刷新
	BlockMonster  bool // 是否阻挡怪物刷新
	AllowMarch    bool // 是否允许行军通过

	// 地形效果
	TerrainEffects map[string]float32 // 地形效果，如{"movement_speed": 0.8, "vision_range": 0.7}
}
