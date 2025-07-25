import React, { useState, useEffect } from 'react';
import ThresholdConfig from './ThresholdConfig';

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
  load_status: string; // 新增：负载状态 'normal' | 'high'
}

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

interface SystemSummary {
  total: number;
  online: number;
  offline: number;
  unknown: number;
}

const ServerMonitor: React.FC = () => {
  const [systems, setSystems] = useState<SystemStats[]>([]);
  const [summary, setSummary] = useState<SystemSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date());
  const [showThresholdConfig, setShowThresholdConfig] = useState(false);
  const [selectedSystem, setSelectedSystem] = useState<SystemStats | null>(null);

  const API_BASE = 'http://localhost:8080/api';

  const fetchData = async () => {
    try {
      const [systemsResponse, summaryResponse] = await Promise.all([
        fetch(`${API_BASE}/systems/stats`),
        fetch(`${API_BASE}/systems/summary`)
      ]);

      if (!systemsResponse.ok || !summaryResponse.ok) {
        throw new Error('Failed to fetch data');
      }

      const systemsData = await systemsResponse.json();
      const summaryData = await summaryResponse.json();

      setSystems(systemsData.systems || []);
      setSummary(summaryData);
      setLastUpdate(new Date());
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    
    // 每5秒刷新一次数据
    const interval = setInterval(fetchData, 5000);
    
    return () => clearInterval(interval);
  }, []);

  const getStatusClass = (status: string) => {
    switch (status) {
      case 'up': return 'status-badge status-up';
      case 'down': return 'status-badge status-down';
      default: return 'status-badge status-unknown';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'up': return '在线';
      case 'down': return '离线';
      default: return '未知';
    }
  };

  const getLoadStatusClass = (loadStatus: string) => {
    switch (loadStatus) {
      case 'high': return 'load-status-high';
      case 'normal': return 'load-status-normal';
      default: return 'load-status-normal';
    }
  };

  const getLoadStatusText = (loadStatus: string) => {
    switch (loadStatus) {
      case 'high': return '高负载';
      case 'normal': return '正常';
      default: return '未知';
    }
  };

  const formatDateTime = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  if (loading) {
    return (
      <div className="loading-container">
        <div>正在加载服务器数据...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-container">
        <div>错误: {error}</div>
      </div>
    );
  }

  return (
    <div className="dashboard-container">
      <div className="dashboard-content">
        <div className="dashboard-header">
          <h1 className="dashboard-title">服务器监控面板</h1>
          <p className="dashboard-subtitle">
            最后更新时间: {lastUpdate.toLocaleTimeString()} (每5秒自动刷新)
          </p>
        </div>

        {/* Summary Cards */}
        {summary && (
          <div className="summary-grid">
            <div className="summary-card">
              <h3>服务器总数</h3>
              <p className="total">{summary.total}</p>
            </div>
            <div className="summary-card">
              <h3>在线</h3>
              <p className="online">{summary.online}</p>
            </div>
            <div className="summary-card">
              <h3>离线</h3>
              <p className="offline">{summary.offline}</p>
            </div>
            <div className="summary-card">
              <h3>未知</h3>
              <p className="unknown">{summary.unknown}</p>
            </div>
          </div>
        )}

        {/* Server List */}
        <div className="servers-table-container">
          <div className="table-header">
            <h2>服务器状态与性能</h2>
          </div>
          
          <div className="table-wrapper">
            <table className="servers-table">
              <thead>
                <tr>
                  <th>服务器</th>
                  <th>状态</th>
                  <th>负载状态</th>
                  <th>CPU (%)</th>
                  <th>内存 (%)</th>
                  <th>网络I/O (Mbps)</th>
                  <th>最后更新</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                {systems.map((system) => (
                  <tr key={system.id}>
                    <td>
                      <div className="server-info">
                        <div className="server-name">{system.name}</div>
                        <div className="server-address">{system.host}:{system.port}</div>
                      </div>
                    </td>
                    <td>
                      <span className={getStatusClass(system.status)}>
                        {getStatusText(system.status)}
                      </span>
                    </td>
                    <td>
                      <span className={getLoadStatusClass(system.load_status)}>
                        {getLoadStatusText(system.load_status)}
                      </span>
                    </td>
                    <td>
                      <div className="progress-container">
                        <div className="progress-bar">
                          <div 
                            className="progress-fill progress-cpu" 
                            style={{width: `${Math.min(system.avg_cpu, 100)}%`}}
                          ></div>
                        </div>
                        <span className="progress-text">{system.avg_cpu.toFixed(1)}%</span>
                      </div>
                    </td>
                    <td>
                      <div className="progress-container">
                        <div className="progress-bar">
                          <div 
                            className="progress-fill progress-memory" 
                            style={{width: `${Math.min(system.avg_mem_pct, 100)}%`}}
                          ></div>
                        </div>
                        <span className="progress-text">{system.avg_mem_pct.toFixed(1)}%</span>
                      </div>
                    </td>
                    <td>
                      <div className="network-metrics">
                        <div>↑ {(system.avg_net_sent * 8).toFixed(2)}</div>
                        <div>↓ {(system.avg_net_recv * 8).toFixed(2)}</div>
                      </div>
                    </td>
                    <td>
                      {system.last_update ? formatDateTime(system.last_update) : '无数据'}
                    </td>
                    <td>
                      <button 
                        className="config-button"
                        onClick={() => {
                          setSelectedSystem(system);
                          setShowThresholdConfig(true);
                        }}
                      >
                        配置阈值
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          
          {systems.length === 0 && (
            <div className="empty-state">
              未找到服务器
            </div>
          )}
        </div>

        {/* 阈值配置模态框 */}
        {showThresholdConfig && selectedSystem && (
          <ThresholdConfig
            system={selectedSystem}
            onClose={() => {
              setShowThresholdConfig(false);
              setSelectedSystem(null);
            }}
            onSave={() => {
              // 保存后刷新数据
              fetchData();
            }}
          />
        )}
      </div>
    </div>
  );
};

export default ServerMonitor;