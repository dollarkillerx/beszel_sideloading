package models

import (
	"time"
)

// System 表示服务器/系统记录
type System struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      string    `json:"port"`
	Status    string    `json:"status"` // up, down, unknown
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SystemStat 表示系统统计记录
type SystemStat struct {
	ID        string    `json:"id"`
	SystemID  string    `json:"system_id"`
	Type      string    `json:"type"` // 1m, 5m, 15m
	CPU       float64   `json:"cpu"`
	Mem       float64   `json:"mem"`       // 总内存 GB
	MemUsed   float64   `json:"mem_used"`  // 使用内存 GB
	MemPct    float64   `json:"mem_pct"`   // 内存使用百分比
	NetSent   float64   `json:"net_sent"`  // 网络发送 MB/s
	NetRecv   float64   `json:"net_recv"`  // 网络接收 MB/s
	CreatedAt time.Time `json:"created_at"`
}

// SystemSummary 服务器摘要
type SystemSummary struct {
	Total   int64 `json:"total"`
	Online  int64 `json:"online"`
	Offline int64 `json:"offline"`
	Unknown int64 `json:"unknown"`
}

// SystemWithAvgStats 带平均统计的系统
type SystemWithAvgStats struct {
	System
	AvgCPU     float64   `json:"avg_cpu"`
	AvgMemPct  float64   `json:"avg_mem_pct"`
	AvgNetSent float64   `json:"avg_net_sent"`
	AvgNetRecv float64   `json:"avg_net_recv"`
	LastUpdate time.Time `json:"last_update"`
}

// AverageStats 平均统计数据
type AverageStats struct {
	AvgCPU     float64   `json:"avg_cpu"`
	AvgMemPct  float64   `json:"avg_mem_pct"`
	AvgNetSent float64   `json:"avg_net_sent"`
	AvgNetRecv float64   `json:"avg_net_recv"`
	LastUpdate time.Time `json:"last_update"`
}