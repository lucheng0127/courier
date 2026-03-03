# unify-api-and-jwt-auth 变更提案

## 概述

本变更旨在统一 API 接口风格并完善权限管理体系，主要包括：

1. **统一 API 路径前缀**：除 `/chat/completions` 外，所有管理接口统一使用 `/api/v1` 前缀
2. **实现 JWT Token 鉴权**：除 `/chat/completions` 使用 API Key 鉴权外，其余接口统一使用 JWT Token 鉴权
3. **完善角色权限管理**：添加用户角色系统（管理员/普通用户），只有管理员可以管理 providers 和创建用户；普通用户可以管理自己的 API Key，管理员可以管理任何用户的 API Key；普通用户可以登录并查询自己的使用统计
4. **统一路由注册风格**：所有 Controller 的路由注册参考 `providerCtrl.RegisterRoutes(api)` 风格

## 动机

### 当前问题

1. **API 路径不统一**：
   - Provider 管理：`/api/v1/providers/*`
   - Provider 运维：`/api/v1/admin/providers/*`
   - 用户管理：`/v1/users/*`
   - 使用统计：`/v1/usage/*`
   - Chat API：`/v1/chat/completions`

2. **鉴权机制混乱**：
   - Chat API 使用 API Key 鉴权（符合预期）
   - 管理接口使用 Admin API Key 鉴权（不灵活）
   - 没有基于用户角色的权限控制

3. **权限管理不足**：
   - 所有管理操作依赖单一的管理员 API Key
   - 无法区分不同管理员的操作
   - 用户表缺少角色字段

4. **代码风格不一致**：
   - Provider Controller 使用 `RegisterRoutes` 方法
   - User Controller 直接在 main.go 中注册路由

### 预期收益

1. **统一的 API 风格**：便于客户端集成和文档编写
2. **灵活的鉴权机制**：JWT Token 支持更细粒度的权限控制
3. **完善的权限管理**：支持多管理员、操作审计
4. **一致的代码风格**：便于维护和扩展

## 影响范围

### 模块变更

1. **数据库层**：
   - 添加 `users` 表的 `role` 字段

2. **模型层** (`internal/model/`)：
   - 更新 `User` 模型添加 `Role` 字段
   - 新增 JWT 相关请求/响应模型

3. **中间件层** (`internal/middleware/`)：
   - 新增 `JWTAuth` 中间件
   - 更新 `RequireAdmin` 中间件（基于角色而非 Admin API Key）
   - 保留 `APIKeyAuth` 中间件（仅用于 `/chat/completions`）

4. **服务层** (`internal/service/`)：
   - 新增 `JWTService` 处理 token 生成和验证
   - 更新 `AuthService` 支持登录（管理员登录）

5. **控制器层** (`internal/controller/`)：
   - 新增 `AuthController` 处理登录
   - 更新 `UserController` 实现 `RegisterRoutes` 方法
   - 更新 `UsageController` 实现 `RegisterRoutes` 方法
   - 更新 `ProviderReloadController` 路由路径

6. **路由层** (`cmd/server/main.go`)：
   - 统一路由注册方式
   - 移除硬编码的路由注册

### API 变更

#### 新增接口

| 方法 | 路径 | 描述 | 鉴权 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录获取 JWT Token（支持 admin 和 user 角色） | 无（使用邮箱+密码） |
| POST | `/api/v1/auth/refresh` | 刷新 JWT Token | JWT（refresh token） |
| POST | `/api/v1/auth/logout` | 登出（可选，标记 token 失效） | JWT |

#### 修改接口

| 原路径 | 新路径 | 变更说明 |
|--------|--------|----------|
| `POST /v1/users` | `POST /api/v1/users` | 路径前缀变更，鉴权改为 JWT + Admin（仅管理员） |
| `GET /v1/users/:id` | `GET /api/v1/users/:id` | 路径前缀变更，鉴权改为 JWT（用户可获取自己的，管理员可获取任何人的） |
| `POST /v1/users/:id/api-keys` | `POST /api/v1/users/:id/api-keys` | 路径前缀变更，鉴权改为 JWT（用户可管理自己的，管理员可管理任何人的） |
| `GET /v1/users/:id/api-keys` | `GET /api/v1/users/:id/api-keys` | 路径前缀变更，鉴权改为 JWT（用户可查看自己的，管理员可查看任何人的） |
| `DELETE /v1/users/:id/api-keys/:key_id` | `DELETE /api/v1/users/:id/api-keys/:key_id` | 路径前缀变更，鉴权改为 JWT（用户可删除自己的，管理员可删除任何人的） |
| `GET /v1/usage` | `GET /api/v1/usage` | 路径前缀变更，鉴权改为 JWT（管理员可查询所有用户，普通用户只能查询自己的） |
| `POST /api/v1/admin/providers/reload` | `POST /api/v1/admin/providers/reload` | 路径不变，鉴权改为 JWT + Admin |
| `POST /api/v1/admin/providers/:name/reload` | `POST /api/v1/admin/providers/:name/reload` | 路径不变，鉴权改为 JWT + Admin |
| `POST /api/v1/admin/providers/:name/enable` | `POST /api/v1/admin/providers/:name/enable` | 路径不变，鉴权改为 JWT + Admin |
| `POST /api/v1/admin/providers/:name/disable` | `POST /api/v1/admin/providers/:name/disable` | 路径不变，鉴权改为 JWT + Admin |

#### 不变接口

| 方法 | 路径 | 描述 | 鉴权 |
|------|------|------|------|
| POST | `/v1/chat/completions` | Chat Completions API | API Key |
| GET/POST/PUT/DELETE | `/api/v1/providers/*` | Provider 管理 | JWT + Admin |

## 实施策略

### 阶段 1：基础设施
1. 数据库迁移：添加 `role` 字段到 `users` 表
2. 实现 `JWTService`
3. 实现 `JWTAuth` 中间件
4. 更新 `RequireAdmin` 中间件

### 阶段 2：认证接口
1. 实现 `AuthController`（登录、刷新）
2. 添加登录路由

### 阶段 3：迁移现有接口
1. 更新 `UserController` 使用 `RegisterRoutes` 风格
2. 更新 `UsageController` 使用 `RegisterRoutes` 风格
3. 修改 `main.go` 中的路由注册

### 阶段 4：测试与文档
1. 更新 API 文档
2. 集成测试

## 兼容性

### 破坏性变更

1. **API 路径变更**：所有 `/v1/*` 路径（除 `/chat/completions`）改为 `/api/v1/*`
2. **鉴权方式变更**：管理接口从 Admin API Key 改为 JWT Token

### 迁移策略

1. **过渡期支持**：可以暂时保留 Admin API Key 鉴权作为降级方案
2. **客户端更新**：需要更新调用管理 API 的客户端代码
3. **文档更新**：及时更新 API 文档

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| JWT Token 泄露 | 可能被滥用 | 设置合理的过期时间，实现 refresh token 机制 |
| 数据库迁移失败 | 服务不可用 | 先在测试环境验证，准备回滚脚本 |
| 现有客户端中断 | 用户无法访问 | 保留过渡期支持，提前通知 |

## 依赖关系

- 需要 Go JWT 库（如 `github.com/golang-jwt/jwt/v5`）
- 可能需要密码哈希库（如果当前没有）

## 后续优化

- 添加 MFA 支持
- 实现 API Key 和 JWT 双重鉴权选项
- 实现 token 黑名单机制（支持主动登出）
