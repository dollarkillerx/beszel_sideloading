package router

import (
	"backend/internal/api/handlers"
	"backend/internal/config"
	"backend/internal/service"
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

	// 静态文件服务 - 使用NoRoute来处理所有未匹配的路由
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
		
		// 服务器别名路由
		systems.PUT("/:id/alias", handlers.SetSystemAlias)      // 设置服务器别名
		systems.GET("/:id/alias", handlers.GetSystemAlias)      // 获取服务器别名
		systems.DELETE("/:id/alias", handlers.DeleteSystemAlias) // 删除服务器别名
		
		// 所有别名
		api.GET("/aliases", handlers.GetAllAliases)             // 获取所有别名
		
		// 节点信息路由
		systems.GET("/:id/nodes", handlers.GetSystemNodes)      // 获取系统节点信息
		api.GET("/nodes", handlers.GetAllSystemsNodes)          // 获取所有系统节点信息
		api.GET("/nodes/search", handlers.SearchNodes)          // 搜索节点
		api.GET("/nodes/load-status", handlers.GetHighLoadNodes) // 获取高负载节点
	}
}

// setupStaticRoutes 设置静态文件路由
func setupStaticRoutes(r *gin.Engine) {
	// 静态文件目录
	staticDir := "./static"
	
	// 根路径返回index.html
	r.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, "index.html"))
	})
	
	// 使用NoRoute处理所有未匹配的路由
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// 如果是API请求，返回404
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}
		
		// 如果是健康检查，跳过（已经被处理了）
		if path == "/health" {
			return
		}
		
		// 检查是否为静态文件请求
		ext := filepath.Ext(path)
		if ext == ".css" || ext == ".js" || ext == ".svg" || ext == ".png" || ext == ".ico" || ext == ".map" {
			// 尝试提供静态文件
			filename := filepath.Base(path)
			c.File(filepath.Join(staticDir, filename))
			return
		}
		
		// 其他所有请求都返回index.html（用于SPA路由）
		c.File(filepath.Join(staticDir, "index.html"))
	})
}