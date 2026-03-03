# Courier LLM Gateway

> 一个统一的 LLM API 网关，提供 OpenAI 兼容接口，支持多 Provider 接入、自动 Fallback、用户管理和使用统计。

## 特性

- **OpenAI 兼容接口**：完全兼容 OpenAI Chat Completions API
- **多 Provider 支持**：支持 OpenAI、通义千问、vLLM 等
- **自动 Fallback**：模型调用失败时自动切换到备用模型
- **用户管理**：基于角色的访问控制（Admin/User）
- **API Key 管理**：为用户生成和管理 API Key
- **使用统计**：记录和查询 API 使用情况
- **JWT 认证**：安全的 Token 认证机制
- **链路追踪**：每个请求唯一 TraceID，方便问题排查

## 快速开始

### 使用 Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f courier
```

### 本地开发

#### 1. 启动数据库

```bash
docker-compose up -d postgres
```

#### 2. 执行数据库迁移

```bash
export DATABASE_URL="host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable"

psql $DATABASE_URL -f migrations/000001_create_providers.up.sql
psql $DATABASE_URL -f migrations/000002_add_fallback_models.up.sql
psql $DATABASE_URL -f migrations/000003_create_users.up.sql
psql $DATABASE_URL -f migrations/000004_create_api_keys.up.sql
psql $DATABASE_URL -f migrations/000005_create_usage_records.up.sql
psql $DATABASE_URL -f migrations/000006_add_user_role.up.sql
psql $DATABASE_URL -f migrations/000007_add_password_hash.up.sql
```

#### 3. 运行服务

```bash
export DATABASE_URL="host=localhost port=5432 user=courier password=courier dbname=courier sslmode=disable"
export JWT_SECRET="your-jwt-secret-key-change-in-production"
export INITIAL_ADMIN_EMAIL="admin@example.com"
export INITIAL_ADMIN_PASSWORD="admin-password-change-me"

go run cmd/server/main.go
```

## 配置

### 环境变量

| 变量 | 描述 | 默认值 | 必需 |
|------|------|--------|------|
| `DATABASE_URL` | PostgreSQL 连接字符串 | - | ✓ |
| `PORT` | HTTP 服务端口 | 8080 | - |
| `JWT_SECRET` | JWT 签名密钥 | - | ✓ |
| `JWT_ACCESS_TOKEN_EXPIRES_IN` | Access Token 有效期 | 15m | - |
| `JWT_REFRESH_TOKEN_EXPIRES_IN` | Refresh Token 有效期 | 168h | - |
| `INITIAL_ADMIN_EMAIL` | 初始管理员邮箱 | - | - |
| `INITIAL_ADMIN_PASSWORD` | 初始管理员密码 | - | - |
| `LOG_LEVEL` | 日志级别 | info | - |

## 使用示例

### 1. 登录获取 JWT Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin-password-change-me"
  }'
```

### 2. 创建用户

```bash
ACCESS_TOKEN="..."  # 从登录响应获取

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com"
  }'
```

### 3. 创建 API Key

```bash
USER_ID="1"  # 用户 ID

curl -X POST "http://localhost:8080/api/v1/users/$USER_ID/api-keys" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "生产环境 Key"
  }'
```

### 4. 创建 Provider

```bash
# OpenAI Provider
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

# 通义千问 Provider
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

# 本地 vLLM Provider
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "local-vllm",
    "type": "vllm",
    "base_url": "http://localhost:8000/v1",
    "timeout": 120,
    "enabled": true
  }'
```

### 5. 调用 Chat API

```bash
API_KEY="sk-..."  # 从创建 API Key 响应获取

curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $API_KEY" \
  -d '{
    "model": "openai-main/gpt-4o",
    "messages": [{"role": "user", "content": "你好！"}]
  }'
```

## 文档

- [完整 API 文档](./docs/api.md)
- [部署文档](./docs/deployment.md)
- [Provider 配置指南](./docs/provider-config.md)
- [认证与鉴权](./docs/authentication.md)
- [API Key 管理](./docs/api-key-management.md)
- [Chat API 使用](./docs/chat-api.md)
- [Fallback 最佳实践](./docs/fallback-best-practices.md)

## 架构

```
┌─────────────────────────────────────────────────────────────┐
│                         Courier Gateway                       │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Auth      │  │   Router    │  │     Retry           │ │
│  │   Service   │  │   Service   │  │     Service         │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Provider  │  │    User     │  │      Usage          │ │
│  │   Service   │  │   Service   │  │      Service        │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                     Adapter Layer                           │
│  ┌─────────────┐  ┌─────────────┐                          │
│  │   OpenAI    │  │    vLLM     │  ┌──────────────┐       │
│  │   Adapter   │  │   Adapter   │  │   ...        │       │
│  └─────────────┘  └─────────────┘  └──────────────┘       │
├─────────────────────────────────────────────────────────────┤
│                     Storage Layer                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   PostgreSQL                         │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 项目结构

```
courier/
├── cmd/
│   └── server/
│       └── main.go              # 服务入口
├── internal/
│   ├── adapter/                 # Provider 适配器层
│   │   ├── openai/
│   │   └── vllm/
│   ├── controller/              # HTTP 控制器
│   ├── middleware/              # 中间件
│   ├── model/                   # 数据模型
│   ├── repository/              # 数据访问层
│   ├── service/                 # 业务逻辑层
│   └── bootstrap/               # 服务初始化
├── migrations/                  # 数据库迁移文件
├── docs/                        # 文档
├── openspec/                    # OpenSpec 变更管理
└── go.mod
```

## 开发

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test ./... -cover

# 运行特定包的测试
go test ./internal/service/...
```

### 代码格式化

```bash
# 格式化代码
go fmt ./...

# 运行 linter
go vet ./...
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f courier
```

### 生产环境建议

1. **安全性**
   - 使用强密码的数据库连接
   - 使用强随机密钥的 JWT_SECRET（至少 32 字符）
   - 配置 HTTPS（使用 Nginx 反向代理）

2. **性能**
   - 配置数据库连接池
   - 启用日志聚合（如 ELK）
   - 监控 Provider 调用延迟

3. **高可用**
   - 部署多实例 + 负载均衡
   - PostgreSQL 主从复制

## 贡献

欢迎提交 Issue 和 Pull Request！

## License

MIT License
