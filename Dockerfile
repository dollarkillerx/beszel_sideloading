# 多阶段构建 Dockerfile
FROM node:18-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app/frontend

# 复制前端文件
COPY frontend/package.json frontend/bun.lockb* ./
COPY frontend/ ./

# 安装bun并构建前端
RUN npm install -g bun
RUN bun install
RUN bun run build

# Go 后端构建阶段
FROM golang:1.21-alpine AS backend-builder

# 安装必要的工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app/backend

# 复制go mod文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源码
COPY backend/ ./

# 构建Go应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# 最终运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates sqlite

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件
COPY --from=backend-builder /app/backend/main .
COPY --from=frontend-builder /app/frontend/dist ./static

# 创建必要的目录
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release
ENV DB_PATH=/app/data/server_monitor.db

# 启动命令
CMD ["./main"]