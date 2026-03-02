# unified-api-routes Specification Delta

## Purpose

统一所有 API 接口的路径前缀和路由注册风格。

## MODIFIED Requirements

### Requirement: API 路径前缀规范

除 `/chat/completions` 外，所有管理接口 SHALL 使用 `/api/v1` 作为路径前缀。

#### Scenario: 管理 API 路径前缀

- **WHEN** 访问任意管理接口
- **THEN** 路径以 `/api/v1` 开头
- **AND** 不存在以 `/v1` 开头的管理接口（除 Chat API）

#### Scenario: Chat API 路径保持不变

- **WHEN** 访问 Chat Completions API
- **THEN** 路径为 `/v1/chat/completions`
- **AND** 保持与 OpenAI 兼容的路径格式

### Requirement: Provider 管理 API 路径

Provider 管理接口 SHALL 统一使用 `/api/v1/providers` 路径前缀。

#### Scenario: Provider CRUD 路径

- **WHEN** 访问 Provider 管理接口
- **THEN** 接口路径如下：
  - `POST /api/v1/providers` - 创建 Provider
  - `GET /api/v1/providers` - 列出所有 Provider
  - `GET /api/v1/providers/:name` - 获取单个 Provider
  - `PUT /api/v1/providers/:name` - 更新 Provider
  - `DELETE /api/v1/providers/:name` - 删除 Provider

#### Scenario: Provider 运维路径

- **WHEN** 访问 Provider 运维接口
- **THEN** 接口路径如下：
  - `POST /api/v1/admin/providers/reload` - 重载所有 Provider
  - `POST /api/v1/admin/providers/:name/reload` - 重载指定 Provider
  - `POST /api/v1/admin/providers/:name/enable` - 启用 Provider
  - `POST /api/v1/admin/providers/:name/disable` - 禁用 Provider

### Requirement: 用户管理 API 路径

用户管理接口 SHALL 统一使用 `/api/v1/users` 路径前缀。

#### Scenario: 用户 CRUD 路径

- **WHEN** 访问用户管理接口
- **THEN** 接口路径如下：
  - `POST /api/v1/users` - 创建用户
  - `GET /api/v1/users` - 列出所有用户
  - `GET /api/v1/users/:id` - 获取单个用户
  - `PUT /api/v1/users/:id` - 更新用户信息
  - `DELETE /api/v1/users/:id` - 删除用户
  - `PATCH /api/v1/users/:id/status` - 更新用户状态

#### Scenario: API Key 管理路径

- **WHEN** 访问 API Key 管理接口
- **THEN** 接口路径如下：
  - `POST /api/v1/users/:id/api-keys` - 为用户创建 API Key
  - `GET /api/v1/users/:id/api-keys` - 获取用户的 API Key 列表
  - `DELETE /api/v1/users/:id/api-keys/:key_id` - 撤销 API Key

### Requirement: 使用统计 API 路径

使用统计接口 SHALL 统一使用 `/api/v1/usage` 路径前缀。

#### Scenario: 使用统计路径

- **WHEN** 访问使用统计接口
- **THEN** 接口路径如下：
  - `GET /api/v1/usage` - 查询使用统计（支持按 user_id、日期范围、分组方式过滤）
  - `GET /api/v1/usage/:user_id` - 查询指定用户的统计（可选，与上面的接口合并）

### Requirement: 认证 API 路径

认证相关接口 SHALL 使用 `/api/v1/auth` 路径前缀。

#### Scenario: 认证接口路径

- **WHEN** 访问认证接口
- **THEN** 接口路径如下：
  - `POST /api/v1/auth/login` - 管理员登录
  - `POST /api/v1/auth/refresh` - 刷新 JWT Token
  - `POST /api/v1/auth/logout` - 登出（可选）

### Requirement: 路由注册风格统一

所有 Controller SHALL 实现 `RegisterRoutes` 方法，统一路由注册风格。

#### Scenario: ProviderController 路由注册

