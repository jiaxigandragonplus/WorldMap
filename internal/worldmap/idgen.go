package worldmap

import (
	"sync"

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
func (g *IDGenerator) GenerateNewID() int64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.useRemote && g.client != nil {
		id, _ := g.client.GenerateID()
		return id
	} else if g.localGenerator != nil {
		id, _ := g.localGenerator.GenerateID()
		return id
	}

	// 回退到简单自增（仅用于测试）
	return 0
}
