package service

import (
	"backend/internal/database"
	"backend/pkg/models"
	"fmt"
)

// NodeTagService 节点标签服务
type NodeTagService struct{}

// NewNodeTagService 创建节点标签服务
func NewNodeTagService() *NodeTagService {
	return &NodeTagService{}
}


// AddTag 添加服务器标签
func (s *NodeTagService) AddTag(systemID string, request *models.NodeTagRequest) error {
	db := database.GetDB()
	
	// 检查是否已存在相同的标签
	var existingTag models.NodeTag
	result := db.Where("system_id = ? AND tag_type = ? AND tag_id = ?", 
		systemID, request.Type, request.ID).First(&existingTag)
	
	if result.Error == nil {
		return fmt.Errorf("标签已存在")
	}
	
	// 创建新标签
	newTag := models.NodeTag{
		SystemID: systemID,
		TagType:  request.Type,
		TagID:    request.ID,
	}
	
	if err := db.Create(&newTag).Error; err != nil {
		return fmt.Errorf("创建标签失败: %w", err)
	}
	
	return nil
}

// RemoveTag 删除服务器标签
func (s *NodeTagService) RemoveTag(systemID string, request *models.NodeTagRequest) error {
	db := database.GetDB()
	
	// 删除特定标签
	result := db.Where("system_id = ? AND tag_type = ? AND tag_id = ?", 
		systemID, request.Type, request.ID).Delete(&models.NodeTag{})
	
	if result.Error != nil {
		return fmt.Errorf("删除标签失败: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("标签不存在")
	}
	
	return nil
}

// GetSystemTags 获取服务器的所有标签
func (s *NodeTagService) GetSystemTags(systemID string) ([]models.NodeTag, error) {
	db := database.GetDB()
	var tags []models.NodeTag
	
	err := db.Where("system_id = ?", systemID).Find(&tags).Error
	
	if err != nil {
		return nil, fmt.Errorf("获取服务器标签失败: %w", err)
	}
	
	return tags, nil
}

// GetNodeLoadStatus 获取节点负载状态
func (s *NodeTagService) GetNodeLoadStatus(requests []models.NodeLoadRequest, systemService *SystemService, thresholdService *ThresholdService) ([]models.NodeLoadResponse, error) {
	db := database.GetDB()
	var responses []models.NodeLoadResponse
	
	// 获取所有系统的统计数据
	systemsWithStats, err := systemService.GetSystemsWithAvgStats()
	if err != nil {
		return nil, fmt.Errorf("获取系统统计数据失败: %w", err)
	}
	
	// 创建系统ID到统计数据的映射
	systemStatsMap := make(map[string]*models.SystemWithAvgStats)
	for _, system := range systemsWithStats {
		systemStatsMap[system.ID] = system
	}
	
	for _, req := range requests {
		response := models.NodeLoadResponse{
			Type:       req.Type,
			ID:         req.ID,
			LoadStatus: "not_found",
		}
		
		// 查找对应的系统ID
		var tag models.NodeTag
		err := db.Where("tag_type = ? AND tag_id = ?", req.Type, req.ID).First(&tag).Error
		if err == nil {
			// 找到对应的系统，检查是否有统计数据
			if systemWithStats, exists := systemStatsMap[tag.SystemID]; exists {
				// 获取阈值配置
				threshold, err := thresholdService.GetThreshold(tag.SystemID)
				if err == nil {
					// 使用现有的负载计算逻辑
					response.LoadStatus = systemService.CalculateLoadStatus(systemWithStats, threshold)
				} else {
					// 没有阈值配置，默认为正常
					response.LoadStatus = "normal"
				}
			}
		}
		
		responses = append(responses, response)
	}
	
	return responses, nil
}

