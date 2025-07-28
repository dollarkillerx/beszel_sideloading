package service

import (
	"backend/internal/database"
	"backend/pkg/models"
	"fmt"
)

// ThresholdService 阈值配置服务
type ThresholdService struct{}

// NewThresholdService 创建阈值配置服务
func NewThresholdService() *ThresholdService {
	return &ThresholdService{}
}

// GetThreshold 获取系统阈值配置
func (s *ThresholdService) GetThreshold(systemID string) (*models.SystemThreshold, error) {
	storage := database.GetStorage()
	
	threshold, err := storage.GetThreshold(systemID)
	if err != nil {
		return nil, fmt.Errorf("获取阈值配置失败: %w", err)
	}
	
	// 如果没有找到，创建默认配置
	if threshold == nil {
		threshold = &models.SystemThreshold{
			SystemID:        systemID,
			CPUAlertLimit:   90.0,
			MemAlertLimit:   90.0,
			NetUpMax:        0,
			NetDownMax:      0,
			NetUpAlert:      80.0,
			NetDownAlert:    80.0,
		}
		
		// 保存默认配置
		if createErr := storage.CreateOrUpdateThreshold(threshold); createErr != nil {
			return nil, fmt.Errorf("创建默认阈值配置失败: %w", createErr)
		}
	}
	
	return threshold, nil
}

// UpdateThreshold 更新系统阈值配置
func (s *ThresholdService) UpdateThreshold(systemID string, threshold *models.SystemThreshold) error {
	storage := database.GetStorage()
	
	// 设置SystemID
	threshold.SystemID = systemID
	
	// 创建或更新阈值
	if err := storage.CreateOrUpdateThreshold(threshold); err != nil {
		return fmt.Errorf("更新阈值配置失败: %w", err)
	}
	
	return nil
}

// UpdateNetworkMax 更新网络最大值（用于动态更新历史极限值）
func (s *ThresholdService) UpdateNetworkMax(systemID string, netUpMbps, netDownMbps float64) error {
	storage := database.GetStorage()
	
	// 获取现有阈值配置
	threshold, err := s.GetThreshold(systemID)
	if err != nil {
		return err
	}
	
	// 更新最大值（只有当新值更大时才更新）
	updated := false
	if netUpMbps > threshold.NetUpMax {
		threshold.NetUpMax = netUpMbps
		updated = true
	}
	if netDownMbps > threshold.NetDownMax {
		threshold.NetDownMax = netDownMbps
		updated = true
	}
	
	if updated {
		return storage.CreateOrUpdateThreshold(threshold)
	}
	
	return nil
}

// GetAllThresholds 获取所有系统的阈值配置
func (s *ThresholdService) GetAllThresholds() ([]*models.SystemThreshold, error) {
	storage := database.GetStorage()
	return storage.ListThresholds()
}

// DeleteThreshold 删除系统阈值配置
func (s *ThresholdService) DeleteThreshold(systemID string) error {
	storage := database.GetStorage()
	
	if err := storage.DeleteThreshold(systemID); err != nil {
		return fmt.Errorf("删除阈值配置失败: %w", err)
	}
	
	return nil
}