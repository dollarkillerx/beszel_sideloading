#!/bin/bash

# 构建脚本
set -e

echo "🏗️  开始构建 Beszel 监控系统..."

# 构建前端
echo "📦 构建前端..."
cd frontend
if command -v bun &> /dev/null; then
    echo "使用 Bun 构建前端..."
    bun install
    bun run build
else
    echo "使用 npm 构建前端..."
    npm install
    npm run build
fi
cd ..

# 构建后端
echo "🔧 构建后端..."
cd backend
echo "下载 Go 依赖..."
go mod download
echo "构建 Go 应用..."
CGO_ENABLED=1 go build -o ../dist/beszel-monitor cmd/main.go
cd ..

# 复制前端构建文件到后端静态目录
echo "📁 复制前端文件..."
mkdir -p dist/static
cp -r frontend/dist/* dist/static/

echo "✅ 构建完成！"
echo "📍 可执行文件位置: ./dist/beszel-monitor"
echo "📍 静态文件位置: ./dist/static/"
echo ""
echo "🚀 运行命令:"
echo "   cd dist && ./beszel-monitor"