# SLG地图资源刷新策略配置指南

## 概述

本文档详细介绍了SLG地图项目中资源刷新策略的配置方法，包括资源点类型、刷新策略、位置配置等。

## 系统架构

### 核心组件

1. **ResourceUnit** (`internal/worldmap/resource.go`)
   - 资源点单位的具体实现
   - 管理资源点的状态、采集和恢复逻辑

2. **ResourceManager** (`internal/worldmap/resourcemgr.go`)
   - 管理所有资源点的生命周期
   - 处理全局刷新和区域刷新逻辑
   - 实现动态平衡算法

3. **资源配置** (`internal/worldmap/config/map.go`)
   - 增强的资源点配置结构
   - 资源区域配置
   - 全局刷新配置

## 资源配置详解

### 1. 资源点类型 (ResourcePointType)

| 类型 | 常量值 | 描述 |
|------|--------|------|
| 固定资源点 | `ResourcePointType_Fixed` | 永久存在，资源被采集后按配置速率恢复 |
| 随机刷新点 | `ResourcePointType_RandomSpawn` | 在指定半径内随机位置刷新，有刷新间隔和概率 |
| 季节性资源点 | `ResourcePointType_Seasonal` | 只在特定月份和时间段活跃 |
| 事件触发资源点 | `ResourcePointType_EventBased` | 由游戏事件触发出现 |

### 2. 刷新策略 (RefreshStrategy)

| 策略 | 常量值 | 描述 |
|------|--------|------|
| 线性恢复 | `RefreshStrategy_Linear` | 每秒按固定速率恢复资源 |
| 指数恢复 | `RefreshStrategy_Exponential` | 恢复速率随时间指数增长 |
| 阶梯式恢复 | `RefreshStrategy_Stepwise` | 每间隔一段时间恢复固定量 |
| 随机恢复 | `RefreshStrategy_Random` | 恢复量在范围内随机 |

### 3. 增强资源点配置 (EnhancedResourcePointConfig)

```go
type EnhancedResourcePointConfig struct {
    PointID        int32              // 资源点ID
    X, Y           int32              // 坐标位置
    ResourceType   string             // 资源类型
    PointType      ResourcePointType  // 资源点类型
    RefreshStrategy RefreshStrategy   // 刷新策略
    
    // 基础属性
    MaxAmount      int32              // 最大资源量
    CurrentAmount  int32              // 当前资源量
    RegenRate      float32            // 基础恢复速率（单位/秒）
    RegenDelay     int32              // 恢复延迟（秒）
    
    // 随机刷新相关
    SpawnRadius    int32              // 刷新半径
    SpawnInterval  int32              // 刷新间隔（秒）
    SpawnChance    float32            // 刷新概率（0.0-1.0）
    
    // 时间相关配置
    ActiveHours    []int32            // 活跃时间段（小时）
    SeasonMonths   []int32            // 活跃月份（1-12）
    
    // 事件触发
    TriggerEventID int32              // 触发事件ID
    DespawnAfter   int32              // 消失时间（秒）
    
    // 玩家限制
    MinPlayerLevel int32              // 最小玩家等级可见
    MaxPlayerLevel int32              // 最大玩家等级可见
    FactionRestrict string            // 阵营限制
}
```

### 4. 资源区域配置 (ResourceZoneConfig)

```go
type ResourceZoneConfig struct {
    ZoneID        int32    // 区域ID
    ZoneName      string   // 区域名称
    ZoneType      string   // 区域类型
    MinX, MinY    int32    // 区域最小坐标
    MaxX, MaxY    int32    // 区域最大坐标
    
    // 资源分布
    ResourceTypes []string // 可能出现的资源类型
    Density       float32  // 资源密度（0.0-1.0）
    MinDistance   int32    // 资源点最小间距
    MaxPoints     int32    // 最大资源点数
    
    // 刷新规则
    RefreshEnabled bool    // 是否启用自动刷新
    RefreshRate    float32 // 刷新速率（点/小时）
    MaxConcurrent  int32   // 最大同时存在资源点数
}
```

### 5. 全局刷新配置 (GlobalRefreshConfig)

```go
type GlobalRefreshConfig struct {
    // 时间相关
    DailyRefreshTime   string  // 每日刷新时间
    WeeklyRefreshDay   int32   // 每周刷新日
    MonthlyRefreshDay  int32   // 每月刷新日
    
    // 资源平衡
    EnableDynamicBalance bool   // 启用动态平衡
    MinResourceRatio     float32 // 最小资源比例
    MaxResourceRatio     float32 // 最大资源比例
    
    // 玩家相关
    PlayerBasedRefresh  bool    // 基于玩家数量的刷新
    PlayersPerResource  int32   // 每个资源点对应的玩家数
    MaxRefreshDistance  int32   // 最大刷新距离
    
    // 性能限制
    MaxRefreshPerTick   int32   // 每Tick最大刷新数量
    RefreshTickInterval int32   // 刷新Tick间隔（秒）
}
```

