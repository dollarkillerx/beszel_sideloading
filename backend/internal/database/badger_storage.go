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

func (s *BadgerStorage) nodeTagKey(systemID string, tagType string, tagID int) []byte {
	return []byte(fmt.Sprintf("nodetag:%s:%s:%d", systemID, tagType, tagID))
}

func (s *BadgerStorage) nodeTagPrefixBySystem(systemID string) []byte {
	return []byte(fmt.Sprintf("nodetag:%s:", systemID))
}

func (s *BadgerStorage) nodeTagPrefixByTypeAndID(tagType string, tagID int) []byte {
	return []byte("nodetag:")
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

// CreateNodeTag 创建节点标签
func (s *BadgerStorage) CreateNodeTag(tag *models.NodeTag) error {
	return s.db.Update(func(txn *badger.Txn) error {
		if tag.ID == 0 {
			tag.ID = uint(time.Now().UnixNano())
		}
		if tag.CreatedAt.IsZero() {
			tag.CreatedAt = time.Now()
		}
		tag.UpdatedAt = time.Now()

		data, err := json.Marshal(tag)
		if err != nil {
			return err
		}

		key := s.nodeTagKey(tag.SystemID, tag.TagType, tag.TagID)
		return txn.Set(key, data)
	})
}

// GetNodeTags 获取系统的所有标签
func (s *BadgerStorage) GetNodeTags(systemID string) ([]*models.NodeTag, error) {
	var tags []*models.NodeTag

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = s.nodeTagPrefixBySystem(systemID)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var tag models.NodeTag

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &tag)
			})
			if err != nil {
				log.Printf("Failed to unmarshal node tag: %v", err)
				continue
			}

			tags = append(tags, &tag)
		}
		return nil
	})

	return tags, err
}

// GetNodeTagsByTypeAndID 根据标签类型和ID获取所有相关节点
func (s *BadgerStorage) GetNodeTagsByTypeAndID(tagType string, tagID int) ([]*models.NodeTag, error) {
	var tags []*models.NodeTag

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = s.nodeTagPrefixByTypeAndID(tagType, tagID)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var tag models.NodeTag

			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &tag)
			})
			if err != nil {
				log.Printf("Failed to unmarshal node tag: %v", err)
				continue
			}

			// 过滤匹配的标签类型和ID
			if tag.TagType == tagType && tag.TagID == tagID {
				tags = append(tags, &tag)
			}
		}
		return nil
	})

	return tags, err
}

// DeleteNodeTag 删除特定标签
func (s *BadgerStorage) DeleteNodeTag(systemID string, tagType string, tagID int) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(s.nodeTagKey(systemID, tagType, tagID))
	})
}

// DeleteAllNodeTags 删除系统的所有标签
func (s *BadgerStorage) DeleteAllNodeTags(systemID string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = s.nodeTagPrefixBySystem(systemID)
		it := txn.NewIterator(opts)
		defer it.Close()

		var keys [][]byte
		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().KeyCopy(nil)
			keys = append(keys, key)
		}

		for _, key := range keys {
			if err := txn.Delete(key); err != nil {
				return err
			}
		}

		return nil
	})
}