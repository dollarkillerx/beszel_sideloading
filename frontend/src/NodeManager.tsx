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

interface NodeTag {
  id: number;
  system_id: string;
  tag_type: string;
  tag_id: number;
  created_at: string;
  updated_at: string;
}

const NodeManager: React.FC = () => {
  const [showTagManager, setShowTagManager] = useState(false);
  const [selectedSystem, setSelectedSystem] = useState<System | null>(null);
  const [systems, setSystems] = useState<System[]>([]);
  const [systemTags, setSystemTags] = useState<Record<string, NodeTag[]>>({});
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
      const systemsList = data.systems || [];
      setSystems(systemsList);
      
      // 获取每个系统的标签
      await fetchAllSystemTags(systemsList);
    } catch (err) {
      console.error('获取服务器列表失败:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchAllSystemTags = async (systemsList: System[]) => {
    const tagsMap: Record<string, NodeTag[]> = {};
    
    // 并行获取所有系统的标签
    await Promise.all(
      systemsList.map(async (system) => {
        try {
          const response = await fetch(`${API_BASE}/systems/${system.id}/tags`);
          if (response.ok) {
            const data = await response.json();
            tagsMap[system.id] = data.tags || [];
          } else {
            tagsMap[system.id] = [];
          }
        } catch (err) {
          console.error(`获取系统 ${system.id} 标签失败:`, err);
          tagsMap[system.id] = [];
        }
      })
    );
    
    setSystemTags(tagsMap);
  };

  const handleManageTags = (system: System) => {
    setSelectedSystem(system);
    setShowTagManager(true);
  };

  const handleCloseTagManager = () => {
    setShowTagManager(false);
    setSelectedSystem(null);
    // 重新获取标签数据以刷新显示
    fetchSystems();
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
              
              {/* 标签显示 */}
              <div className="node-tags-section">
                <div className="tags-header">
                  <span>标签 ({(systemTags[system.id] || []).length})</span>
                </div>
                <div className="tags-display">
                  {(systemTags[system.id] || []).length === 0 ? (
                    <span className="no-tags-text">暂无标签</span>
                  ) : (
                    <div className="tags-list-inline">
                      {(systemTags[system.id] || []).map((tag) => (
                        <span key={tag.id} className="tag-badge">
                          {tag.tag_type}:{tag.tag_id}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
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