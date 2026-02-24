# api-key-management Specification

## Purpose
TBD - created by archiving change init-user-management. Update Purpose after archive.
## Requirements
### Requirement: API Key 生成

系统 MUST 支持为指定用户生成新的 API Key。

#### Scenario: 成功生成 API Key

- **WHEN** 客户端发送 POST 请求到 `/api/v1/users/:id/apikeys`，`id` 为已存在的用户 ID
- **THEN** 系统返回 HTTP 201 状态码
- **AND** 返回 API Key 信息，包含 `id`、`key`、`user_id`、`status`（默认为 `active`）、`created_at`
- **AND** `key` 格式为 `ck_` 后跟 24 字节的十六进制编码随机字符串
- **AND** API Key 数据被持久化到 SQLite 数据库

#### Scenario: 用户不存在

- **WHEN** 客户端发送 POST 请求到 `/api/v1/users/:id/apikeys`，`id` 为不存在的用户 ID
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示用户不存在

#### Scenario: 无效的用户 ID

- **WHEN** 客户端发送 POST 请求到 `/api/v1/users/:id/apikeys`，`id` 为无效格式（如非数字）
- **THEN** 系统返回 HTTP 400 状态码
- **AND** 返回错误信息提示参数无效

### Requirement: API Key 列表查询

系统 MUST 支持查询指定用户的所有 API Key。

#### Scenario: 成功查询 API Key 列表

- **WHEN** 客户端发送 GET 请求到 `/api/v1/users/:id/apikeys`，`id` 为已存在的用户 ID
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 返回该用户的 API Key 列表数组，包含 `id`、`key`、`status`、`last_used_at`、`created_at`、`updated_at`
- **AND** 返回空数组当该用户没有 API Key 时

#### Scenario: 用户不存在

- **WHEN** 客户端发送 GET 请求到 `/api/v1/users/:id/apikeys`，`id` 为不存在的用户 ID
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示用户不存在

#### Scenario: 无效的用户 ID

- **WHEN** 客户端发送 GET 请求到 `/api/v1/users/:id/apikeys`，`id` 为无效格式（如非数字）
- **THEN** 系统返回 HTTP 400 状态码
- **AND** 返回错误信息提示参数无效

### Requirement: API Key 删除

系统 MUST 支持删除指定的 API Key。

#### Scenario: 成功删除 API Key

- **WHEN** 客户端发送 DELETE 请求到 `/api/v1/users/:id/apikeys/:keyid`，`id` 为已存在的用户 ID，`keyid` 为该用户的 API Key ID
- **THEN** 系统返回 HTTP 204 状态码
- **AND** 该 API Key 从数据库中被物理删除

#### Scenario: 用户不存在

- **WHEN** 客户端发送 DELETE 请求到 `/api/v1/users/:id/apikeys/:keyid`，`id` 为不存在的用户 ID
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示用户不存在

#### Scenario: API Key 不存在

- **WHEN** 客户端发送 DELETE 请求到 `/api/v1/users/:id/apikeys/:keyid`，`id` 为已存在的用户 ID，`keyid` 为不存在的 API Key ID
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示 API Key 不存在

#### Scenario: API Key 不属于该用户

- **WHEN** 客户端发送 DELETE 请求到 `/api/v1/users/:id/apikeys/:keyid`，`keyid` 存在但不属于用户 `id`
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示 API Key 不存在

#### Scenario: 无效的 ID 参数

- **WHEN** 客户端发送 DELETE 请求到 `/api/v1/users/:id/apikeys/:keyid`，`id` 或 `keyid` 为无效格式（如非数字）
- **THEN** 系统返回 HTTP 400 状态码
- **AND** 返回错误信息提示参数无效

### Requirement: API Key 禁用

系统 MUST 支持禁用指定的 API Key。

#### Scenario: 成功禁用 API Key

- **WHEN** 客户端发送 PUT 请求到 `/api/v1/users/:id/apikeys/:keyid/disable`，`id` 为已存在的用户 ID，`keyid` 为该用户的 API Key ID
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 返回更新后的 API Key 信息，`status` 字段值为 `disabled`
- **AND** 数据库中该 API Key 的 `status` 被更新为 `disabled`

#### Scenario: 用户不存在

- **WHEN** 客户端发送 PUT 请求到 `/api/v1/users/:id/apikeys/:keyid/disable`，`id` 为不存在的用户 ID
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示用户不存在

#### Scenario: API Key 不存在

- **WHEN** 客户端发送 PUT 请求到 `/api/v1/users/:id/apikeys/:keyid/disable`，`id` 为已存在的用户 ID，`keyid` 为不存在的 API Key ID
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示 API Key 不存在

#### Scenario: API Key 不属于该用户

- **WHEN** 客户端发送 PUT 请求到 `/api/v1/users/:id/apikeys/:keyid/disable`，`keyid` 存在但不属于用户 `id`
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示 API Key 不存在

#### Scenario: 禁用已禁用的 API Key

- **WHEN** 客户端发送 PUT 请求到 `/api/v1/users/:id/apikeys/:keyid/disable`，`keyid` 为状态已是 `disabled` 的 API Key ID
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 返回 API Key 信息，`status` 仍为 `disabled`

#### Scenario: 无效的 ID 参数

- **WHEN** 客户端发送 PUT 请求到 `/api/v1/users/:id/apikeys/:keyid/disable`，`id` 或 `keyid` 为无效格式（如非数字）
- **THEN** 系统返回 HTTP 400 状态码
- **AND** 返回错误信息提示参数无效

