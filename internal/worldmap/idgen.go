package worldmap

import (
	"sync"
	"time"

	"github.com/GooLuck/WorldMap/internal/idgen"
)

// IDGenerator 全局ID生成器管理器
type IDGenerator struct {
	localGenerator *idgen.LocalGenerator
	client         *idgen.Client
	useRemote      bool
	mu             sync.RWMutex
}

var (
	globalIDGenerator *IDGenerator
	once              sync.Once
)

// InitIDGenerator 初始化全局ID生成器
// machineID: 机器ID，范围0-1023
// useRemote: 是否使用远程ID生成服务
// remoteURL: 远程服务URL（如果useRemote为true）
func InitIDGenerator(machineID int64, useRemote bool, remoteURL string) error {
	var initErr error
	once.Do(func() {
		if useRemote {
			// 使用远程ID生成服务
			client := idgen.NewClient(remoteURL)
			globalIDGenerator = &IDGenerator{
				client:    client,
				useRemote: true,
			}
		} else {
			// 使用本地ID生成器
			localGen, err := idgen.NewLocalGenerator(machineID)
			if err != nil {
				initErr = err
				return
			}
			globalIDGenerator = &IDGenerator{
				localGenerator: localGen,
				useRemote:      false,
			}
		}
	})
	return initErr
}

// GetIDGenerator 获取全局ID生成器实例
func GetIDGenerator() *IDGenerator {
	if globalIDGenerator == nil {
		// 使用默认配置初始化（机器ID为1，本地模式）
		InitIDGenerator(1, false, "")
	}
	return globalIDGenerator
}

// GenerateUnitID 生成Unit ID
func (g *IDGenerator) GenerateUnitID() (int64, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.useRemote && g.client != nil {
		return g.client.GenerateID()
	} else if g.localGenerator != nil {
		return g.localGenerator.GenerateID()
	}

	// 回退到简单自增（仅用于测试）
	return 0, nil
}

// GenerateUnitIDs 批量生成Unit ID
func (g *IDGenerator) GenerateUnitIDs(count int) ([]int64, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.useRemote && g.client != nil {
		return g.client.GenerateIDs(count)
	} else if g.localGenerator != nil {
		return g.localGenerator.GenerateIDs(count)
	}

	// 回退到简单自增（仅用于测试）
	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		ids[i] = int64(i + 1)
	}
	return ids, nil
}

// ParseUnitID 解析Unit ID
func (g *IDGenerator) ParseUnitID(id int64) (timestamp int64, machineID int64, sequence int64) {
	if g.localGenerator != nil {
		return g.localGenerator.ParseID(id)
	}

	// 使用idgen包的Parse函数
	return idgen.Parse(id)
}

// GetUnitIDTimestamp 获取Unit ID的时间戳
func (g *IDGenerator) GetUnitIDTimestamp(id int64) string {
	if g.localGenerator != nil {
		return g.localGenerator.GetTimestamp(id).Format("2006-01-02 15:04:05")
	}

	// 使用idgen包的GetTimestamp函数
	return idgen.GetTimestamp(id).Format("2006-01-02 15:04:05")
}

// IsValidUnitID 检查Unit ID是否有效
func (g *IDGenerator) IsValidUnitID(id int64) bool {
	if id <= 0 {
		return false
	}

	// 检查时间戳是否在合理范围内（2024年之后）
	timestamp, _, _ := g.ParseUnitID(id)
	if timestamp < 1704067200000 { // 2024-01-01之前
		return false
	}

	// 检查时间戳是否在未来（允许一定的时钟偏移）
	if timestamp > time.Now().UnixMilli()+3600000 { // 1小时未来
		return false
	}

	return true
}

// 辅助函数

// GenerateUnitID 生成单个Unit ID（便捷函数）
func GenerateUnitID() (int64, error) {
	return GetIDGenerator().GenerateUnitID()
}

// MustGenerateUnitID 生成单个Unit ID，如果失败则panic
func MustGenerateUnitID() int64 {
	id, err := GenerateUnitID()
	if err != nil {
		panic("Failed to generate unit ID: " + err.Error())
	}
	return id
}

// GenerateUnitIDs 批量生成Unit ID（便捷函数）
func GenerateUnitIDs(count int) ([]int64, error) {
	return GetIDGenerator().GenerateUnitIDs(count)
}

// ParseUnitID 解析Unit ID（便捷函数）
func ParseUnitID(id int64) (timestamp int64, machineID int64, sequence int64) {
	return GetIDGenerator().ParseUnitID(id)
}
