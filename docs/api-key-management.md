# API Key 和用户管理 API 文档

本文档描述 Courier LLM Gateway 的用户管理、API Key 管理和使用统计 API。

## 前置条件

所有管理 API 需要通过 `X-Admin-API-Key` Header 进行管理员鉴权。

```bash
export ADMIN_API_KEY="your-admin-api-key"
```

## 用户管理

### 创建用户

创建新用户，用于管理和追踪 API 使用。

**请求**

```bash
POST /v1/users
Content-Type: application/json
X-Admin-API-Key: $ADMIN_API_KEY

{
  "name": "张三",
  "email": "zhangsan@example.com"
}
```

**响应**

```json
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com",
  "status": "active",
  "created_at": "2026-03-02T12:00:00Z",
  "updated_at": "2026-03-02T12:00:00Z"
}
```

### 获取用户信息

获取指定用户的详细信息。

**请求**

```bash
GET /v1/users/{user_id}
X-Admin-API-Key: $ADMIN_API_KEY
```

**响应**

```json
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com",
  "status": "active",
  "created_at": "2026-03-02T12:00:00Z",
  "updated_at": "2026-03-02T12:00:00Z"
}
```

## API Key 管理

### 创建 API Key

为指定用户创建新的 API Key。

**请求**

```bash
POST /v1/users/{user_id}/api-keys
Content-Type: application/json
X-Admin-API-Key: $ADMIN_API_KEY

{
  "name": "生产环境 Key",
  "expires_at": "2027-03-02T12:00:00Z"
}
```

**参数说明**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | API Key 的名称（用于识别） |
| expires_at | string | 否 | 过期时间（ISO 8601 格式） |

**响应**

```json
{
  "id": 1,
  "key": "sk-a1b2c3d4e5f67890abcdef1234567890",
  "key_prefix": "sk-a1b2c3d4",
  "name": "生产环境 Key",
  "status": "active",
  "expires_at": "2027-03-02T12:00:00Z",
  "created_at": "2026-03-02T12:00:00Z"
}
```

> **注意**: 完整的 API Key (`key` 字段) 仅在创建时返回一次，请妥善保存。

### 获取用户的 API Key 列表

获取指定用户的所有 API Key。

**请求**

```bash
GET /v1/users/{user_id}/api-keys
X-Admin-API-Key: $ADMIN_API_KEY
```

**响应**

```json
{
  "api_keys": [
    {
      "id": 1,
      "key_prefix": "sk-a1b2c3d4",
      "name": "生产环境 Key",
      "status": "active",
      "last_used_at": "2026-03-02T14:30:00Z",
      "expires_at": "2027-03-02T12:00:00Z",
      "created_at": "2026-03-02T12:00:00Z"
    }
  ]
}
```

### 撤销 API Key

撤销（删除）指定的 API Key。

**请求**

```bash
DELETE /v1/users/{user_id}/api-keys/{key_id}
X-Admin-API-Key: $ADMIN_API_KEY
```

**响应**

```
204 No Content
```

## 使用统计

### 查询使用统计

查询指定用户的使用统计信息。

**请求**

```bash
GET /v1/usage?user_id={user_id}&start_date={date}&end_date={date}&group_by={dimension}
X-Admin-API-Key: $ADMIN_API_KEY
```

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| user_id | integer | 是 | - | 用户 ID |
| start_date | string | 否 | 30 天前 | 开始日期（RFC3339 格式） |
| end_date | string | 否 | 当前时间 | 结束日期（RFC3339 格式） |
| group_by | string | 否 | day | 聚合维度：`day` 或 `model` |

**响应（按天聚合）**

```json
{
  "user_id": 1,
  "period": {
    "start": "2026-02-01T00:00:00Z",
    "end": "2026-03-02T00:00:00Z"
  },
  "summary": {
    "total_requests": 1250,
    "total_tokens": 1250000,
    "total_prompt_tokens": 800000,
    "total_completion_tokens": 450000,
    "average_latency_ms": 1250.5
  },
  "daily_breakdown": [
    {
      "date": "2026-03-01",
      "requests": 120,
      "tokens": 120000,
      "prompt_tokens": 80000,
      "completion_tokens": 40000,
      "average_latency_ms": 1100.0
    }
  ]
}
```

**响应（按模型聚合）**

```json
{
  "user_id": 1,
  "period": {
    "start": "2026-02-01T00:00:00Z",
    "end": "2026-03-02T00:00:00Z"
  },
  "summary": {
    "total_requests": 1250,
    "total_tokens": 1250000,
    "total_prompt_tokens": 800000,
    "total_completion_tokens": 450000,
    "average_latency_ms": 1250.5
  },
  "model_breakdown": [
    {
      "model": "openai-main/gpt-4o",
      "requests": 800,
      "tokens": 800000,
      "prompt_tokens": 500000,
      "completion_tokens": 300000,
      "average_latency_ms": 1200.0
    }
  ]
}
```

## 使用 API Key 调用 Chat API

获取到 API Key 后，使用标准的 `Authorization: Bearer` Header 调用 Chat API。

```bash
export API_KEY="sk-a1b2c3d4e5f67890abcdef1234567890"

curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $API_KEY" \
  -d '{
    "model": "openai-main/gpt-4o",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

## 错误码

| HTTP 状态码 | 错误类型 | 说明 |
|-------------|----------|------|
| 400 | invalid_request_error | 请求参数错误 |
| 401 | invalid_request_error | API Key 无效或缺失 |
| 403 | permission_error | 权限不足（如用户已禁用） |
| 404 | invalid_request_error | 资源不存在 |
| 409 | invalid_request_error | 资源冲突（如邮箱已存在） |
| 500 | api_error | 服务器内部错误 |

## API Key 格式说明

- **格式**: `sk-<32位随机十六进制字符>`
- **长度**: 35 个字符
- **示例**: `sk-a1b2c3d4e5f67890abcdef1234567890`
- **前缀**: 前 10 位用于识别（如 `sk-a1b2c3d4`）

## 最佳实践

1. **API Key 安全**
   - API Key 仅在创建时完整显示一次，请妥善保存
   - 不要在代码中硬编码 API Key
   - 定期轮换 API Key

2. **用户管理**
   - 一个用户可以创建多个 API Key，用于不同的应用场景
   - 为不同的环境（开发、测试、生产）创建不同的 Key
   - 为过期时间设置合理的值

3. **使用统计**
   - 使用按天聚合查看每日使用趋势
   - 使用按模型聚合了解不同模型的使用情况
   - 监控 `average_latency_ms` 关注性能变化
