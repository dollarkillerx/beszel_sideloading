package handlers

import (
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var nodeService *service.NodeService

// InitNodeHandler 初始化节点处理器
func InitNodeHandler(ns *service.NodeService) {
	nodeService = ns
}

// GetSystemNodes 获取系统的节点信息
func GetSystemNodes(c *gin.Context) {
	// 检查nodeService是否初始化
	if nodeService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "节点服务不可用，Redis连接失败"})
		return
	}

	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	// 获取系统基本信息
	systems, err := systemService.GetSystems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统信息失败"})
		return
	}

	var systemName string
	for _, system := range systems {
		if system.ID == systemID {
			systemName = system.Name
			break
		}
	}

	if systemName == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "系统不存在"})
		return
	}

	// 获取节点信息
	nodeInfo, err := nodeService.GetSystemNodeInfo(systemID, systemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nodeInfo)
}

// GetAllSystemsNodes 获取所有系统的节点信息
func GetAllSystemsNodes(c *gin.Context) {
	// 检查nodeService是否初始化
	if nodeService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "节点服务不可用，Redis连接失败"})
		return
	}

	// 获取所有系统
	systems, err := systemService.GetSystems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统列表失败"})
		return
	}

	// 获取所有系统的节点信息
	allNodeInfo, err := nodeService.GetAllSystemsNodeInfo(systems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"systems": allNodeInfo,
	})
}

// SearchNodes 搜索节点
func SearchNodes(c *gin.Context) {
	// 检查nodeService是否初始化
	if nodeService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "节点服务不可用，Redis连接失败"})
		return
	}

	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	nodes, err := nodeService.SearchNodesByKeyword(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
		"count": len(nodes),
	})
}

// GetHighLoadNodes 获取所有高负载节点
func GetHighLoadNodes(c *gin.Context) {
	// 检查nodeService是否初始化
	if nodeService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "节点服务不可用，Redis连接失败",
			"data":  []map[string]interface{}{},
		})
		return
	}

	// 获取带负载状态的系统列表
	systems, err := systemService.GetSystemsWithLoadStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统负载状态失败"})
		return
	}

	var highLoadNodes []map[string]interface{}

	// 遍历所有系统，找出高负载的
	for _, system := range systems {
		// 只处理高负载或离线的服务器
		if system.LoadStatus == "high" || system.Status != "up" {
			// 获取该系统的节点信息
			nodeInfo, err := nodeService.GetSystemNodeInfo(system.ID, system.Name)
			if err != nil {
				continue // 跳过获取失败的系统
			}

			// 将该系统的所有节点添加到高负载节点列表
			for _, node := range nodeInfo.Nodes {
				highLoadNodes = append(highLoadNodes, map[string]interface{}{
					"name":   node.Name,
					"type":   node.Type,
					"id":     node.ID,
					"online": node.Online,
				})
			}
		}
	}

	c.JSON(http.StatusOK, highLoadNodes)
}