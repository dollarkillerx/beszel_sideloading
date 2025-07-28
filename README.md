# Beszel 监控系统

一个集成了前端和后端的服务器监控系统，支持服务器负载监控、标签管理和负载状态查询。

##  功能特性

- **服务器监控**: 实时监控服务器CPU、内存、网络状态
- **标签管理**: 为服务器添加自定义标签进行分类管理
- **负载状态**: 基于阈值配置的智能负载状态判断
- **批量查询**: 支持批量查询节点负载状态API
- **实时更新**: 基于PocketBase的实时数据更新
- **Docker部署**: 一键Docker部署，前后端集成

## 快速开始

### 使用 Docker（推荐）

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd beszel_sideloading
   ```

2. **配置环境变量**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件，配置 PocketBase 连接信息
   ```
   
   配置示例：
   ```env
   POCKETBASE_BASE_URL=https://bz.baidua.top
   POCKETBASE_EMAIL=your-email@example.com
   POCKETBASE_PASSWORD=your-password
   ```

3. **使用 Docker Compose 启动**
   ```bash
   docker-compose up -d
   ```

4. **访问应用**
   - 打开浏览器访问: http://localhost:8080
   - 包含三个主要页面：
     - **服务器监控** - 查看服务器状态和负载
     - **节点管理** - 管理服务器标签
     - **负载测试** - 测试批量负载查询API

### 手动构建和部署

1. **前置要求**
   - Go 1.21+
   - Node.js 18+ 或 Bun
   - SQLite

2. **使用构建脚本**
   ```bash
   chmod +x scripts/build.sh
   ./scripts/build.sh
   ```

3. **或者手动构建**
   ```bash
   # 构建前端
   cd frontend
   bun install && bun run build  # 或使用 npm install && npm run build
   cd ..
   
   # 构建后端
   cd backend
   go build -o ../dist/beszel-monitor cmd/main.go
   cd ..
   
   # 复制前端文件
   mkdir -p dist/static
   cp -r frontend/dist/* dist/static/
   ```

4. **运行应用**
   ```bash
   cd dist
   ./beszel-monitor
   ```

## API 文档

### 节点负载状态查询 API

这是系统的核心API，支持批量查询节点负载状态。

**端点**: `POST /api/nodes/load-status`

**请求示例**:
```json
[
  {"type": "ss", "id": 1},
  {"type": "v2ray", "id": 2},
  {"type": "trojan", "id": 3}
]
```

**响应示例**:
```json
[
  {"type": "ss", "id": 1, "load_status": "normal"},
  {"type": "v2ray", "id": 2, "load_status": "high"},
  {"type": "trojan", "id": 3, "load_status": "not_found"}
]
```

**负载状态说明**:
- `normal`: 负载正常
- `high`: 负载过高（包括离线服务器）
- `not_found`: 未找到对应标签
- `no_data`: 无统计数据

### 服务器管理 API

- `GET /api/systems` - 获取所有服务器列表
- `GET /api/systems/summary` - 获取服务器摘要
- `GET /api/systems/stats` - 获取服务器统计数据
- `GET /api/systems/:id/stats` - 获取特定服务器统计

### 标签管理 API

- `GET /api/systems/:id/tags` - 获取服务器标签
- `POST /api/systems/:id/tags` - 添加服务器标签
- `DELETE /api/systems/:id/tags` - 删除服务器标签

**标签操作示例**:
```bash
# 添加标签
curl -X POST "http://localhost:8080/api/systems/server-id/tags" \
  -H "Content-Type: application/json" \
  -d '{"type": "ss", "id": 1}'

# 删除标签
curl -X DELETE "http://localhost:8080/api/systems/server-id/tags" \
  -H "Content-Type: application/json" \
  -d '{"type": "ss", "id": 1}'
```

### 阈值配置 API

- `GET /api/systems/:id/threshold` - 获取服务器阈值配置
- `PUT /api/systems/:id/threshold` - 更新服务器阈值配置
- `DELETE /api/systems/:id/threshold` - 删除服务器阈值配置
- `GET /api/thresholds` - 获取所有阈值配置

## ⚙️ 配置

### 环境变量

| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `POCKETBASE_BASE_URL` | PocketBase 服务地址 | - | ✅ |
| `POCKETBASE_EMAIL` | PocketBase 登录邮箱 | - | ✅ |
| `POCKETBASE_PASSWORD` | PocketBase 登录密码 | - | ✅ |
| `DB_PATH` | SQLite 数据库路径 | `./server_monitor.db` | ❌ |
| `GIN_MODE` | Gin 运行模式 | `debug` | ❌ |
| `PORT` | 服务端口 | `8080` | ❌ |

### 阈值配置

系统支持为每台服务器设置独立的负载阈值：

- **CPU阈值**: CPU使用率告警百分比（默认90%）
- **内存阈值**: 内存使用率告警百分比（默认90%）
- **网络上行**: 上行带宽最大值和告警百分比
- **网络下行**: 下行带宽最大值和告警百分比

### Docker 卷挂载

- `/app/data` - 数据库文件存储
- `/app/config` - 配置文件存储（可选）

## 系统架构

### 技术栈

**前端**:
- React 18 + TypeScript
- Bun 运行时
- CSS Modules

**后端**:
- Go 1.21+
- Gin Web框架

**部署**:
- Docker + Docker Compose
- 多阶段构建
- Alpine Linux

### 目录结构

```
beszel_sideloading/
├── frontend/              # React + TypeScript 前端
│   ├── src/
│   │   ├── App.tsx       # 主应用组件
│   │   ├── ServerMonitor.tsx    # 服务器监控页面
│   │   ├── NodeManager.tsx      # 节点管理页面
│   │   ├── NodeTagManager.tsx   # 标签管理组件
│   │   ├── LoadStatusTest.tsx   # 负载测试页面
│   │   └── index.css     # 样式文件
│   └── package.json
├── backend/               # Go 后端
│   ├── cmd/
│   │   └── main.go       # 应用入口
│   ├── internal/
│   │   ├── api/          # API路由和处理器
│   │   ├── service/      # 业务逻辑层
│   │   ├── database/     # 数据库配置
│   │   └── config/       # 配置管理
│   ├── pkg/
│   │   └── models/       # 数据模型
│   └── go.mod
├── scripts/
│   └── build.sh          # 构建脚本
├── Dockerfile             # Docker 构建文件
├── docker-compose.yml     # Docker Compose 配置
├── .env.example          # 环境变量示例
└── README.md             # 项目文档
```

## 🔧 开发指南

### 开发环境设置

1. **启动后端开发服务器**
   ```bash
   cd backend
   go run cmd/main.go
   ```

2. **启动前端开发服务器**
   ```bash
   cd frontend
   bun run dev  # 或 npm run dev
   ```

3. **访问开发环境**
   - 前端: http://localhost:3000
   - 后端API: http://localhost:8080/api

