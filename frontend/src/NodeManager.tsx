import React, { useState, useEffect } from 'react';
import AliasManager from './AliasManager';
import { API_BASE } from './utils/api';

interface System {
  id: string;
  name: string;
  host: string;
  port: string;
  status: string;
  created_at: string;
  updated_at: string;
}

interface SystemAlias {
  id: number;
  system_id: string;
  alias: string;
  created_at: string;
  updated_at: string;
}

interface V2boardNode {
  name: string;
  id: number;
  type: string;
  online: number;
  last_update: number;
}

interface SystemNodeInfo {
  system_id: string;
  system_name: string;
  alias?: string;
  nodes: V2boardNode[];
  total_online: number;
}

const NodeManager: React.FC = () => {
  const [showAliasManager, setShowAliasManager] = useState(false);
  const [selectedSystem, setSelectedSystem] = useState<System | null>(null);
  const [systems, setSystems] = useState<System[]>([]);
  const [systemAliases, setSystemAliases] = useState<Record<string, SystemAlias | null>>({});
  const [systemNodes, setSystemNodes] = useState<Record<string, SystemNodeInfo>>({});
  const [loading, setLoading] = useState(true);

  // API_BASE is now imported from utils/api.ts

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
      
      // 获取每个系统的别名
      await fetchAllSystemAliases(systemsList);
      
      // 获取每个系统的节点信息
      await fetchAllSystemNodes(systemsList);
    } catch (err) {
      console.error('获取服务器列表失败:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchAllSystemAliases = async (systemsList: System[]) => {
    const aliasesMap: Record<string, SystemAlias | null> = {};
    
    // 并行获取所有系统的别名
    await Promise.all(
      systemsList.map(async (system) => {
        try {
          const response = await fetch(`${API_BASE}/systems/${system.id}/alias`);
          if (response.ok) {
            const data = await response.json();
            aliasesMap[system.id] = data.alias || null;
          } else {
            aliasesMap[system.id] = null;
          }
        } catch (err) {
          console.error(`获取系统 ${system.id} 别名失败:`, err);
          aliasesMap[system.id] = null;
        }
      })
    );
    
    setSystemAliases(aliasesMap);
  };

  const fetchAllSystemNodes = async (systemsList: System[]) => {
    const nodesMap: Record<string, SystemNodeInfo> = {};
    
    // 并行获取所有系统的节点信息
    await Promise.all(
      systemsList.map(async (system) => {
        try {
          const response = await fetch(`${API_BASE}/systems/${system.id}/nodes`);
          if (response.ok) {
            const nodeInfo = await response.json();
            nodesMap[system.id] = nodeInfo;
          } else {
            nodesMap[system.id] = {
              system_id: system.id,
              system_name: system.name,
              nodes: [],
              total_online: 0
            };
          }
        } catch (err) {
          console.error(`获取系统 ${system.id} 节点信息失败:`, err);
          nodesMap[system.id] = {
            system_id: system.id,
            system_name: system.name,
            nodes: [],
            total_online: 0
          };
        }
      })
    );
    
    setSystemNodes(nodesMap);
  };

  const handleManageAlias = (system: System) => {
    setSelectedSystem(system);
    setShowAliasManager(true);
  };

  const handleCloseAliasManager = () => {
    setShowAliasManager(false);
    setSelectedSystem(null);
    // 重新获取别名和节点数据以刷新显示
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
              
              {/* 别名显示 */}
              <div className="node-alias-section">
                <div className="alias-header">
                  <span>别名</span>
                </div>
                <div className="alias-display">
                  {systemAliases[system.id] ? (
                    <span className="alias-text">{systemAliases[system.id]?.alias}</span>
                  ) : (
                    <span className="no-alias-text">暂无别名</span>
                  )}
                </div>
              </div>
              
              {/* 节点信息显示 */}
              <div className="node-info-section">
                <div className="node-info-header">
                  <span>节点信息</span>
                  {systemNodes[system.id]?.nodes?.length > 0 && (
                    <span className="total-online">(总在线: {systemNodes[system.id].total_online})</span>
                  )}
                </div>
                <div className="node-info-display">
                  {systemNodes[system.id]?.nodes?.length > 0 ? (
                    <div className="nodes-list">
                      {systemNodes[system.id].nodes.map((node) => (
                        <div key={node.id} className="node-item">
                          <span className="node-name">{node.name}</span>
                          <span className="node-online">在线: {node.online}</span>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <span className="no-nodes-text">未找到匹配节点</span>
                  )}
                </div>
              </div>
              
              <div className="node-card-actions">
                <button
                  className="manage-alias-button"
                  onClick={() => handleManageAlias(system)}
                >
                  管理别名
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* 别名管理模态框 */}
      {showAliasManager && selectedSystem && (
        <AliasManager
          systemId={selectedSystem.id}
          systemName={selectedSystem.name}
          onClose={handleCloseAliasManager}
        />
      )}
    </div>
  );
};

export default NodeManager;