package idgen

import (
	"context"
	"errors"
	"sync"
)

// DatabaseIDGenerator 基于数据库的ID生成器接口
// 支持Redis、MongoDB、MySQL等数据库实现
type DatabaseIDGenerator interface {
	// GenerateID 生成单个ID
	GenerateID(ctx context.Context, key string) (int64, error)

	// GenerateIDs 批量生成ID
	GenerateIDs(ctx context.Context, key string, count int) ([]int64, error)

	// SetInitialValue 设置初始值（如果键不存在）
	SetInitialValue(ctx context.Context, key string, initialValue int64) error

	// GetCurrentValue 获取当前值
	GetCurrentValue(ctx context.Context, key string) (int64, error)

	// Close 关闭连接
	Close() error
}

// MemoryIDGenerator 内存实现的ID生成器（用于测试和单机场景）
type MemoryIDGenerator struct {
	counters map[string]int64
	mutex    sync.RWMutex
}

// NewMemoryIDGenerator 创建内存ID生成器
func NewMemoryIDGenerator() *MemoryIDGenerator {
	return &MemoryIDGenerator{
		counters: make(map[string]int64),
	}
}

// GenerateID 生成单个ID
func (m *MemoryIDGenerator) GenerateID(ctx context.Context, key string) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.counters[key]++
	return m.counters[key], nil
}

// GenerateIDs 批量生成ID
func (m *MemoryIDGenerator) GenerateIDs(ctx context.Context, key string, count int) ([]int64, error) {
	if count <= 0 {
		return nil, errors.New("count must be positive")
	}

	if count > 10000 {
		return nil, errors.New("count cannot exceed 10000")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	startID := m.counters[key] + 1
	m.counters[key] += int64(count)

	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		ids[i] = startID + int64(i)
	}

	return ids, nil
}

// SetInitialValue 设置初始值
func (m *MemoryIDGenerator) SetInitialValue(ctx context.Context, key string, initialValue int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if current, exists := m.counters[key]; exists {
		if current < initialValue {
			m.counters[key] = initialValue
		}
	} else {
		m.counters[key] = initialValue
	}

	return nil
}

// GetCurrentValue 获取当前值
func (m *MemoryIDGenerator) GetCurrentValue(ctx context.Context, key string) (int64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if value, exists := m.counters[key]; exists {
		return value, nil
	}

	return 0, nil
}

// Close 关闭（内存实现无需关闭）
func (m *MemoryIDGenerator) Close() error {
	return nil
}

// DatabaseIDGeneratorFactory 数据库ID生成器工厂
type DatabaseIDGeneratorFactory struct{}

// CreateGenerator 创建ID生成器
func (f *DatabaseIDGeneratorFactory) CreateGenerator(dbType string, config interface{}) (DatabaseIDGenerator, error) {
	switch dbType {
	case "memory":
		return NewMemoryIDGenerator(), nil
	case "redis":
		return f.createRedisGenerator(config)
	case "mongodb":
		return f.createMongoDBGenerator(config)
	case "mysql":
		return f.createMySQLGenerator(config)
	default:
		return nil, errors.New("unsupported database type")
	}
}

// createRedisGenerator 创建Redis ID生成器（需要用户实现）
func (f *DatabaseIDGeneratorFactory) createRedisGenerator(config interface{}) (DatabaseIDGenerator, error) {
	// 这里返回一个错误，提示用户需要实现Redis客户端
	// 实际使用时可以添加github.com/redis/go-redis/v9依赖
	return nil, errors.New("Redis implementation requires redis client dependency. " +
		"Add github.com/redis/go-redis/v9 to go.mod and implement RedisIDGenerator")
}

// createMongoDBGenerator 创建MongoDB ID生成器（需要用户实现）
func (f *DatabaseIDGeneratorFactory) createMongoDBGenerator(config interface{}) (DatabaseIDGenerator, error) {
	// 这里返回一个错误，提示用户需要实现MongoDB客户端
	return nil, errors.New("MongoDB implementation requires mongo client dependency. " +
		"Add go.mongodb.org/mongo-driver to go.mod and implement MongoDBIDGenerator")
}

// createMySQLGenerator 创建MySQL ID生成器（需要用户实现）
func (f *DatabaseIDGeneratorFactory) createMySQLGenerator(config interface{}) (DatabaseIDGenerator, error) {
	// 这里返回一个错误，提示用户需要实现MySQL客户端
	return nil, errors.New("MySQL implementation requires database/sql and MySQL driver. " +
		"Add github.com/go-sql-driver/mysql to go.mod and implement MySQLIDGenerator")
}

// 数据库ID生成器使用示例

// DatabaseIDGenWrapper 数据库ID生成器包装器，提供便捷方法
type DatabaseIDGenWrapper struct {
	generator DatabaseIDGenerator
	keyPrefix string
}

