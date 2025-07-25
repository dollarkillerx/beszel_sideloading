package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PocketBase API 客户端
type PocketBaseClient struct {
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

// NewPocketBaseClient 创建新的PocketBase客户端
func NewPocketBaseClient(baseURL string) *PocketBaseClient {
	return &PocketBaseClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest 向PocketBase API发送HTTP请求
func (pb *PocketBaseClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
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
func (pb *PocketBaseClient) Login(email, password string) error {
	fmt.Println("🔐 正在进行用户认证...")

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
	fmt.Printf("✅ 登录成功，用户: %s\n", authResp.Record.Email)
	return nil
}

// CheckAuth 检查认证状态
func (pb *PocketBaseClient) CheckAuth() error {
	if pb.AuthToken == "" {
		return fmt.Errorf("未认证: 请先登录")
	}

	// 测试认证token是否有效
	resp, err := pb.makeRequest("POST", "/api/collections/users/auth-refresh", nil)
	if err != nil {
		return fmt.Errorf("认证检查失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ 认证状态正常")
		return nil
	} else {
		pb.AuthToken = "" // 清除无效token
		return fmt.Errorf("认证已过期，请重新登录")
	}
}

// ListSystems 获取所有系统/服务器
func (pb *PocketBaseClient) ListSystems() (*ListResponse[System], error) {
	fmt.Println("🔍 正在获取所有服务器...")

	// 构建查询参数
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
func (pb *PocketBaseClient) GetSystemLoadAverage(systemID string, count int) (*ListResponse[SystemStats], error) {
	fmt.Printf("🔍 正在获取系统 %s 最近 %d 条负载数据...\n", systemID, count)

	// 构建查询参数 - 获取最近的N条记录
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

// PrintSystems 显示系统信息
func PrintSystems(systems *ListResponse[System]) {
	fmt.Printf("✅ 找到 %d 台服务器:\n", systems.TotalItems)
	fmt.Println(strings.Repeat("-", 60))

	for i, system := range systems.Items {
		fmt.Printf("%d. 服务器: %s\n", i+1, system.Name)
		fmt.Printf("   主机: %s:%s\n", system.Host, system.Port)
		fmt.Printf("   状态: %s\n", system.Status)
		fmt.Printf("   ID: %s\n", system.ID)
		fmt.Printf("   创建时间: %s\n", formatTime(system.Created))
		fmt.Printf("   更新时间: %s\n", formatTime(system.Updated))
		fmt.Println("   " + strings.Repeat("-", 40))
	}
}

// formatTime 格式化时间字符串
func formatTime(timeStr string) string {
	// 尝试解析PocketBase的时间格式
	layouts := []string{
		"2006-01-02 15:04:05.999Z",
		"2006-01-02T15:04:05.999Z",
		"2006-01-02 15:04:05Z",
		"2006-01-02T15:04:05Z",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}

	// 如果都解析失败，返回原始字符串
	return timeStr
}

// PrintSystemStats 显示系统统计数据
func PrintSystemStats(stats *ListResponse[SystemStats]) {
	fmt.Printf("✅ 找到 %d 条系统统计记录:\n", stats.TotalItems)
	fmt.Println(strings.Repeat("-", 60))

	for i, stat := range stats.Items {
		fmt.Printf("%d. 时间: %s\n", i+1, formatTime(stat.Created))

		if stat.Stats.CPU > 0 {
			fmt.Printf("   CPU使用率: %.2f%%\n", stat.Stats.CPU)
		}

		// 内存使用率
		memPct := stat.Stats.MemPct
		if memPct == 0 && stat.Stats.Mem > 0 && stat.Stats.MemUsed > 0 {
			memPct = (stat.Stats.MemUsed / stat.Stats.Mem) * 100
		}
		fmt.Printf("   内存: %.2f%% (%.2fGB/%.2fGB)\n", memPct, stat.Stats.MemUsed, stat.Stats.Mem)

		// 网络带宽（数据已经是 MB/s，转换为 Mbps）
		netSentMbps := stat.Stats.NetworkSent * 8
		netRecvMbps := stat.Stats.NetworkRecv * 8
		fmt.Printf("   网络: ↑ %.2f Mbps / ↓ %.2f Mbps\n", netSentMbps, netRecvMbps)

		fmt.Println("   " + strings.Repeat("-", 40))
	}
}

// GetServersSummary 返回服务器状态摘要
func GetServersSummary(systems *ListResponse[System]) map[string]int {
	summary := map[string]int{
		"total":   systems.TotalItems,
		"online":  0,
		"offline": 0,
		"unknown": 0,
	}

	for _, system := range systems.Items {
		switch system.Status {
		case "up":
			summary["online"]++
		case "down":
			summary["offline"]++
		default:
			summary["unknown"]++
		}
	}

	return summary
}

// PrintSummary 显示服务器摘要
func PrintSummary(summary map[string]int) {
	fmt.Println("✅ 服务器摘要:")
	fmt.Printf("   总计: %d\n", summary["total"])
	fmt.Printf("   在线: %d\n", summary["online"])
	fmt.Printf("   离线: %d\n", summary["offline"])
	fmt.Printf("   未知: %d\n", summary["unknown"])
}

func main() {
	// 创建PocketBase客户端
	client := NewPocketBaseClient("https://bz.baidua.top")

	fmt.Println("🚀 开始执行Beszel API测试: https://bz.baidua.top")
	fmt.Println(strings.Repeat("=", 60))

	// 步骤1: 用户认证 (请替换为实际的用户名和密码)
	email := "Spike.wook@gmail.com" // 请替换为实际邮箱
	password := "adadmin/1213"      // 请替换为实际密码

	if err := client.Login(email, password); err != nil {
		log.Fatalf("❌ 认证失败: %v\n请检查邮箱和密码是否正确", err)
	}
	fmt.Println()

	// 步骤2: 检查认证状态
	if err := client.CheckAuth(); err != nil {
		log.Fatalf("❌ 认证检查失败: %v", err)
	}
	fmt.Println()

	// 步骤3: 列出所有服务器
	systems, err := client.ListSystems()
	if err != nil {
		log.Fatalf("❌ 获取系统列表失败: %v", err)
	}

	PrintSystems(systems)
	fmt.Println()

	// 步骤4: 获取服务器摘要
	summary := GetServersSummary(systems)
	PrintSummary(summary)
	fmt.Println()

	// 步骤5: 遍历所有服务器，测试负载平均值
	if len(systems.Items) > 0 {
		fmt.Println("\n📊 开始测试每台服务器的负载数据...")
		fmt.Println(strings.Repeat("=", 60))

		for i, system := range systems.Items {
			fmt.Printf("\n🔸 服务器 %d/%d: %s (ID: %s)\n", i+1, len(systems.Items), system.Name, system.ID)

			// 先尝试获取所有数据（不限时间）
			fmt.Println("   1️⃣ 尝试获取所有统计数据（不限时间）...")
			testParams := url.Values{}
			testParams.Set("page", "1")
			testParams.Set("perPage", "3")
			testParams.Set("sort", "-created")
			testParams.Set("filter", fmt.Sprintf(`system = "%s"`, system.ID))

			testEndpoint := "/api/collections/system_stats/records?" + testParams.Encode()

			resp, err := client.makeRequest("GET", testEndpoint, nil)
			if err == nil {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)

				// 解析响应查看是否有数据
				var testResult map[string]interface{}
				if err := json.Unmarshal(body, &testResult); err == nil {
					if items, ok := testResult["items"].([]interface{}); ok {
						fmt.Printf("   ✅ 找到 %d 条记录\n", len(items))
						if len(items) > 0 {
							// 显示第一条记录的信息
							if firstItem, ok := items[0].(map[string]interface{}); ok {
								fmt.Printf("   最新记录: 类型=%v, 时间=%v\n", firstItem["type"], firstItem["created"])
							}
						}
					}
				}
			} else {
				fmt.Printf("   ❌ 查询失败: %v\n", err)
			}

			// 获取最近5条1m类型数据并计算平均值
			fmt.Println("   2️⃣ 获取最近5条系统数据并计算平均值...")
			loadAvgData, err := client.GetSystemLoadAverage(system.ID, 5)
			if err != nil {
				fmt.Printf("   ❌ 获取数据失败: %v\n", err)
			} else {
				if loadAvgData.TotalItems > 0 {
					fmt.Printf("   ✅ 找到 %d 条记录\n", loadAvgData.TotalItems)

					// 计算平均值
					var avgCPU, avgMemPct float64
					var avgNetSent, avgNetRecv float64
					var count float64

					// 显示每条记录并累加
					for j, stat := range loadAvgData.Items {
						fmt.Printf("      📍 记录 %d - 时间: %s\n", j+1, formatTime(stat.Created))
						fmt.Printf("         CPU: %.2f%%\n", stat.Stats.CPU)

						// 内存使用率
						memPct := stat.Stats.MemPct
						if memPct == 0 && stat.Stats.Mem > 0 && stat.Stats.MemUsed > 0 {
							memPct = (stat.Stats.MemUsed / stat.Stats.Mem) * 100
						}
						fmt.Printf("         内存: %.2f%% (%.2fGB/%.2fGB)\n", memPct, stat.Stats.MemUsed, stat.Stats.Mem)

						// 网络带宽（数据已经是 MB/s，转换为 Mbps）
						netSentMbps := stat.Stats.NetworkSent * 8 // MB/s 转 Mbps
						netRecvMbps := stat.Stats.NetworkRecv * 8
						fmt.Printf("         网络: ↑ %.2f Mbps / ↓ %.2f Mbps\n", netSentMbps, netRecvMbps)

						avgCPU += stat.Stats.CPU
						avgMemPct += memPct
						avgNetSent += netSentMbps
						avgNetRecv += netRecvMbps
						count++

						if j >= 4 { // 最多5条
							break
						}
					}

					// 计算并显示平均值
					if count > 0 {
						avgCPU /= count
						avgMemPct /= count
						avgNetSent /= count
						avgNetRecv /= count

						fmt.Println("\n   📊 平均值汇总:")
						fmt.Printf("      - CPU平均: %.2f%%\n", avgCPU)
						fmt.Printf("      - 内存平均: %.2f%%\n", avgMemPct)
						fmt.Printf("      - 网络平均: ↑ %.2f Mbps / ↓ %.2f Mbps\n", avgNetSent, avgNetRecv)
					}
				} else {
					fmt.Println("   ⚠️  没有找到数据")
				}
			}

			fmt.Println("   " + strings.Repeat("-", 40))
		}
	} else {
		fmt.Println("ℹ️  未找到服务器，无法测试负载平均值")
	}

	fmt.Println("🎉 所有测试完成!")
}
