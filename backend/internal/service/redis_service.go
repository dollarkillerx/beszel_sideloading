package service

import (
	"backend/internal/config"
	"backend/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

// RedisService Redis服务
type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService 创建Redis服务
func NewRedisService(cfg *config.Config) (*RedisService, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password, // 支持空密码
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()

	// 测试Redis连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("❌ Redis连接失败 [%s] 数据库:%d - %v", addr, cfg.Redis.DB, err)
		rdb.Close()
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	log.Printf("✅ Redis连接成功 [%s] 数据库:%d", addr, cfg.Redis.DB)

	return &RedisService{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisService) Close() error {
	return r.client.Close()
}

// GetNodesByAlias 根据别名模糊匹配获取节点信息
func (r *RedisService) GetNodesByAlias(alias string) ([]models.V2boardNode, error) {
	var nodes []models.V2boardNode

	// 使用SCAN命令获取所有v2board_database_AGENT_*的key
	pattern := "v2board_database_AGENT_*"
	iter := r.client.Scan(r.ctx, 0, pattern, 0).Iterator()

	for iter.Next(r.ctx) {
		key := iter.Val()
		
		// 获取key的值
		val, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			log.Printf("获取Redis key %s 失败: %v", key, err)
			continue
		}

		// 解析JSON
		var node models.V2boardNode
		if err := json.Unmarshal([]byte(val), &node); err != nil {
			log.Printf("解析Redis key %s 的JSON失败: %v", key, err)
			continue
		}

		// 模糊匹配节点名称
		if strings.Contains(node.Name, alias) {
			nodes = append(nodes, node)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("扫描Redis keys失败: %w", err)
	}

	return nodes, nil
}

// GetAllNodes 获取所有节点信息（用于调试）
func (r *RedisService) GetAllNodes() ([]models.V2boardNode, error) {
	var nodes []models.V2boardNode

	pattern := "v2board_database_AGENT_*"
	iter := r.client.Scan(r.ctx, 0, pattern, 0).Iterator()

	for iter.Next(r.ctx) {
		key := iter.Val()
		
		val, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			log.Printf("获取Redis key %s 失败: %v", key, err)
			continue
		}

		var node models.V2boardNode
		if err := json.Unmarshal([]byte(val), &node); err != nil {
			log.Printf("解析Redis key %s 的JSON失败: %v", key, err)
			continue
		}

		nodes = append(nodes, node)
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("扫描Redis keys失败: %w", err)
	}

	return nodes, nil
}