# role-based-access Specification

## Purpose
TBD - created by archiving change unify-api-and-jwt-auth. Update Purpose after archive.
## Requirements
### Requirement: 用户角色

系统 SHALL 支持用户角色区分，包含 `admin` 和 `user` 两种角色。

#### Scenario: Admin 角色定义

- **GIVEN** 用户角色为 `admin`
- **THEN** 该用户可以访问所有管理接口
- **AND** 可以读取、更新、删除其他用户
- **AND** 可以管理 Providers（创建、读取、更新、删除、重载、启用、禁用）
- **AND** 可以查询所有用户的使用统计
- **AND** 可以为任何用户创建 API Key

#### Scenario: User 角色定义

- **GIVEN** 用户角色为 `user`
- **THEN** 该用户只能访问自己的资源
- **AND** 可以使用自己的 API Key 调用 Chat API
- **AND** 可以查询自己的使用统计
- **AND** 可以查看自己的用户信息
- **AND** 不能访问管理接口
- **AND** 不能管理其他用户

### Requirement: 数据库角色字段

系统 SHALL 在 `users` 表中添加 `role` 字段。

#### Scenario: 角色字段定义

- **WHEN** 查看 `users` 表结构
- **THEN** 包含 `role` 字段
- **AND** 字段类型为 `VARCHAR(20)`
- **AND** 字段非空
- **AND** 默认值为 `'user'`
- **AND** 包含检查约束 `CHECK (role IN ('user', 'admin'))`

#### Scenario: 角色字段索引

- **WHEN** 查看 `users` 表索引
- **THEN** 包含 `idx_users_role` 索引
- **AND** 索引字段为 `role`

### Requirement: RequireAdmin 中间件

系统 SHALL 提供 `RequireAdmin` 中间件，用于验证用户是否具有管理员权限。

#### Scenario: 管理员用户访问

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `admin`
- **WHEN** 访问使用了 `RequireAdmin` 中间件的接口
- **THEN** 中间件从上下文获取 `user_role`
- **AND** 验证角色为 `admin`
- **AND** 请求继续传递到后续处理函数

#### Scenario: 非管理员用户访问

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 访问使用了 `RequireAdmin` 中间件的接口
- **THEN** 中间件从上下文获取 `user_role`
- **AND** 检测到角色不是 `admin`
- **AND** 返回 403 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Admin privileges required"
  - `type`: "permission_error"
- **AND** 终止请求处理

#### Scenario: 未通过 JWT 鉴权

- **GIVEN** 用户未通过 JWT 鉴权（缺少用户信息）
- **WHEN** 访问使用了 `RequireAdmin` 中间件的接口
- **THEN** 中间件检测到上下文中缺少 `user_role`
- **AND** 返回 401 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Authentication required"
  - `type`: "authentication_error"

### Requirement: Provider 管理权限

系统 SHALL 允许管理员用户管理 Provider，所有认证用户可查询 Provider 列表和模型列表。

#### Scenario: 普通用户查询 Provider 列表

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 发送 GET 请求到 `/api/v1/providers`
- **THEN** 返回 200 状态码
- **AND** 响应体只包含非敏感信息：
  - `name`: Provider 名称
  - `type`: Provider 类型
  - `base_url`: API 地址
  - `enabled`: 启用状态
  - `fallback_models`: 支持的模型列表
- **AND** 响应体不包含敏感信息：
  - `api_key`: API 密钥
  - `timeout`: 超时配置
  - `is_running`: 运行状态
  - `extra_config`: 额外配置

#### Scenario: 管理员用户查询 Provider 列表

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `admin`
- **WHEN** 发送 GET 请求到 `/api/v1/providers`
- **THEN** 返回 200 状态码
- **AND** 响应体包含完整的 Provider 信息（包括敏感信息）

#### Scenario: 普通用户查询 Provider 模型列表

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 发送 GET 请求到 `/api/v1/providers/:name/models`
- **THEN** 返回 200 状态码
- **AND** 响应体包含 Provider 名称、类型和模型列表

### Requirement: 用户管理权限

