import React, { useState } from 'react';
import ServerMonitor from './ServerMonitor';
import NodeManager from './NodeManager';
import "./index.css";

export function App() {
  const [currentView, setCurrentView] = useState<'monitor' | 'nodes'>('monitor');

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
        </div>
      </nav>
      
      <div className="app-content">
        {currentView === 'monitor' && <ServerMonitor />}
        {currentView === 'nodes' && <NodeManager />}
      </div>
    </div>
  );
}

export default App;
