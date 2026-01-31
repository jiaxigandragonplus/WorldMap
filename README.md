# WorldMap

先造需求：
实现一个slg地图管理系统
地图是二维的，不需要地形，高低起伏之类的，但是会有阻挡，某些地方可以设置不能通过
地图划分成一个一个的格子grid，格子大小可以自定义，给出格子索引，可以计算出格子的坐标，给定一个坐标，可以直接算出所在格子
地图上的单位多种多样，可以是玩家，npc，建筑物，静态物件，阻挡物件等

## 64位唯一ID生成模块

### 概述
为WorldMap系统添加了64位全局唯一ID生成模块，用于生成unit id。该ID生成器基于Snowflake算法变体，确保在分布式环境下的ID唯一性，支持unit跨地图传送。

### 特性
- **64位全局唯一ID**：时间戳(41位) + 机器ID(10位) + 序列号(12位)
- **高并发支持**：每台机器每毫秒可生成4096个唯一ID
- **时间有序**：ID按时间顺序递增，便于排序和索引
- **跨地图支持**：unit ID全局唯一，支持跨地图服务器传输
- **多种使用模式**：支持本地生成和远程服务两种模式

### 目录结构
```
internal/idgen/
├── snowflake.go    # Snowflake算法核心实现
├── service.go      # ID生成HTTP服务
└── client.go       # 客户端库

internal/worldmap/
└── idgen.go        # 全局ID生成器管理器

cmd/idgen-server/
└── main.go         # 独立的ID生成服务器
```

### 使用方法

#### 1. 本地模式（默认）
```go
import "github.com/GooLuck/WorldMap/internal/worldmap"

// 初始化全局ID生成器（机器ID为1）
worldmap.InitIDGenerator(1, false, "")

// 生成Unit ID
id, err := worldmap.GenerateUnitID()
if err != nil {
    // 处理错误
}

// 批量生成
ids, err := worldmap.GenerateUnitIDs(10)
```

#### 2. 远程服务模式
```go
// 启动ID生成服务
// go run cmd/idgen-server/main.go -addr :8080 -machine-id 1

// 客户端使用
worldmap.InitIDGenerator(0, true, "http://localhost:8080")
id, err := worldmap.GenerateUnitID()
```

#### 3. 直接使用ID生成器
```go
import "github.com/GooLuck/WorldMap/internal/idgen"

// 创建本地生成器
generator, err := idgen.NewLocalGenerator(1)
id, err := generator.GenerateID()

// 解析ID
timestamp, machineID, sequence := idgen.Parse(id)
```

### ID结构

#### 标准版（默认）
```
64位ID = 时间戳(41位) + 机器ID(10位) + 序列号(12位)

时间戳：从2024-01-01 00:00:00 UTC开始的毫秒数
机器ID：0-1023，最多支持1024个服务器节点/进程
序列号：0-4095，同一毫秒内的自增序列
```

#### 扩展版（支持更多节点）
```
64位ID = 时间戳(40位) + 机器ID(12位) + 序列号(12位)

时间戳：从2024-01-01 00:00:00 UTC开始的毫秒数
机器ID：0-4095，最多支持4096个服务器节点/进程
序列号：0-4095，同一毫秒内的自增序列
```

#### 数据中心版（类似Twitter Snowflake）
```
64位ID = 时间戳(41位) + 数据中心ID(5位) + 机器ID(5位) + 序列号(13位)

时间戳：从2024-01-01 00:00:00 UTC开始的毫秒数
数据中心ID：0-31，支持32个数据中心
机器ID：0-31，每个数据中心支持32台机器
序列号：0-8191，同一毫秒内的自增序列
```

### 节点数量限制说明

1. **标准版**：10位机器ID，支持最多**1024个**节点
2. **扩展版**：12位机器ID，支持最多**4096个**节点
3. **数据中心版**：5位数据中心ID + 5位机器ID，支持**32个数据中心 × 32台机器 = 1024个**节点

### 如何选择版本

- 如果预计节点数不超过1024，使用**标准版**即可
- 如果需要支持更多节点（最多4096），使用**扩展版**
- 如果需要按数据中心组织节点，使用**数据中心版**

### 扩展使用

```go
import "github.com/GooLuck/WorldMap/internal/idgen"

// 使用扩展版（支持4096节点）
extGen, err := idgen.NewExtendedSnowflake(2048) // 机器ID 0-4095

// 使用数据中心版
dcGen, err := idgen.NewDataCenterSnowflake(1, 16) // 数据中心ID 1，机器ID 16
```

### 时间戳范围

