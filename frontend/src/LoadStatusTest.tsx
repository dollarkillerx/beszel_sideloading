import React, { useState } from 'react';

interface NodeLoadRequest {
  type: string;
  id: number;
}

interface NodeLoadResponse {
  type: string;
  id: number;
  load_status: string;
}

const LoadStatusTest: React.FC = () => {
  const [requests, setRequests] = useState<NodeLoadRequest[]>([
    { type: 'proxy', id: 1 },
    { type: 'cache', id: 2 }
  ]);
  const [responses, setResponses] = useState<NodeLoadResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [newType, setNewType] = useState('');
  const [newId, setNewId] = useState('');

  const API_BASE = 'http://localhost:8080/api';

  const addRequest = () => {
    if (!newType.trim() || !newId.trim()) {
      setError('类型和ID都不能为空');
      return;
    }

    const id = parseInt(newId);
    if (isNaN(id)) {
      setError('ID必须是数字');
      return;
    }

    const newRequest: NodeLoadRequest = {
      type: newType.trim(),
      id: id
    };

    setRequests([...requests, newRequest]);
    setNewType('');
    setNewId('');
    setError(null);
  };

  const removeRequest = (index: number) => {
    setRequests(requests.filter((_, i) => i !== index));
  };

  const testLoadStatus = async () => {
    if (requests.length === 0) {
      setError('请至少添加一个查询项');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await fetch(`${API_BASE}/nodes/load-status`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requests),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '查询失败');
      }

      const data: NodeLoadResponse[] = await response.json();
      setResponses(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : '查询失败');
      setResponses([]);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'normal':
        return '#4CAF50'; // 绿色
      case 'high':
        return '#f44336'; // 红色
      case 'not_found':
        return '#9E9E9E'; // 灰色
      default:
        return '#FF9800'; // 橙色
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'normal':
        return '正常';
      case 'high':
        return '高负载';
      case 'not_found':
        return '未找到';
      default:
        return '未知';
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
      <h1>节点负载状态测试</h1>
      <p>测试批量查询节点负载状态API</p>

      {error && (
        <div style={{ 
          background: '#ffebee', 
          color: '#c62828', 
          padding: '10px', 
          borderRadius: '4px', 
          marginBottom: '20px' 
        }}>
          {error}
        </div>
      )}

      {/* 添加查询项 */}
      <div style={{ marginBottom: '20px', padding: '15px', border: '1px solid #ddd', borderRadius: '4px' }}>
        <h3>添加查询项</h3>
        <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
          <input
            type="text"
            placeholder="类型 (如: proxy, cache, server)"
            value={newType}
            onChange={(e) => setNewType(e.target.value)}
            style={{ padding: '8px', border: '1px solid #ccc', borderRadius: '4px', flex: 1 }}
          />
          <input
            type="number"
            placeholder="ID"
            value={newId}
            onChange={(e) => setNewId(e.target.value)}
            style={{ padding: '8px', border: '1px solid #ccc', borderRadius: '4px', width: '100px' }}
          />
          <button
            onClick={addRequest}
            style={{
              padding: '8px 16px',
              background: '#2196F3',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer'
            }}
          >
            添加
          </button>
        </div>
      </div>

      {/* 当前查询项列表 */}
      <div style={{ marginBottom: '20px' }}>
        <h3>查询列表 ({requests.length})</h3>
        {requests.length === 0 ? (
          <p style={{ color: '#666' }}>暂无查询项</p>
        ) : (
          <div>
            {requests.map((req, index) => (
              <div
                key={index}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '10px',
                  background: '#f5f5f5',
                  marginBottom: '5px',
                  borderRadius: '4px'
                }}
              >
                <span>类型: {req.type}, ID: {req.id}</span>
                <button
                  onClick={() => removeRequest(index)}
                  style={{
                    padding: '4px 8px',
                    background: '#f44336',
                    color: 'white',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: 'pointer'
                  }}
                >
                  删除
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* 测试按钮 */}
      <div style={{ textAlign: 'center', marginBottom: '20px' }}>
        <button
          onClick={testLoadStatus}
          disabled={loading || requests.length === 0}
          style={{
            padding: '12px 24px',
            background: loading ? '#ccc' : '#4CAF50',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: loading ? 'not-allowed' : 'pointer',
            fontSize: '16px'
          }}
        >
          {loading ? '查询中...' : '测试负载状态查询'}
        </button>
      </div>

      {/* 查询结果 */}
      {responses.length > 0 && (
        <div>
          <h3>查询结果</h3>
          <div style={{ border: '1px solid #ddd', borderRadius: '4px' }}>
            {responses.map((resp, index) => (
              <div
                key={index}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '15px',
                  borderBottom: index < responses.length - 1 ? '1px solid #eee' : 'none'
                }}
              >
                <div>
                  <strong>类型:</strong> {resp.type} | <strong>ID:</strong> {resp.id}
                </div>
                <div
                  style={{
                    padding: '6px 12px',
                    borderRadius: '20px',
                    color: 'white',
                    background: getStatusColor(resp.load_status),
                    fontWeight: 'bold'
                  }}
                >
                  {getStatusText(resp.load_status)}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* API 信息 */}
      <div style={{ marginTop: '30px', padding: '15px', background: '#f9f9f9', borderRadius: '4px' }}>
        <h4>API 信息</h4>
        <p><strong>端点:</strong> POST /api/nodes/load-status</p>
        <p><strong>请求格式:</strong></p>
        <pre style={{ background: '#fff', padding: '10px', borderRadius: '4px', overflow: 'auto' }}>
{JSON.stringify(requests, null, 2)}
        </pre>
        {responses.length > 0 && (
          <>
            <p><strong>响应格式:</strong></p>
            <pre style={{ background: '#fff', padding: '10px', borderRadius: '4px', overflow: 'auto' }}>
{JSON.stringify(responses, null, 2)}
            </pre>
          </>
        )}
      </div>
    </div>
  );
};

export default LoadStatusTest;