- **WHEN** 查看 ProviderController
- **THEN** 实现了 `RegisterRoutes(r *gin.RouterGroup)` 方法
- **AND** 方法内注册所有 Provider 相关路由
- **AND** 不在 main.go 中硬编码路由

#### Scenario: UserController 路由注册

- **WHEN** 查看 UserController
- **THEN** 实现了 `RegisterRoutes(r *gin.RouterGroup)` 方法
- **AND** 方法内注册所有用户和 API Key 相关路由
- **AND** 包括：创建、列表、获取、更新、删除用户
- **AND** 包括：创建、列表、撤销 API Key

#### Scenario: UsageController 路由注册

- **WHEN** 查看 UsageController
- **THEN** 实现了 `RegisterRoutes(r *gin.RouterGroup)` 方法
- **AND** 方法内注册所有使用统计相关路由

#### Scenario: ChatController 路由注册

- **WHEN** 查看 ChatController
- **THEN** 实现了 `RegisterRoutes(r *gin.RouterGroup)` 方法
- **AND** 方法内注册 Chat API 相关路由

#### Scenario: AuthController 路由注册

- **WHEN** 查看 AuthController（新增）
- **THEN** 实现了 `RegisterRoutes(r *gin.RouterGroup)` 方法
- **AND** 方法内注册登录、刷新、登出路由

### Requirement: main.go 路由组织

main.go 中的路由注册 SHALL 遵循统一的组织方式。

#### Scenario: 路由组结构

- **WHEN** 查看 main.go 中的路由设置
- **THEN** 按以下结构组织：
  1. 创建 `/api/v1` 路由组（管理接口）
  2. 在 `/api/v1` 下注册认证路由（无需鉴权）
  3. 创建需要 JWT 鉴权的子组
  4. 在 JWT 子组下创建需要 Admin 角色的子组
  5. 在相应组下调用各 Controller 的 `RegisterRoutes` 方法
  6. 创建 `/v1` 路由组（Chat API）
  7. 在 `/v1` 下注册需要 API Key 鉴权的路由

#### Scenario: 中间件应用顺序

- **WHEN** 查看路由中间件
- **THEN** 中间件按以下顺序应用：
  1. 日志/恢复中间件（全局）
  2. CORS 中间件（全局）
  3. JWT 鉴权中间件（管理 API）
  4. RequireAdmin 中间件（需要管理员权限的接口）
  5. API Key 鉴权中间件（Chat API）
  6. TraceID 中间件（Chat API）

## ADDED Requirements

### Requirement: API 版本管理

系统 SHALL 支持通过 URL 路径进行版本控制。

#### Scenario: 当前版本

- **WHEN** 访问任何 API
- **THEN** 所有路径包含版本号 `v1`
- **AND** 未来可以通过 `/api/v2` 引入不兼容的变更

### Requirement: 路由文档注释

每个路由处理函数 SHALL 包含文档注释说明其路径和方法。

#### Scenario: 路由注释格式

- **WHEN** 查看路由处理函数
- **THEN** 函数上方包含注释：
  ```go
  // CreateUser 创建用户
  // POST /api/v1/users
  func (c *UserController) CreateUser(ctx *gin.Context) {
      // ...
  }
  ```

## REMOVED Requirements

### Requirement: 移除 /v1 管理接口

原有的 `/v1/*` 管理接口路径 SHALL 被移除。

#### Scenario: 旧路径不再可用

- **GIVEN** 系统已完成迁移
- **WHEN** 访问旧的管理接口路径
  - `/v1/users`
  - `/v1/usage`
- **THEN** 返回 404 状态码
- **AND** 响应体包含错误信息：
  - `message`: "API endpoint not found"
  - `type`: "invalid_request_error"

### Requirement: 移除 main.go 硬编码路由

原有的在 main.go 中硬编码注册路由的方式 SHALL 被移除。

#### Scenario: 不再有直接路由注册

- **WHEN** 查看 main.go
- **THEN** 不出现类似 `userGroup.POST("", userCtrl.CreateUser)` 的直接路由注册
- **AND** 所有路由注册通过 Controller 的 `RegisterRoutes` 方法完成
