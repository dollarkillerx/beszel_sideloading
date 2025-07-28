package database

import (
	"backend/pkg/models"
)

// Storage 定义存储接口
type Storage interface {
	// 系统阈值相关
	CreateOrUpdateThreshold(threshold *models.SystemThreshold) error
	GetThreshold(systemID string) (*models.SystemThreshold, error)
	ListThresholds() ([]*models.SystemThreshold, error)
	DeleteThreshold(systemID string) error

	// 节点标签相关
	CreateNodeTag(tag *models.NodeTag) error
	GetNodeTags(systemID string) ([]*models.NodeTag, error)
	GetNodeTagsByTypeAndID(tagType string, tagID int) ([]*models.NodeTag, error)
	DeleteNodeTag(systemID string, tagType string, tagID int) error
	DeleteAllNodeTags(systemID string) error

	// 关闭存储
	Close() error
}