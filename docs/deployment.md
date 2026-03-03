# Courier LLM Gateway - 部署文档

## 前置要求

- Docker 和 Docker Compose
- Go 1.23+（本地开发）
- PostgreSQL 16+（本地开发）

## 本地开发

### 1. 启动数据库

```bash
docker-compose up -d postgres
```

### 2. 运行服务

系统启动时会自动执行数据库迁移（GORM AutoMigrate）。

```bash
export DATABASE_URL="host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable"
export JWT_SECRET="your-jwt-secret-key-change-in-production"
export INITIAL_ADMIN_EMAIL="admin@example.com"
export INITIAL_ADMIN_PASSWORD="admin-password-change-me"
go run cmd/server/main.go
```

### 环境变量

| 变量 | 说明 | 默认值 | 必需 |
|------|------|--------|------|
| DATABASE_URL | PostgreSQL 连接字符串 | - | ✓ |
| PORT | HTTP 服务端口 | 8080 | - |
| JWT_SECRET | JWT 签名密钥 | - | ✓ |
| JWT_ACCESS_TOKEN_EXPIRES_IN | Access Token 有效期 | 15m | - |
| JWT_REFRESH_TOKEN_EXPIRES_IN | Refresh Token 有效期 | 168h | - |
| JWT_ISSUER | Token 发行者标识 | courier-gateway | - |
| INITIAL_ADMIN_EMAIL | 初始管理员邮箱 | - | - |
| INITIAL_ADMIN_PASSWORD | 初始管理员密码 | - | - |
| LOG_LEVEL | 日志级别（debug/info/warn/error） | info | - |
| ENV | 运行环境（development/production） | production | - |
| AUTO_MIGRATE | 是否自动执行数据库迁移 | true | - |

### 日志配置

系统使用 uber-go/zap 结构化日志：

- **开发环境**（ENV=development）：彩色 console 格式，便于调试
- **生产环境**（ENV=production）：JSON 格式，便于日志聚合

设置日志级别：
```bash
export LOG_LEVEL=debug  # debug/info/warn/error
```

### 3. 测试 API

```bash
# 用户注册（无需认证）
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "password": "user-password-123"
  }')
USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.id')

# 用户登录获取 JWT Token
USER_ACCESS_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "zhangsan@example.com",
    "password": "user-password-123"
  }' | jq -r '.access_token')

# 为自己创建 API Key
API_KEY_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/users/$USER_ID/api-keys" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN" \
  -d '{
    "name": "生产环境 Key"
  }')
API_KEY=$(echo $API_KEY_RESPONSE | jq -r '.key')
# 返回的 key 仅此一次可见，请妥善保存

# 使用 API Key 调用 Chat API
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $API_KEY" \
  -d '{
    "model": "openai-main/gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# 查询自己的使用统计
curl "http://localhost:8080/api/v1/usage" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN"

# 管理员登录获取 JWT Token
ADMIN_ACCESS_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin-password-change-me"
  }' | jq -r '.access_token')

# 查询所有用户（仅管理员）
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN"
```

### 4. 管理 Provider