系统 SHALL 只允许管理员用户查看和管理用户账户状态，用户通过自主注册创建。

#### Scenario: 查看用户列表

- **GIVEN** 用户已登录
- **AND** 用户角色为 `admin`
- **WHEN** 发送 GET 请求到 `/api/v1/users`
- **THEN** 返回所有用户的列表

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 GET 请求到 `/api/v1/users`
- **THEN** 返回 403 状态码

#### Scenario: 管理员查看任意用户信息

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 GET 请求到 `/api/v1/users/:id`
- **THEN** 可以查看任何用户的信息

#### Scenario: 普通用户查看自己的信息

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/123`
- **THEN** 返回该用户的信息

#### Scenario: 普通用户尝试查看他人信息

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/456`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

#### Scenario: 更新用户信息

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 PUT 请求到 `/api/v1/users/:id`
- **THEN** 可以更新用户信息

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 PUT 请求到 `/api/v1/users/:id`
- **THEN** 返回 403 状态码

#### Scenario: 删除用户

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 DELETE 请求到 `/api/v1/users/:id`
- **THEN** 可以删除用户

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 DELETE 请求到 `/api/v1/users/:id`
- **THEN** 返回 403 状态码

#### Scenario: 更新用户状态

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 PATCH 请求到 `/api/v1/users/:id/status`
- **THEN** 可以更新用户状态（active/disabled）

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 PATCH 请求到 `/api/v1/users/:id/status`
- **THEN** 返回 403 状态码

#### Scenario: 管理员为用户创建 API Key

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 POST 请求到 `/api/v1/users/:id/api-keys`
- **THEN** 可以为任何用户创建 API Key

#### Scenario: 普通用户为自己创建 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 POST 请求到 `/api/v1/users/123/api-keys`
- **THEN** 可以创建 API Key

#### Scenario: 普通用户尝试为他人创建 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 POST 请求到 `/api/v1/users/456/api-keys`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

#### Scenario: 普通用户查询自己的 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/123/api-keys`
- **THEN** 返回该用户的 API Key 列表

#### Scenario: 普通用户尝试查询他人 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/456/api-keys`
- **THEN** 返回 403 状态码

#### Scenario: 普通用户撤销自己的 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 DELETE 请求到 `/api/v1/users/123/api-keys/:key_id`
- **THEN** 可以撤销该 API Key

#### Scenario: 普通用户尝试撤销他人 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 DELETE 请求到 `/api/v1/users/456/api-keys/:key_id`
- **THEN** 返回 403 状态码

### Requirement: 使用统计查询权限

系统 SHALL 允许管理员查询所有用户的使用统计，普通用户查询自己的使用统计。

#### Scenario: 管理员查询任意用户统计

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 GET 请求到 `/api/v1/usage?user_id=123`
- **THEN** 可以查询任意用户的使用统计
- **AND** 可以不传 `user_id` 查询所有用户的统计数据

#### Scenario: 普通用户查询自己的统计

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/usage`
- **THEN** 返回该用户自己的使用统计
- **AND** 自动过滤 `user_id=123` 的数据
- **AND** 如果尝试传递其他 `user_id` 参数，忽略该参数

#### Scenario: 普通用户查询所有统计

- **GIVEN** 普通用户已登录
- **WHEN** 发送 GET 请求到 `/api/v1/usage`（不带 user_id）
- **THEN** 只返回该用户自己的使用统计
- **AND** 不返回其他用户的数据

### Requirement: 初始管理员用户

系统 SHALL 在首次部署时创建初始管理员用户。

#### Scenario: 初始管理员创建

- **GIVEN** 数据库初始化时
- **WHEN** 系统首次启动
- **THEN** 检查是否存在任何管理员用户
- **AND** 如果不存在，则创建默认管理员账户
- **AND** 默认管理员邮箱从环境变量 `INITIAL_ADMIN_EMAIL` 读取
- **AND** 默认管理员密码从环境变量 `INITIAL_ADMIN_PASSWORD` 读取
- **AND** 如果环境变量未设置，在日志中警告并跳过创建

