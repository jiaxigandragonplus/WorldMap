package idgen

import (
	"errors"
	"sync"
	"time"
)

const (
	// Snowflake算法位分配
	// 64位ID结构：时间戳(41位) + 机器ID(10位) + 序列号(12位)
	// 时间戳从自定义epoch开始（2024-01-01 00:00:00 UTC）

	timestampBits = 41
	machineIDBits = 10
	sequenceBits  = 12

	maxTimestamp = (1 << timestampBits) - 1
	maxMachineID = (1 << machineIDBits) - 1
	maxSequence  = (1 << sequenceBits) - 1

	// 自定义epoch：2024-01-01 00:00:00 UTC
	customEpoch int64 = 1704067200000 // 毫秒

	// 位偏移
	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits
)

// Snowflake ID生成器
type Snowflake struct {
	machineID     int64
	sequence      int64
	lastTimestamp int64
	mutex         sync.Mutex
}

// NewSnowflake 创建新的Snowflake ID生成器
// machineID: 机器ID，范围0-1023
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machine ID out of range")
	}

	return &Snowflake{
		machineID:     machineID,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

// Generate 生成一个新的64位唯一ID
func (s *Snowflake) Generate() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixMilli() - customEpoch

	if timestamp < 0 {
		return 0, errors.New("clock moved backwards")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				time.Sleep(time.Microsecond * 100)
				timestamp = time.Now().UnixMilli() - customEpoch
			}
		}
	} else {
		s.sequence = 0
	}

	if timestamp > maxTimestamp {
		return 0, errors.New("timestamp overflow")
	}

	s.lastTimestamp = timestamp

	id := (timestamp << timestampShift) |
		(s.machineID << machineIDShift) |
		s.sequence

	return id, nil
}

// GenerateBatch 批量生成多个ID
func (s *Snowflake) GenerateBatch(count int) ([]int64, error) {
	if count <= 0 || count > 10000 {
		return nil, errors.New("count must be between 1 and 10000")
	}

	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		id, err := s.Generate()
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}

	return ids, nil
}

// Parse 解析Snowflake ID，返回时间戳、机器ID和序列号
func Parse(id int64) (timestamp int64, machineID int64, sequence int64) {
	sequence = id & maxSequence
	machineID = (id >> machineIDShift) & maxMachineID
	timestamp = (id >> timestampShift) & maxTimestamp

	// 转换回实际时间戳（毫秒）
	timestamp += customEpoch

	return timestamp, machineID, sequence
}

// GetTimestamp 从ID中提取时间戳
func GetTimestamp(id int64) time.Time {
	timestamp, _, _ := Parse(id)
	return time.UnixMilli(timestamp)
}

// GetMachineID 从ID中提取机器ID
func GetMachineID(id int64) int64 {
	_, machineID, _ := Parse(id)
	return machineID
}
