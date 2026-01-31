package idgen

import (
	"errors"
	"sync"
	"time"
)

// ExtendedSnowflake 扩展的Snowflake ID生成器，支持更多机器节点
// 64位ID结构：时间戳(40位) + 机器ID(12位) + 序列号(12位)
// 支持最多4096个机器节点，每毫秒4096个ID
type ExtendedSnowflake struct {
	machineID     int64
	sequence      int64
	lastTimestamp int64
	mutex         sync.Mutex
}

const (
	// 扩展版位分配
	extTimestampBits = 40
	extMachineIDBits = 12
	extSequenceBits  = 12

	extMaxTimestamp = (1 << extTimestampBits) - 1
	extMaxMachineID = (1 << extMachineIDBits) - 1
	extMaxSequence  = (1 << extSequenceBits) - 1

	// 位偏移
	extMachineIDShift = extSequenceBits
	extTimestampShift = extSequenceBits + extMachineIDBits
)

// NewExtendedSnowflake 创建新的扩展Snowflake ID生成器
// machineID: 机器ID，范围0-4095
func NewExtendedSnowflake(machineID int64) (*ExtendedSnowflake, error) {
	if machineID < 0 || machineID > extMaxMachineID {
		return nil, errors.New("machine ID out of range (0-4095)")
	}

	return &ExtendedSnowflake{
		machineID:     machineID,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

// Generate 生成一个新的64位唯一ID
func (es *ExtendedSnowflake) Generate() (int64, error) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	timestamp := time.Now().UnixMilli() - customEpoch

	if timestamp < 0 {
		return 0, errors.New("clock moved backwards")
	}

	if timestamp == es.lastTimestamp {
		es.sequence = (es.sequence + 1) & extMaxSequence
		if es.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for timestamp <= es.lastTimestamp {
				time.Sleep(time.Microsecond * 100)
				timestamp = time.Now().UnixMilli() - customEpoch
			}
		}
	} else {
		es.sequence = 0
	}

	if timestamp > extMaxTimestamp {
		return 0, errors.New("timestamp overflow")
	}

	es.lastTimestamp = timestamp

	id := (timestamp << extTimestampShift) |
		(es.machineID << extMachineIDShift) |
		es.sequence

	return id, nil
}

// GenerateBatch 批量生成多个ID
func (es *ExtendedSnowflake) GenerateBatch(count int) ([]int64, error) {
	if count <= 0 || count > 10000 {
		return nil, errors.New("count must be between 1 and 10000")
	}

	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		id, err := es.Generate()
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}

	return ids, nil
}

// ParseExtended 解析扩展版Snowflake ID
func ParseExtended(id int64) (timestamp int64, machineID int64, sequence int64) {
	sequence = id & extMaxSequence
	machineID = (id >> extMachineIDShift) & extMaxMachineID
	timestamp = (id >> extTimestampShift) & extMaxTimestamp

	// 转换回实际时间戳（毫秒）
	timestamp += customEpoch

	return timestamp, machineID, sequence
}

// GetExtendedTimestamp 从扩展版ID中提取时间戳
func GetExtendedTimestamp(id int64) time.Time {
	timestamp, _, _ := ParseExtended(id)
	return time.UnixMilli(timestamp)
}

// GetExtendedMachineID 从扩展版ID中提取机器ID
func GetExtendedMachineID(id int64) int64 {
	_, machineID, _ := ParseExtended(id)
	return machineID
}

// 数据中心版Snowflake（类似Twitter原始设计）
// 64位ID结构：时间戳(41位) + 数据中心ID(5位) + 机器ID(5位) + 序列号(13位)
// 支持32个数据中心，每个数据中心32台机器，共1024个节点

type DataCenterSnowflake struct {
	dataCenterID  int64
	machineID     int64
	sequence      int64
	lastTimestamp int64
	mutex         sync.Mutex
}

