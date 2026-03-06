# 设计文档：API Key 启用/禁用/删除和 Chat 接口双重认证

## 架构概览

本变更涉及两个主要模块：
1. **API Key 状态管理** - 添加启用、禁用和删除操作
2. **Chat 接口认证** - 支持双重认证方式

## 模块 1：API Key 状态管理

### 现有实现

当前 API Key 模型支持三种状态：
- `active` - 激活状态，可以通过鉴权
- `disabled` - 禁用状态，无法通过鉴权
- `revoked` - 撤销状态，无法通过鉴权

现有接口：
- `POST /api/v1/users/:id/api-keys` - 创建 API Key
- `GET /api/v1/users/:id/api-keys` - 列出 API Key
- `DELETE /api/v1/users/:id/api-keys/:key_id` - 撤销 API Key（设置状态为 `revoked`）

### 新增接口

#### 1. 启用 API Key

**端点**：`PATCH /api/v1/users/:id/api-keys/:key_id/enable`

**请求**：无请求体

**响应**：200 OK
```json
{
  "id": 123,
  "key_prefix": "sk-courier1234",
  "name": "My API Key",
  "status": "active",
  "last_used_at": "2026-03-05T10:30:00Z",
  "expires_at": null,
  "created_at": "2026-03-01T08:00:00Z"
}
```

**实现逻辑**：
1. 验证用户权限（所有者或管理员）
2. 验证 API Key 属于该用户
3. 更新状态为 `active`
4. 返回更新后的 API Key 信息

#### 2. 禁用 API Key

**端点**：`PATCH /api/v1/users/:id/api-keys/:key_id/disable`

**请求**：无请求体

**响应**：200 OK（同启用接口）

**实现逻辑**：
1. 验证用户权限（所有者或管理员）
2. 验证 API Key 属于该用户
3. 更新状态为 `disabled`
4. 返回更新后的 API Key 信息

#### 3. 删除 API Key（硬删除）

**端点**：`DELETE /api/v1/users/:id/api-keys/:key_id`

**响应**：204 No Content

**实现逻辑**：
1. 验证用户权限（所有者或管理员）
2. 验证 API Key 属于该用户
3. 从数据库直接删除记录
4. 返回 204 No Content

**注意**：这是真正的删除操作，不可恢复。不同于现有的撤销操作（软删除）。

### 数据层变更

#### Repository 接口

新增 `DeleteAPIKey` 方法：

```go
// DeleteAPIKey 删除 API Key（硬删除）
DeleteAPIKey(ctx context.Context, id int64) error
```

#### SQL 实现

```sql
DELETE FROM api_keys WHERE id = $1
```

## 模块 2：Chat 接口双重认证

### 现有实现

当前 Chat 接口使用 `middleware.APIKeyAuth` 中间件进行认证：
- 从 `Authorization: Bearer <api_key>` 提取 API Key
- 验证 API Key 有效性
- 注入用户信息到 Context

### 新实现

#### 认证流程

创建新的组合认证中间件 `DualAuth`：

```
1. 尝试 JWT 认证
   ↓ 成功
   注入用户信息
   检查用户状态
   ↓
   处理请求

2. JWT 失败，尝试 API Key 认证
   ↓ 成功
   注入用户信息
   检查用户和 API Key 状态
   ↓
   处理请求

3. 两者都失败
   ↓
   返回 401 Unauthorized
```

#### 中间件实现

```go
// DualAuth 双重认证中间件（JWT 或 API Key）
func DualAuth(authService *service.AuthService, jwtSvc service.JWTService) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")

        // 尝试 JWT 认证
        if tryJWTAuth(ctx, jwtSvc, authHeader) {
            ctx.Next()
            return
        }

        // JWT 失败，尝试 API Key 认证
        if tryAPIKeyAuth(ctx, authService, authHeader) {
            ctx.Next()
            return
        }

        // 两者都失败
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "message": "Authentication failed. Please provide a valid API key or JWT token.",
                "type":    "authentication_error",
            },
        })
        ctx.Abort()
    }
}
```

