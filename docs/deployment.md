# Courier LLM Gateway - 部署文档

## 前置要求

- Docker 和 Docker Compose
- Go 1.22+（本地开发）
- PostgreSQL 16+（本地开发）

## 本地开发

### 1. 启动数据库

```bash
docker-compose up -d postgres
```

### 2. 执行数据库迁移

```bash
psql "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable" -f migrations/000001_create_providers.up.sql
```

### 3. 运行服务

```bash
export DATABASE_URL="host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable"
go run cmd/server/main.go
```

### 4. 测试 API

```bash
# 创建 Provider
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "openai-main",
    "type": "openai",
    "base_url": "https://api.openai.com/v1",
    "timeout": 300,
    "api_key": "sk-xxx",
    "enabled": true
  }'

# 查询 Provider 列表
curl http://localhost:8080/api/v1/providers

# 重载 Provider
curl -X POST http://localhost:8080/api/v1/admin/providers/reload
```

## Docker 部署

### 构建镜像

```bash
docker-compose build
```

### 启动服务

```bash
docker-compose up -d
```

### 查看日志

```bash
docker-compose logs -f courier
```

### 停止服务

```bash
docker-compose down
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| DATABASE_URL | PostgreSQL 连接字符串 | - |
| PORT | HTTP 服务端口 | 8080 |
| ADMIN_API_KEY | 管理员 API Key（可选） | - |

## 数据库迁移

### 创建迁移文件

```bash
# 格式: VERSION_description.up.sql / VERSION_description.down.sql
# 例如: 000002_add_models_table.up.sql
```

### 执行迁移

```bash
# 手动执行
psql $DATABASE_URL -f migrations/000001_create_providers.up.sql

# 使用 migrate 工具（推荐）
migrate -path migrations -database "$DATABASE_URL" up
```

## 生产部署注意事项

1. **安全性**
   - 设置强密码的 DATABASE_URL
   - 配置 ADMIN_API_KEY 启用管理员认证
   - 使用 HTTPS（配置反向代理如 Nginx）

2. **性能**
   - 配置数据库连接池
   - 启用日志聚合（如 ELK）
   - 监控 Provider 调用延迟

3. **高可用**
   - 部署多实例 + 负载均衡
   - PostgreSQL 主从复制
   - 健康检查 `/health`

## 回滚

```bash
# 回滚代码
git checkout <previous-commit>

# 回滚数据库迁移
psql $DATABASE_URL -f migrations/000001_create_providers.down.sql
```
