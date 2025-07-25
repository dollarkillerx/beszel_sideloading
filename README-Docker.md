# Beszel 监控系统 Docker 部署指南

## 📦 Docker 一键部署

这个项目已经配置好了完整的 Docker 部署方案，将前端和后端打包在一个容器中运行。

## 🚀 快速启动

### 1. 克隆项目
```bash
git clone <your-repo-url>
cd beszel_sideloading
```

### 2. 配置环境变量
```bash
cp .env.example .env
vim .env  # 编辑配置文件
```

在 `.env` 文件中配置：
```env
POCKETBASE_BASE_URL=https://bz.baidua.top
POCKETBASE_EMAIL=your-email@example.com
POCKETBASE_PASSWORD=your-password
```

### 3. 使用 Docker Compose 启动
```bash
docker-compose up -d
```

### 4. 访问应用
打开浏览器访问: http://localhost:8080

## 🛠️ 手动 Docker 构建

### 构建镜像
```bash
docker build -t beszel-monitor .
```

### 运行容器
```bash
docker run -d \
  --name beszel-monitor \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e POCKETBASE_BASE_URL=https://bz.baidua.top \
  -e POCKETBASE_EMAIL=your-email@example.com \
  -e POCKETBASE_PASSWORD=your-password \
  beszel-monitor
```

## 📁 文件结构

构建后的 Docker 镜像包含：

```
/app/
├── main                 # Go 后端可执行文件
├── static/              # 前端构建文件
│   ├── index.html
│   ├── assets/
│   └── favicon.ico
└── data/                # 数据库存储目录
    └── server_monitor.db
```

## 🔧 Docker 配置说明

### Dockerfile 特性
- **多阶段构建**: 分别构建前端和后端，最终合并到一个轻量镜像
- **前端构建**: 使用 Bun 构建 React 应用
- **后端构建**: 使用 Go 构建可执行文件
- **运行环境**: 基于 Alpine Linux，包含 SQLite 支持
- **安全性**: 使用非 root 用户运行

### 端口映射
- **容器端口**: 8080
- **主机端口**: 8080 (可在 docker-compose.yml 中修改)

### 数据持久化
- **数据库**: `/app/data/server_monitor.db`
- **配置**: `/app/config/` (可选)

### 环境变量
| 变量 | 描述 | 默认值 |
|------|------|--------|
| `GIN_MODE` | Gin 框架运行模式 | `release` |
| `DB_PATH` | SQLite 数据库路径 | `/app/data/server_monitor.db` |
| `POCKETBASE_BASE_URL` | PocketBase API 地址 | - |
| `POCKETBASE_EMAIL` | PocketBase 登录邮箱 | - |
| `POCKETBASE_PASSWORD` | PocketBase 登录密码 | - |

## 🔍 故障排除

### 查看日志
```bash
# 查看容器日志
docker logs beszel-monitor

# 实时查看日志
docker logs -f beszel-monitor

# Docker Compose 日志
docker-compose logs -f
```

### 进入容器调试
```bash
docker exec -it beszel-monitor sh
```

### 重启服务
```bash
docker-compose restart
```

### 完全重建
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 📊 监控和维护

### 容器状态检查
```bash
docker ps
docker stats beszel-monitor
```

### 数据备份
```bash
# 备份数据库
docker cp beszel-monitor:/app/data/server_monitor.db ./backup/

# 恢复数据库
docker cp ./backup/server_monitor.db beszel-monitor:/app/data/
```

### 更新应用
```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up -d --build
```

## 🌐 生产环境建议

### 1. 使用 nginx 反向代理
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. 配置 HTTPS
```bash
# 使用 Let's Encrypt
certbot --nginx -d your-domain.com
```

### 3. 配置自动备份
```bash
# 添加到 crontab
0 2 * * * docker cp beszel-monitor:/app/data/server_monitor.db /backup/beszel_$(date +\%Y\%m\%d).db
```

### 4. 资源限制
```yaml
# docker-compose.yml
services:
  beszel-monitor:
    # ...
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

## 🆘 常见问题

**Q: 容器启动失败**
A: 检查环境变量配置，确保 PocketBase 连接信息正确

**Q: 前端页面显示空白**
A: 检查静态文件是否正确复制，查看容器日志

**Q: 数据库权限错误**
A: 确保数据目录有正确的写权限：`chmod 755 ./data`

**Q: API 请求失败**
A: 检查 CORS 配置，确保前后端在同一域名下运行