```bash
# 创建 OpenAI Provider（带 Fallback 配置）
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "openai-main",
    "type": "openai",
    "base_url": "https://api.openai.com/v1",
    "timeout": 60,
    "api_key": "sk-xxx",
    "enabled": true,
    "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
  }'

# 创建通义千问 Provider
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "qwen-main",
    "type": "openai",
    "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1",
    "timeout": 60,
    "api_key": "your-api-key",
    "enabled": true,
    "extra_config": {
      "temperature": 0.8,
      "max_tokens": 1500
    },
    "fallback_models": ["qwen-max", "qwen-plus", "qwen-turbo"]
  }'

# 查询 Provider 列表
curl http://localhost:8080/api/v1/providers \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 获取单个 Provider 信息
curl http://localhost:8080/api/v1/providers/openai-main \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 更新 Provider（只更新需要修改的字段）
curl -X PUT http://localhost:8080/api/v1/providers/openai-main \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "timeout": 120,
    "enabled": false
  }'

# 删除 Provider
curl -X DELETE http://localhost:8080/api/v1/providers/old-provider \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 重载 Provider
curl -X POST http://localhost:8080/api/v1/providers/reload \
  -H "Authorization: Bearer $ACCESS_TOKEN"
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

### 完全清理（包括数据库卷）

```bash
docker-compose down -v
```

## 数据库迁移

系统使用 GORM AutoMigrate 进行自动数据库迁移，无需手动执行 SQL 文件。

### 迁移机制

1. **启动时自动执行**：服务启动时自动检查并执行迁移
2. **Schema 版本跟踪**：使用 `schema_migrations` 表记录版本和 hash
3. **变更检测**：检测 struct 定义变化并自动同步
4. **环境变量控制**：可通过 `AUTO_MIGRATE=false` 禁用自动迁移

### 迁移日志

```json
{"level":"info","msg":"Starting database auto-migration..."}
{"level":"info","msg":"Database auto-migration completed successfully","version":"v1.0.0","hash":"..."}
```

### 禁用自动迁移

```bash
export AUTO_MIGRATE=false
```

### Model 定义

数据库表结构由 Go struct 定义，位于 `internal/model/` 目录：

- `provider.go` - Provider 表
- `user.go` - User 和 APIKey 表
- `usage.go` - UsageRecord 表

## 生产部署注意事项

### 安全性

1. **密钥管理**
   - 设置强随机密钥的 `JWT_SECRET`（至少 32 字符）
   - 设置强密码的 `DATABASE_URL`
   - **初始管理员**：通过 `INITIAL_ADMIN_EMAIL` 和 `INITIAL_ADMIN_PASSWORD` 环境变量创建初始管理员账户
     - 这是创建管理员用户的唯一方式
     - 用户注册只能创建普通用户（role="user"）
     - 如需创建更多管理员，需通过数据库直接操作用户角色

2. **HTTPS**
   - 生产环境必须使用 HTTPS
   - 配置反向代理（Nginx、Caddy）

3. **API Key 保护**
   - 不要在代码中硬编码 API Key
   - 使用环境变量或密钥管理服务

### 性能优化

1. **数据库**
   - 配置合适的连接池大小
   - 使用 PostgreSQL 连接池（PgBouncer）

2. **日志**
   - 设置 `LOG_LEVEL=info` 或 `warn` 减少日志量
   - 启用日志聚合（ELK、Loki）

3. **监控**
   - 监控 Provider 调用延迟
   - 监控 Fallback 频率
   - 设置告警规则

### 高可用

1. **多实例部署**
   - 部署多个服务实例
   - 使用负载均衡（Nginx、HAProxy）

2. **数据库**
   - PostgreSQL 主从复制
   - 连接池管理

3. **健康检查**

健康检查端点：
```bash
curl http://localhost:8080/health
```

响应：
```json
{"status":"ok"}
```

## API 接口说明

### 认证接口（无需鉴权）

| 方法 | 路径 | 说明 | 速率限制 |
|------|------|------|----------|
| POST | `/api/v1/auth/register` | 用户注册 | IP 级别，5 次/小时 |
| POST | `/api/v1/auth/login` | 用户登录获取 JWT Token | - |
| POST | `/api/v1/auth/refresh` | 刷新 JWT Token | - |

### 管理接口（需要 JWT Token，Admin 角色）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/users` | 查询用户列表 |
| GET | `/api/v1/users/:id` | 获取用户信息 |
| PUT | `/api/v1/users/:id` | 更新用户 |
| PATCH | `/api/v1/users/:id/status` | 更新用户状态 |
| DELETE | `/api/v1/users/:id` | 删除用户 |
| POST | `/api/v1/providers` | 创建 Provider |
| GET | `/api/v1/providers` | 查询 Provider 列表 |
| GET | `/api/v1/providers/:name` | 获取单个 Provider 信息 |
| PUT | `/api/v1/providers/:name` | 更新 Provider 配置 |
| DELETE | `/api/v1/providers/:name` | 删除 Provider |
| POST | `/api/v1/providers/reload` | 重载所有 Provider |
| POST | `/api/v1/admin/providers/:name/reload` | 重载指定 Provider |
| POST | `/api/v1/admin/providers/:name/enable` | 启用 Provider |
| POST | `/api/v1/admin/providers/:name/disable` | 禁用 Provider |

### 用户接口（需要 JWT Token）

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/users/:id` | 获取用户信息 | 管理员可获取任意用户，普通用户只能获取自己 |
| POST | `/api/v1/users/:id/api-keys` | 为用户创建 API Key | 管理员可为任意用户创建，普通用户只能为自己创建 |
| GET | `/api/v1/users/:id/api-keys` | 获取用户的 API Key 列表 | 管理员可查看任意用户，普通用户只能查看自己 |
| DELETE | `/api/v1/users/:id/api-keys/:key_id` | 撤销 API Key | 管理员可撤销任意用户的，普通用户只能撤销自己的 |
| GET | `/api/v1/usage` | 查询使用统计 | 管理员可查询所有用户，普通用户只能查询自己 |

### Chat API（需要用户 API Key）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/chat/completions` | Chat Completions（OpenAI 兼容） |

### 认证方式

- **JWT Token**：通过 `Authorization: Bearer <token>` 传递，用于管理接口和用户接口
- **用户 API Key**：格式为 `sk-<32位随机字符>`，通过 `Authorization: Bearer <key>` 传递，用于 Chat API

## 日志说明

### 结构化日志格式

所有日志使用 JSON 格式输出：

```json
{
  "level": "info",
  "ts": 1772523484.036684,
  "caller": "migrate/migrator.go:40",
  "msg": "Starting database auto-migration..."
}
```

### Chat 请求日志

每个 Chat 请求都会记录详细日志：

```json
{
  "trace_id": "trace-550e8400-e29b-41d4-a716-446655440000",
  "request_id": "chatcmpl-123",
  "api_key": "sk-...key",
  "model": "openai-main/gpt-4o",
  "provider_name": "openai-main",
  "model_name": "gpt-4o",
  "fallback_count": 1,
  "final_model": "gpt-4o-mini",
  "prompt_tokens": 10,
  "completion_tokens": 9,
  "total_tokens": 19,
  "latency_ms": 1250,
  "status": "success",
  "timestamp": "2026-03-03T12:00:00Z"
}
```

### TraceID

每个请求都会生成唯一的 TraceID，可通过响应头获取：

```http
X-Trace-ID: trace-550e8400-e29b-41d4-a716-446655440000
```

TraceID 用于追踪请求链路，方便排查问题。

## 故障排查

### 数据库连接失败

```
Failed to connect to database: connection refused
```

**解决方案**：
1. 检查 PostgreSQL 是否运行
2. 检查 `DATABASE_URL` 是否正确
3. 检查网络连接

### 迁移失败

```
Database migration failed
```

**解决方案**：
1. 检查数据库权限
2. 查看详细错误日志
3. 尝试手动删除 `schema_migrations` 表后重启

### Provider 不可用

```bash
# 检查 Provider 状态
curl http://localhost:8080/api/v1/providers \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

检查响应中的 `enabled` 字段。

### 认证失败

- 检查 JWT Token 是否过期
- 验证 `JWT_SECRET` 配置
- 确认用户状态为 `active`
