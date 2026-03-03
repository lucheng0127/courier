# 认证与鉴权

本文档说明 Courier LLM Gateway 的认证和鉴权机制。

## 概述

系统支持两种认证方式：

1. **JWT Token 认证**：用于管理接口，支持基于角色的访问控制
2. **API Key 认证**：用于 Chat API

## JWT Token 认证

### 登录获取 Token

**请求**：
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your-password"
}
```

**成功响应**：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**错误响应**：
```json
{
  "message": "Invalid email or password",
  "type": "authentication_error"
}
```

### 使用 Token 访问接口

在请求头中添加 Authorization：

```bash
GET /api/v1/users/123
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 刷新 Token

**请求**：
```bash
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

## API Key 认证

Chat API 使用 API Key 进行认证。

### 生成 API Key

**请求**：
```bash
POST /api/v1/users/123/api-keys
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "My API Key"
}
```

**响应**：
```json
{
  "id": 1,
  "key": "sk-abc123...",
  "key_prefix": "sk-abc123",
  "name": "My API Key",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z"
}
```

> **注意**：完整的 `key` 只在创建时返回一次，请妥善保存。

### 使用 API Key

```bash
POST /v1/chat/completions
Authorization: Bearer sk-abc123...
Content-Type: application/json

{
  "model": "openai/gpt-4",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ]
}
```

## 用户角色

系统支持两种用户角色：

### Admin（管理员）

- 可以访问所有管理接口
- 可以管理 Providers
- 可以创建和管理所有用户
- 可以查看所有用户的使用统计

### User（普通用户）

- 可以登录获取 Token
- 可以管理自己的 API Key
- 可以查看自己的使用统计
- 可以查看自己的用户信息

## 权限控制

### 接口权限表

| 接口 | Admin | User |
|------|-------|------|
| `/api/v1/auth/login` | ✓ | ✓ |
| `/api/v1/auth/refresh` | ✓ | ✓ |
| `/api/v1/providers/*` | ✓ | ✗ |
| `/api/v1/admin/providers/*` | ✓ | ✗ |
| `/api/v1/users` (列出) | ✓ | ✗ |
| `/api/v1/users` (创建) | ✓ | ✗ |
| `/api/v1/users/:id` (查看) | 任意用户 | 仅自己 |
| `/api/v1/users/:id` (更新) | ✓ | ✗ |
| `/api/v1/users/:id` (删除) | ✓ | ✗ |
| `/api/v1/users/:id/api-keys` | 任意用户 | 仅自己 |
| `/api/v1/usage` | 所有用户 | 仅自己 |
| `/v1/chat/completions` | API Key | API Key |

### 错误响应

**权限不足** (403)：
```json
{
  "message": "Admin privileges required",
  "type": "permission_error"
}
```

**未认证** (401)：
```json
{
  "message": "Invalid or expired access token",
  "type": "authentication_error"
}
```

## Token 过期处理

Access Token 有效期为 15 分钟。当 Token 过期时，客户端应该：

1. 捕获 401 响应
2. 使用 Refresh Token 调用 `/api/v1/auth/refresh`
3. 使用新的 Access Token 重试原请求

## 安全建议

1. **HTTPS**：生产环境必须使用 HTTPS
2. **Token 存储**：客户端应安全存储 Token（如使用 httpOnly cookie）
3. **密码强度**：密码至少 8 个字符，包含大小写字母和数字
4. **API Key 保护**：API Key 应该保密，不要提交到版本控制
5. **定期轮换**：定期更换 API Key 和密码

## 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `JWT_SECRET` | JWT 签名密钥（必须设置） | - |
| `JWT_ACCESS_TOKEN_EXPIRES_IN` | Access Token 有效期 | 15m |
| `JWT_REFRESH_TOKEN_EXPIRES_IN` | Refresh Token 有效期 | 168h |
| `INITIAL_ADMIN_EMAIL` | 初始管理员邮箱 | - |
| `INITIAL_ADMIN_PASSWORD` | 初始管理员密码 | - |
