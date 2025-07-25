package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client PocketBase API 客户端
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// System 表示服务器/系统记录
type System struct {
	ID             string   `json:"id"`
	CollectionID   string   `json:"collectionId"`
	CollectionName string   `json:"collectionName"`
	Created        string   `json:"created"`
	Updated        string   `json:"updated"`
	Name           string   `json:"name"`
	Host           string   `json:"host"`
	Port           string   `json:"port"`
	Status         string   `json:"status"`
	Users          []string `json:"users"`
}

// SystemStats 表示系统统计记录
type SystemStats struct {
	ID             string    `json:"id"`
	CollectionID   string    `json:"collectionId"`
	CollectionName string    `json:"collectionName"`
	Created        string    `json:"created"`
	Updated        string    `json:"updated"`
	System         string    `json:"system"`
	Type           string    `json:"type"`
	Stats          StatsData `json:"stats"`
}

// StatsData 包含实际的指标数据
type StatsData struct {
	CPU         float64 `json:"cpu"`
	Mem         float64 `json:"m"`  // 总内存 GB
	MemUsed     float64 `json:"mu"` // 使用内存 GB
	MemPct      float64 `json:"mp"` // 内存使用百分比
	NetworkSent float64 `json:"ns"` // 网络发送（字节/秒）
	NetworkRecv float64 `json:"nr"` // 网络接收（字节/秒）
}

// 自定义JSON解码，处理stats字段
func (s *SystemStats) UnmarshalJSON(data []byte) error {
	type Alias SystemStats
	aux := &struct {
		Stats interface{} `json:"stats"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 处理stats字段 - 可能是字符串或对象
	switch v := aux.Stats.(type) {
	case string:
		// 如果是字符串，解析为StatsData
		if err := json.Unmarshal([]byte(v), &s.Stats); err != nil {
			return err
		}
	case map[string]interface{}:
		// 如果是对象，转换为StatsData
		statsJSON, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(statsJSON, &s.Stats); err != nil {
			return err
		}
	}

	return nil
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token  string `json:"token"`
	Record struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"record"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

// ListResponse PocketBase API 响应包装器
type ListResponse[T any] struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Items      []T `json:"items"`
}

// NewClient 创建新的PocketBase客户端
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest 向PocketBase API发送HTTP请求
func (pb *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, pb.BaseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if pb.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+pb.AuthToken)
	}

	resp, err := pb.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// Login 用户登录认证
func (pb *Client) Login(email, password string) error {
	loginReq := LoginRequest{
		Identity: email,
		Password: password,
	}

	resp, err := pb.makeRequest("POST", "/api/collections/users/auth-with-password", loginReq)
	if err != nil {
		return fmt.Errorf("登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("登录失败，状态码 %d: %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("解析认证响应失败: %w", err)
	}

	pb.AuthToken = authResp.Token
	return nil
}

// ListSystems 获取所有系统/服务器
func (pb *Client) ListSystems() (*ListResponse[System], error) {
	params := url.Values{}
	params.Set("page", "1")
	params.Set("perPage", "50")
	params.Set("sort", "-created")

	endpoint := "/api/collections/systems/records?" + params.Encode()

	resp, err := pb.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch systems: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ListResponse[System]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetSystemLoadAverage 获取指定系统的负载平均值数据
func (pb *Client) GetSystemLoadAverage(systemID string, count int) (*ListResponse[SystemStats], error) {
	params := url.Values{}
	params.Set("page", "1")
	params.Set("perPage", fmt.Sprintf("%d", count))
	params.Set("sort", "-created")

	// 只过滤该系统的1m类型数据
	filter := fmt.Sprintf(`system = "%s" && type = "1m"`, systemID)
	params.Set("filter", filter)

	endpoint := "/api/collections/system_stats/records?" + params.Encode()

	resp, err := pb.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch system stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ListResponse[SystemStats]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}