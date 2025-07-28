import React, { useState, useEffect } from 'react';
import { API_BASE } from './utils/api';

interface NodeTag {
  id: number;
  system_id: string;
  tag_type: string;
  tag_id: number;
  created_at: string;
  updated_at: string;
}

interface NodeTagRequest {
  type: string;
  id: number;
}

interface NodeTagManagerProps {
  systemId: string;
  systemName?: string;
  onClose: () => void;
}

const NodeTagManager: React.FC<NodeTagManagerProps> = ({ 
  systemId, 
  systemName, 
  onClose 
}) => {
  const [tags, setTags] = useState<NodeTag[]>([]);
  const [newTagType, setNewTagType] = useState('');
  const [newTagId, setNewTagId] = useState('');
  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // API_BASE is now imported from utils/api.ts

  useEffect(() => {
    fetchSystemTags();
  }, [systemId]);

  const fetchSystemTags = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE}/systems/${systemId}/tags`);
      if (!response.ok) {
        throw new Error('获取标签失败');
      }
      const data = await response.json();
      setTags(data.tags || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取标签失败');
    } finally {
      setLoading(false);
    }
  };

  const addTag = async () => {
    if (!newTagType.trim()) {
      setError('标签类型不能为空');
      return;
    }

    if (!newTagId.trim()) {
      setError('标签ID不能为空');
      return;
    }

    const tagId = parseInt(newTagId);
    if (isNaN(tagId)) {
      setError('标签ID必须是数字');
      return;
    }

    try {
      setAdding(true);
      const request: NodeTagRequest = {
        type: newTagType.trim(),
        id: tagId,
      };

      const response = await fetch(`${API_BASE}/systems/${systemId}/tags`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '添加标签失败');
      }

      setNewTagType('');
      setNewTagId('');
      setError(null);
      await fetchSystemTags();
    } catch (err) {
      setError(err instanceof Error ? err.message : '添加标签失败');
    } finally {
      setAdding(false);
    }
  };

  const removeTag = async (tag: NodeTag) => {
    try {
      const request: NodeTagRequest = {
        type: tag.tag_type,
        id: tag.tag_id,
      };

      const response = await fetch(`${API_BASE}/systems/${systemId}/tags`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '删除标签失败');
      }

      setError(null);
      await fetchSystemTags();
    } catch (err) {
      setError(err instanceof Error ? err.message : '删除标签失败');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !adding) {
      addTag();
    }
  };

  return (
    <div className="tag-modal">
      <div className="tag-modal-content">
        <div className="tag-header">
          <div>
            <h3>标签管理 - {systemName || systemId}</h3>
            <p className="node-info">服务器ID: {systemId}</p>
          </div>
          <button className="close-button" onClick={onClose}>✕</button>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {/* 添加标签 */}
        <div className="add-tag-section">
          <div className="add-tag-input-group">
            <input
              type="text"
              placeholder="Type: (如: ss, v2ray, trojan...)"
              value={newTagType}
              onChange={(e) => setNewTagType(e.target.value)}
              onKeyPress={handleKeyPress}
              disabled={adding}
            />
            <input
              type="number"
              placeholder="标签ID (数字)"
              value={newTagId}
              onChange={(e) => setNewTagId(e.target.value)}
              onKeyPress={handleKeyPress}
              disabled={adding}
            />
            <button 
              onClick={addTag} 
              disabled={adding || !newTagType.trim() || !newTagId.trim()}
              className="add-tag-button"
            >
              {adding ? '添加中...' : '添加'}
            </button>
          </div>
        </div>

        {/* 当前标签列表 */}
        <div className="current-tags-section">
          <h4>当前标签 ({tags.length})</h4>
          {loading ? (
            <div className="loading-tags">加载中...</div>
          ) : tags.length === 0 ? (
            <div className="no-tags">暂无标签</div>
          ) : (
            <div className="tags-list">
              {tags.map((tag) => (
                <div key={tag.id} className="tag-item">
                  <span className="tag-content">{tag.tag_type}:{tag.tag_id}</span>
                  <button
                    className="remove-tag-button"
                    onClick={() => removeTag(tag)}
                    title="删除标签"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="tag-actions">
          <button className="close-modal-button" onClick={onClose}>
            关闭
          </button>
        </div>
      </div>
    </div>
  );
};

export default NodeTagManager;