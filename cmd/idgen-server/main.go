package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GooLuck/WorldMap/internal/idgen"
)

func main() {
	// 解析命令行参数
	var (
		addr      = flag.String("addr", ":8080", "HTTP server address")
		machineID = flag.Int64("machine-id", 1, "Machine ID (0-1023)")
	)
	flag.Parse()

	// 验证机器ID
	if *machineID < 0 || *machineID > 1023 {
		log.Fatalf("Machine ID must be between 0 and 1023, got %d", *machineID)
	}

	fmt.Printf("Starting ID Generation Service\n")
	fmt.Printf("  Machine ID: %d\n", *machineID)
	fmt.Printf("  Server Address: %s\n", *addr)
	fmt.Printf("  Custom Epoch: 2024-01-01 00:00:00 UTC\n")
	fmt.Printf("  Max IDs per millisecond: %d\n", 1<<12) // 序列号最大值

	// 创建ID生成服务
	service, err := idgen.NewIDService(*machineID)
	if err != nil {
		log.Fatalf("Failed to create ID service: %v", err)
	}

	// 设置优雅关闭
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// 在goroutine中启动HTTP服务器
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting HTTP server on %s", *addr)
		if err := service.StartHTTPServer(*addr); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// 等待停止信号
	select {
	case err := <-serverErr:
		log.Fatalf("Server error: %v", err)
	case <-stop:
		log.Println("Shutdown signal received")
	}

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.Stop(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("ID generation service stopped")
}

// 简单的健康检查端点（已包含在service.go中，这里只是示例）
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "timestamp": %d}`, time.Now().Unix())
}
