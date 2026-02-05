package idgen

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Client ID生成服务客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient 创建新的ID生成服务客户端
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GenerateID 从远程服务生成单个ID
func (c *Client) GenerateID() (int64, error) {
	url := c.baseURL + "/id"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var result struct {
		ID        int64  `json:"id"`
		Timestamp int64  `json:"timestamp"`
		Hex       string `json:"hex"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.ID, nil
}

// GenerateIDs 从远程服务批量生成ID
func (c *Client) GenerateIDs(count int) ([]int64, error) {
	if count < 1 || count > 1000 {
		return nil, errors.New("count must be between 1 and 1000")
	}

	url := c.baseURL + "/ids?count=" + strconv.Itoa(count)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var result struct {
		IDs       []int64 `json:"ids"`
		Count     int     `json:"count"`
		Timestamp int64   `json:"timestamp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.IDs, nil
}

// LocalGenerator 本地ID生成器（不依赖远程服务）
type LocalGenerator struct {
	generator *Snowflake
}

// NewLocalGenerator 创建本地ID生成器
func NewLocalGenerator(machineID int64) (*LocalGenerator, error) {
	generator, err := NewSnowflake(machineID)
	if err != nil {
		return nil, err
	}

	return &LocalGenerator{
		generator: generator,
	}, nil
}

// GenerateID 生成单个ID
func (lg *LocalGenerator) GenerateID() (int64, error) {
	return lg.generator.Generate()
}

// GenerateIDs 批量生成ID
func (lg *LocalGenerator) GenerateIDs(count int) ([]int64, error) {
	return lg.generator.GenerateBatch(count)
}

// ParseID 解析ID
func (lg *LocalGenerator) ParseID(id int64) (timestamp int64, machineID int64, sequence int64) {
	return Parse(id)
}

// GetTimestamp 获取ID的时间戳
func (lg *LocalGenerator) GetTimestamp(id int64) time.Time {
	return GetTimestamp(id)
}

// GetMachineID 获取ID的机器ID
func (lg *LocalGenerator) GetMachineID(id int64) int64 {
	return GetMachineID(id)
}
