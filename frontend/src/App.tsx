import React, { useState } from 'react';
import ServerMonitor from './ServerMonitor';
import NodeManager from './NodeManager';
import LoadStatusTest from './LoadStatusTest';
import "./index.css";

export function App() {
  const [currentView, setCurrentView] = useState<'monitor' | 'nodes' | 'test'>('monitor');

  return (
    <div className="app">
      <nav className="app-nav">
        <div className="nav-links">
          <button 
            className={`nav-link ${currentView === 'monitor' ? 'active' : ''}`}
            onClick={() => setCurrentView('monitor')}
          >
            服务器监控
          </button>
          <button 
            className={`nav-link ${currentView === 'nodes' ? 'active' : ''}`}
            onClick={() => setCurrentView('nodes')}
          >
            节点管理
          </button>
          <button 
            className={`nav-link ${currentView === 'test' ? 'active' : ''}`}
            onClick={() => setCurrentView('test')}
          >
            负载测试
          </button>
        </div>
      </nav>
      
      <div className="app-content">
        {currentView === 'monitor' && <ServerMonitor />}
        {currentView === 'nodes' && <NodeManager />}
        {currentView === 'test' && <LoadStatusTest />}
      </div>
    </div>
  );
}

export default App;
