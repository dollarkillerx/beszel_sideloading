package server

import (
	"backend/internal/api/handlers"
	"backend/internal/api/router"
	"backend/internal/config"
	"backend/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Server 服务器结构
type Server struct {
	config        *config.Config
	httpServer    *http.Server
	router        *gin.Engine
	systemService *service.SystemService
	redisService  *service.RedisService
	nodeService   *service.NodeService
}

// New 创建新的服务器实例
func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 初始化服务
	if err := s.initServices(); err != nil {
		return err
	}

	// 设置路由
	s.router = router.SetupRouter(s.config, s.systemService)

	// 创建HTTP服务器
	s.httpServer = &http.Server{
		Addr:         s.config.GetAddress(),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	log.Printf("Server starting on %s", s.config.GetAddress())
	
	// 在goroutine中启动服务器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号优雅关闭
	s.waitForShutdown()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	log.Println("Shutting down server...")
	
	// 关闭Redis连接
	if s.redisService != nil {
		if err := s.redisService.Close(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
		} else {
			log.Println("Redis connection closed")
		}
	}

	// 创建5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server stopped gracefully")
	return nil
}

// initServices 初始化服务
func (s *Server) initServices() error {
	log.Println("Initializing services...")
	
	// 初始化系统服务
	s.systemService = service.NewSystemService(s.config)
	
	// 初始化Redis服务
	var err error
	s.redisService, err = service.NewRedisService(s.config)
	if err != nil {
		log.Printf("Redis服务初始化失败: %v", err)
		// Redis失败不应该阻止服务启动，但会影响节点查询功能
		log.Println("节点查询功能将不可用")
	}
	
	// 初始化节点服务
	if s.redisService != nil {
		s.nodeService = service.NewNodeService(s.redisService)
		// 设置SystemService的NodeService引用
		s.systemService.SetNodeService(s.nodeService)
		// 初始化节点处理器
		handlers.InitNodeHandler(s.nodeService)
		log.Println("节点服务初始化成功")
	}
	
	log.Println("Services initialized successfully")
	return nil
}

// waitForShutdown 等待关闭信号
func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit
	log.Println("Received shutdown signal")
	
	if err := s.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}