const (
	dcTimestampBits  = 41
	dcDataCenterBits = 5
	dcMachineIDBits  = 5
	dcSequenceBits   = 13

	dcMaxTimestamp    = (1 << dcTimestampBits) - 1
	dcMaxDataCenterID = (1 << dcDataCenterBits) - 1
	dcMaxMachineID    = (1 << dcMachineIDBits) - 1
	dcMaxSequence     = (1 << dcSequenceBits) - 1

	dcMachineIDShift  = dcSequenceBits
	dcDataCenterShift = dcSequenceBits + dcMachineIDBits
	dcTimestampShift  = dcSequenceBits + dcMachineIDBits + dcDataCenterBits
)

// NewDataCenterSnowflake 创建数据中心版Snowflake ID生成器
// dataCenterID: 数据中心ID，范围0-31
// machineID: 机器ID，范围0-31
func NewDataCenterSnowflake(dataCenterID, machineID int64) (*DataCenterSnowflake, error) {
	if dataCenterID < 0 || dataCenterID > dcMaxDataCenterID {
		return nil, errors.New("data center ID out of range (0-31)")
	}
	if machineID < 0 || machineID > dcMaxMachineID {
		return nil, errors.New("machine ID out of range (0-31)")
	}

	return &DataCenterSnowflake{
		dataCenterID:  dataCenterID,
		machineID:     machineID,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

// Generate 生成一个新的64位唯一ID
func (dcs *DataCenterSnowflake) Generate() (int64, error) {
	dcs.mutex.Lock()
	defer dcs.mutex.Unlock()

	timestamp := time.Now().UnixMilli() - customEpoch

	if timestamp < 0 {
		return 0, errors.New("clock moved backwards")
	}

	if timestamp == dcs.lastTimestamp {
		dcs.sequence = (dcs.sequence + 1) & dcMaxSequence
		if dcs.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for timestamp <= dcs.lastTimestamp {
				time.Sleep(time.Microsecond * 100)
				timestamp = time.Now().UnixMilli() - customEpoch
			}
		}
	} else {
		dcs.sequence = 0
	}

	if timestamp > dcMaxTimestamp {
		return 0, errors.New("timestamp overflow")
	}

	dcs.lastTimestamp = timestamp

	id := (timestamp << dcTimestampShift) |
		(dcs.dataCenterID << dcDataCenterShift) |
		(dcs.machineID << dcMachineIDShift) |
		dcs.sequence

	return id, nil
}

// ParseDataCenter 解析数据中心版Snowflake ID
func ParseDataCenter(id int64) (timestamp int64, dataCenterID int64, machineID int64, sequence int64) {
	sequence = id & dcMaxSequence
	machineID = (id >> dcMachineIDShift) & dcMaxMachineID
	dataCenterID = (id >> dcDataCenterShift) & dcMaxDataCenterID
	timestamp = (id >> dcTimestampShift) & dcMaxTimestamp

	// 转换回实际时间戳（毫秒）
	timestamp += customEpoch

	return timestamp, dataCenterID, machineID, sequence
}

// 根据需求选择合适的生成器
type GeneratorType int

const (
	GeneratorTypeStandard   GeneratorType = iota // 标准版：10位机器ID，1024节点
	GeneratorTypeExtended                        // 扩展版：12位机器ID，4096节点
	GeneratorTypeDataCenter                      // 数据中心版：5位数据中心+5位机器，1024节点
)

// CreateGenerator 根据类型创建ID生成器
func CreateGenerator(genType GeneratorType, param1, param2 int64) (interface{}, error) {
	switch genType {
	case GeneratorTypeStandard:
		return NewSnowflake(param1)
	case GeneratorTypeExtended:
		return NewExtendedSnowflake(param1)
	case GeneratorTypeDataCenter:
		return NewDataCenterSnowflake(param1, param2)
	default:
		return nil, errors.New("unknown generator type")
	}
}
