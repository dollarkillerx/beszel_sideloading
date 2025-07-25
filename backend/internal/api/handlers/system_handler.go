package handlers

import (
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var systemService *service.SystemService
var thresholdHandler *ThresholdHandler

// InitHandlers 初始化处理器
func InitHandlers(svc *service.SystemService) {
	systemService = svc
	thresholdHandler = NewThresholdHandler()
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

// GetSystemsWithAvgStats 获取所有系统及其平均统计数据（包含负载状态）
func GetSystemsWithAvgStats(c *gin.Context) {
	systems, err := systemService.GetSystemsWithLoadStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统统计数据失败", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"systems": systems})
}

// 阈值配置相关的全局函数包装器
func GetThreshold(c *gin.Context) {
	thresholdHandler.GetThreshold(c)
}

func UpdateThreshold(c *gin.Context) {
	thresholdHandler.UpdateThreshold(c)
}

func DeleteThreshold(c *gin.Context) {
	thresholdHandler.DeleteThreshold(c)
}

func GetAllThresholds(c *gin.Context) {
	thresholdHandler.GetAllThresholds(c)
}