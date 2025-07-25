import React, { useState, useEffect } from 'react';

interface SystemThreshold {
  id: number;
  system_id: string;
  cpu_alert_limit: number;
  mem_alert_limit: number;
  net_up_max: number;
  net_down_max: number;
  net_up_alert: number;
  net_down_alert: number;
  created_at: string;
  updated_at: string;
}

interface SystemStats {
  id: number;
  name: string;
  host: string;
  port: string;
  status: string;
  avg_cpu: number;
  avg_mem_pct: number;
  avg_net_sent: number;
  avg_net_recv: number;
  last_update: string;
  load_status: string;
}

interface ThresholdConfigProps {
  system: SystemStats;
  onClose: () => void;
  onSave: () => void;
}

const ThresholdConfig: React.FC<ThresholdConfigProps> = ({ system, onClose, onSave }) => {
  const [threshold, setThreshold] = useState<SystemThreshold | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const API_BASE = 'http://localhost:8080/api';

  useEffect(() => {
    fetchThreshold();
  }, [system.id]);

  const fetchThreshold = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/systems/${system.id}/threshold`);
      if (!response.ok) {
        throw new Error('获取阈值配置失败');
      }
      const data = await response.json();
      setThreshold(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取阈值配置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!threshold) return;

    try {
      setSaving(true);
      
      // 处理网络最大值：如果小于1则设置为2000
      const processedThreshold = {
        ...threshold,
        net_up_max: threshold.net_up_max < 1 ? 2000 : threshold.net_up_max,
        net_down_max: threshold.net_down_max < 1 ? 2000 : threshold.net_down_max,
      };
      
      // 调试信息：显示即将发送的数据
      console.log('原始阈值配置:', threshold);
      console.log('处理后阈值配置:', processedThreshold);
      
      const response = await fetch(`${API_BASE}/systems/${system.id}/threshold`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(processedThreshold),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '保存阈值配置失败');
      }

      setError(null);
      
      // 如果有值被自动调整，显示提示信息
      if (threshold.net_up_max > 0 && threshold.net_up_max < 1) {
        console.log(`上行最大值从 ${threshold.net_up_max} 自动调整为 2000 Mbps`);
      }
      if (threshold.net_down_max > 0 && threshold.net_down_max < 1) {
        console.log(`下行最大值从 ${threshold.net_down_max} 自动调整为 2000 Mbps`);
      }
      
      onSave();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : '保存阈值配置失败');
    } finally {
      setSaving(false);
    }
  };

  const handleInputChange = (field: keyof SystemThreshold, value: number) => {
    if (threshold) {
      // 处理NaN值，如果解析失败则设为0
      let validValue = isNaN(value) ? 0 : value;
      
      // 确保值不为负数
      if (validValue < 0) {
        validValue = 0;
      }
      
      setThreshold({
        ...threshold,
        [field]: validValue,
      });
    }
  };

  if (loading) {
    return (
      <div className="threshold-modal">
        <div className="threshold-modal-content">
          <div className="loading">正在加载配置...</div>
        </div>
      </div>
    );
  }

  if (!threshold) {
    return (
      <div className="threshold-modal">
        <div className="threshold-modal-content">
          <div className="error">无法加载阈值配置</div>
          <button onClick={onClose}>关闭</button>
        </div>
      </div>
    );
  }

  return (
    <div className="threshold-modal">
      <div className="threshold-modal-content">
        <div className="threshold-header">
          <div>
            <h3>阈值配置 - {system.name}</h3>
            <p className="system-info">{system.host}:{system.port}</p>
          </div>
          <button className="close-button" onClick={onClose}>✕</button>
        </div>
        
        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        <div className="threshold-form">
          <div className="form-group">
            <label>CPU 告警阈值 (%)</label>
            <input
              type="number"
              min="0"
              max="100"
              step="0.1"
              value={threshold.cpu_alert_limit}
              onChange={(e) => {
                const value = e.target.value === '' ? 0 : parseFloat(e.target.value);
                handleInputChange('cpu_alert_limit', value);
              }}
            />
          </div>

          <div className="form-group">
            <label>内存告警阈值 (%)</label>
            <input
              type="number"
              min="0"
              max="100"
              step="0.1"
              value={threshold.mem_alert_limit}
              onChange={(e) => {
                const value = e.target.value === '' ? 0 : parseFloat(e.target.value);
                handleInputChange('mem_alert_limit', value);
              }}
            />
          </div>

          <div className="form-group">
            <label>上行最大值 (Mbps)</label>
            <input
              type="number"
              min="0"
              step="0.1"
              value={threshold.net_up_max}
              onChange={(e) => {
                const value = e.target.value === '' ? 0 : parseFloat(e.target.value);
                handleInputChange('net_up_max', value);
              }}
            />
            <small>
              历史记录的最大上行速度，系统会自动更新。
              {threshold.net_up_max === 0 && (
                <span className="warning-text"> ⚠️ 当前为0，不会检查上行网络负载</span>
              )}
              {threshold.net_up_max > 0 && threshold.net_up_max < 1 && (
                <span className="info-text"> ℹ️ 小于1将自动设置为2000 Mbps</span>
              )}
            </small>
          </div>

          <div className="form-group">
            <label>下行最大值 (Mbps)</label>
            <input
              type="number"
              min="0"
              step="0.1"
              value={threshold.net_down_max}
              onChange={(e) => {
                const value = e.target.value === '' ? 0 : parseFloat(e.target.value);
                handleInputChange('net_down_max', value);
              }}
            />
            <small>
              历史记录的最大下行速度，系统会自动更新。
              {threshold.net_down_max === 0 && (
                <span className="warning-text"> ⚠️ 当前为0，不会检查下行网络负载</span>
              )}
              {threshold.net_down_max > 0 && threshold.net_down_max < 1 && (
                <span className="info-text"> ℹ️ 小于1将自动设置为2000 Mbps</span>
              )}
            </small>
          </div>

          <div className="form-group">
            <label>上行告警阈值 (%)</label>
            <input
              type="number"
              min="0"
              max="100"
              step="0.1"
              value={threshold.net_up_alert}
              onChange={(e) => {
                const value = e.target.value === '' ? 0 : parseFloat(e.target.value);
                handleInputChange('net_up_alert', value);
              }}
            />
            <small>当上行速度达到最大值的该百分比时告警</small>
          </div>

          <div className="form-group">
            <label>下行告警阈值 (%)</label>
            <input
              type="number"
              min="0"
              max="100"
              step="0.1"
              value={threshold.net_down_alert}
              onChange={(e) => {
                const value = e.target.value === '' ? 0 : parseFloat(e.target.value);
                handleInputChange('net_down_alert', value);
              }}
            />
            <small>当下行速度达到最大值的该百分比时告警</small>
          </div>
        </div>

        <div className="threshold-actions">
          <button className="cancel-button" onClick={onClose}>
            取消
          </button>
          <button 
            className="save-button" 
            onClick={handleSave}
            disabled={saving}
          >
            {saving ? '保存中...' : '保存'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ThresholdConfig;