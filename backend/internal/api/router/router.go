package router

import (
	"backend/internal/api/handlers"
	"backend/internal/config"
	"backend/internal/service"

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
	}
}