import React, { useState, useEffect } from 'react';
import { API_BASE } from './utils/api';

interface SystemAlias {
  id: number;
  system_id: string;
  alias: string;
  created_at: string;
  updated_at: string;
}

interface AliasManagerProps {
  systemId: string;
  systemName?: string;
  onClose: () => void;
}

const AliasManager: React.FC<AliasManagerProps> = ({ 
  systemId, 
  systemName, 
  onClose 
}) => {
  const [alias, setAlias] = useState<SystemAlias | null>(null);
  const [newAlias, setNewAlias] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchSystemAlias();
  }, [systemId]);

  const fetchSystemAlias = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/systems/${systemId}/alias`);
      if (!response.ok) {
        throw new Error('获取别名失败');
      }
      const data = await response.json();
      setAlias(data.alias || null);
      setNewAlias(data.alias?.alias || '');
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取别名失败');
    } finally {
      setLoading(false);
    }
  };

  const saveAlias = async () => {
    if (!newAlias.trim()) {
      setError('别名不能为空');
      return;
    }

    try {
      setSaving(true);
      const response = await fetch(`${API_BASE}/systems/${systemId}/alias`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ alias: newAlias.trim() }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '设置别名失败');
      }

      setError(null);
      await fetchSystemAlias();
    } catch (err) {
      setError(err instanceof Error ? err.message : '设置别名失败');
    } finally {
      setSaving(false);
    }
  };

  const deleteAlias = async () => {
    try {
      setSaving(true);
      const response = await fetch(`${API_BASE}/systems/${systemId}/alias`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '删除别名失败');
      }

      setAlias(null);
      setNewAlias('');
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : '删除别名失败');
    } finally {
      setSaving(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !saving) {
      saveAlias();
    }
  };

  return (
    <div className="tag-modal">
      <div className="tag-modal-content">
        <div className="tag-header">
          <div>
            <h3>别名管理 - {systemName || systemId}</h3>
            <p className="node-info">服务器ID: {systemId}</p>
          </div>
          <button className="close-button" onClick={onClose}>✕</button>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {/* 别名设置 */}
        <div className="alias-section">
          <div className="alias-input-group">
            <label>服务器别名：</label>
            <input
              type="text"
              placeholder="请输入服务器别名"
              value={newAlias}
              onChange={(e) => setNewAlias(e.target.value)}
              onKeyPress={handleKeyPress}
              disabled={saving || loading}
              className="alias-input"
            />
          </div>
          
          <div className="alias-info">
            {loading ? (
              <span className="loading-text">加载中...</span>
            ) : alias ? (
              <span className="current-alias">当前别名: {alias.alias}</span>
            ) : (
              <span className="no-alias">暂未设置别名</span>
            )}
          </div>
        </div>

        <div className="tag-actions">
          <button 
            className="save-button" 
            onClick={saveAlias}
            disabled={saving || loading || !newAlias.trim()}
          >
            {saving ? '保存中...' : '保存别名'}
          </button>
          {alias && (
            <button 
              className="delete-button" 
              onClick={deleteAlias}
              disabled={saving || loading}
            >
              删除别名
            </button>
          )}
          <button className="close-modal-button" onClick={onClose}>
            关闭
          </button>
        </div>
      </div>
    </div>
  );
};

export default AliasManager;