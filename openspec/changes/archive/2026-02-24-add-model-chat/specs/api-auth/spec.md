## ADDED Requirements

### Requirement: API Key 认证

系统 MUST 支持 API Key 认证机制保护受保护的 API 端点。

#### Scenario: 使用有效的 API Key

- **WHEN** 客户端发送请求到受保护的端点，请求头包含 `Authorization: Bearer <valid_api_key>`
- **AND** 该 API Key 在数据库中存在且状态为 `active`
- **THEN** 系统返回对应的业务响应
- **AND** 该 API Key 的 `last_used_at` 被更新为当前时间

#### Scenario: 缺少 Authorization 头

- **WHEN** 客户端发送请求到受保护的端点，请求头不包含 `Authorization`
- **THEN** 系统返回 HTTP 401 状态码
- **AND** 返回错误信息提示未提供认证信息

#### Scenario: API Key 不存在

- **WHEN** 客户端发送请求到受保护的端点，请求头中的 API Key 在数据库中不存在
- **THEN** 系统返回 HTTP 401 状态码
- **AND** 返回错误信息提示 API Key 无效

#### Scenario: API Key 已禁用

- **WHEN** 客户端发送请求到受保护的端点，请求头中的 API Key 状态为 `disabled`
- **THEN** 系统返回 HTTP 401 状态码
- **AND** 返回错误信息提示 API Key 已禁用

#### Scenario: Authorization 格式错误

- **WHEN** 客户端发送请求到受保护的端点，`Authorization` 头格式不是 `Bearer <key>`
- **THEN** 系统返回 HTTP 401 状态码
- **AND** 返回错误信息提示认证格式无效

### Requirement: 用户上下文传递

系统 MUST 在认证成功后将用户信息传递给后续处理。

#### Scenario: 用户信息可用

- **WHEN** API Key 认证成功
- **THEN** 用户 ID 和用户信息被设置到请求上下文中
- **AND** 后续 Handler 可以从上下文中获取用户信息

### Requirement: 公开端点

系统 MUST 对某些端点不要求 API Key 认证。

#### Scenario: 访问公开端点

- **WHEN** 客户端发送请求到公开端点（如 `/api/v1/models`）
- **THEN** 系统不验证 API Key
- **AND** 直接返回对应的业务响应
