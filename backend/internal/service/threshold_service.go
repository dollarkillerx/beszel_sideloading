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
	db := database.GetDB()
	var threshold models.SystemThreshold
	
	err := db.Where("system_id = ?", systemID).First(&threshold).Error
	if err != nil {
		// 如果没有找到，创建默认配置
		threshold = models.SystemThreshold{
			SystemID:        systemID,
			CPUAlertLimit:   90.0,
			MemAlertLimit:   90.0,
			NetUpMax:        0,
			NetDownMax:      0,
			NetUpAlert:      80.0,
			NetDownAlert:    80.0,
		}
		
		// 保存默认配置
		if createErr := db.Create(&threshold).Error; createErr != nil {
			return nil, fmt.Errorf("创建默认阈值配置失败: %w", createErr)
		}
	}
	
	return &threshold, nil
}

// UpdateThreshold 更新系统阈值配置
func (s *ThresholdService) UpdateThreshold(systemID string, threshold *models.SystemThreshold) error {
	db := database.GetDB()
	
	// 方案：使用 map 明确指定要更新的字段，包括零值
	updates := map[string]interface{}{
		"cpu_alert_limit": threshold.CPUAlertLimit,
		"mem_alert_limit": threshold.MemAlertLimit,
		"net_up_max":      threshold.NetUpMax,
		"net_down_max":    threshold.NetDownMax,
		"net_up_alert":    threshold.NetUpAlert,
		"net_down_alert":  threshold.NetDownAlert,
	}
	
	result := db.Model(&models.SystemThreshold{}).Where("system_id = ?", systemID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新阈值配置失败: %w", result.Error)
	}
	
	// 如果没有找到记录，创建新记录
	if result.RowsAffected == 0 {
		threshold.SystemID = systemID
		if err := db.Create(threshold).Error; err != nil {
			return fmt.Errorf("创建阈值配置失败: %w", err)
		}
	}
	
	return nil
}

// UpdateNetworkMax 更新网络最大值（用于动态更新历史极限值）
func (s *ThresholdService) UpdateNetworkMax(systemID string, netUpMbps, netDownMbps float64) error {
	db := database.GetDB()
	
	var threshold models.SystemThreshold
	err := db.Where("system_id = ?", systemID).First(&threshold).Error
	if err != nil {
		// 如果没有配置，先创建默认配置
		if _, err := s.GetThreshold(systemID); err != nil {
			return err
		}
		// 重新获取
		if err := db.Where("system_id = ?", systemID).First(&threshold).Error; err != nil {
			return err
		}
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
		return db.Save(&threshold).Error
	}
	
	return nil
}

// GetAllThresholds 获取所有系统的阈值配置
func (s *ThresholdService) GetAllThresholds() ([]models.SystemThreshold, error) {
	db := database.GetDB()
	var thresholds []models.SystemThreshold
	
	err := db.Find(&thresholds).Error
	return thresholds, err
}

// DeleteThreshold 删除系统阈值配置
func (s *ThresholdService) DeleteThreshold(systemID string) error {
	db := database.GetDB()
	
	result := db.Where("system_id = ?", systemID).Delete(&models.SystemThreshold{})
	if result.Error != nil {
		return fmt.Errorf("删除阈值配置失败: %w", result.Error)
	}
	
	return nil
}