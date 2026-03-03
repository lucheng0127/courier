# user-registration Specification

## Purpose
TBD - created by archiving change add-user-registration. Update Purpose after archive.
## Requirements
### Requirement: 用户注册接口

系统 SHALL 提供用户注册接口，允许用户自主创建账户。

#### Scenario: 成功注册

- **GIVEN** 系统正常运行
- **AND** 提供的邮箱尚未被注册
- **AND** 密码符合强度要求
- **WHEN** 发送 POST 请求到 `/api/v1/auth/register`
- **AND** 请求体包含：
  - `name`: 用户名称
  - `email`: 邮箱地址
  - `password`: 密码
- **THEN** 返回 201 状态码
- **AND** 响应体包含创建的用户信息：
  - `id`: 用户 ID
  - `name`: 用户名称
  - `email`: 邮箱地址
  - `role`: 固定为 "user"
  - `status`: 固定为 "active"
  - `created_at`: 创建时间
- **AND** 密码被哈希存储
- **AND** 响应体不包含密码或密码哈希

#### Scenario: 邮箱已存在

- **GIVEN** 系统中已存在使用相同邮箱的用户
- **WHEN** 发送 POST 请求到 `/api/v1/auth/register`
- **AND** 请求体中的邮箱与已有用户相同
- **THEN** 返回 409 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Email already exists"
  - `type`: "invalid_request_error"
- **AND** 不创建新用户

#### Scenario: 请求参数无效

- **GIVEN** 系统正常运行
- **WHEN** 发送 POST 请求到 `/api/v1/auth/register`
- **AND** 请求体缺少必填字段或格式不正确
- **THEN** 返回 400 状态码
- **AND** 响应体包含错误信息：
  - `message`: 描述具体错误（如 "Name is required"）
  - `type`: "invalid_request_error"

### Requirement: 注册接口无需鉴权

用户注册接口 SHALL 不需要 JWT 鉴权或 API Key 鉴权。

#### Scenario: 未登录用户访问注册接口

- **GIVEN** 用户未登录
- **AND** 用户未提供任何认证信息
- **WHEN** 发送 POST 请求到 `/api/v1/auth/register`
- **THEN** 请求被正常处理
- **AND** 不返回 401 或 403 状态码

### Requirement: 密码强度要求

系统 SHALL 对注册密码进行强度验证。

#### Scenario: 密码过短

- **GIVEN** 系统正常运行
- **WHEN** 发送 POST 请求到 `/api/v1/auth/register`
- **AND** 密码少于 8 个字符
- **THEN** 返回 400 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Password must be at least 8 characters"
  - `type`: "invalid_request_error"

#### Scenario: 密码符合要求

- **GIVEN** 系统正常运行
- **WHEN** 发送 POST 请求到 `/api/v1/auth/register`
- **AND** 密码至少 8 个字符
- **THEN** 注册成功
- **AND** 密码被使用 bcrypt 哈希存储

### Requirement: 注册速率限制

系统 SHALL 对注册接口实施速率限制，防止恶意批量注册。

#### Scenario: 正常速率注册

- **GIVEN** IP 地址在 1 小时内注册次数少于 5 次
- **WHEN** 发送 POST 请求到 `/api/v1/auth/register`
- **THEN** 请求被正常处理

#### Scenario: 超过速率限制

- **GIVEN** 同一 IP 地址在 1 小时内已尝试注册 5 次
- **WHEN** 再次发送 POST 请求到 `/api/v1/auth/register`
- **THEN** 返回 429 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Too many registration attempts, please try again later"
  - `type`: "rate_limit_error"

### Requirement: 新用户默认角色和状态

新注册的用户 SHALL 具有默认的角色和状态。

#### Scenario: 默认角色为 user

- **GIVEN** 用户成功注册
- **WHEN** 查询创建的用户信息
- **THEN** `role` 字段为 "user"
- **AND** 用户不具有管理员权限

#### Scenario: 默认状态为 active

- **GIVEN** 用户成功注册
- **WHEN** 查询创建的用户信息
- **THEN** `status` 字段为 "active"
- **AND** 用户可以立即使用登录接口进行认证

