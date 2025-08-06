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
	pbClient         *pocketbase.Client
	config           *config.Config
	thresholdService *ThresholdService
	nodeService      *NodeService
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
	
	service := &SystemService{
		pbClient:         client,
		config:           cfg,
		thresholdService: NewThresholdService(),
	}
	
	// 启动token刷新定时器（每12天刷新一次）
	go service.startTokenRefreshTimer()
	
	return service
}

// SetNodeService 设置节点服务（避免循环依赖）
func (s *SystemService) SetNodeService(nodeService *NodeService) {
	s.nodeService = nodeService
}

// startTokenRefreshTimer 启动token刷新定时器
func (s *SystemService) startTokenRefreshTimer() {
	ticker := time.NewTicker(12 * 24 * time.Hour) // 每12天刷新一次
	defer ticker.Stop()
	
	for range ticker.C {
		if err := s.pbClient.RefreshAuth(); err != nil {
			log.Printf("刷新PocketBase认证失败: %v", err)
			// 如果刷新失败，尝试重新登录
			if err := s.pbClient.Login(s.config.PocketBase.Email, s.config.PocketBase.Password); err != nil {
				log.Printf("重新登录PocketBase失败: %v", err)
			} else {
				log.Printf("成功重新登录PocketBase")
			}
		} else {
			log.Printf("成功刷新PocketBase认证token")
		}
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
			onlineUsers := 0
			if s.nodeService != nil {
				if nodeInfo, err := s.nodeService.GetSystemNodeInfo(system.ID, system.Name); err == nil {
					onlineUsers = nodeInfo.TotalOnline
				}
			}
			
			systemWithStats := &models.SystemWithAvgStats{
				System:      *system,
				AvgCPU:      0,
				AvgMemPct:   0,
				AvgNetSent:  0,
				AvgNetRecv:  0,
				OnlineUsers: onlineUsers,
				LastUpdate:  time.Now(),
			}
			result = append(result, systemWithStats)
			continue
		}
		
		// 计算平均值
		avgStats := calculateAverageStats(pbStats.Items)
		
		// 获取在线人数
		onlineUsers := 0
		if s.nodeService != nil {
			if nodeInfo, err := s.nodeService.GetSystemNodeInfo(system.ID, system.Name); err == nil {
				onlineUsers = nodeInfo.TotalOnline
			}
		}
		
		systemWithStats := &models.SystemWithAvgStats{
			System:      *system,
			AvgCPU:      avgStats.AvgCPU,
			AvgMemPct:   avgStats.AvgMemPct,
			AvgNetSent:  avgStats.AvgNetSent,
			AvgNetRecv:  avgStats.AvgNetRecv,
			OnlineUsers: onlineUsers,
			LastUpdate:  avgStats.LastUpdate,
		}
		
		result = append(result, systemWithStats)
	}
	
	return result, nil
}

// GetSystemsWithLoadStatus 获取带负载状态的系统列表
func (s *SystemService) GetSystemsWithLoadStatus() ([]*models.SystemWithLoadStatus, error) {
	systems, err := s.GetSystemsWithAvgStats()
	if err != nil {
		return nil, err
	}
	
	var result []*models.SystemWithLoadStatus
	
	for _, system := range systems {
		// 获取阈值配置
		threshold, err := s.thresholdService.GetThreshold(system.ID)
		if err != nil {
			log.Printf("获取系统 %s 阈值配置失败: %v", system.Name, err)
			// 使用默认配置继续处理
			threshold = &models.SystemThreshold{
				SystemID:        system.ID,
				CPUAlertLimit:   90.0,
				MemAlertLimit:   90.0,
				NetUpMax:        0,
				NetDownMax:      0,
				NetUpAlert:      80.0,
				NetDownAlert:    80.0,
			}
		}
		
		// 计算负载状态
		loadStatus := s.CalculateLoadStatus(system, threshold)
		
		// 更新网络最大值（动态更新历史极限值）
		netUpMbps := system.AvgNetSent * 8  // 转换为 Mbps
		netDownMbps := system.AvgNetRecv * 8 // 转换为 Mbps
		if err := s.thresholdService.UpdateNetworkMax(system.ID, netUpMbps, netDownMbps); err != nil {
			log.Printf("更新系统 %s 网络最大值失败: %v", system.Name, err)
		}
		
		systemWithLoadStatus := &models.SystemWithLoadStatus{
			SystemWithAvgStats: *system,
			LoadStatus:         loadStatus,
		}
		
		result = append(result, systemWithLoadStatus)
	}
	
	return result, nil
}

// CalculateLoadStatus 计算负载状态
func (s *SystemService) CalculateLoadStatus(system *models.SystemWithAvgStats, threshold *models.SystemThreshold) string {
	// 检查CPU使用率
	if system.AvgCPU >= threshold.CPUAlertLimit {
		log.Printf("系统 %s CPU负载过高: %.2f%% >= %.2f%%", system.Name, system.AvgCPU, threshold.CPUAlertLimit)
		return "high"
	}
	
	// 检查内存使用率
	if system.AvgMemPct >= threshold.MemAlertLimit {
		log.Printf("系统 %s 内存负载过高: %.2f%% >= %.2f%%", system.Name, system.AvgMemPct, threshold.MemAlertLimit)
		return "high"
	}
	
	// 检查网络上行 - 只有当设置了最大值且大于0时才检查
	if threshold.NetUpMax > 0 {
		netUpMbps := system.AvgNetSent * 8 // 转换为 Mbps
		upThreshold := threshold.NetUpMax * (threshold.NetUpAlert / 100)
		if netUpMbps >= upThreshold {
			log.Printf("系统 %s 上行网络负载过高: %.2f Mbps >= %.2f Mbps (%.1f%% of %.2f Mbps)", 
				system.Name, netUpMbps, upThreshold, threshold.NetUpAlert, threshold.NetUpMax)
			return "high"
		}
	}
	
	// 检查网络下行 - 只有当设置了最大值且大于0时才检查
	if threshold.NetDownMax > 0 {
		netDownMbps := system.AvgNetRecv * 8 // 转换为 Mbps
		downThreshold := threshold.NetDownMax * (threshold.NetDownAlert / 100)
		if netDownMbps >= downThreshold {
			log.Printf("系统 %s 下行网络负载过高: %.2f Mbps >= %.2f Mbps (%.1f%% of %.2f Mbps)", 
				system.Name, netDownMbps, downThreshold, threshold.NetDownAlert, threshold.NetDownMax)
			return "high"
		}
	}
	
	// 检查在线人数 - 如果设置了阈值且大于0
	if threshold.OnlineUsersLimit > 0 && system.OnlineUsers >= threshold.OnlineUsersLimit {
		log.Printf("系统 %s 在线人数过多: %d >= %d", system.Name, system.OnlineUsers, threshold.OnlineUsersLimit)
		return "high"
	}
	
	return "normal"
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