package database

import (
	"backend/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

// BadgerStorage BadgerDB存储实现
type BadgerStorage struct {
	db *badger.DB
}

// NewBadgerStorage 创建BadgerDB存储实例
func NewBadgerStorage(dbPath string) (*BadgerStorage, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil // 禁用BadgerDB的日志输出，避免日志过多
	
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger db: %w", err)
	}

	log.Println("BadgerDB 初始化成功")
	return &BadgerStorage{db: db}, nil
}

// Close 关闭数据库
func (s *BadgerStorage) Close() error {
	return s.db.Close()
}

// 生成键前缀
func (s *BadgerStorage) thresholdKey(systemID string) []byte {
	return []byte(fmt.Sprintf("threshold:%s", systemID))
}

func (s *BadgerStorage) aliasKey(systemID string) []byte {
	return []byte(fmt.Sprintf("alias:%s", systemID))
}

func (s *BadgerStorage) aliasPrefix() []byte {
	return []byte("alias:")
}

// CreateOrUpdateThreshold 创建或更新系统阈值
func (s *BadgerStorage) CreateOrUpdateThreshold(threshold *models.SystemThreshold) error {
	return s.db.Update(func(txn *badger.Txn) error {
		// 如果是更新，先获取旧数据保留ID和创建时间
		key := s.thresholdKey(threshold.SystemID)
		item, err := txn.Get(key)
		if err == nil {
			var existing models.SystemThreshold
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &existing)
			})
			if err == nil {
				if threshold.ID == 0 {
					threshold.ID = existing.ID
				}
				if threshold.CreatedAt.IsZero() {
					threshold.CreatedAt = existing.CreatedAt
				}
			}
		} else {
			// 新建记录，生成ID和创建时间
			if threshold.ID == 0 {
				threshold.ID = uint(time.Now().UnixNano())
			}
			if threshold.CreatedAt.IsZero() {
				threshold.CreatedAt = time.Now()
			}
		}

		threshold.UpdatedAt = time.Now()

		data, err := json.Marshal(threshold)
		if err != nil {
			return err
		}

		return txn.Set(key, data)
	})
}

// GetThreshold 获取系统阈值
func (s *BadgerStorage) GetThreshold(systemID string) (*models.SystemThreshold, error) {
	var threshold models.SystemThreshold

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.thresholdKey(systemID))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &threshold)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}

	return &threshold, err
}

// ListThresholds 列出所有系统阈值
func (s *BadgerStorage) ListThresholds() ([]*models.SystemThreshold, error) {
	var thresholds []*models.SystemThreshold

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("threshold:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var threshold models.SystemThreshold

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &threshold)
			})
			if err != nil {
				log.Printf("Failed to unmarshal threshold: %v", err)
				continue
			}

			thresholds = append(thresholds, &threshold)
		}
		return nil
	})

	return thresholds, err
}

// DeleteThreshold 删除系统阈值
func (s *BadgerStorage) DeleteThreshold(systemID string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(s.thresholdKey(systemID))
	})
}

// SetSystemAlias 设置系统别名（创建或更新）
func (s *BadgerStorage) SetSystemAlias(alias *models.SystemAlias) error {
	return s.db.Update(func(txn *badger.Txn) error {
		// 先检查是否已存在
		key := s.aliasKey(alias.SystemID)
		item, err := txn.Get(key)
		if err == nil {
			// 已存在，保留ID和创建时间
			var existing models.SystemAlias
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &existing)
			})
			if err == nil {
				if alias.ID == 0 {
					alias.ID = existing.ID
				}
				if alias.CreatedAt.IsZero() {
					alias.CreatedAt = existing.CreatedAt
				}
			}
		} else {
			// 新建记录
			if alias.ID == 0 {
				alias.ID = uint(time.Now().UnixNano())
			}
			if alias.CreatedAt.IsZero() {
				alias.CreatedAt = time.Now()
			}
		}

		alias.UpdatedAt = time.Now()

		data, err := json.Marshal(alias)
		if err != nil {
			return err
		}

		return txn.Set(key, data)
	})
}

// GetSystemAlias 获取系统别名
func (s *BadgerStorage) GetSystemAlias(systemID string) (*models.SystemAlias, error) {
	var alias models.SystemAlias

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.aliasKey(systemID))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &alias)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}

	return &alias, err
}

// GetAllSystemAliases 获取所有系统别名
func (s *BadgerStorage) GetAllSystemAliases() ([]*models.SystemAlias, error) {
	var aliases []*models.SystemAlias

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = s.aliasPrefix()
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var alias models.SystemAlias

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &alias)
			})
			if err != nil {
				log.Printf("Failed to unmarshal alias: %v", err)
				continue
			}

			aliases = append(aliases, &alias)
		}
		return nil
	})

	return aliases, err
}

// DeleteSystemAlias 删除系统别名
func (s *BadgerStorage) DeleteSystemAlias(systemID string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(s.aliasKey(systemID))
	})
}