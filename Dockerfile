# 多阶段构建 Dockerfile
FROM node:18-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app/frontend

# 复制前端文件
COPY frontend/package.json frontend/bun.lockb* ./
COPY frontend/ ./

# 安装bun并构建前端
# 使用 curl 直接安装 bun（支持多架构）
RUN apk add --no-cache curl bash && \
    curl -fsSL https://bun.sh/install | bash && \
    mv /root/.bun/bin/bun /usr/local/bin/

# 安装依赖并构建
RUN bun install
RUN bun run build

# Go 后端构建阶段
FROM golang:1.23.0 AS backend-builder

# 设置工作目录
WORKDIR /app/backend

# 复制go mod文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源码
COPY backend/ ./

# 构建Go应用 - 不再需要CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

# 最终运行阶段 - 使用alpine更小的镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件
COPY --from=backend-builder /app/backend/main .
COPY --from=frontend-builder /app/frontend/dist ./static

# 创建必要的目录
RUN mkdir -p /app/data

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./main"]