package main

import (
	"fmt"
	"time"

	"github.com/GooLuck/WorldMap/internal/idgen"
	"github.com/GooLuck/WorldMap/internal/worldmap"
)

func main() {
	fmt.Println("=== 测试ID生成模块 ===")

	// 测试1: 本地ID生成器
	fmt.Println("\n1. 测试本地ID生成器:")
	testLocalGenerator()

	// 测试2: 全局ID生成器
	fmt.Println("\n2. 测试全局ID生成器:")
	testGlobalGenerator()

	// 测试3: ID唯一性测试
	fmt.Println("\n3. ID唯一性测试:")
	testUniqueness()

	// 测试4: ID解析测试
	fmt.Println("\n4. ID解析测试:")
	testParsing()

	// 测试5: 跨地图传输模拟
	fmt.Println("\n5. 跨地图传输模拟:")
	testCrossMapTransfer()

	fmt.Println("\n=== 测试完成 ===")
}

func testLocalGenerator() {
	// 创建本地ID生成器
	generator, err := idgen.NewLocalGenerator(1)
	if err != nil {
		fmt.Printf("创建本地生成器失败: %v\n", err)
		return
	}

	// 生成单个ID
	id1, err := generator.GenerateID()
	if err != nil {
		fmt.Printf("生成ID失败: %v\n", err)
		return
	}
	fmt.Printf("生成的ID: %d (0x%016x)\n", id1, id1)

	// 批量生成ID
	ids, err := generator.GenerateIDs(5)
	if err != nil {
		fmt.Printf("批量生成ID失败: %v\n", err)
		return
	}
	fmt.Printf("批量生成的ID: %v\n", ids)

	// 解析ID
	timestamp, machineID, sequence := generator.ParseID(id1)
	fmt.Printf("解析ID: 时间戳=%d, 机器ID=%d, 序列号=%d\n", timestamp, machineID, sequence)
	fmt.Printf("创建时间: %s\n", generator.GetTimestamp(id1).Format("2006-01-02 15:04:05"))
}

func testGlobalGenerator() {
	// 初始化全局ID生成器（本地模式）
	err := worldmap.InitIDGenerator(2, false, "")
	if err != nil {
		fmt.Printf("初始化全局ID生成器失败: %v\n", err)
		return
	}

	// 使用便捷函数生成ID
	id, err := worldmap.GenerateUnitID()
	if err != nil {
		fmt.Printf("生成Unit ID失败: %v\n", err)
		return
	}
	fmt.Printf("生成的Unit ID: %d\n", id)

	// 批量生成
	ids, err := worldmap.GenerateUnitIDs(3)
	if err != nil {
		fmt.Printf("批量生成Unit ID失败: %v\n", err)
		return
	}
	fmt.Printf("批量生成的Unit IDs: %v\n", ids)

	// 验证ID有效性
	valid := worldmap.GetIDGenerator().IsValidUnitID(id)
	fmt.Printf("ID有效性检查: %v\n", valid)
}

func testUniqueness() {
	// 测试ID唯一性
	generator, err := idgen.NewLocalGenerator(1)
	if err != nil {
		fmt.Printf("创建生成器失败: %v\n", err)
		return
	}

	// 生成1000个ID，检查是否重复
	idSet := make(map[int64]bool)
	duplicates := 0

	start := time.Now()
	for i := 0; i < 1000; i++ {
		id, err := generator.GenerateID()
		if err != nil {
			fmt.Printf("生成ID失败: %v\n", err)
			return
		}

		if idSet[id] {
			duplicates++
			fmt.Printf("发现重复ID: %d (第%d次)\n", id, i)
		}
		idSet[id] = true
	}
	elapsed := time.Since(start)

	fmt.Printf("生成1000个ID，耗时: %v\n", elapsed)
	fmt.Printf("重复ID数量: %d\n", duplicates)
	fmt.Printf("唯一ID数量: %d\n", len(idSet))

	if duplicates == 0 {
		fmt.Println("✓ ID唯一性测试通过")
	} else {
		fmt.Println("✗ ID唯一性测试失败")
	}
}

