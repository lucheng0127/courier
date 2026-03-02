# apikey-auth Specification

## Purpose
TBD - created by archiving change add-apikey-auth-and-usage. Update Purpose after archive.
## Requirements
### Requirement: 用户数据模型

系统 SHALL 提供用户数据模型用于管理 API Key 的所有者。

#### Scenario: 创建用户

- **WHEN** 管理员通过 API 创建用户时
- **THEN** 系统生成唯一的用户 ID
- **AND** 存储用户的姓名、邮箱和状态
- **AND** 默认状态为 `active`
- **AND** 记录创建时间

#### Scenario: 用户状态管理

- **WHEN** 用户被禁用时
- **THEN** 系统将用户状态设置为 `disabled`
- **AND** 该用户的所有 API Key 无法通过鉴权

#### Scenario: 用户唯一性约束

- **WHEN** 创建用户时邮箱已存在
- **THEN** 返回 409 Conflict 错误
- **AND** 错误信息说明邮箱已被使用

### Requirement: API Key 数据模型

系统 SHALL 提供用于 API 访问鉴权的 API Key 数据模型。

#### Scenario: API Key 格式

- **WHEN** 创建新 API Key 时
- **THEN** API Key 格式为 `sk-<32位随机字符>`
- **AND** 使用加密安全的随机数生成器
- **AND** 仅在创建响应中返回完整 Key 一次

#### Scenario: API Key 存储

- **WHEN** 存储 API Key 时
- **THEN** 使用 SHA256 哈希存储完整 Key
- **AND** 存储前 10 位作为 `key_prefix` 用于识别
- **AND** 哈希值在数据库中唯一

#### Scenario: API Key 关联用户

- **WHEN** 创建 API Key 时
- **THEN** 必须指定所属用户 ID
- **AND** 一个用户可以拥有多个 API Key
- **AND** 删除用户时级联删除其所有 API Key

#### Scenario: API Key 生命周期状态

- **WHEN** API Key 状态变更时
- **THEN** 状态可以是 `active`、`disabled`、`revoked`
- **AND** 只有 `active` 状态的 Key 能通过鉴权
- **AND** 支持通过 API 禁用或撤销 Key

#### Scenario: API Key 过期时间

- **WHEN** 创建 API Key 时可选设置过期时间
- **THEN** 过期的 Key 无法通过鉴权
- **AND** 返回明确的过期错误信息

### Requirement: API Key 鉴权中间件

系统 SHALL 提供中间件从 API Key 解析用户信息并注入到请求上下文。

#### Scenario: 从请求头提取 API Key

- **WHEN** 客户端发送请求时
- **THEN** 从 `Authorization: Bearer <api_key>` Header 提取 API Key
- **AND** 验证 Bearer token 格式

#### Scenario: 验证 API Key 有效性

- **WHEN** API Key 提取后
- **THEN** 从数据库查询对应的 Key 记录
- **AND** 验证 Key 状态为 `active`
- **AND** 验证 Key 未过期（如果设置了过期时间）

#### Scenario: 获取关联用户信息

- **WHEN** API Key 验证通过后
- **THEN** 根据 Key 关联的用户 ID 查询用户信息
- **AND** 验证用户状态为 `active`
- **AND** 用户被禁用时返回 403 Forbidden

#### Scenario: 注入用户信息到 Context

- **WHEN** API Key 和用户验证都通过后
- **THEN** 在 Gin Context 中设置 `user_id`（用户 ID）
- **AND** 设置 `user_email`（用户邮箱）
- **AND** 设置 `api_key_id`（API Key 记录 ID）
- **AND** 设置 `api_key_masked`（脱敏后的 Key，格式 `sk-cour...xxxx`）

#### Scenario: 更新最后使用时间

- **WHEN** API Key 验证通过后
- **THEN** 异步更新该 Key 的 `last_used_at` 字段
- **AND** 更新失败不影响主请求流程

#### Scenario: 鉴权失败响应

- **WHEN** API Key 无效、已禁用或已过期时
- **THEN** 返回 401 Unauthorized
- **AND** 响应格式为 `{"error": "错误描述"}`
- **AND** 不泄露系统内部信息

### Requirement: 用户管理 API

系统 SHALL 提供用户管理的 API 端点。

#### Scenario: 创建用户

- **WHEN** POST 请求到 `/v1/users` 时
- **THEN** 验证 `X-Admin-API-Key` Header
- **AND** 请求体包含 `name` 和 `email` 字段
- **AND** 成功时返回 201 Created 和用户信息

#### Scenario: 获取用户信息

- **WHEN** GET 请求到 `/v1/users/:id` 时
- **THEN** 验证管理员权限
- **AND** 返回指定用户的完整信息

### Requirement: API Key 管理 API

系统 SHALL 提供 API Key 管理的 API 端点。

#### Scenario: 为用户创建 API Key

- **WHEN** POST 请求到 `/v1/users/:id/api-keys` 时
- **THEN** 验证用户存在且状态为 `active`
- **AND** 请求体包含 `name` 字段（用户定义的 Key 名称）
- **AND** 可选包含 `expires_at` 字段（过期时间）
- **AND** 返回完整的 API Key（仅此一次）

#### Scenario: 获取用户的 API Key 列表

- **WHEN** GET 请求到 `/v1/users/:id/api-keys` 时
- **THEN** 返回该用户所有 API Key
- **AND** 不返回完整 Key，只返回 `key_prefix` 和元数据
- **AND** 包含状态、最后使用时间、创建时间等信息

#### Scenario: 删除/禁用 API Key

- **WHEN** DELETE 请求到 `/v1/users/:id/api-keys/:key_id` 时
- **THEN** 将 API Key 状态设置为 `revoked`
- **AND** 返回 204 No Content
- **AND** 该 Key 立即无法通过鉴权

### Requirement: 管理员鉴权

系统 SHALL 对用户和 API Key 管理接口进行管理员鉴权。

#### Scenario: 管理员 API Key 验证

- **WHEN** 调用管理接口时
- **THEN** 从 `X-Admin-API-Key` Header 获取管理员 Key
- **AND** 验证 Key 有效性
- **AND** 验证失败返回 401 Unauthorized

#### Scenario: 管理接口权限控制

- **WHEN** 管理员 Key 验证通过后
- **THEN** 允许访问用户和 API Key 管理接口
- **AND** 允许访问使用统计接口

