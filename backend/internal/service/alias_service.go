package service

import (
	"backend/internal/database"
	"backend/pkg/models"
	"fmt"
)

// AliasService 别名服务
type AliasService struct{}

// NewAliasService 创建别名服务
func NewAliasService() *AliasService {
	return &AliasService{}
}

// SetAlias 设置服务器别名（每个服务器只能有一个别名）
func (s *AliasService) SetAlias(systemID string, request *models.SystemAliasRequest) error {
	storage := database.GetStorage()
	
	// 创建或更新别名
	alias := &models.SystemAlias{
		SystemID: systemID,
		Alias:    request.Alias,
	}
	
	if err := storage.SetSystemAlias(alias); err != nil {
		return fmt.Errorf("设置别名失败: %w", err)
	}
	
	return nil
}

// GetAlias 获取服务器别名
func (s *AliasService) GetAlias(systemID string) (*models.SystemAlias, error) {
	storage := database.GetStorage()
	
	alias, err := storage.GetSystemAlias(systemID)
	if err != nil {
		return nil, fmt.Errorf("获取别名失败: %w", err)
	}
	
	return alias, nil
}

// DeleteAlias 删除服务器别名
func (s *AliasService) DeleteAlias(systemID string) error {
	storage := database.GetStorage()
	
	if err := storage.DeleteSystemAlias(systemID); err != nil {
		return fmt.Errorf("删除别名失败: %w", err)
	}
	
	return nil
}

// GetAllAliases 获取所有别名
func (s *AliasService) GetAllAliases() ([]*models.SystemAlias, error) {
	storage := database.GetStorage()
	
	aliases, err := storage.GetAllSystemAliases()
	if err != nil {
		return nil, fmt.Errorf("获取所有别名失败: %w", err)
	}
	
	return aliases, nil
}