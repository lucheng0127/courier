# API Key 鉴权和用户使用统计 - 设计文档

## 概述

本设计文档描述 API Key 鉴权和用户使用统计功能的详细实现方案，包括数据模型、API 设计、中间件行为和统计查询逻辑。

## 数据模型设计

### 用户 (User)

```go
type User struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    Status    string    `json:"status" db:"status"` // active, disabled
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### API Key

```go
type APIKey struct {
    ID           string    `json:"id" db:"id"`
    UserID       string    `json:"user_id" db:"user_id"`
    KeyHash      string    `json:"-" db:"key_hash"` // SHA256 哈希存储
    KeyPrefix    string    `json:"key_prefix" db:"key_prefix"` // 前8位用于识别
    Name         string    `json:"name" db:"name"` // 用户定义的名称
    Status       string    `json:"status" db:"status"` // active, disabled, revoked
    LastUsedAt   *time.Time `json:"last_used_at" db:"last_used_at"`
    ExpiresAt    *time.Time `json:"expires_at" db:"expires_at"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
```

**设计决策**：
- API Key 使用 SHA256 哈希存储，原始 Key 只在创建时返回一次
- 存储前 8 位（key_prefix）用于日志和管理界面识别
- 支持过期时间和状态管理（active/disabled/revoked）

### 使用记录 (UsageRecord)

```go
type UsageRecord struct {
    ID              string    `json:"id" db:"id"`
    UserID          string    `json:"user_id" db:"user_id"`
    APIKeyID        string    `json:"api_key_id" db:"api_key_id"`
    RequestID       string    `json:"request_id" db:"request_id"`
    TraceID         string    `json:"trace_id" db:"trace_id"`
    Model           string    `json:"model" db:"model"`
    ProviderName    string    `json:"provider_name" db:"provider_name"`
    PromptTokens    int       `json:"prompt_tokens" db:"prompt_tokens"`
    CompletionTokens int      `json:"completion_tokens" db:"completion_tokens"`
    TotalTokens     int       `json:"total_tokens" db:"total_tokens"`
    LatencyMs       int64     `json:"latency_ms" db:"latency_ms"`
    Status          string    `json:"status" db:"status"` // success, error
    ErrorType       string    `json:"error_type,omitempty" db:"error_type"`
    Timestamp       time.Time `json:"timestamp" db:"timestamp"`
}
```

## 数据库表设计

### users 表

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

### api_keys 表

```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(64) UNIQUE NOT NULL,
    key_prefix VARCHAR(8) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT check_status CHECK (status IN ('active', 'disabled', 'revoked'))
);

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_status ON api_keys(status);
```

### usage_records 表

```sql
CREATE TABLE usage_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    api_key_id UUID NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    request_id VARCHAR(255) NOT NULL,
    trace_id VARCHAR(255),
    model VARCHAR(255) NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    prompt_tokens INTEGER NOT NULL DEFAULT 0,
    completion_tokens INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    latency_ms BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_type VARCHAR(100),
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX idx_usage_records_api_key_id ON usage_records(api_key_id);
CREATE INDEX idx_usage_records_timestamp ON usage_records(timestamp);
CREATE INDEX idx_usage_records_user_timestamp ON usage_records(user_id, timestamp);
```

## API 设计

### 用户管理 API

#### 创建用户
```
POST /v1/users
Content-Type: application/json
X-Admin-API-Key: <admin_key>

{
  "name": "张三",
  "email": "zhangsan@example.com"
}

Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "张三",
  "email": "zhangsan@example.com",
  "status": "active",
  "created_at": "2026-03-02T12:00:00Z"
}
```

#### 获取用户信息
```
GET /v1/users/:id
X-Admin-API-Key: <admin_key>

Response:
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "张三",
  "email": "zhangsan@example.com",
  "status": "active",
  "created_at": "2026-03-02T12:00:00Z"
}
```

### API Key 管理 API

#### 创建 API Key
```
POST /v1/users/:id/api-keys
X-Admin-API-Key: <admin_key>

{
  "name": "生产环境 Key"
}

Response:
{
  "id": "key-id",
  "key": "sk-courier-xxxxx...",  // 仅在创建时返回一次
  "key_prefix": "sk-cour",
  "name": "生产环境 Key",
  "status": "active",
  "created_at": "2026-03-02T12:00:00Z"
}
```

**设计决策**：API Key 格式为 `sk-courier-<随机32字符>`，使用 crypto/rand 生成。

#### 获取用户的 API Key 列表
```
GET /v1/users/:id/api-keys
X-Admin-API-Key: <admin_key>

Response:
{
  "api_keys": [
    {
      "id": "key-id",
      "key_prefix": "sk-cour",
      "name": "生产环境 Key",
      "status": "active",
      "last_used_at": "2026-03-02T13:00:00Z",
      "created_at": "2026-03-02T12:00:00Z"
    }
  ]
}
```

#### 删除/禁用 API Key
```
DELETE /v1/users/:id/api-keys/:key_id
X-Admin-API-Key: <admin_key>

Response: 204 No Content
```

### 使用统计 API

#### 查询使用统计
```
GET /v1/usage?user_id=<id>&start_date=<date>&end_date=<date>&group_by=<field>
X-Admin-API-Key: <admin_key>

Query Parameters:
- user_id (required): 用户 ID
- start_date (optional): 开始日期，默认 30 天前
- end_date (optional): 结束日期，默认今天
- group_by (optional): 聚合维度，可选 day/model，默认按天聚合

Response (按天聚合):
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "period": {
    "start": "2026-02-01T00:00:00Z",
    "end": "2026-03-02T00:00:00Z"
  },
  "summary": {
    "total_requests": 1250,
    "total_tokens": 1250000,
    "total_prompt_tokens": 800000,
    "total_completion_tokens": 450000,
    "average_latency_ms": 1250
  },
  "daily_breakdown": [
    {
      "date": "2026-03-01",
      "requests": 120,
      "tokens": 120000,
      "prompt_tokens": 80000,
      "completion_tokens": 40000,
      "average_latency_ms": 1100
    }
  ]
}
```

## 中间件设计

### 鉴权中间件行为

```go
func APIKeyAuth(authService AuthService) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // 1. 从 Authorization Header 提取 Bearer token
        apiKey := extractBearerToken(ctx)

        // 2. 从数据库查询 API Key
        keyRecord, err := authService.ValidateAPIKey(ctx, apiKey)
        if err != nil {
            ctx.JSON(401, gin.H{"error": "Invalid API key"})
            ctx.Abort()
            return
        }

        // 3. 检查 API Key 状态
        if keyRecord.Status != "active" {
            ctx.JSON(401, gin.H{"error": "API key is disabled"})
            ctx.Abort()
            return
        }

        // 4. 检查过期时间
        if keyRecord.ExpiresAt != nil && keyRecord.ExpiresAt.Before(time.Now()) {
            ctx.JSON(401, gin.H{"error": "API key has expired"})
            ctx.Abort()
            return
        }

        // 5. 获取用户信息
        user, err := authService.GetUserByID(ctx, keyRecord.UserID)
        if err != nil {
            ctx.JSON(401, gin.H{"error": "User not found"})
            ctx.Abort()
            return
        }

        if user.Status != "active" {
            ctx.JSON(403, gin.H{"error": "User account is disabled"})
            ctx.Abort()
            return
        }

        // 6. 注入到 Context
        ctx.Set("user_id", user.ID)
        ctx.Set("user_email", user.Email)
        ctx.Set("api_key_id", keyRecord.ID)
        ctx.Set("api_key_masked", maskAPIKey(apiKey))

        // 7. 异步更新 last_used_at（后台任务）
        go authService.UpdateKeyLastUsed(keyRecord.ID)

        ctx.Next()
    }
}
```

### 使用量记录中间件

在 Chat 请求完成后，从 Context 中提取用户信息并记录使用量：

```go
func (c *ChatController) logRequestWithRetry(...) {
    // 从 Context 获取用户信息
    userID, _ := ctx.Get("user_id")
    apiKeyID, _ := ctx.Get("api_key_id")

    // ... 现有日志逻辑 ...

    // 新增：记录到数据库
    if err := c.usageService.RecordUsage(context.Background(), &model.UsageRecord{
        UserID:          userID.(string),
        APIKeyID:        apiKeyID.(string),
        RequestID:       requestID,
        TraceID:         traceID,
        Model:           req.Model,
        ProviderName:    modelInfo.ProviderName,
        PromptTokens:    log.PromptTokens,
        CompletionTokens: log.CompletionTokens,
        TotalTokens:     log.TotalTokens,
        LatencyMs:       latencyMs,
        Status:          status,
        Timestamp:       time.Now(),
    }); err != nil {
        // 记录失败不影响主流程，只记录错误日志
        logger.Error("Failed to record usage", map[string]any{"error": err.Error()})
    }
}
```

## 性能考虑

1. **API Key 查询缓存**：使用 Redis 缓存活跃的 API Key，减少数据库查询
2. **异步写入使用记录**：使用 channel 批量写入数据库，避免阻塞主请求
3. **使用记录分区**：按时间分区存储 usage_records 表，便于查询和清理历史数据

## 安全考虑

1. **API Key 哈希**：使用 SHA256 哈希存储，salt 使用固定字符串 + key_prefix
2. **管理员 API Key**：管理员接口使用单独的 `X-Admin-API-Key` Header 验证
3. **最小权限原则**：普通用户只能查询自己的使用统计（未来扩展）
4. **审计日志**：记录所有 API Key 创建/删除操作

## 未来扩展

1. **配额管理**：为用户设置每日/每月 Token 配额
2. **计费统计**：按不同模型设置不同价格，计算费用
3. **自助管理**：用户自行管理 API Key（需要用户登录系统）
4. **使用报告**：定期生成使用报告并发送给用户
