# 使用多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache gcc musl-dev

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/main.go

# 使用轻量级基础镜像
FROM alpine:latest

# 设置时区
RUN apk add --no-cache tzdata wget && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 创建非root用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# 从builder阶段复制编译好的应用
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
COPY --from=builder /app/config ./config

# 设置适当的权限
RUN chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --spider -q http://localhost:8080/health || exit 1

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./main"] 