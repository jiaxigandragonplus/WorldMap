package idgen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ID生成服务
type IDService struct {
	generator *Snowflake
	server    *http.Server
	mu        sync.RWMutex
	stats     ServiceStats
}

// ServiceStats 服务统计信息
type ServiceStats struct {
	TotalIDsGenerated int64     `json:"total_ids_generated"`
	StartTime         time.Time `json:"start_time"`
	LastRequestTime   time.Time `json:"last_request_time"`
}

// NewIDService 创建新的ID生成服务
func NewIDService(machineID int64) (*IDService, error) {
	generator, err := NewSnowflake(machineID)
	if err != nil {
		return nil, err
	}

	return &IDService{
		generator: generator,
		stats: ServiceStats{
			StartTime: time.Now(),
		},
	}, nil
}

// GenerateID 生成单个ID
func (s *IDService) GenerateID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := s.generator.Generate()
	if err != nil {
		return 0, err
	}

	s.stats.TotalIDsGenerated++
	s.stats.LastRequestTime = time.Now()

	return id, nil
}

// GenerateIDs 批量生成ID
func (s *IDService) GenerateIDs(count int) ([]int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ids, err := s.generator.GenerateBatch(count)
	if err != nil {
		return nil, err
	}

	s.stats.TotalIDsGenerated += int64(count)
	s.stats.LastRequestTime = time.Now()

	return ids, nil
}

// GetStats 获取服务统计信息
func (s *IDService) GetStats() ServiceStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.stats
}

// HTTP处理函数

// GenerateIDHandler 生成单个ID的HTTP处理函数
func (s *IDService) GenerateIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := s.GenerateID()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate ID: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":        id,
		"timestamp": time.Now().UnixMilli(),
		"hex":       fmt.Sprintf("%016x", id),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GenerateIDsHandler 批量生成ID的HTTP处理函数
func (s *IDService) GenerateIDsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := 1
	if r.URL.Query().Get("count") != "" {
		_, err := fmt.Sscanf(r.URL.Query().Get("count"), "%d", &count)
		if err != nil || count < 1 || count > 1000 {
			http.Error(w, "Count must be between 1 and 1000", http.StatusBadRequest)
			return
		}
	}

	ids, err := s.GenerateIDs(count)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate IDs: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"ids":       ids,
		"count":     len(ids),
		"timestamp": time.Now().UnixMilli(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StatsHandler 获取服务统计信息的HTTP处理函数
func (s *IDService) StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := s.GetStats()
	uptime := time.Since(stats.StartTime)

	response := map[string]interface{}{
		"stats":          stats,
		"uptime_seconds": uptime.Seconds(),
		"uptime_human":   uptime.String(),
		"machine_id":     s.generator.machineID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ParseIDHandler 解析ID的HTTP处理函数
func (s *IDService) ParseIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	timestamp, machineID, sequence := Parse(id)
	createTime := time.UnixMilli(timestamp)

	response := map[string]interface{}{
		"id":               id,
		"hex":              fmt.Sprintf("%016x", id),
		"timestamp":        timestamp,
		"machine_id":       machineID,
		"sequence":         sequence,
		"created_at":       createTime.Format(time.RFC3339),
		"created_at_unix":  timestamp,
		"time_since_epoch": timestamp - customEpoch,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StartHTTPServer 启动HTTP服务器
func (s *IDService) StartHTTPServer(addr string) error {
	mux := http.NewServeMux()

	// 注册路由
	mux.HandleFunc("/id", s.GenerateIDHandler)
	mux.HandleFunc("/ids", s.GenerateIDsHandler)
	mux.HandleFunc("/stats", s.StatsHandler)
	mux.HandleFunc("/parse", s.ParseIDHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("ID generation service started on %s\n", addr)
	fmt.Printf("Endpoints:\n")
	fmt.Printf("  GET /id      - Generate a single ID\n")
	fmt.Printf("  GET /ids?count=N - Generate N IDs (1-1000)\n")
	fmt.Printf("  GET /stats   - Get service statistics\n")
	fmt.Printf("  GET /parse?id=ID - Parse an ID\n")
	fmt.Printf("  GET /health  - Health check\n")

	return s.server.ListenAndServe()
}

// Stop 停止服务
func (s *IDService) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