// NewDatabaseIDGenWrapper 创建数据库ID生成器包装器
func NewDatabaseIDGenWrapper(generator DatabaseIDGenerator, keyPrefix string) *DatabaseIDGenWrapper {
	if keyPrefix == "" {
		keyPrefix = "idgen:"
	}

	return &DatabaseIDGenWrapper{
		generator: generator,
		keyPrefix: keyPrefix,
	}
}

// GenerateUnitID 生成Unit ID
func (w *DatabaseIDGenWrapper) GenerateUnitID() (int64, error) {
	ctx := context.Background()
	return w.generator.GenerateID(ctx, w.keyPrefix+"unit")
}

// GenerateUnitIDs 批量生成Unit ID
func (w *DatabaseIDGenWrapper) GenerateUnitIDs(count int) ([]int64, error) {
	ctx := context.Background()
	return w.generator.GenerateIDs(ctx, w.keyPrefix+"unit", count)
}

// GenerateResourceID 生成Resource ID
func (w *DatabaseIDGenWrapper) GenerateResourceID() (int64, error) {
	ctx := context.Background()
	return w.generator.GenerateID(ctx, w.keyPrefix+"resource")
}

// GeneratePlayerID 生成Player ID
func (w *DatabaseIDGenWrapper) GeneratePlayerID() (int64, error) {
	ctx := context.Background()
	return w.generator.GenerateID(ctx, w.keyPrefix+"player")
}

// SetUnitIDInitialValue 设置Unit ID初始值
func (w *DatabaseIDGenWrapper) SetUnitIDInitialValue(initialValue int64) error {
	ctx := context.Background()
	return w.generator.SetInitialValue(ctx, w.keyPrefix+"unit", initialValue)
}

// GetCurrentUnitID 获取当前Unit ID值
func (w *DatabaseIDGenWrapper) GetCurrentUnitID() (int64, error) {
	ctx := context.Background()
	return w.generator.GetCurrentValue(ctx, w.keyPrefix+"unit")
}

// Close 关闭生成器
func (w *DatabaseIDGenWrapper) Close() error {
	return w.generator.Close()
}

// 数据库ID生成方案的优势和适用场景说明

/*
数据库ID生成方案的优势：

1. 绝对唯一性：由数据库的原子操作保证，无重复风险
2. 无节点限制：不受机器ID位数限制，支持无限扩展
3. 顺序递增：ID严格递增，便于数据库索引和范围查询
4. 集中管理：所有ID由数据库统一管理，便于监控和调整
5. 持久化：ID状态持久化，服务重启后可以继续

适用场景：

1. Redis方案：
   - 高性能需求，需要低延迟ID生成
   - 已有Redis基础设施
   - 需要支持高并发ID生成
   - 适合：游戏服务器、高并发Web应用

2. MongoDB方案：
   - 已有MongoDB基础设施
   - 需要文档存储的灵活性
   - 适合：微服务架构、文档型应用

3. MySQL方案：
   - 已有MySQL/PostgreSQL基础设施
   - 需要强一致性保证
   - 适合：传统企业应用、金融系统

4. 内存方案：
   - 单机测试环境
   - 开发调试
   - 小型单机应用

实现建议：

1. 对于Redis实现，使用INCR/INCRBY原子命令
2. 对于MongoDB实现，使用findAndModify原子操作
3. 对于MySQL实现，使用事务+SELECT FOR UPDATE
4. 考虑批量预取机制减少数据库压力
5. 添加适当的重试机制处理网络故障
*/

// 批量预取缓存机制示例
type BatchPrefetchIDGenerator struct {
	baseGenerator DatabaseIDGenerator
	cache         map[string][]int64
	cacheMutex    sync.RWMutex
	batchSize     int
	keyPrefix     string
}

// NewBatchPrefetchIDGenerator 创建批量预取ID生成器
func NewBatchPrefetchIDGenerator(baseGenerator DatabaseIDGenerator, batchSize int) *BatchPrefetchIDGenerator {
	return &BatchPrefetchIDGenerator{
		baseGenerator: baseGenerator,
		cache:         make(map[string][]int64),
		batchSize:     batchSize,
		keyPrefix:     "idgen:",
	}
}

// GenerateID 生成ID（带缓存预取）
func (b *BatchPrefetchIDGenerator) GenerateID(ctx context.Context, key string) (int64, error) {
	fullKey := b.keyPrefix + key

	// 尝试从缓存获取
	b.cacheMutex.Lock()
	if ids, exists := b.cache[fullKey]; exists && len(ids) > 0 {
		id := ids[0]
		b.cache[fullKey] = ids[1:]
		b.cacheMutex.Unlock()
		return id, nil
	}
	b.cacheMutex.Unlock()

	// 缓存为空，批量获取
	ids, err := b.baseGenerator.GenerateIDs(ctx, fullKey, b.batchSize)
	if err != nil {
		return 0, err
	}

	// 取第一个，其余放入缓存
	if len(ids) == 0 {
		return 0, errors.New("no IDs generated")
	}

	id := ids[0]
	if len(ids) > 1 {
		b.cacheMutex.Lock()
		b.cache[fullKey] = ids[1:]
		b.cacheMutex.Unlock()
	}

	return id, nil
}