func testParsing() {
	// 测试ID解析功能
	generator, err := idgen.NewLocalGenerator(3)
	if err != nil {
		fmt.Printf("创建生成器失败: %v\n", err)
		return
	}

	id, err := generator.GenerateID()
	if err != nil {
		fmt.Printf("生成ID失败: %v\n", err)
		return
	}

	// 使用idgen包的Parse函数
	timestamp, machineID, sequence := idgen.Parse(id)
	createTime := time.UnixMilli(timestamp)

	fmt.Printf("原始ID: %d\n", id)
	fmt.Printf("十六进制: 0x%016x\n", id)
	fmt.Printf("解析结果:\n")
	fmt.Printf("  时间戳: %d (%s)\n", timestamp, createTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  机器ID: %d\n", machineID)
	fmt.Printf("  序列号: %d\n", sequence)

	// 验证解析的正确性
	expectedMachineID := int64(3)
	if machineID == expectedMachineID {
		fmt.Printf("✓ 机器ID正确: %d\n", machineID)
	} else {
		fmt.Printf("✗ 机器ID错误: 期望 %d, 实际 %d\n", expectedMachineID, machineID)
	}

	// 检查时间戳是否合理
	now := time.Now().UnixMilli()
	if timestamp > 1704067200000 && timestamp <= now+3600000 { // 2024年之后，1小时未来之内
		fmt.Printf("✓ 时间戳合理: %d\n", timestamp)
	} else {
		fmt.Printf("✗ 时间戳异常: %d\n", timestamp)
	}
}

func testCrossMapTransfer() {
	// 模拟跨地图传输场景
	fmt.Println("模拟跨地图传输场景:")

	// 初始化两个不同的地图服务器（使用不同的机器ID）
	err := worldmap.InitIDGenerator(10, false, "") // 地图服务器1
	if err != nil {
		fmt.Printf("初始化地图服务器1失败: %v\n", err)
		return
	}

	// 在地图服务器1上创建unit
	unitID1, err := worldmap.GenerateUnitID()
	if err != nil {
		fmt.Printf("生成Unit ID失败: %v\n", err)
		return
	}
	fmt.Printf("地图服务器1生成的Unit ID: %d\n", unitID1)

	// 模拟unit传输到地图服务器2
	// 在实际场景中，地图服务器2会有不同的机器ID
	// 但ID本身是全局唯一的，所以可以跨服务器识别

	// 解析ID信息
	timestamp, machineID, sequence := worldmap.ParseUnitID(unitID1)
	createTime := time.UnixMilli(timestamp)

	fmt.Printf("Unit ID解析结果:\n")
	fmt.Printf("  生成时间: %s\n", createTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("  生成服务器(机器ID): %d\n", machineID)
	fmt.Printf("  序列号: %d\n", sequence)

	// 验证ID在不同服务器上的有效性
	// 假设地图服务器2初始化了不同的机器ID
	err = worldmap.InitIDGenerator(20, false, "") // 重新初始化（实际应该用不同的实例）
	if err != nil {
		fmt.Printf("初始化地图服务器2失败: %v\n", err)
		return
	}

	valid := worldmap.GetIDGenerator().IsValidUnitID(unitID1)
	fmt.Printf("在地图服务器2上验证ID有效性: %v\n", valid)

	if valid {
		fmt.Println("✓ Unit ID可以跨地图服务器识别")
	} else {
		fmt.Println("✗ Unit ID跨地图服务器识别失败")
	}

	// 生成新的unit ID（在地图服务器2上）
	unitID2, err := worldmap.GenerateUnitID()
	if err != nil {
		fmt.Printf("在地图服务器2上生成Unit ID失败: %v\n", err)
		return
	}
	fmt.Printf("地图服务器2生成的Unit ID: %d\n", unitID2)

	// 检查两个ID是否不同
	if unitID1 != unitID2 {
		fmt.Println("✓ 不同服务器生成的ID不同")
	} else {
		fmt.Println("✗ 不同服务器生成的ID相同（冲突）")
	}
}
