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
# 初始迁移
psql "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable" -f migrations/000001_create_providers.up.sql

# Fallback 重试迁移
psql "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable" -f migrations/000002_add_fallback_models.up.sql

# 用户和 API Key 管理迁移
psql "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable" -f migrations/000003_create_users.up.sql
psql "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable" -f migrations/000004_create_api_keys.up.sql
psql "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable" -f migrations/000005_create_usage_records.up.sql

# 角色和密码迁移（JWT 认证）
psql "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable" -f migrations/000006_add_user_role.up.sql
psql "host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable" -f migrations/000007_add_password_hash.up.sql
```

### 3. 运行服务

```bash
export DATABASE_URL="host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable"
export JWT_SECRET="your-jwt-secret-key-change-in-production"
export INITIAL_ADMIN_EMAIL="admin@example.com"
export INITIAL_ADMIN_PASSWORD="admin-password-change-me"
go run cmd/server/main.go
```

### 4. 测试 API

```bash
# 登录获取 JWT Token
ACCESS_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin-password-change-me"
  }' | jq -r '.access_token')

# 创建用户
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com"
  }')
USER_ID=$(echo $USER_RESPONSE | jq -r '.id')

# 为用户创建 API Key
API_KEY_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/users/$USER_ID/api-keys" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
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

# 查询使用统计
curl "http://localhost:8080/api/v1/usage?user_id=$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### 5. 管理 Provider

```bash
# 创建 Provider（带 Fallback 配置）
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

# 查询 Provider 列表
curl http://localhost:8080/api/v1/providers \
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

## 环境变量

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

## Provider 配置

### Fallback 模型配置

为提高服务可用性，可为每个 Provider 配置 Fallback 模型列表：

```json
{
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "timeout": 60,
  "api_key": "sk-xxx",
  "enabled": true,
  "fallback_models": [
    "gpt-4o",
    "gpt-4o-mini",
    "gpt-3.5-turbo"
  ]
}
```

### Fallback 工作原理

1. 请求优先使用列表中的第一个模型（主模型）
2. 当主模型失败时（超时、网络错误、5xx 错误），自动尝试下一个模型
3. 直到成功或所有模型都失败

### 可观测性

所有请求日志包含以下字段：

- `trace_id` - 链路追踪 ID
- `fallback_count` - Fallback 次数
- `final_model` - 最终使用的模型
- `attempt_details` - 每次尝试的详情

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
   - 设置强随机密钥的 JWT_SECRET（至少 32 字符）
   - 配置 INITIAL_ADMIN_EMAIL 和 INITIAL_ADMIN_PASSWORD 创建初始管理员
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

# 回滚数据库迁移（按相反顺序）
psql $DATABASE_URL -f migrations/000007_add_password_hash.down.sql
psql $DATABASE_URL -f migrations/000006_add_user_role.down.sql
psql $DATABASE_URL -f migrations/000005_create_usage_records.down.sql
psql $DATABASE_URL -f migrations/000004_create_api_keys.down.sql
psql $DATABASE_URL -f migrations/000003_create_users.down.sql
psql $DATABASE_URL -f migrations/000002_add_fallback_models.down.sql
psql $DATABASE_URL -f migrations/000001_create_providers.down.sql
```

## API 接口说明

### 认证接口（无需鉴权）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录获取 JWT Token |
| POST | `/api/v1/auth/refresh` | 刷新 JWT Token |

### 管理接口（需要 JWT Token，Admin 角色）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/users` | 创建用户 |
| GET | `/api/v1/users` | 查询用户列表 |
| DELETE | `/api/v1/users/:id` | 删除用户 |
| PUT | `/api/v1/users/:id` | 更新用户 |
| PATCH | `/api/v1/users/:id/status` | 更新用户状态 |
| POST | `/api/v1/providers` | 创建 Provider |
| GET | `/api/v1/providers` | 查询 Provider 列表 |
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

详细认证说明请参考 [authentication.md](./authentication.md)
