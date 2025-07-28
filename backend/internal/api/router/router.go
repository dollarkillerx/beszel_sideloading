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

	// 静态文件服务 - 必须在API路由之前
	setupStaticRoutes(r)

	// API路由组
	setupAPIRoutes(r)

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
	
	// 为所有具体的静态文件类型提供服务
	r.GET("/*.css", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, filepath.Base(c.Request.URL.Path)))
	})
	r.GET("/*.js", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, filepath.Base(c.Request.URL.Path)))
	})
	r.GET("/*.js.map", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, filepath.Base(c.Request.URL.Path)))
	})
	r.GET("/*.svg", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, filepath.Base(c.Request.URL.Path)))
	})
	r.GET("/*.png", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, filepath.Base(c.Request.URL.Path)))
	})
	r.GET("/*.ico", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, filepath.Base(c.Request.URL.Path)))
	})
	
	// 根路径返回index.html
	r.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(staticDir, "index.html"))
	})
}