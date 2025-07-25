package service

import (
	"backend/pkg/models"
	"testing"
	"time"
)

func TestCalculateLoadStatus(t *testing.T) {
	s := &SystemService{}

	// 测试用例1: CPU超过阈值，应该返回high
	system1 := &models.SystemWithAvgStats{
		System: models.System{
			ID:   "test1",
			Name: "TestServer1",
		},
		AvgCPU:     95.0, // 超过90%阈值
		AvgMemPct:  50.0, // 正常
		AvgNetSent: 1.0,  // 正常
		AvgNetRecv: 1.0,  // 正常
		LastUpdate: time.Now(),
	}
	
	threshold1 := &models.SystemThreshold{
		CPUAlertLimit: 90.0,
		MemAlertLimit: 90.0,
		NetUpMax:      0,    // 未设置网络最大值
		NetDownMax:    0,    // 未设置网络最大值
		NetUpAlert:    80.0,
		NetDownAlert:  80.0,
	}
	
	result1 := s.calculateLoadStatus(system1, threshold1)
	if result1 != "high" {
		t.Errorf("测试用例1失败: 期望 'high', 得到 '%s'", result1)
	}

	// 测试用例2: 内存超过阈值，应该返回high
	system2 := &models.SystemWithAvgStats{
		System: models.System{
			ID:   "test2",
			Name: "TestServer2",
		},
		AvgCPU:     50.0, // 正常
		AvgMemPct:  95.0, // 超过90%阈值
		AvgNetSent: 1.0,  // 正常
		AvgNetRecv: 1.0,  // 正常
		LastUpdate: time.Now(),
	}
	
	result2 := s.calculateLoadStatus(system2, threshold1)
	if result2 != "high" {
		t.Errorf("测试用例2失败: 期望 'high', 得到 '%s'", result2)
	}

	// 测试用例3: 所有指标正常，网络最大值为0，应该返回normal
	system3 := &models.SystemWithAvgStats{
		System: models.System{
			ID:   "test3",
			Name: "TestServer3",
		},
		AvgCPU:     50.0, // 正常
		AvgMemPct:  50.0, // 正常
		AvgNetSent: 100.0, // 虽然很高，但NetUpMax=0，不检查
		AvgNetRecv: 100.0, // 虽然很高，但NetDownMax=0，不检查
		LastUpdate: time.Now(),
	}
	
	result3 := s.calculateLoadStatus(system3, threshold1)
	if result3 != "normal" {
		t.Errorf("测试用例3失败: 期望 'normal', 得到 '%s'", result3)
	}

	// 测试用例4: 网络上行超过阈值，应该返回high
	system4 := &models.SystemWithAvgStats{
		System: models.System{
			ID:   "test4",
			Name: "TestServer4",
		},
		AvgCPU:     50.0, // 正常
		AvgMemPct:  50.0, // 正常
		AvgNetSent: 10.0, // 10 MB/s * 8 = 80 Mbps，超过80%阈值(64 Mbps)
		AvgNetRecv: 1.0,  // 正常
		LastUpdate: time.Now(),
	}
	
	threshold4 := &models.SystemThreshold{
		CPUAlertLimit: 90.0,
		MemAlertLimit: 90.0,
		NetUpMax:      100.0, // 100 Mbps最大值
		NetDownMax:    100.0, // 100 Mbps最大值
		NetUpAlert:    80.0,  // 80%阈值 = 80 Mbps
		NetDownAlert:  80.0,  // 80%阈值 = 80 Mbps
	}
	
	result4 := s.calculateLoadStatus(system4, threshold4)
	if result4 != "high" {
		t.Errorf("测试用例4失败: 期望 'high', 得到 '%s'", result4)
	}

	// 测试用例5: 所有指标正常，包括网络，应该返回normal
	system5 := &models.SystemWithAvgStats{
		System: models.System{
			ID:   "test5",
			Name: "TestServer5",
		},
		AvgCPU:     50.0, // 正常
		AvgMemPct:  50.0, // 正常
		AvgNetSent: 5.0,  // 5 MB/s * 8 = 40 Mbps，低于80%阈值(64 Mbps)
		AvgNetRecv: 5.0,  // 5 MB/s * 8 = 40 Mbps，低于80%阈值(64 Mbps)
		LastUpdate: time.Now(),
	}
	
	result5 := s.calculateLoadStatus(system5, threshold4)
	if result5 != "normal" {
		t.Errorf("测试用例5失败: 期望 'normal', 得到 '%s'", result5)
	}
}