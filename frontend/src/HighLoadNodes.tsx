import React, { useState, useEffect } from 'react';
import { API_BASE } from './utils/api';

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
  online_users: number;
  last_update: string;
  load_status: string;
}

interface HighLoadNode {
  name: string;
  type: string;
  id: number;
  online: number;
}

interface SystemSummary {
  total: number;
  online: number;
  offline: number;
  unknown: number;
  high_load: number; // 高负载节点数量
}

const HighLoadNodes: React.FC = () => {
  const [systems, setSystems] = useState<SystemStats[]>([]);
  const [summary, setSummary] = useState<SystemSummary | null>(null);
  const [highLoadNodes, setHighLoadNodes] = useState<HighLoadNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date());

  const fetchData = async () => {
    try {
      setLoading(true);

      // 获取所有系统统计数据
      const statsResponse = await fetch(`${API_BASE}/systems/stats`);
      if (!statsResponse.ok) {
        throw new Error('获取系统统计失败');
      }
      const statsData = await statsResponse.json();

      // 过滤高负载服务器（load_status为'high'或离线的服务器）
      const highLoadSystems = (statsData.systems || []).filter((system: SystemStats) => 
        system.load_status === 'high' || system.status !== 'up'
      );

      setSystems(highLoadSystems);

      // 获取高负载节点数据
      const nodesResponse = await fetch(`${API_BASE}/nodes/load-status`);
      if (nodesResponse.ok) {
        const nodesData = await nodesResponse.json();
        setHighLoadNodes(nodesData);
      } else {
        setHighLoadNodes([]);
      }

      // 计算高负载摘要
      const totalHighLoad = highLoadSystems.length;
      const offlineCount = highLoadSystems.filter(s => s.status !== 'up').length;
      const highCpuMemCount = highLoadSystems.filter(s => s.status === 'up' && s.load_status === 'high').length;

      setSummary({
        total: totalHighLoad,
        online: highCpuMemCount,
        offline: offlineCount,
        unknown: 0,
        high_load: totalHighLoad,
      });

      setError(null);
      setLastUpdate(new Date());
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取数据失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();

    // 每10秒刷新一次数据
    const interval = setInterval(fetchData, 10000);
    return () => clearInterval(interval);
  }, []);

  const formatDateTime = (dateString: string) => {
    try {
      return new Date(dateString).toLocaleString('zh-CN');
    } catch {
      return '无效时间';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'up': return '在线';
      case 'down': return '离线';
      default: return '未知';
    }
  };

  const getStatusClass = (status: string) => {
    switch (status) {
      case 'up': return 'status-up';
      case 'down': return 'status-down';
      default: return 'status-unknown';
    }
  };

  const getLoadStatusText = (loadStatus: string, systemStatus: string) => {
    if (systemStatus !== 'up') return '离线';
    switch (loadStatus) {
      case 'high': return '高负载';
      case 'normal': return '正常';
      default: return '未知';
    }
  };

  const getLoadStatusClass = (loadStatus: string, systemStatus: string) => {
    if (systemStatus !== 'up') return 'load-offline';
    switch (loadStatus) {
      case 'high': return 'load-high';
      case 'normal': return 'load-normal';
      default: return 'load-unknown';
    }
  };

  const getLoadReasonText = (system: SystemStats) => {
    const reasons = [];
    
    if (system.status !== 'up') {
      reasons.push('服务器离线');
    } else {
      if (system.avg_cpu > 90) reasons.push(`CPU: ${system.avg_cpu.toFixed(1)}%`);
      if (system.avg_mem_pct > 90) reasons.push(`内存: ${system.avg_mem_pct.toFixed(1)}%`);
      // 可以根据在线人数阈值判断
      if (system.online_users > 300) reasons.push(`在线人数: ${system.online_users}`);
    }
    
    return reasons.length > 0 ? reasons.join(', ') : '其他原因';
  };

  return (
    <div className="dashboard-container high-load-dashboard">
      <div className="dashboard-header">
        <h1 className="dashboard-title high-load-title">🔥 高负载服务器监控</h1>
        <p className="dashboard-description">显示当前处于高负载状态或离线的服务器</p>
        <div className="last-update">
          最后更新: {lastUpdate.toLocaleString('zh-CN')}
          <button className="refresh-button" onClick={fetchData} disabled={loading}>
            {loading ? '更新中...' : '刷新'}
          </button>
        </div>
      </div>

      {error && (
        <div className="error-banner">
          错误: {error}
        </div>
      )}


      {/* 高负载服务器表格 */}
      <div className="dashboard-content" style={{ maxWidth: '90rem', margin: '0 auto' }}>
        {loading && systems.length === 0 ? (
          <div className="loading">正在加载高负载服务器数据...</div>
        ) : systems.length === 0 ? (
          <div className="no-data">
            🎉 暂无高负载服务器，所有服务器运行正常！
          </div>
        ) : (
          <div className="servers-table-container">
            <div className="table-header">
              <h2>🖥️ 高负载服务器详情 ({systems.length}个)</h2>
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
                    <th>在线人数</th>
                    <th>网络I/O (Mbps)</th>
                    <th>负载原因</th>
                    <th>最后更新</th>
                  </tr>
                </thead>
                <tbody>
                  {systems.map((system) => (
                    <tr key={system.id} className="high-load-row">
                      <td>
                        <div className="server-info">
                          <div className="server-name">{system.name}</div>
                          <div className="server-address">{system.host}:{system.port}</div>
                        </div>
                      </td>
                      <td>
                        <span className={`status-badge ${getStatusClass(system.status)}`}>
                          {getStatusText(system.status)}
                        </span>
                      </td>
                      <td>
                        <span className={getLoadStatusClass(system.load_status, system.status)}>
                          {getLoadStatusText(system.load_status, system.status)}
                        </span>
                      </td>
                      <td>
                        <div className="progress-container">
                          <div className="progress-bar">
                            <div 
                              className={`progress-fill ${system.avg_cpu > 90 ? 'progress-cpu-danger' : 'progress-cpu'}`}
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
                              className={`progress-fill ${system.avg_mem_pct > 90 ? 'progress-memory-danger' : 'progress-memory'}`}
                              style={{width: `${Math.min(system.avg_mem_pct, 100)}%`}}
                            ></div>
                          </div>
                          <span className="progress-text">{system.avg_mem_pct.toFixed(1)}%</span>
                        </div>
                      </td>
                      <td className="online-users-cell">
                        <span className={`online-users-count ${system.online_users > 300 ? 'high-users' : ''}`}>
                          {system.online_users}
                        </span>
                      </td>
                      <td>
                        <div className="network-metrics">
                          <div>↑ {(system.avg_net_sent * 8).toFixed(2)}</div>
                          <div>↓ {(system.avg_net_recv * 8).toFixed(2)}</div>
                        </div>
                      </td>
                      <td className="load-reason">
                        <span className="reason-text">
                          {getLoadReasonText(system)}
                        </span>
                      </td>
                      <td>
                        <small style={{ color: '#6b7280', fontSize: '0.875rem' }}>
                          {system.last_update ? formatDateTime(system.last_update) : '无数据'}
                        </small>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>

      {/* 高负载节点列表 */}
      {highLoadNodes.length > 0 && (
        <div className="dashboard-content" style={{ marginTop: '2rem', maxWidth: '90rem', margin: '2rem auto 0' }}>
          <div className="servers-table-container">
            <div className="table-header">
              <h2>🔗 高负载节点详情 ({highLoadNodes.length}个)</h2>
              <p style={{ margin: '0.5rem 0 0 0', color: '#6b7280', fontSize: '0.875rem' }}>
                对应高负载服务器的所有节点信息
              </p>
            </div>
            <div className="table-wrapper">
              <table className="servers-table">
                <thead>
                  <tr>
                    <th style={{ width: '70%' }}>节点名称</th>
                    <th style={{ width: '30%', textAlign: 'center' }}>在线人数</th>
                  </tr>
                </thead>
                <tbody>
                  {highLoadNodes.map((node, index) => (
                    <tr key={`${node.id}-${index}`} className="high-load-row">
                      <td>
                        <div className="node-name-display">
                          <strong>{node.name}</strong>
                          <div style={{ 
                            fontSize: '0.75rem', 
                            color: '#6b7280', 
                            marginTop: '0.25rem',
                            textTransform: 'uppercase'
                          }}>
                            {node.type} · ID: {node.id}
                          </div>
                        </div>
                      </td>
                      <td className="online-users-cell" style={{ textAlign: 'center' }}>
                        <span className="online-users-count high-users" style={{
                          fontSize: '1rem',
                          fontWeight: '600',
                          padding: '0.25rem 0.75rem',
                          borderRadius: '9999px',
                          backgroundColor: '#fee2e2',
                          color: '#dc2626'
                        }}>
                          {node.online}人
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default HighLoadNodes;