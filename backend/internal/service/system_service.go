package service

import (
	"backend/internal/config"
	"backend/internal/pocketbase"
	"backend/pkg/models"
	"fmt"
	"log"
	"time"
)

// SystemService 系统服务
type SystemService struct {
	pbClient *pocketbase.Client
	config   *config.Config
}

// NewSystemService 创建系统服务
func NewSystemService(cfg *config.Config) *SystemService {
	client := pocketbase.NewClient(cfg.PocketBase.BaseURL)
	
	// 登录认证
	if err := client.Login(cfg.PocketBase.Email, cfg.PocketBase.Password); err != nil {
		log.Printf("PocketBase 登录失败: %v", err)
	} else {
		log.Printf("PocketBase 登录成功，连接到: %s", cfg.PocketBase.BaseURL)
	}
	
	return &SystemService{
		pbClient: client,
		config:   cfg,
	}
}

// GetSystems 获取所有系统
func (s *SystemService) GetSystems() ([]*models.System, error) {
	pbSystems, err := s.pbClient.ListSystems()
	if err != nil {
		return nil, fmt.Errorf("获取系统列表失败: %w", err)
	}
	
	var systems []*models.System
	for _, pbSystem := range pbSystems.Items {
		system := &models.System{
			ID:     pbSystem.ID,
			Name:   pbSystem.Name,
			Host:   pbSystem.Host,
			Port:   pbSystem.Port,
			Status: pbSystem.Status,
			CreatedAt: parseTime(pbSystem.Created),
			UpdatedAt: parseTime(pbSystem.Updated),
		}
		systems = append(systems, system)
	}
	
	return systems, nil
}

// GetSystemSummary 获取系统摘要
func (s *SystemService) GetSystemSummary() (*models.SystemSummary, error) {
	systems, err := s.GetSystems()
	if err != nil {
		return nil, err
	}
	
	summary := &models.SystemSummary{
		Total: int64(len(systems)),
	}
	
	for _, system := range systems {
		switch system.Status {
		case "up":
			summary.Online++
		case "down":
			summary.Offline++
		default:
			summary.Unknown++
		}
	}
	
	return summary, nil
}

// GetSystemsWithAvgStats 获取所有系统及其平均统计数据
func (s *SystemService) GetSystemsWithAvgStats() ([]*models.SystemWithAvgStats, error) {
	systems, err := s.GetSystems()
	if err != nil {
		return nil, err
	}
	
	var result []*models.SystemWithAvgStats
	
	for _, system := range systems {
		// 获取最近5条1分钟数据
		pbStats, err := s.pbClient.GetSystemLoadAverage(system.ID, 5)
		if err != nil {
			log.Printf("获取系统 %s 统计数据失败: %v", system.Name, err)
			// 如果获取失败，仍然添加系统信息，但统计数据为0
			systemWithStats := &models.SystemWithAvgStats{
				System:     *system,
				AvgCPU:     0,
				AvgMemPct:  0,
				AvgNetSent: 0,
				AvgNetRecv: 0,
				LastUpdate: time.Now(),
			}
			result = append(result, systemWithStats)
			continue
		}
		
		// 计算平均值
		avgStats := calculateAverageStats(pbStats.Items)
		
		systemWithStats := &models.SystemWithAvgStats{
			System:     *system,
			AvgCPU:     avgStats.AvgCPU,
			AvgMemPct:  avgStats.AvgMemPct,
			AvgNetSent: avgStats.AvgNetSent,
			AvgNetRecv: avgStats.AvgNetRecv,
			LastUpdate: avgStats.LastUpdate,
		}
		
		result = append(result, systemWithStats)
	}
	
	return result, nil
}

// GetSystemStats 获取指定系统的统计数据
func (s *SystemService) GetSystemStats(systemID string, limit int) ([]*models.SystemStat, error) {
	pbStats, err := s.pbClient.GetSystemLoadAverage(systemID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取系统统计数据失败: %w", err)
	}
	
	var stats []*models.SystemStat
	for _, pbStat := range pbStats.Items {
		// 计算内存使用百分比
		memPct := pbStat.Stats.MemPct
		if memPct == 0 && pbStat.Stats.Mem > 0 && pbStat.Stats.MemUsed > 0 {
			memPct = (pbStat.Stats.MemUsed / pbStat.Stats.Mem) * 100
		}
		
		stat := &models.SystemStat{
			ID:       pbStat.ID,
			SystemID: pbStat.System,
			Type:     pbStat.Type,
			CPU:      pbStat.Stats.CPU,
			Mem:      pbStat.Stats.Mem,
			MemUsed:  pbStat.Stats.MemUsed,
			MemPct:   memPct,
			NetSent:  pbStat.Stats.NetworkSent,
			NetRecv:  pbStat.Stats.NetworkRecv,
			CreatedAt: parseTime(pbStat.Created),
		}
		stats = append(stats, stat)
	}
	
	return stats, nil
}

// calculateAverageStats 计算平均统计数据
func calculateAverageStats(pbStats []pocketbase.SystemStats) *models.AverageStats {
	if len(pbStats) == 0 {
		return &models.AverageStats{
			LastUpdate: time.Now(),
		}
	}
	
	var avgCPU, avgMemPct, avgNetSent, avgNetRecv float64
	var lastUpdate time.Time
	
	for _, stat := range pbStats {
		avgCPU += stat.Stats.CPU
		
		// 计算内存使用百分比
		memPct := stat.Stats.MemPct
		if memPct == 0 && stat.Stats.Mem > 0 && stat.Stats.MemUsed > 0 {
			memPct = (stat.Stats.MemUsed / stat.Stats.Mem) * 100
		}
		avgMemPct += memPct
		
		// 网络数据转换为 MB/s
		avgNetSent += stat.Stats.NetworkSent
		avgNetRecv += stat.Stats.NetworkRecv
		
		// 记录最新时间
		statTime := parseTime(stat.Created)
		if statTime.After(lastUpdate) {
			lastUpdate = statTime
		}
	}
	
	count := float64(len(pbStats))
	return &models.AverageStats{
		AvgCPU:     avgCPU / count,
		AvgMemPct:  avgMemPct / count,
		AvgNetSent: avgNetSent / count,
		AvgNetRecv: avgNetRecv / count,
		LastUpdate: lastUpdate,
	}
}

// parseTime 解析时间字符串
func parseTime(timeStr string) time.Time {
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
			return t
		}
	}
	
	return time.Now()
}