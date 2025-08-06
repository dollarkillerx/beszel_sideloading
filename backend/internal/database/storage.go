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

	// 系统别名相关
	SetSystemAlias(alias *models.SystemAlias) error
	GetSystemAlias(systemID string) (*models.SystemAlias, error)
	GetAllSystemAliases() ([]*models.SystemAlias, error)
	DeleteSystemAlias(systemID string) error

	// 关闭存储
	Close() error
}