## ADDED Requirements

### Requirement: 用户创建

系统 MUST 支持通过 HTTP API 创建新用户。

#### Scenario: 成功创建用户

- **WHEN** 客户端发送 POST 请求到 `/api/v1/users`，请求体包含有效的 `name` 和可选的 `email`
- **THEN** 系统返回 HTTP 201 状态码
- **AND** 返回用户信息，包含自动生成的 `id`、`name`、`email`、`created_at`、`updated_at`
- **AND** 用户数据被持久化到 SQLite 数据库

#### Scenario: 用户名重复

- **WHEN** 客户端发送 POST 请求到 `/api/v1/users`，请求体包含已存在的 `name`
- **THEN** 系统返回 HTTP 409 状态码
- **AND** 返回错误信息提示用户名已存在

#### Scenario: 请求体无效

- **WHEN** 客户端发送 POST 请求到 `/api/v1/users`，请求体缺少 `name` 字段
- **THEN** 系统返回 HTTP 400 状态码
- **AND** 返回错误信息提示请求参数无效

### Requirement: 用户列表查询

系统 MUST 支持查询所有用户列表。

#### Scenario: 成功查询用户列表

- **WHEN** 客户端发送 GET 请求到 `/api/v1/users`
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 返回用户列表数组，包含所有用户的 `id`、`name`、`email`、`created_at`、`updated_at`
- **AND** 返回空数组当没有用户时

### Requirement: 用户详情查询

系统 MUST 支持通过用户 ID 查询单个用户详情。

#### Scenario: 成功查询用户详情

- **WHEN** 客户端发送 GET 请求到 `/api/v1/users/:id`，`id` 为已存在的用户 ID
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 返回用户详细信息，包含 `id`、`name`、`email`、`created_at`、`updated_at`

#### Scenario: 用户不存在

- **WHEN** 客户端发送 GET 请求到 `/api/v1/users/:id`，`id` 为不存在的用户 ID
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示用户不存在

#### Scenario: 无效的用户 ID

- **WHEN** 客户端发送 GET 请求到 `/api/v1/users/:id`，`id` 为无效格式（如非数字）
- **THEN** 系统返回 HTTP 400 状态码
- **AND** 返回错误信息提示参数无效