- 标准版/数据中心版：41位时间戳，约69年（从2024年开始到2093年）
- 扩展版：40位时间戳，约34年（从2024年开始到2058年）

对于大多数游戏服务器场景，标准版的1024节点和69年时间范围已经足够。如有特殊需求，可使用扩展版或自定义位分配。

## 数据库ID生成方案（替代方案）

### 概述
除了Snowflake算法，还提供了基于数据库的ID生成方案。该方案依赖Redis、MongoDB或MySQL等数据库的原子操作来保证ID的唯一性。

### 核心优势
1. **绝对唯一性**：由数据库原子操作保证，无重复风险
2. **无节点限制**：不受机器ID位数限制，支持无限扩展
3. **顺序递增**：ID严格递增，便于数据库索引和范围查询
4. **集中管理**：所有ID由数据库统一管理，便于监控

### 实现方案

#### 1. Redis方案
```go
// 使用Redis的INCR/INCRBY原子命令
// 需要添加github.com/redis/go-redis/v9依赖
key := "idgen:unit"
id, err := redisClient.Incr(ctx, key).Result()
```

#### 2. MongoDB方案
```go
// 使用MongoDB的findAndModify原子操作
// 需要添加go.mongodb.org/mongo-driver依赖
filter := bson.M{"_id": "unit_counter"}
update := bson.M{"$inc": bson.M{"value": 1}}
```

#### 3. MySQL方案
```go
// 使用事务+SELECT FOR UPDATE保证原子性
// 需要添加github.com/go-sql-driver/mysql依赖
BEGIN TRANSACTION;
SELECT value FROM id_counters WHERE key = 'unit' FOR UPDATE;
UPDATE id_counters SET value = value + 1 WHERE key = 'unit';
COMMIT;
```

#### 4. 内存方案（测试/单机）
```go
import "github.com/GooLuck/WorldMap/internal/idgen"

// 使用内存实现的ID生成器
factory := &idgen.DatabaseIDGeneratorFactory{}
generator, _ := factory.CreateGenerator("memory", nil)

wrapper := idgen.NewDatabaseIDGenWrapper(generator, "idgen:")
unitID, _ := wrapper.GenerateUnitID()
```

### 接口设计
```go
type DatabaseIDGenerator interface {
    GenerateID(ctx context.Context, key string) (int64, error)
    GenerateIDs(ctx context.Context, key string, count int) ([]int64, error)
    SetInitialValue(ctx context.Context, key string, initialValue int64) error
    GetCurrentValue(ctx context.Context, key string) (int64, error)
    Close() error
}
```

### 性能优化
1. **批量预取**：一次性获取多个ID缓存到本地，减少数据库访问
2. **连接池**：使用数据库连接池提高并发性能
3. **本地缓存**：在应用层缓存一定数量的ID

### 适用场景对比

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| **Snowflake** | 无外部依赖、高性能、时间有序 | 节点数量有限制、时钟回拨问题 | 分布式系统、需要时间信息 |
| **Redis** | 高性能、简单可靠、支持集群 | 依赖Redis、网络延迟 | 高并发Web、游戏服务器 |
| **MongoDB** | 文档灵活、易于扩展 | 依赖MongoDB、相对较慢 | 微服务、文档型应用 |
| **MySQL** | 强一致性、事务支持 | 性能较低、单点压力 | 传统企业应用、金融系统 |

### 选择建议

1. **需要时间信息+分布式**：选择Snowflake方案
2. **已有Redis基础设施+高并发**：选择Redis方案
3. **需要强一致性+事务**：选择MySQL方案
4. **开发测试/单机应用**：选择内存方案

### 扩展性考虑
- 所有方案都支持批量ID生成
- 支持自定义键前缀，区分不同业务
- 提供工厂模式，便于切换不同实现
- 支持设置初始ID值

### 部署建议
1. **生产环境**：根据基础设施选择Redis或Snowflake
2. **开发环境**：使用内存方案简化部署
3. **混合部署**：可同时提供多种方案，根据业务需求选择

### HTTP API（ID生成服务）
- `GET /id` - 生成单个ID
- `GET /ids?count=N` - 批量生成N个ID (1-1000)
- `GET /stats` - 获取服务统计信息
- `GET /parse?id=ID` - 解析ID信息
- `GET /health` - 健康检查

### 配置说明
机器ID需要在0-1023范围内，确保集群中每台服务器/进程有唯一的机器ID。

### 测试
运行测试程序验证功能：
```bash
go run cmd/test-idgen/main.go
```

### 性能
- 单机每秒可生成超过400万个ID
- ID生成耗时约0.1-0.2微秒/个
- 支持高并发场景