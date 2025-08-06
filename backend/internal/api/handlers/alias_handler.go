package handlers

import (
	"backend/internal/service"
	"backend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var aliasService *service.AliasService

// InitAliasHandler 初始化别名处理器
func InitAliasHandler() {
	aliasService = service.NewAliasService()
}

// SetSystemAlias 设置服务器别名
func SetSystemAlias(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	var request models.SystemAliasRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
		return
	}

	// 设置别名
	err := aliasService.SetAlias(systemID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SystemAliasResponse{
		Success: "别名设置成功",
	})
}

// GetSystemAlias 获取服务器别名
func GetSystemAlias(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	alias, err := aliasService.GetAlias(systemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if alias == nil {
		c.JSON(http.StatusOK, models.SystemAliasResponse{
			Alias: nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.SystemAliasResponse{
		Alias: alias,
	})
}

// DeleteSystemAlias 删除服务器别名
func DeleteSystemAlias(c *gin.Context) {
	systemID := c.Param("id")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统ID不能为空"})
		return
	}

	err := aliasService.DeleteAlias(systemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SystemAliasResponse{
		Success: "别名删除成功",
	})
}

// GetAllAliases 获取所有别名
func GetAllAliases(c *gin.Context) {
	aliases, err := aliasService.GetAllAliases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aliases": aliases,
	})
}