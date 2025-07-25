package database

import (
	"backend/pkg/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 初始化数据库连接和迁移
func Init(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// 检查是否需要迁移node_tags表结构
	if DB.Migrator().HasTable("node_tags") {
		// 检查是否有旧的列结构
		if DB.Migrator().HasColumn(&models.NodeTag{}, "node_type") {
			log.Println("检测到旧的node_tags表结构，正在迁移...")
			DB.Migrator().DropTable("node_tags")
		}
	}
	
	// 自动迁移数据库模型
	err = DB.AutoMigrate(&models.SystemThreshold{}, &models.NodeTag{})
	if err != nil {
		return err
	}
	
	log.Println("数据库表结构迁移完成")

	log.Println("数据库初始化成功")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}