package config

import (
	"os"
)

// Config 应用配置
type Config struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	CORS       CORSConfig       `json:"cors"`
	PocketBase PocketBaseConfig `json:"pocketbase"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `json:"path"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowOrigins []string `json:"allow_origins"`
	AllowMethods []string `json:"allow_methods"`
	AllowHeaders []string `json:"allow_headers"`
}

// PocketBaseConfig PocketBase配置
type PocketBaseConfig struct {
	BaseURL  string `json:"base_url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", ""),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "badger_data"),
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"*"},
		},
		PocketBase: PocketBaseConfig{
			BaseURL:  getEnv("POCKETBASE_URL", "https://bz.baidua.top"),
			Email:    getEnv("POCKETBASE_EMAIL", "Spike.wook@gmail.com"),
			Password: getEnv("POCKETBASE_PASSWORD", "adadmin/1213"),
		},
	}
}

// GetAddress 获取服务器地址
func (c *Config) GetAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

