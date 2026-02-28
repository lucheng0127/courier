# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安装必要工具
RUN apk add --no-cache git make

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .

# 复制迁移文件
COPY --from=builder /app/migrations ./migrations

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./server"]
