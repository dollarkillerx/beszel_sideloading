package database

import (
	"log"
)

var storage Storage

// Init 初始化数据库存储
func Init(dbPath string) error {
	// 使用BadgerDB作为存储后端
	badgerStorage, err := NewBadgerStorage(dbPath)
	if err != nil {
		return err
	}
	
	storage = badgerStorage
	log.Println("数据库存储初始化成功")
	return nil
}

// GetStorage 获取存储实例
func GetStorage() Storage {
	return storage
}

// Close 关闭数据库
func Close() error {
	if storage != nil {
		return storage.Close()
	}
	return nil
}