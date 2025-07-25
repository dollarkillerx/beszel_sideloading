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

// SystemThreshold 系统阈值配置（本地SQLite存储）
type SystemThreshold struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	SystemID        string  `gorm:"uniqueIndex;not null" json:"system_id"`
	CPUAlertLimit   float64 `gorm:"default:90.0" json:"cpu_alert_limit"`   // CPU告警阈值（%）
	MemAlertLimit   float64 `gorm:"default:90.0" json:"mem_alert_limit"`   // 内存告警阈值（%）
	NetUpMax        float64 `gorm:"default:0" json:"net_up_max"`           // 上行最大Mbps（历史极限值）
	NetDownMax      float64 `gorm:"default:0" json:"net_down_max"`         // 下行最大Mbps（历史极限值）
	NetUpAlert      float64 `gorm:"default:80.0" json:"net_up_alert"`      // 上行告警阈值（百分比）
	NetDownAlert    float64 `gorm:"default:80.0" json:"net_down_alert"`    // 下行告警阈值（百分比）
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SystemWithLoadStatus 带负载状态的系统统计
type SystemWithLoadStatus struct {
	SystemWithAvgStats
	LoadStatus string `json:"load_status"` // normal, high
}

// NodeTag 服务器标签（本地SQLite存储）
type NodeTag struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	SystemID string `gorm:"not null" json:"system_id"`        // 服务器ID
	TagType  string `gorm:"not null" json:"tag_type"`         // 标签类型
	TagID    int    `gorm:"not null" json:"tag_id"`           // 标签ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NodeTagRequest 创建/删除标签的请求结构
type NodeTagRequest struct {
	Type string `json:"type" binding:"required"` // 标签类型
	ID   int    `json:"id" binding:"required"`   // 标签ID
}

// NodeTagsResponse 节点标签响应
type NodeTagsResponse struct {
	Success string    `json:"success"`
	Tags    []NodeTag `json:"tags,omitempty"`
}

// NodeLoadRequest 节点负载查询请求
type NodeLoadRequest struct {
	Type string `json:"type" binding:"required"` // 标签类型
	ID   int    `json:"id" binding:"required"`   // 标签ID
}

// NodeLoadResponse 节点负载查询响应
type NodeLoadResponse struct {
	Type       string `json:"type"`        // 标签类型
	ID         int    `json:"id"`          // 标签ID
	LoadStatus string `json:"load_status"` // 负载状态: normal, high, not_found
}