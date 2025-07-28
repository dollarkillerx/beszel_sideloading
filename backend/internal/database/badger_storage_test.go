package database

import (
	"backend/pkg/models"
	"os"
	"testing"
)

func TestBadgerStorage(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "badger-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建存储实例
	storage, err := NewBadgerStorage(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	defer storage.Close()

	// 测试系统阈值操作
	t.Run("SystemThreshold", func(t *testing.T) {
		threshold := &models.SystemThreshold{
			SystemID:      "test-system-1",
			CPUAlertLimit: 90.0,
			MemAlertLimit: 85.0,
		}

		// 创建阈值
		err := storage.CreateOrUpdateThreshold(threshold)
		if err != nil {
			t.Errorf("Failed to create threshold: %v", err)
		}

		// 获取阈值
		got, err := storage.GetThreshold("test-system-1")
		if err != nil {
			t.Errorf("Failed to get threshold: %v", err)
		}
		if got.CPUAlertLimit != 90.0 {
			t.Errorf("Expected CPU limit 90.0, got %f", got.CPUAlertLimit)
		}

		// 更新阈值
		threshold.CPUAlertLimit = 95.0
		err = storage.CreateOrUpdateThreshold(threshold)
		if err != nil {
			t.Errorf("Failed to update threshold: %v", err)
		}

		// 验证更新
		got, err = storage.GetThreshold("test-system-1")
		if err != nil {
			t.Errorf("Failed to get updated threshold: %v", err)
		}
		if got.CPUAlertLimit != 95.0 {
			t.Errorf("Expected updated CPU limit 95.0, got %f", got.CPUAlertLimit)
		}

		// 列出所有阈值
		all, err := storage.ListThresholds()
		if err != nil {
			t.Errorf("Failed to list thresholds: %v", err)
		}
		if len(all) != 1 {
			t.Errorf("Expected 1 threshold, got %d", len(all))
		}
	})

	// 测试节点标签操作
	t.Run("NodeTag", func(t *testing.T) {
		tag := &models.NodeTag{
			SystemID: "test-system-1",
			TagType:  "group",
			TagID:    1,
		}

		// 创建标签
		err := storage.CreateNodeTag(tag)
		if err != nil {
			t.Errorf("Failed to create tag: %v", err)
		}

		// 获取系统标签
		tags, err := storage.GetNodeTags("test-system-1")
		if err != nil {
			t.Errorf("Failed to get tags: %v", err)
		}
		if len(tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(tags))
		}

		// 根据类型和ID获取标签
		tags, err = storage.GetNodeTagsByTypeAndID("group", 1)
		if err != nil {
			t.Errorf("Failed to get tags by type and ID: %v", err)
		}
		if len(tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(tags))
		}

		// 删除标签
		err = storage.DeleteNodeTag("test-system-1", "group", 1)
		if err != nil {
			t.Errorf("Failed to delete tag: %v", err)
		}

		// 验证删除
		tags, err = storage.GetNodeTags("test-system-1")
		if err != nil {
			t.Errorf("Failed to get tags after delete: %v", err)
		}
		if len(tags) != 0 {
			t.Errorf("Expected 0 tags after delete, got %d", len(tags))
		}
	})
}