## 配置示例

### 示例1：固定金矿点
```json
{
  "point_id": 1001,
  "x": 300,
  "y": 300,
  "resource_type": "gold",
  "point_type": "fixed",
  "refresh_strategy": "linear",
  "max_amount": 5000,
  "regen_rate": 2.5,
  "regen_delay": 600
}
```

### 示例2：随机刷新木材点
```json
{
  "point_id": 1002,
  "x": 800,
  "y": 800,
  "resource_type": "wood",
  "point_type": "random_spawn",
  "refresh_strategy": "random",
  "max_amount": 1000,
  "regen_rate": 0.5,
  "regen_delay": 1800,
  "spawn_radius": 200,
  "spawn_interval": 7200,
  "spawn_chance": 0.6
}
```

### 示例3：季节性草药点
```json
{
  "point_id": 1003,
  "x": 1200,
  "y": 200,
  "resource_type": "herb",
  "point_type": "seasonal",
  "refresh_strategy": "exponential",
  "max_amount": 3000,
  "regen_rate": 1.0,
  "regen_delay": 1200,
  "active_hours": [6, 18],
  "season_months": [3, 4, 5, 9, 10, 11]
}
```

## 使用流程

### 1. 初始化资源管理器
```go
// 创建网格管理器
gridMgr := NewGridManager(width, height, xGridLen, yGridLen)

// 创建全局配置
globalConfig := &config.GlobalRefreshConfig{
    DailyRefreshTime: "04:00",
    EnableDynamicBalance: true,
    // ... 其他配置
}

// 创建资源管理器
resourceMgr := NewResourceManager(gridMgr, globalConfig)

// 加载地图配置
resourceMgr.LoadConfig(mapConfig)
```

### 2. 更新资源状态
```go
// 在游戏主循环中调用
func gameLoop() {
    for {
        now := time.Now()
        resourceMgr.Update(now)
        time.Sleep(time.Second)
    }
}
```

### 3. 玩家采集资源
```go
func playerHarvestResource(player *Player, resourceId int64, amount int32) {
    harvested := resourceMgr.HarvestResource(
        resourceId, 
        amount, 
        player.Level, 
        player.Faction
    )
    
    if harvested > 0 {
        player.AddResource(resourceType, harvested)
    }
}
```

## 最佳实践

### 1. 资源点布局策略
- **固定资源点**：放置在重要战略位置，如据点附近
- **随机刷新点**：分布在野外区域，增加探索乐趣
- **季节性资源点**：模拟现实世界季节变化，增加游戏真实感
- **事件触发点**：用于特殊活动或任务奖励

### 2. 刷新参数调优
- **恢复速率**：根据资源稀有度调整，稀有资源恢复慢
- **刷新间隔**：平衡玩家等待时间和资源可用性
- **刷新概率**：控制资源点密度，避免过于密集或稀疏

### 3. 性能优化
- 使用资源区域管理大片区域
- 限制每Tick刷新数量，避免性能峰值
- 实现动态平衡，根据玩家数量调整资源总量

### 4. 游戏平衡
- 监控资源消耗与产出的比例
- 根据玩家反馈调整资源配置
- 定期分析资源分布数据，优化游戏体验

## 扩展功能

### 1. 自定义刷新策略
可以通过实现 `RefreshStrategy` 接口创建自定义刷新策略：
```go
type CustomRefreshStrategy struct {
    // 自定义参数
}

func (s *CustomRefreshStrategy) CalculateRegen(elapsedTime float64, config *ResourceConfig) int32 {
    // 自定义计算逻辑
}
```

### 2. 资源点事件系统
可以扩展资源点支持更多事件：
- `OnHarvest`：资源被采集时触发
- `OnRegenComplete`：资源恢复完成时触发
- `OnSpawn`：资源点生成时触发
- `OnDespawn`：资源点消失时触发

### 3. 数据分析与监控
集成数据分析功能：
- 资源采集统计
- 资源分布热图
- 玩家行为分析
- 平衡性指标监控

## 故障排除

### 常见问题

1. **资源点不刷新**
   - 检查 `RefreshEnabled` 配置
   - 验证时间限制条件（ActiveHours, SeasonMonths）
   - 检查刷新概率设置

2. **资源恢复速度异常**
   - 验证 `RegenRate` 和 `RegenDelay` 配置
   - 检查刷新策略计算逻辑
   - 确认时间单位是否正确

3. **性能问题**
   - 调整 `MaxRefreshPerTick` 限制
   - 优化资源点查询算法
   - 考虑分区域更新策略

### 调试工具

建议实现以下调试工具：
- 资源点状态查看器
- 刷新日志记录
- 实时配置热重载
- 性能监控面板

## 总结

本资源刷新策略系统提供了高度可配置的资源管理方案，支持多种资源点类型和刷新策略。通过合理的配置，可以实现丰富的游戏体验和良好的游戏平衡性。

详细配置示例请参考 `config_example.json` 文件。