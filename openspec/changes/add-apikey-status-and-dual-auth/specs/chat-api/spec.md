# chat-api 规范变更

## MODIFIED Requirements

### Requirement: API Key 鉴权

系统 SHALL 支持 API Key 和 JWT 两种鉴权方式访问 Chat API。

#### Scenario: 使用 API Key 鉴权

- **WHEN** 客户端发送 API 请求时
- **THEN** 必须在 Header 中提供 `Authorization: Bearer <api_key>`
- **AND** API Key 格式为 `sk-` 开头
- **AND** 系统验证 API Key 的有效性和状态
- **AND** 在上下文中注入认证类型 `auth_type=apikey`

#### Scenario: 使用 JWT Token 鉴权（新增）

- **WHEN** 客户端发送 API 请求时
- **THEN** 可以在 Header 中提供 `Authorization: Bearer <jwt_token>`
- **AND** JWT Token 必须是有效的 Access Token
- **AND** 系统验证 Token 的签名和过期时间
- **AND** 系统验证关联的用户状态为 `active`
- **AND** 在上下文中注入认证类型 `auth_type=jwt`

#### Scenario: 两种鉴权方式都失败

- **WHEN** 客户端提供的 Authorization Header 既不是有效的 API Key 也不是有效的 JWT Token
- **THEN** 返回 401 Unauthorized
- **AND** 响应体包含 `{"error": {"message": "Authentication failed. Please provide a valid API key or JWT token.", "type": "authentication_error"}}`

#### Scenario: 缺少 Authorization Header

- **WHEN** 客户端未提供 Authorization Header
- **THEN** 返回 401 Unauthorized
- **AND** 响应体包含 `{"error": {"message": "Missing authorization header", "type": "authentication_error"}}`

### Requirement: 鉴权上下文注入

系统 SHALL 在鉴权成功后向请求上下文注入用户信息和认证类型。

#### Scenario: API Key 鉴权成功后的上下文注入

- **WHEN** API Key 鉴权成功后
- **THEN** 在上下文中设置：
  - `user_id`: 用户 ID
  - `user_email`: 用户邮箱
  - `api_key_id`: API Key 记录 ID
  - `api_key_masked`: 脱敏后的 API Key
  - `auth_type`: 认证类型（值为 `"apikey"`）

#### Scenario: JWT 鉴权成功后的上下文注入

- **WHEN** JWT 鉴权成功后
- **THEN** 在上下文中设置：
  - `user_id`: 用户 ID
  - `user_email`: 用户邮箱
  - `user_role`: 用户角色
  - `auth_type`: 认证类型（值为 `"jwt"`）

### Requirement: 使用量记录认证类型

系统 SHALL 在使用量记录中区分不同的认证类型。

#### Scenario: 记录 API Key 认证的使用量

- **WHEN** 使用 API Key 认证成功处理 Chat 请求
- **THEN** 使用量记录包含：
  - `user_id`: 用户 ID
  - `api_key_id`: API Key ID
  - `auth_type`: 认证类型（值为 `"apikey"`）

#### Scenario: 记录 JWT 认证的使用量

- **WHEN** 使用 JWT 认证成功处理 Chat 请求
- **THEN** 使用量记录包含：
  - `user_id`: 用户 ID
  - `api_key_id`: NULL（不关联 API Key）
  - `auth_type`: 认证类型（值为 `"jwt"`）

#### Scenario: 日志中包含认证类型

- **WHEN** 记录 Chat 请求日志时
- **THEN** 日志包含认证类型字段
- **AND** 方便追踪和分析不同认证方式的使用情况

## MODIFIED Requirements

### Requirement: 请求日志

系统 SHALL 记录所有 Chat API 请求的日志。

#### Scenario: 日志内容（扩展）

- **WHEN** 处理 Chat API 请求时
- **THEN** 记录：
  - 请求 ID、API Key（脱敏）、认证类型
  - 模型名称、Provider
  - Token 使用量、耗时、状态
- **AND** 日志格式为 JSON
- **AND** 新增 `auth_type` 字段区分认证方式

#### Scenario: 错误日志

(无变更，保持原样)
