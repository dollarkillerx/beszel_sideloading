package handlers

import (
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var systemService *service.SystemService

// InitHandlers 初始化处理器
func InitHandlers(svc *service.SystemService) {
	systemService = svc
}

// GetSystems 获取所有系统列表
func GetSystems(c *gin.Context) {
	systems, err := systemService.GetSystems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统列表失败", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"systems": systems})
}

// GetSystemSummary 获取系统摘要
func GetSystemSummary(c *gin.Context) {
	summary, err := systemService.GetSystemSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统摘要失败", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, summary)
}

// GetSystemStats 获取指定系统的统计数据
func GetSystemStats(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}
	
	// 获取条数参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	
	stats, err := systemService.GetSystemStats(systemID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统统计数据失败", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"stats": stats, "total": len(stats)})
}

// GetSystemsWithAvgStats 获取所有系统及其平均统计数据
func GetSystemsWithAvgStats(c *gin.Context) {
	systems, err := systemService.GetSystemsWithAvgStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统统计数据失败", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"systems": systems})
}