#### Context 注入

**JWT 认证成功时注入**：
- `user_id` - 用户 ID
- `user_email` - 用户邮箱
- `user_role` - 用户角色
- `auth_type` - 认证类型（`jwt`）

**API Key 认证成功时注入**：
- `user_id` - 用户 ID
- `user_email` - 用户邮箱
- `api_key_id` - API Key 记录 ID
- `api_key_masked` - 脱敏后的 API Key
- `auth_type` - 认证类型（`apikey`）

#### 使用量记录

**JWT 认证**：
- 使用量记录只包含 `user_id`，不关联 API Key
- `api_key_id` 字段为 NULL

**API Key 认证**：
- 使用量记录包含 `user_id` 和 `api_key_id`（现有行为保持不变）

## 安全考虑

### 1. 权限验证

所有 API Key 操作都需要验证：
- 操作者是 Key 的所有者
- 或者操作者是管理员

### 2. 状态转换

API Key 状态转换规则：
- `active` ↔ `disabled` - 可以互相切换
- `active` → `revoked` - 可以撤销
- `disabled` → `revoked` - 可以撤销
- `revoked` - 不可逆，不能重新激活
- **删除操作** - 可以删除任何状态的 Key

### 3. 审计日志

建议记录以下操作：
- API Key 启用/禁用操作
- API Key 删除操作
- Chat 请求的认证类型（JWT 或 API Key）

## 兼容性

### 向后兼容

1. **现有 API Key 接口**：
   - 创建、列出、撤销接口保持不变
   - 客户端无需修改

2. **Chat 接口**：
   - API Key 认证方式保持完全兼容
   - 新增 JWT 认证为可选方式

### API 版本

所有变更都在 `/api/v1` 和 `/v1` 路径下，属于同一版本。

## 错误处理

### API Key 操作错误

| 错误场景 | HTTP 状态码 | 错误类型 | 错误信息 |
|---------|------------|---------|---------|
| API Key 不存在 | 404 | `invalid_request_error` | API key not found |
| 无权限操作 | 403 | `permission_error` | Permission denied |
| 状态已是目标状态 | 400 | `invalid_request_error` | API key is already {status} |

### 认证错误

| 错误场景 | HTTP 状态码 | 错误类型 | 错误信息 |
|---------|------------|---------|---------|
| 缺少 Authorization Header | 401 | `authentication_error` | Missing authorization header |
| Authorization 格式错误 | 401 | `authentication_error` | Invalid authorization header format |
| JWT 无效或过期 | 401 | `authentication_error` | Invalid or expired access token |
| API Key 无效 | 401 | `invalid_request_error` | Invalid API key |
| 用户被禁用 | 403 | `permission_error` | User account is disabled |

## 性能考虑

### 数据库查询优化

双重认证不会增加额外的数据库查询：
- JWT 认证：只验证 Token 签名和过期时间，无需查询 API Key
- API Key 认证：查询 API Key 进行验证（现有行为）

### 并发控制

API Key 状态更新的并发控制：
- 使用数据库级别的乐观锁
- 或使用 `UPDATE ... WHERE id = ? AND status = ?` 确保状态一致性

## 测试策略

### 单元测试

1. **API Key 操作**：
   - 启用 API Key
   - 禁用 API Key
   - 删除 API Key
   - 权限验证

2. **双重认证**：
   - JWT 认证成功
   - API Key 认证成功
   - 两者都失败
   - 认证类型注入

### 集成测试

1. **端到端流程**：
   - 用户创建 API Key
   - 禁用 API Key
   - 使用 JWT 访问 Chat 接口
   - 使用 API Key 访问 Chat 接口
   - 删除 API Key

2. **使用量记录**：
   - JWT 认证时的使用量记录（只记录用户 ID）
   - API Key 认证时的使用量记录（记录用户 ID 和 API Key ID）
