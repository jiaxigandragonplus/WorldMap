package idgen

import (
	"errors"
	"sync"
	"time"
)

// CustomIDGenerator 自定义ID生成器
// 编码规则: 32位时间戳(秒) + 15位ServerId + 17位自增ID
// 支持每秒生成 2^17 = 131072 个ID
// ServerId范围: 0-32767
// 借ID机制: 允许向未来借1小时ID
type CustomIDGenerator struct {
	serverID      int64      // 15位，范围0-32767
	lastTimestamp int64      // 上次生成ID的时间戳(秒)
	sequence      int64      // 17位自增序列，范围0-131071
	borrowedCount int64      // 已借用的未来ID数量
	mutex         sync.Mutex // 互斥锁
	startTime     int64      // 起始时间戳(秒)
}

const (
	// 位分配
	customTimestampBits = 32
	customServerIDBits  = 15
	customSequenceBits  = 17

	// 最大值
	customMaxTimestamp = (1 << customTimestampBits) - 1
	customMaxServerID  = (1 << customServerIDBits) - 1 // 32767
	customMaxSequence  = (1 << customSequenceBits) - 1 // 131071

	// 位偏移
	customServerIDShift  = customSequenceBits
	customTimestampShift = customSequenceBits + customServerIDBits

	// 时间相关
	customEpochSeconds     int64 = 1704067200                     // 2024-01-01 00:00:00 UTC (秒)
	customMaxBorrowSeconds int64 = 3600                           // 最大允许借1小时
	customMaxBorrowIDs     int64 = 3600 * (customMaxSequence + 1) // 3600 * 131072
)

// NewCustomIDGenerator 创建自定义ID生成器
// serverID: 服务器ID，范围0-32767
func NewCustomIDGenerator(serverID int64) (*CustomIDGenerator, error) {
	if serverID < 0 || serverID > customMaxServerID {
		return nil, errors.New("server ID out of range (0-32767)")
	}

	return &CustomIDGenerator{
		serverID:      serverID,
		lastTimestamp: -1,
		sequence:      0,
		borrowedCount: 0,
		startTime:     time.Now().Unix(),
	}, nil
}

// Generate 生成一个新的64位唯一ID
func (c *CustomIDGenerator) Generate() (int64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 获取当前时间戳(秒)
	now := time.Now().Unix()
	currentTimestamp := now - customEpochSeconds

	if currentTimestamp < 0 {
		return 0, errors.New("clock moved backwards before custom epoch")
	}

	// 检查时间戳是否溢出
	if currentTimestamp > customMaxTimestamp {
		return 0, errors.New("timestamp overflow")
	}

	// 处理时间戳变化
	if currentTimestamp > c.lastTimestamp {
		// 新的时间戳，重置序列号和借用计数
		c.sequence = 0
		c.borrowedCount = 0
		c.lastTimestamp = currentTimestamp
	} else if currentTimestamp == c.lastTimestamp {
		// 同一秒内，递增序列号
		c.sequence++

		// 检查序列号是否溢出
		if c.sequence > customMaxSequence {
			// 序列号溢出，尝试借用未来ID
			return c.borrowFutureID(currentTimestamp)
		}
	} else {
		// 时间回拨（当前时间小于上次记录的时间）
		// 使用借ID机制继续生成
		return c.handleClockBackwards(currentTimestamp)
	}

	// 生成ID
	id := (currentTimestamp << customTimestampShift) |
		(c.serverID << customServerIDShift) |
		c.sequence

	return id, nil
}

// borrowFutureID 借用未来ID
func (c *CustomIDGenerator) borrowFutureID(currentTimestamp int64) (int64, error) {
	// 增加借用计数
	c.borrowedCount++

	// 检查是否超过最大借用限制
	if c.borrowedCount > customMaxBorrowIDs {
		return 0, errors.New("exceeded maximum borrow limit (1 hour)")
	}

	// 计算借用后的时间戳
	// 每借用131072个ID，时间戳增加1秒
	borrowedSeconds := c.borrowedCount / (customMaxSequence + 1)
	borrowedTimestamp := currentTimestamp + borrowedSeconds

	// 检查是否超过最大借用时间
	if borrowedSeconds > customMaxBorrowSeconds {
		return 0, errors.New("exceeded maximum borrow time (1 hour)")
	}

	// 计算借用后的序列号
	borrowedSequence := c.borrowedCount % (customMaxSequence + 1)

	// 生成借用ID
	id := (borrowedTimestamp << customTimestampShift) |
		(c.serverID << customServerIDShift) |
		borrowedSequence

	return id, nil
}

