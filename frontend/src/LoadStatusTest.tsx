import React, { useState, useEffect } from 'react';
import { API_BASE } from './utils/api';

interface HighLoadNode {
  name: string;
  type: string;
  id: number;
  online: number;
}

const LoadStatusTest: React.FC = () => {
  const [highLoadNodes, setHighLoadNodes] = useState<HighLoadNode[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date());

  const fetchHighLoadNodes = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch(`${API_BASE}/nodes/load-status`);

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '获取高负载节点失败');
      }

      const data = await response.json();
      // 如果是服务不可用的情况，使用返回的data字段
      if (Array.isArray(data)) {
        setHighLoadNodes(data);
      } else if (data.data && Array.isArray(data.data)) {
        setHighLoadNodes(data.data);
      } else {
        setHighLoadNodes([]);
      }
      setLastUpdate(new Date());
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取高负载节点失败');
      setHighLoadNodes([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHighLoadNodes();
  }, []);

  const getTypeColor = (type: string) => {
    switch (type.toLowerCase()) {
      case 'vmess':
        return '#2196F3'; // 蓝色
      case 'vless':
        return '#4CAF50'; // 绿色
      case 'trojan':
        return '#FF9800'; // 橙色
      case 'shadowsocks':
      case 'ss':
        return '#9C27B0'; // 紫色
      case 'hysteria2':
        return '#E91E63'; // 粉红色
      default:
        return '#757575'; // 灰色
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ marginBottom: '30px' }}>
        <h1 style={{ margin: '0 0 10px 0', color: '#333' }}>🔥 高负载节点API测试</h1>
        <p style={{ margin: '0', color: '#666', fontSize: '16px' }}>
          测试 GET /api/nodes/load-status API - 获取所有高负载节点信息
        </p>
      </div>

      {error && (
        <div style={{ 
          background: '#fee2e2', 
          color: '#dc2626', 
          padding: '15px', 
          borderRadius: '8px', 
          marginBottom: '20px',
          border: '1px solid #fecaca'
        }}>
          ❌ {error}
        </div>
      )}

      {/* 控制面板 */}
      <div style={{ 
        background: '#f8fafc', 
        padding: '20px', 
        borderRadius: '8px', 
        border: '1px solid #e5e7eb',
        marginBottom: '30px' 
      }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h3 style={{ margin: '0 0 5px 0', color: '#374151' }}>API测试控制面板</h3>
            <p style={{ margin: '0', color: '#6b7280', fontSize: '14px' }}>
              最后更新: {lastUpdate.toLocaleString('zh-CN')}
            </p>
          </div>
          <button
            onClick={fetchHighLoadNodes}
            disabled={loading}
            style={{
              padding: '12px 24px',
              background: loading ? '#d1d5db' : '#3b82f6',
              color: 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: loading ? 'not-allowed' : 'pointer',
              fontSize: '14px',
              fontWeight: '500',
              display: 'flex',
              alignItems: 'center',
              gap: '8px'
            }}
          >
            {loading ? '🔄 查询中...' : '🔄 刷新数据'}
          </button>
        </div>
      </div>

      {/* 查询结果 */}
      <div style={{ marginBottom: '30px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '15px' }}>
          <h2 style={{ margin: '0', color: '#374151' }}>📊 高负载节点列表</h2>
          <div style={{
            padding: '6px 12px',
            background: highLoadNodes.length > 0 ? '#fee2e2' : '#f0fdf4',
            color: highLoadNodes.length > 0 ? '#dc2626' : '#059669',
            borderRadius: '20px',
            fontSize: '14px',
            fontWeight: '600'
          }}>
            {highLoadNodes.length}个节点
          </div>
        </div>

        {loading && highLoadNodes.length === 0 ? (
          <div style={{
            padding: '60px 20px',
            textAlign: 'center',
            background: '#f9fafb',
            borderRadius: '8px',
            color: '#6b7280'
          }}>
            🔄 正在获取高负载节点数据...
          </div>
        ) : highLoadNodes.length === 0 ? (
          <div style={{
            padding: '60px 20px',
            textAlign: 'center',
            background: '#f0fdf4',
            borderRadius: '8px',
            color: '#059669',
            border: '2px dashed #bbf7d0'
          }}>
            🎉 当前暂无高负载节点，所有节点运行正常！
          </div>
        ) : (
          <div style={{ 
            background: '#fff', 
            borderRadius: '8px', 
            border: '1px solid #e5e7eb',
            overflow: 'hidden'
          }}>
            {highLoadNodes.map((node, index) => (
              <div
                key={`${node.id}-${index}`}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '20px',
                  borderBottom: index < highLoadNodes.length - 1 ? '1px solid #f3f4f6' : 'none',
                  background: '#fef2f2',
                  borderLeft: '4px solid #dc2626'
                }}
              >
                <div style={{ flex: 1 }}>
                  <div style={{ 
                    fontSize: '16px', 
                    fontWeight: '600', 
                    color: '#111827',
                    marginBottom: '8px'
                  }}>
                    {node.name}
                  </div>
                  <div style={{ display: 'flex', gap: '20px', alignItems: 'center' }}>
                    <span
                      style={{
                        padding: '4px 12px',
                        borderRadius: '20px',
                        fontSize: '12px',
                        fontWeight: '600',
                        textTransform: 'uppercase',
                        background: getTypeColor(node.type),
                        color: 'white'
                      }}
                    >
                      {node.type}
                    </span>
                    <span style={{ color: '#6b7280', fontSize: '14px' }}>
                      ID: {node.id}
                    </span>
                  </div>
                </div>
                <div style={{ textAlign: 'right' }}>
                  <div style={{
                    fontSize: '24px',
                    fontWeight: '700',
                    color: '#dc2626',
                    marginBottom: '4px'
                  }}>
                    {node.online}人
                  </div>
                  <div style={{ 
                    fontSize: '12px', 
                    color: '#6b7280',
                    textTransform: 'uppercase',
                    fontWeight: '500'
                  }}>
                    在线用户
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* API 信息 */}
      <div style={{ 
        background: '#f8fafc', 
        padding: '20px', 
        borderRadius: '8px', 
        border: '1px solid #e5e7eb' 
      }}>
        <h3 style={{ margin: '0 0 15px 0', color: '#374151' }}>🔧 API 技术信息</h3>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
          <div>
            <p style={{ margin: '0 0 8px 0', fontWeight: '600', color: '#374151' }}>请求信息:</p>
            <div style={{ background: '#fff', padding: '15px', borderRadius: '6px', border: '1px solid #d1d5db' }}>
              <code style={{ fontSize: '14px', color: '#059669' }}>GET /api/nodes/load-status</code>
            </div>
          </div>
          <div>
            <p style={{ margin: '0 0 8px 0', fontWeight: '600', color: '#374151' }}>响应状态:</p>
            <div style={{ background: '#fff', padding: '15px', borderRadius: '6px', border: '1px solid #d1d5db' }}>
              <code style={{ fontSize: '14px', color: '#dc2626' }}>
                {error ? 'Error' : highLoadNodes.length > 0 ? `${highLoadNodes.length} nodes` : 'No data'}
              </code>
            </div>
          </div>
        </div>
        
        {highLoadNodes.length > 0 && (
          <div style={{ marginTop: '20px' }}>
            <p style={{ margin: '0 0 10px 0', fontWeight: '600', color: '#374151' }}>响应数据示例:</p>
            <pre style={{ 
              background: '#fff', 
              padding: '15px', 
              borderRadius: '6px', 
              overflow: 'auto',
              border: '1px solid #d1d5db',
              fontSize: '12px',
              color: '#374151',
              maxHeight: '200px'
            }}>
{JSON.stringify(highLoadNodes.slice(0, 2), null, 2)}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
};

export default LoadStatusTest;