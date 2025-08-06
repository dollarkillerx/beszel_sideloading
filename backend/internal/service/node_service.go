package service

import (
	"backend/pkg/models"
	"fmt"
	"strings"
)

// NodeService 节点服务
type NodeService struct {
	redisService *RedisService
	aliasService *AliasService
}

// NewNodeService 创建节点服务
func NewNodeService(redisService *RedisService) *NodeService {
	return &NodeService{
		redisService: redisService,
		aliasService: NewAliasService(),
	}
}

// GetSystemNodeInfo 获取系统的节点信息
func (s *NodeService) GetSystemNodeInfo(systemID, systemName string) (*models.SystemNodeInfo, error) {
	// 获取系统别名
	alias, err := s.aliasService.GetAlias(systemID)
	if err != nil {
		return nil, fmt.Errorf("获取系统别名失败: %w", err)
	}

	result := &models.SystemNodeInfo{
		SystemID:   systemID,
		SystemName: systemName,
		Nodes:      []models.V2boardNode{},
	}

	// 如果没有别名，返回空的节点信息
	if alias == nil || alias.Alias == "" {
		return result, nil
	}

	result.Alias = alias.Alias

	// 根据别名模糊匹配节点
	nodes, err := s.redisService.GetNodesByAlias(alias.Alias)
	if err != nil {
		return nil, fmt.Errorf("查询节点信息失败: %w", err)
	}

	result.Nodes = nodes

	// 计算总在线人数
	totalOnline := 0
	for _, node := range nodes {
		totalOnline += node.Online
	}
	result.TotalOnline = totalOnline

	return result, nil
}

// GetAllSystemsNodeInfo 获取所有系统的节点信息
func (s *NodeService) GetAllSystemsNodeInfo(systems []*models.System) ([]*models.SystemNodeInfo, error) {
	var results []*models.SystemNodeInfo

	for _, system := range systems {
		nodeInfo, err := s.GetSystemNodeInfo(system.ID, system.Name)
		if err != nil {
			// 记录错误但继续处理其他系统
			fmt.Printf("获取系统 %s 节点信息失败: %v\n", system.ID, err)
			continue
		}
		results = append(results, nodeInfo)
	}

	return results, nil
}

// SearchNodesByKeyword 根据关键词搜索节点
func (s *NodeService) SearchNodesByKeyword(keyword string) ([]models.V2boardNode, error) {
	if strings.TrimSpace(keyword) == "" {
		return []models.V2boardNode{}, nil
	}

	nodes, err := s.redisService.GetNodesByAlias(keyword)
	if err != nil {
		return nil, fmt.Errorf("搜索节点失败: %w", err)
	}

	return nodes, nil
}