// handleClockBackwards 处理时钟回拨
func (c *CustomIDGenerator) handleClockBackwards(currentTimestamp int64) (int64, error) {
	// 计算时间回拨量
	backwardsSeconds := c.lastTimestamp - currentTimestamp

	// 如果回拨时间在可接受范围内（比如5秒内），继续使用上次时间戳
	// 这样可以避免因时钟微调导致的ID重复
	if backwardsSeconds <= 5 {
		// 使用上次时间戳继续生成
		c.sequence++

		if c.sequence > customMaxSequence {
			// 序列号溢出，借用未来ID
			return c.borrowFutureID(c.lastTimestamp)
		}

		id := (c.lastTimestamp << customTimestampShift) |
			(c.serverID << customServerIDShift) |
			c.sequence

		return id, nil
	}

	// 严重时钟回拨，返回错误
	return 0, errors.New("clock moved backwards significantly")
}

// GenerateBatch 批量生成多个ID
func (c *CustomIDGenerator) GenerateBatch(count int) ([]int64, error) {
	if count <= 0 || count > 10000 {
		return nil, errors.New("count must be between 1 and 10000")
	}

	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		id, err := c.Generate()
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}

	return ids, nil
}

// ParseCustom 解析自定义ID
func ParseCustom(id int64) (timestamp int64, serverID int64, sequence int64) {
	sequence = id & customMaxSequence
	serverID = (id >> customServerIDShift) & customMaxServerID
	timestamp = (id >> customTimestampShift) & customMaxTimestamp

	// 转换回实际时间戳（秒）
	timestamp += customEpochSeconds

	return timestamp, serverID, sequence
}

// GetCustomTimestamp 从ID中提取时间戳
func GetCustomTimestamp(id int64) time.Time {
	timestamp, _, _ := ParseCustom(id)
	return time.Unix(timestamp, 0)
}

// GetCustomServerID 从ID中提取服务器ID
func GetCustomServerID(id int64) int64 {
	_, serverID, _ := ParseCustom(id)
	return serverID
}

// GetCustomSequence 从ID中提取序列号
func GetCustomSequence(id int64) int64 {
	_, _, sequence := ParseCustom(id)
	return sequence
}

// GetBorrowedCount 获取已借用的ID数量
func (c *CustomIDGenerator) GetBorrowedCount() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.borrowedCount
}

// GetStats 获取生成器统计信息
func (c *CustomIDGenerator) GetStats() map[string]interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now().Unix()
	uptime := now - c.startTime

	return map[string]interface{}{
		"server_id":        c.serverID,
		"last_timestamp":   c.lastTimestamp + customEpochSeconds,
		"current_sequence": c.sequence,
		"borrowed_count":   c.borrowedCount,
		"max_borrow_ids":   customMaxBorrowIDs,
		"max_sequence":     customMaxSequence,
		"uptime_seconds":   uptime,
		"start_time":       time.Unix(c.startTime, 0).Format(time.RFC3339),
	}
}

// 性能对比函数

// CompareWithSnowflake 与Snowflake算法对比
func CompareWithSnowflake() map[string]interface{} {
	snowflake := map[string]interface{}{
		"algorithm":        "Snowflake",
		"timestamp_bits":   41,
		"timestamp_unit":   "millisecond",
		"time_range":       "~69 years",
		"machine_id_bits":  10,
		"max_nodes":        1024,
		"sequence_bits":    12,
		"ids_per_ms":       4096,
		"ids_per_second":   4096000,
		"borrow_mechanism": "No",
		"clock_backwards":  "Need handling",
	}

	custom := map[string]interface{}{
		"algorithm":        "Custom (32+15+17)",
		"timestamp_bits":   32,
		"timestamp_unit":   "second",
		"time_range":       "~136 years",
		"server_id_bits":   15,
		"max_nodes":        32768,
		"sequence_bits":    17,
		"ids_per_second":   131072,
		"borrow_mechanism": "Yes (1 hour)",
		"clock_backwards":  "Tolerant (5 seconds)",
	}

	return map[string]interface{}{
		"snowflake": snowflake,
		"custom":    custom,
		"summary": map[string]string{
			"for_large_clusters":  "Custom方案更好（32768节点 vs 1024节点）",
			"for_high_throughput": "Snowflake更好（409.6万/秒 vs 13.1万/秒）",
			"for_time_precision":  "Snowflake更好（毫秒 vs 秒）",
			"for_time_range":      "Custom更好（136年 vs 69年）",
			"recommendation":      "游戏服务器推荐Custom方案，节点多且秒级精度足够",
		},
	}
}

// IsValidCustomID 验证ID有效性
func IsValidCustomID(id int64) bool {
	if id <= 0 {
		return false
	}

	timestamp, serverID, sequence := ParseCustom(id)

	// 检查服务器ID范围
	if serverID < 0 || serverID > customMaxServerID {
		return false
	}

	// 检查序列号范围
	if sequence < 0 || sequence > customMaxSequence {
		return false
	}

	// 检查时间戳范围
	now := time.Now().Unix()
	maxFutureTime := now + customMaxBorrowSeconds // 允许借用1小时

	if timestamp < customEpochSeconds {
		return false
	}

	if timestamp > maxFutureTime {
		return false
	}

	return true
}
