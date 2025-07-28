package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/server"
	"log"
)

func main() {
	// 加载配置
	cfg := config.Load()
	log.Println("Configuration loaded successfully")

	// 初始化数据库
	if err := database.Init(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close() // 确保程序退出时关闭数据库
	log.Println("Database initialized successfully")

	// 创建服务器实例
	srv := server.New(cfg)

	// 启动服务器
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}