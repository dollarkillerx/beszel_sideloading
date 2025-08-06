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
	AvgCPU      float64   `json:"avg_cpu"`
	AvgMemPct   float64   `json:"avg_mem_pct"`
	AvgNetSent  float64   `json:"avg_net_sent"`
	AvgNetRecv  float64   `json:"avg_net_recv"`
	OnlineUsers int       `json:"online_users"`     // 在线人数
	LastUpdate  time.Time `json:"last_update"`
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
	ID                uint    `gorm:"primaryKey" json:"id"`
	SystemID          string  `gorm:"uniqueIndex;not null" json:"system_id"`
	CPUAlertLimit     float64 `gorm:"default:90.0" json:"cpu_alert_limit"`     // CPU告警阈值（%）
	MemAlertLimit     float64 `gorm:"default:90.0" json:"mem_alert_limit"`     // 内存告警阈值（%）
	NetUpMax          float64 `gorm:"default:0" json:"net_up_max"`             // 上行最大Mbps（历史极限值）
	NetDownMax        float64 `gorm:"default:0" json:"net_down_max"`           // 下行最大Mbps（历史极限值）
	NetUpAlert        float64 `gorm:"default:80.0" json:"net_up_alert"`        // 上行告警阈值（百分比）
	NetDownAlert      float64 `gorm:"default:80.0" json:"net_down_alert"`      // 下行告警阈值（百分比）
	OnlineUsersLimit  int     `gorm:"default:300" json:"online_users_limit"`   // 在线人数告警阈值（默认300人）
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// SystemWithLoadStatus 带负载状态的系统统计
type SystemWithLoadStatus struct {
	SystemWithAvgStats
	LoadStatus string `json:"load_status"` // normal, high
}

// SystemAlias 服务器别名（本地存储）
type SystemAlias struct {
	ID       uint   `json:"id"`
	SystemID string `json:"system_id"`  // 服务器ID，唯一索引
	Alias    string `json:"alias"`      // 别名
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SystemAliasRequest 创建/更新别名的请求结构
type SystemAliasRequest struct {
	Alias string `json:"alias" binding:"required"` // 别名
}

// SystemAliasResponse 别名响应
type SystemAliasResponse struct {
	Success string       `json:"success,omitempty"`
	Alias   *SystemAlias `json:"alias,omitempty"`
}

// V2boardNode V2board节点信息
type V2boardNode struct {
	Name       string `json:"name"`
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Online     int    `json:"online"`
	LastUpdate int64  `json:"last_update"`
}

// SystemNodeInfo 服务器节点信息
type SystemNodeInfo struct {
	SystemID    string        `json:"system_id"`
	SystemName  string        `json:"system_name"`
	Alias       string        `json:"alias,omitempty"`
	Nodes       []V2boardNode `json:"nodes"`
	TotalOnline int           `json:"total_online"`
}