### Requirement: 注册后可登录

新注册的用户 SHALL 能够使用注册时设置的邮箱和密码登录。

#### Scenario: 注册后登录成功

- **GIVEN** 用户成功注册
- **AND** 注册时使用邮箱 `test@example.com` 和密码 `password123`
- **WHEN** 发送 POST 请求到 `/api/v1/auth/login`
- **AND** 请求体包含相同的邮箱和密码
- **THEN** 返回 200 状态码
- **AND** 响应体包含 JWT 访问令牌和刷新令牌

### Requirement: 密码哈希存储

系统 SHALL 使用安全的哈希算法存储注册用户密码。

#### Scenario: 密码哈希

- **WHEN** 用户注册时提供密码
- **THEN** 使用 bcrypt 算法进行哈希
- **AND** 使用 cost factor 12
- **AND** 只存储哈希值，不存储明文密码
- **AND** 数据库 `password_hash` 字段包含哈希值

### Requirement: 与现有认证体系兼容

用户注册功能 SHALL 与现有的 JWT 认证和角色权限控制体系完全兼容。

#### Scenario: 注册用户使用 API Key

- **GIVEN** 用户已注册并登录
- **AND** 用户创建了 API Key
- **WHEN** 使用该 API Key 调用 Chat API
- **THEN** 请求成功处理
- **AND** 用户身份正确识别

#### Scenario: 注册用户查询自己的信息

- **GIVEN** 用户已注册并登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/123`
- **THEN** 返回该用户的信息
- **AND** 不返回其他用户的信息

#### Scenario: 注册用户不能访问管理员接口

- **GIVEN** 用户已注册并登录
- **AND** 用户角色为 "user"
- **WHEN** 发送 POST 请求到 `/api/v1/providers`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

### Requirement: 移除管理员创建用户功能

系统 SHALL 移除管理员创建普通用户的功能。

#### Scenario: POST /api/v1/users 接口不存在

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 POST 请求到 `/api/v1/users`
- **THEN** 返回 404 或 405 状态码
- **AND** 不创建新用户

#### Scenario: 保留用户管理功能

- **GIVEN** 管理员用户已登录
- **WHEN** 查看用户管理接口列表
- **THEN** 仍保留以下功能：
  - `GET /api/v1/users` - 查看用户列表
  - `GET /api/v1/users/:id` - 查看用户信息
  - `PUT /api/v1/users/:id` - 更新用户信息
  - `DELETE /api/v1/users/:id` - 删除用户
  - `PATCH /api/v1/users/:id/status` - 更新用户状态

### Requirement: 管理员用户创建方式

系统 SHALL 仅支持特定方式创建管理员用户。

#### Scenario: 初始管理员

- **GIVEN** 数据库为空
- **AND** 环境变量 `INITIAL_ADMIN_EMAIL` 和 `INITIAL_ADMIN_PASSWORD` 已设置
- **WHEN** 系统首次启动
- **THEN** 创建初始管理员用户
- **AND** 该用户角色为 "admin"

#### Scenario: 用户角色升级

- **GIVEN** 系统已运行
- **AND** 存在普通用户
- **WHEN** 需要将普通用户升级为管理员
- **THEN** 通过数据库直接操作更新用户角色
- **AND** 此操作仅限运维场景

### Requirement: 请求和响应模型

系统 SHALL 定义明确的注册请求和响应模型。

#### Scenario: RegisterRequest 模型

- **WHEN** 定义注册请求模型
- **THEN** 包含以下字段：
  - `name`: string，必填，用户名称
  - `email`: string，必填，邮箱格式，唯一
  - `password`: string，必填，最少 8 个字符

#### Scenario: RegisterResponse 模型

- **WHEN** 定义注册响应模型
- **THEN** 包含以下字段：
  - `id`: int64，用户 ID
  - `name`: string，用户名称
  - `email`: string，邮箱地址
  - `role`: string，固定为 "user"
  - `status`: string，固定为 "active"
  - `created_at`: time，创建时间
  - 不包含 `password` 或 `password_hash` 字段

