package handlers

import (
	"backend/internal/service"
	"backend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NodeTagHandler 节点标签处理器
type NodeTagHandler struct {
	nodeTagService *service.NodeTagService
}

// NewNodeTagHandler 创建节点标签处理器
func NewNodeTagHandler() *NodeTagHandler {
	return &NodeTagHandler{
		nodeTagService: service.NewNodeTagService(),
	}
}

// AddSystemTag 添加服务器标签
// POST /api/systems/:id/tags
func (h *NodeTagHandler) AddSystemTag(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	var request models.NodeTagRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误: " + err.Error()})
		return
	}

	err := h.nodeTagService.AddTag(systemID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.NodeTagsResponse{Success: "ok"})
}

// RemoveSystemTag 删除服务器标签
// DELETE /api/systems/:id/tags
func (h *NodeTagHandler) RemoveSystemTag(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	var request models.NodeTagRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误: " + err.Error()})
		return
	}

	err := h.nodeTagService.RemoveTag(systemID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.NodeTagsResponse{Success: "ok"})
}

// GetSystemTags 获取服务器的所有标签
// GET /api/systems/:id/tags
func (h *NodeTagHandler) GetSystemTags(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	tags, err := h.nodeTagService.GetSystemTags(systemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.NodeTagsResponse{
		Success: "ok",
		Tags:    tags,
	})
}

// GetNodeLoadStatus 获取节点负载状态
// POST /api/nodes/load-status
func (h *NodeTagHandler) GetNodeLoadStatus(c *gin.Context) {
	var requests []models.NodeLoadRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数格式错误: " + err.Error()})
		return
	}

	responses, err := h.nodeTagService.GetNodeLoadStatus(requests, systemService, thresholdHandler.thresholdService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

