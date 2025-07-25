package router

import (
	"backend/internal/api/handlers"
	"backend/internal/config"
	"backend/internal/service"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, systemService *service.SystemService) *gin.Engine {
	// 初始化处理器
	handlers.InitHandlers(systemService)
	r := gin.Default()

	// 配置CORS中间件
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.CORS.AllowOrigins
	corsConfig.AllowMethods = cfg.CORS.AllowMethods
	corsConfig.AllowHeaders = cfg.CORS.AllowHeaders
	r.Use(cors.New(corsConfig))

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
	})

	// API路由组
	setupAPIRoutes(r)

	// 静态文件服务
	setupStaticRoutes(r)

	return r
}

// setupAPIRoutes 设置API路由
func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 系统相关路由
		systems := api.Group("/systems")
		{
			systems.GET("", handlers.GetSystems)
			systems.GET("/summary", handlers.GetSystemSummary)
			systems.GET("/stats", handlers.GetSystemsWithAvgStats)
			systems.GET("/:id/stats", handlers.GetSystemStats)
			
			// 阈值配置路由
			systems.GET("/:id/threshold", handlers.GetThreshold)
			systems.PUT("/:id/threshold", handlers.UpdateThreshold)
			systems.DELETE("/:id/threshold", handlers.DeleteThreshold)
		}
		
		// 全局阈值配置路由
		api.GET("/thresholds", handlers.GetAllThresholds)
		
		// 服务器标签路由
		systems.POST("/:id/tags", handlers.AddSystemTag)        // 添加服务器标签
		systems.DELETE("/:id/tags", handlers.RemoveSystemTag)   // 删除服务器标签
		systems.GET("/:id/tags", handlers.GetSystemTags)        // 获取服务器标签
		
		// 节点负载状态查询
		api.POST("/nodes/load-status", handlers.GetNodeLoadStatus) // 批量查询节点负载状态
	}
}

// setupStaticRoutes 设置静态文件路由
func setupStaticRoutes(r *gin.Engine) {
	// 静态文件目录
	staticDir := "./static"
	
	// 提供静态文件服务
	r.Static("/assets", filepath.Join(staticDir, "assets"))
	r.StaticFile("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))
	
	// SPA 路由处理 - 所有非API路由都返回index.html
	r.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		
		// 其他所有请求都返回index.html（用于SPA路由）
		c.File(filepath.Join(staticDir, "index.html"))
	})
}