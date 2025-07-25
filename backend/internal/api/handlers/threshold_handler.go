package handlers

import (
	"backend/internal/service"
	"backend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ThresholdHandler 阈值配置处理器
type ThresholdHandler struct {
	thresholdService *service.ThresholdService
}

// NewThresholdHandler 创建阈值配置处理器
func NewThresholdHandler() *ThresholdHandler {
	return &ThresholdHandler{
		thresholdService: service.NewThresholdService(),
	}
}

// GetThreshold 获取系统阈值配置
// GET /api/systems/:id/threshold
func (h *ThresholdHandler) GetThreshold(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	threshold, err := h.thresholdService.GetThreshold(systemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, threshold)
}

// UpdateThreshold 更新系统阈值配置
// PUT /api/systems/:id/threshold
func (h *ThresholdHandler) UpdateThreshold(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	var threshold models.SystemThreshold
	if err := c.ShouldBindJSON(&threshold); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误: " + err.Error()})
		return
	}

	// 验证阈值范围
	if threshold.CPUAlertLimit < 0 || threshold.CPUAlertLimit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CPU阈值必须在0-100之间"})
		return
	}
	if threshold.MemAlertLimit < 0 || threshold.MemAlertLimit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内存阈值必须在0-100之间"})
		return
	}
	if threshold.NetUpAlert < 0 || threshold.NetUpAlert > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上行告警阈值必须在0-100之间"})
		return
	}
	if threshold.NetDownAlert < 0 || threshold.NetDownAlert > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "下行告警阈值必须在0-100之间"})
		return
	}

	err := h.thresholdService.UpdateThreshold(systemID, &threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的配置
	updatedThreshold, err := h.thresholdService.GetThreshold(systemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedThreshold)
}

// GetAllThresholds 获取所有系统的阈值配置
// GET /api/thresholds
func (h *ThresholdHandler) GetAllThresholds(c *gin.Context) {
	thresholds, err := h.thresholdService.GetAllThresholds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"thresholds": thresholds})
}

// DeleteThreshold 删除系统阈值配置
// DELETE /api/systems/:id/threshold
func (h *ThresholdHandler) DeleteThreshold(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	err := h.thresholdService.DeleteThreshold(systemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "阈值配置删除成功"})
}