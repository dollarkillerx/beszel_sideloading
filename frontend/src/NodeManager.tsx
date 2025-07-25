import React, { useState, useEffect } from 'react';
import NodeTagManager from './NodeTagManager';

interface System {
  id: string;
  name: string;
  host: string;
  port: string;
  status: string;
  created_at: string;
  updated_at: string;
}

const NodeManager: React.FC = () => {
  const [showTagManager, setShowTagManager] = useState(false);
  const [selectedSystem, setSelectedSystem] = useState<System | null>(null);
  const [systems, setSystems] = useState<System[]>([]);
  const [loading, setLoading] = useState(true);

  const API_BASE = 'http://localhost:8080/api';

  useEffect(() => {
    fetchSystems();
  }, []);

  const fetchSystems = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/systems`);
      if (!response.ok) {
        throw new Error('获取服务器列表失败');
      }
      const data = await response.json();
      setSystems(data.systems || []);
    } catch (err) {
      console.error('获取服务器列表失败:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleManageTags = (system: System) => {
    setSelectedSystem(system);
    setShowTagManager(true);
  };

  const handleCloseTagManager = () => {
    setShowTagManager(false);
    setSelectedSystem(null);
  };

  return (
    <div className="node-manager-container">
      <div className="node-manager-header">
        <h2>服务器管理</h2>
        <p>管理服务器标签，便于分类和查找</p>
      </div>

      {loading ? (
        <div className="loading">加载服务器列表中...</div>
      ) : (
        <div className="nodes-grid">
          {systems.map((system) => (
            <div key={system.id} className="node-card">
              <div className="node-card-header">
                <h3>{system.name}</h3>
                <span className={`status-badge status-${system.status}`}>
                  {system.status === 'up' ? '在线' : '离线'}
                </span>
              </div>
              <div className="node-card-info">
                <p>服务器ID: {system.id}</p>
                <p>地址: {system.host}:{system.port}</p>
              </div>
              <div className="node-card-actions">
                <button
                  className="manage-tags-button"
                  onClick={() => handleManageTags(system)}
                >
                  管理标签
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* 标签管理模态框 */}
      {showTagManager && selectedSystem && (
        <NodeTagManager
          systemId={selectedSystem.id}
          systemName={selectedSystem.name}
          onClose={handleCloseTagManager}
        />
      )}
    </div>
  );
};

export default NodeManager;