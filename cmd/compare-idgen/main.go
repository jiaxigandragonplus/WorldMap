package main

import (
	"fmt"
	"time"

	"github.com/GooLuck/WorldMap/internal/idgen"
)

func main() {
	fmt.Println("=== ID生成算法对比测试 ===")

	// 1. 显示算法对比
	fmt.Println("\n1. 算法参数对比:")
	comparison := idgen.CompareWithSnowflake()

	snowflake := comparison["snowflake"].(map[string]interface{})
	custom := comparison["custom"].(map[string]interface{})
	summary := comparison["summary"].(map[string]string)

	fmt.Println("\nSnowflake算法:")
	fmt.Printf("  时间戳: %d位 (%s), 范围: %s\n",
		snowflake["timestamp_bits"], snowflake["timestamp_unit"], snowflake["time_range"])
	fmt.Printf("  机器ID: %d位, 最大节点: %d\n",
		snowflake["machine_id_bits"], snowflake["max_nodes"])
	fmt.Printf("  序列号: %d位, ID/毫秒: %d, ID/秒: %d\n",
		snowflake["sequence_bits"], snowflake["ids_per_ms"], snowflake["ids_per_second"])

	fmt.Println("\n自定义算法 (32+15+17):")
	fmt.Printf("  时间戳: %d位 (%s), 范围: %s\n",
		custom["timestamp_bits"], custom["timestamp_unit"], custom["time_range"])
	fmt.Printf("  服务器ID: %d位, 最大节点: %d\n",
		custom["server_id_bits"], custom["max_nodes"])
	fmt.Printf("  序列号: %d位, ID/秒: %d\n",
		custom["sequence_bits"], custom["ids_per_second"])
	fmt.Printf("  借ID机制: %s\n", custom["borrow_mechanism"])

	fmt.Println("\n对比总结:")
	for key, value := range summary {
		fmt.Printf("  %s: %s\n", key, value)
	}

	// 2. 实际生成测试
	fmt.Println("\n2. 实际生成测试:")

	// 测试Snowflake
	fmt.Println("\nSnowflake算法测试:")
	snowflakeGen, err := idgen.NewSnowflake(1)
	if err != nil {
		fmt.Printf("创建Snowflake生成器失败: %v\n", err)
	} else {
		start := time.Now()
		snowflakeIDs := make([]int64, 1000)
		for i := 0; i < 1000; i++ {
			id, err := snowflakeGen.Generate()
			if err != nil {
				fmt.Printf("生成ID失败: %v\n", err)
				break
			}
			snowflakeIDs[i] = id
		}
		snowflakeTime := time.Since(start)

		// 检查唯一性
		snowflakeSet := make(map[int64]bool)
		for _, id := range snowflakeIDs {
			snowflakeSet[id] = true
		}

		fmt.Printf("  生成1000个ID耗时: %v\n", snowflakeTime)
		fmt.Printf("  唯一ID数量: %d\n", len(snowflakeSet))

		// 解析示例ID
		if len(snowflakeIDs) > 0 {
			timestamp, machineID, sequence := idgen.Parse(snowflakeIDs[0])
			fmt.Printf("  示例ID解析: 时间戳=%d, 机器ID=%d, 序列号=%d\n",
				timestamp, machineID, sequence)
		}
	}

	// 测试自定义算法
	fmt.Println("\n自定义算法测试:")
	customGen, err := idgen.NewCustomIDGenerator(100)
	if err != nil {
		fmt.Printf("创建自定义生成器失败: %v\n", err)
	} else {
		start := time.Now()
		customIDs := make([]int64, 1000)
		for i := 0; i < 1000; i++ {
			id, err := customGen.Generate()
			if err != nil {
				fmt.Printf("生成ID失败: %v\n", err)
				break
			}
			customIDs[i] = id
		}
		customTime := time.Since(start)

		// 检查唯一性
		customSet := make(map[int64]bool)
		for _, id := range customIDs {
			customSet[id] = true
		}

		fmt.Printf("  生成1000个ID耗时: %v\n", customTime)
		fmt.Printf("  唯一ID数量: %d\n", len(customSet))
		fmt.Printf("  已借用ID数量: %d\n", customGen.GetBorrowedCount())

		// 解析示例ID
		if len(customIDs) > 0 {
			timestamp, serverID, sequence := idgen.ParseCustom(customIDs[0])
			fmt.Printf("  示例ID解析: 时间戳=%d, 服务器ID=%d, 序列号=%d\n",
				timestamp, serverID, sequence)
			fmt.Printf("  时间: %s\n", time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"))
		}

		// 显示统计信息
		stats := customGen.GetStats()
		fmt.Printf("  统计信息: 最后时间戳=%d, 当前序列号=%d\n",
			stats["last_timestamp"], stats["current_sequence"])
	}

	// 3. 性能对比
	fmt.Println("\n3. 性能对比分析:")

	// 模拟高并发场景
	fmt.Println("\n模拟高并发场景 (10万ID生成):")

	// Snowflake性能估算
	snowflakePerMs := 4096
	snowflakePerSecond := snowflakePerMs * 1000
	snowflakeTimeFor100k := float64(100000) / float64(snowflakePerSecond) * 1000

	// 自定义算法性能估算
	customPerSecond := 131072
	customTimeFor100k := float64(100000) / float64(customPerSecond) * 1000

	fmt.Printf("  Snowflake估算: %.2f 毫秒 (理论值)\n", snowflakeTimeFor100k)
	fmt.Printf("  自定义算法估算: %.2f 毫秒 (理论值)\n", customTimeFor100k)

	// 4. 适用场景建议
	fmt.Println("\n4. 适用场景建议:")

	fmt.Println("\n选择Snowflake当:")
	fmt.Println("  - 需要毫秒级时间精度")
	fmt.Println("  - 每秒需要生成超过13万个ID")
	fmt.Println("  - 节点数量不超过1024个")
	fmt.Println("  - 需要时间有序的ID")

	fmt.Println("\n选择自定义算法当:")
	fmt.Println("  - 秒级时间精度足够")
	fmt.Println("  - 需要支持大量节点（最多32768个）")
	fmt.Println("  - 每秒ID需求在13万个以内")
	fmt.Println("  - 需要借ID机制应对瞬时高峰")
	fmt.Println("  - 需要更长的时间范围（136年）")

	// 5. 游戏服务器场景分析
	fmt.Println("\n5. 游戏服务器场景分析:")

	// 典型游戏服务器参数
	playersPerServer := 5000                            // 每台服务器5000玩家
	unitsPerPlayer := 50                                // 每个玩家50个单位
	unitsPerServer := playersPerServer * unitsPerPlayer // 25万单位

	// ID生成需求
	idGenerationRate := unitsPerServer / 3600 // 每秒新生成单位数（假设平均）

	fmt.Printf("  典型游戏服务器参数:\n")
	fmt.Printf("    - 每台服务器玩家: %d\n", playersPerServer)
	fmt.Printf("    - 每个玩家单位数: %d\n", unitsPerPlayer)
	fmt.Printf("    - 每台服务器总单位: %d\n", unitsPerServer)
	fmt.Printf("    - 平均每秒新单位: %d\n", idGenerationRate)

	fmt.Printf("\n  算法满足度:\n")
	fmt.Printf("    Snowflake: 每秒409.6万ID > 需求 %d ID ✓\n", idGenerationRate)
	fmt.Printf("    自定义算法: 每秒13.1万ID > 需求 %d ID ✓\n", idGenerationRate)

	fmt.Printf("\n  节点支持:\n")
	fmt.Printf("    Snowflake: 1024节点 × 5000玩家 = 512万总玩家\n")
	fmt.Printf("    自定义算法: 32768节点 × 5000玩家 = 1.64亿总玩家\n")

	fmt.Println("\n=== 测试完成 ===")
}
