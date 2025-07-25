package main

import (
	"backend/internal/config"
	"backend/internal/server"
	"log"
)

func main() {
	// 加载配置
	cfg := config.Load()
	log.Println("Configuration loaded successfully")

	// 创建服务器实例
	srv := server.New(cfg)

	// 启动服务器
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}