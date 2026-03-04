# Courier LLM Gateway - API 文档

## 目录

- [概述](#概述)
- [认证](#认证)
- [认证接口](#认证接口)
  - [用户注册](#用户注册)
  - [登录](#登录)
  - [刷新 Token](#刷新-token)
- [用户管理](#用户管理)
- [API Key 管理](#api-key-管理)
- [Provider 管理](#provider-管理)
- [Chat API](#chat-api)
- [使用统计](#使用统计)
- [错误处理](#错误处理)

---

## 概述

Courier LLM Gateway 是一个统一的 LLM API 网关，提供 OpenAI 兼容的接口，支持对接多个 LLM Provider。

**Base URL**: `http://localhost:8080`

**API 版本**: `v1`

---

## 认证

系统支持两种认证方式：

### 1. JWT Token 认证

用于管理接口，通过 `Authorization: Bearer <token>` 传递。

**获取 Token**：调用 `/api/v1/auth/login` 接口

### 2. API Key 认证

用于 Chat API，格式为 `sk-<32位随机字符>`，通过 `Authorization: Bearer <key>` 传递。

---

## 认证接口

### 用户注册

**描述**：用户自主注册账户，注册成功后默认角色为 `user`，状态为 `active`。

**权限**：无需认证

**速率限制**：同一 IP 每小时最多 5 次注册请求

**请求**：
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "张三",
  "email": "zhangsan@example.com",
  "password": "your-password"
}
```

**请求参数**：
| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| name | string | 是 | 用户名称 |
| email | string | 是 | 邮箱地址，必须唯一 |
| password | string | 是 | 密码，至少 8 个字符 |

**响应**（成功）：
```json
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com",
  "role": "user",
  "status": "active",
  "created_at": "2026-03-03T00:00:00Z"
}
```

**响应**（邮箱已存在）：
```json
{
  "message": "Email already exists",
  "type": "invalid_request_error"
}
```

**响应**（密码过短）：
```json
{
  "message": "Password must be at least 8 characters",
  "type": "invalid_request_error"
}
```

**响应**（超过速率限制）：
```json
{
  "message": "Too many registration attempts, please try again later",
  "type": "rate_limit_error"
}
```

> **注意**：注册成功后，用户需要调用登录接口获取 JWT Token。

### 登录

**请求**：
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your-password"
}
```

**响应**：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

### 刷新 Token

**请求**：
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**响应**：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

---

## 用户管理

> **注意**：用户创建已改为自主注册模式。管理员不再通过 API 创建普通用户。新用户通过 `POST /api/v1/auth/register` 接口自行注册。

### 查询用户列表

**权限**: Admin

**请求**：
```http
GET /api/v1/users
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "users": [
    {
      "id": 1,
      "name": "张三",
      "email": "zhangsan@example.com",
      "role": "user",
      "status": "active",
      "created_at": "2026-03-03T00:00:00Z"
    }
  ],
  "total": 1
}
```

### 获取用户信息

**权限**: Admin（可查看任意用户），User（仅可查看自己）

**请求**：
```http
GET /api/v1/users/:id
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "id": 1,
  "name": "张三",
  "email": "zhangsan@example.com",
  "role": "user",
  "status": "active",
  "created_at": "2026-03-03T00:00:00Z"
}
```

### 更新用户

**权限**: Admin

**请求**：
```http
PUT /api/v1/users/:id
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "李四"
}
```

**响应**：
```json
{
  "id": 1,
  "name": "李四",
  "email": "zhangsan@example.com",
  "role": "user",
  "status": "active",
  "updated_at": "2026-03-03T00:00:00Z"
}
```

### 更新用户状态

**权限**: Admin

**请求**：
```http
PATCH /api/v1/users/:id/status
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "status": "disabled"
}
```

**响应**：
```json
{
  "id": 1,
  "status": "disabled",
  "updated_at": "2026-03-03T00:00:00Z"
}
```

### 删除用户

**权限**: Admin

**请求**：
```http
DELETE /api/v1/users/:id
Authorization: Bearer <jwt-token>
```

**响应**: `204 No Content`

---

## API Key 管理

### 创建 API Key

**权限**: Admin（可为任意用户创建），User（仅可为自己创建）

**请求**：
```http
POST /api/v1/users/:id/api-keys
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "生产环境 Key"
}
```

**响应**：
```json
{
  "id": 1,
  "key": "sk-abc123def456...",
  "key_prefix": "sk-abc123",
  "name": "生产环境 Key",
  "status": "active",
  "created_at": "2026-03-03T00:00:00Z"
}
```

> **注意**: 完整的 `key` 仅在创建时返回一次，请妥善保存。

### 获取 API Key 列表

**权限**: Admin（可查看任意用户），User（仅可查看自己）

**请求**：
```http
GET /api/v1/users/:id/api-keys
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "api_keys": [
    {
      "id": 1,
      "key_prefix": "sk-abc123",
      "name": "生产环境 Key",
      "status": "active",
      "created_at": "2026-03-03T00:00:00Z",
      "last_used_at": "2026-03-03T12:00:00Z"
    }
  ]
}
```

### 撤销 API Key

**权限**: Admin（可撤销任意用户的），User（仅可撤销自己的）

**请求**：
```http
DELETE /api/v1/users/:id/api-keys/:key_id
Authorization: Bearer <jwt-token>
```

**响应**: `204 No Content`

---

## Provider 管理

### 创建 Provider

**权限**: Admin

**请求**：
```http
POST /api/v1/providers
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "timeout": 60,
  "api_key": "sk-xxx",
  "enabled": true,
  "extra_config": {
    "temperature": 0.7,
    "max_tokens": 2000
  },
  "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
}
```

**响应**：
```json
{
  "id": 1,
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "timeout": 60,
  "enabled": true,
  "created_at": "2026-03-03T00:00:00Z"
}
```

### 查询 Provider 列表

**权限**: 所有认证用户

**查询参数**:
- `enabled` (可选): 过滤条件
  - `true`: 只返回已启用的 Provider
  - `false`: 只返回已禁用的 Provider
  - 不传参数: 返回所有 Provider

**请求**：
```http
GET /api/v1/providers
Authorization: Bearer <jwt-token>
```

**管理员响应**（完整信息）：
```json
{
  "providers": [
    {
      "provider": {
        "id": 1,
        "name": "openai-main",
        "type": "openai",
        "base_url": "https://api.openai.com/v1",
        "timeout": 60,
        "enabled": true,
        "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"],
        "created_at": "2026-03-03T00:00:00Z",
        "updated_at": "2026-03-03T00:00:00Z"
      },
      "is_running": true
    }
  ]
}
```

**普通用户响应**（简化信息）：
```json
{
  "providers": [
    {
      "name": "openai-main",
      "type": "openai",
      "base_url": "https://api.openai.com/v1",
      "enabled": true,
      "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
    }
  ]
}
```

> **注意**: 普通用户响应不包含敏感信息如 `api_key`、`timeout`、`is_running` 等。

### 获取单个 Provider

**权限**: Admin

**请求**：
```http
GET /api/v1/providers/:name
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "name": "openai-main",
  "type": "openai"
}
```

### 获取 Provider 模型列表

**权限**: 所有认证用户

**请求**：
```http
GET /api/v1/providers/:name/models
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "name": "openai-main",
  "type": "openai",
  "models": [
    "gpt-4o",
    "gpt-4o-mini",
    "gpt-3.5-turbo"
  ]
}
```

> **说明**: 模型列表来源于 Provider 配置中的 `fallback_models` 字段。

### 更新 Provider

**权限**: Admin

**请求**：
```http
PUT /api/v1/providers/:name
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "timeout": 120,
  "enabled": false
}
```

**响应**：
```json
{
  "id": 1,
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "timeout": 120,
  "enabled": false,
  "updated_at": "2026-03-03T00:00:00Z"
}
```

### 删除 Provider

**权限**: Admin

**请求**：
```http
DELETE /api/v1/providers/:name
Authorization: Bearer <jwt-token>
```

**响应**: `204 No Content`

### 重载所有 Provider

**权限**: Admin

**请求**：
```http
POST /api/v1/admin/providers/reload
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "message": "All providers reloaded successfully",
  "count": 2
}
```

### 重载指定 Provider

**权限**: Admin

**请求**：
```http
POST /api/v1/admin/providers/:name/reload
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "message": "Provider openai-main reloaded successfully"
}
```

### 启用 Provider

**权限**: Admin

**请求**：
```http
POST /api/v1/admin/providers/:name/enable
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "message": "Provider openai-main enabled successfully"
}
```

### 禁用 Provider

**权限**: Admin

**请求**：
```http
POST /api/v1/admin/providers/:name/disable
Authorization: Bearer <jwt-token>
```

**响应**：
```json
{
  "message": "Provider openai-main disabled successfully"
}
```

---

## Chat API

### Chat Completions

OpenAI 兼容的 Chat Completions API。

**端点**：`POST /v1/chat/completions`

**认证**: API Key

**请求**：
```json
{
  "model": "provider-name/model-name",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello!"}
  ],
  "temperature": 0.7,
  "max_tokens": 1000,
  "stream": false
}
```

**响应（非流式）**：
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "openai-main/gpt-4o",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! How can I help you today?"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 9,
    "total_tokens": 19
  }
}
```

**响应（流式）**：
```
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"openai-main/gpt-4o","choices":[{"index":0,"delta":{"content":"Hello"}}]}

data: [DONE]
```

### 模型格式

采用 **OpenRouter 风格**：`provider/model_name`

- `provider` - Provider 实例名称
- `model_name` - 模型名称

**示例**：
- `openai-main/gpt-4o`
- `qwen-main/qwen-turbo`
- `vllm-local/llama-2-7b`

### Fallback 机制

当模型调用失败时，系统会自动尝试 Fallback 列表中的下一个模型。

**触发条件**：
- 超时错误
- 网络错误
- 5xx 服务器错误
- 连接拒绝

**不触发 Fallback**：
- 4xx 客户端错误
- 认证失败
- 模型不存在

---

## 使用统计

### 查询使用统计

**权限**: Admin（可查询所有用户），User（仅可查询自己）

**请求**：
```http
GET /api/v1/usage?user_id=1&start_date=2026-03-01&end_date=2026-03-03
Authorization: Bearer <jwt-token>
```

**参数**：
| 参数 | 类型 | 描述 |
|------|------|------|
| user_id | int | 用户 ID（Admin 必填，User 可选） |
| start_date | string | 开始日期（可选） |
| end_date | string | 结束日期（可选） |
| page | int | 页码（默认 1） |
| page_size | int | 每页数量（默认 20） |

**响应**：
```json
{
  "records": [
    {
      "id": 1,
      "user_id": 1,
      "model": "openai-main/gpt-4o",
      "provider_name": "openai-main",
      "prompt_tokens": 100,
      "completion_tokens": 50,
      "total_tokens": 150,
      "latency_ms": 1250,
      "status": "success",
      "timestamp": "2026-03-03T12:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

---

## 错误处理

所有错误响应遵循统一格式：

```json
{
  "error": {
    "message": "错误描述",
    "type": "error_type"
  }
}
```

### 错误类型

| HTTP 状态 | 错误类型 | 描述 |
|-----------|----------|------|
| 400 | `invalid_request_error` | 请求参数错误 |
| 401 | `authentication_error` | 认证失败 |
| 403 | `permission_error` | 权限不足 |
| 404 | `not_found_error` | 资源不存在 |
| 429 | `rate_limit_error` | 请求频率限制 |
| 500 | `api_error` | 服务器内部错误 |
| 503 | `service_unavailable` | 服务不可用 |

### 错误示例

**认证失败**：
```json
{
  "error": {
    "message": "Invalid or expired access token",
    "type": "authentication_error"
  }
}
```

**权限不足**：
```json
{
  "error": {
    "message": "Admin privileges required",
    "type": "permission_error"
  }
}
```

**模型格式错误**：
```json
{
  "error": {
    "message": "invalid model format: gpt-4 (expected format: provider/model_name)",
    "type": "invalid_request_error"
  }
}
```

**Fallback 耗尽**：
```json
{
  "error": {
    "message": "All models failed after 3 attempts. Last error: timeout",
    "type": "service_unavailable",
    "details": [
      {
        "model": "gpt-4o",
        "error_type": "timeout",
        "duration_ms": 30000
      }
    ]
  }
}
```

---

## 请求头

| 请求头 | 描述 |
|--------|------|
| `Authorization` | Bearer Token（JWT 或 API Key） |
| `Content-Type` | application/json |
| `X-Trace-ID` | 链路追踪 ID（响应返回） |

---

## 速率限制

- **注册接口**：同一 IP 每小时最多 5 次注册请求
- **JWT Token 认证接口**：无限制
- **Chat API**：根据用户配置限制

---

## SDK 示例

### Python

```python
import requests

# Chat API
response = requests.post(
    "http://localhost:8080/v1/chat/completions",
    headers={
        "Authorization": "Bearer sk-your-key",
        "Content-Type": "application/json"
    },
    json={
        "model": "openai-main/gpt-4o",
        "messages": [{"role": "user", "content": "Hello!"}]
    }
)
print(response.json())
```

### JavaScript/Node.js

```javascript
const response = await fetch('http://localhost:8080/v1/chat/completions', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer sk-your-key',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    model: 'openai-main/gpt-4o',
    messages: [{ role: 'user', content: 'Hello!' }]
  })
});

const data = await response.json();
console.log(data);
```

### cURL

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer sk-your-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "openai-main/gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```
