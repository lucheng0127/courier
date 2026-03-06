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

(无变更，保持原样)

#### Scenario: 获取用户的 API Key 列表

(无变更，保持原样)

#### Scenario: 删除/禁用 API Key

此场景已被新增的需求替代和扩展。原有的 "删除/禁用" 实际上是 "撤销"（设置状态为 `revoked`）。

- **WHEN** DELETE 请求到 `/v1/users/:id/api-keys/:key_id` 时（旧行为）
- **THEN** 将 API Key 状态设置为 `revoked`
- **AND** 该 Key 立即无法通过鉴权

**注意**：此场景保留向后兼容，但新的删除接口将执行真正的删除操作。

**修改后的行为**：
- 现有的 DELETE 接口行为保持不变（软删除，状态设置为 `revoked`）
- 新增的删除接口将执行硬删除（从数据库删除记录）
- API 路由可能需要调整以区分两种操作

**设计决策**：为避免破坏现有兼容性，建议：
- 保留现有 `DELETE /api/v1/users/:id/api-keys/:key_id` 接口的行为（撤销）
- 新增操作通过新的路由实现（启用/禁用）
- 如需真正的硬删除，可在后续版本中通过新的路由或参数实现

**修正**：根据用户需求，DELETE 接口应该执行真正的删除。但这会影响现有客户端。

**最终决策**：
- 为了向后兼容，保留现有 RevokeAPIKey 接口和路由
- 新增启用/禁用接口
- 新增真正的删除接口，可以通过查询参数或新路由实现

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

### Requirement: API Key 启用接口

系统 SHALL 提供 API Key 启用接口，允许用户或管理员重新激活被禁用的 API Key。

#### Scenario: 启用被禁用的 API Key

- **GIVEN** 用户有一个状态为 `disabled` 的 API Key
- **WHEN** 用户发送 PATCH 请求到 `/api/v1/users/:id/api-keys/:key_id/enable`
- **THEN** 系统验证请求者是 Key 的所有者或管理员
- **AND** 系统将 API Key 状态更新为 `active`
- **AND** 返回 200 OK 和更新后的 API Key 信息
- **AND** 该 Key 可以正常通过鉴权

#### Scenario: 启用已激活的 API Key

- **GIVEN** 用户有一个状态为 `active` 的 API Key
- **WHEN** 用户尝试启用该 API Key
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明 API Key 已是激活状态

#### Scenario: 无权限启用 API Key

- **GIVEN** 用户 A 有一个 API Key
- **WHEN** 用户 B 尝试启用用户 A 的 API Key（用户 B 不是管理员）
- **THEN** 返回 403 Forbidden
- **AND** 错误信息说明权限不足

### Requirement: API Key 禁用接口

系统 SHALL 提供 API Key 禁用接口，允许用户或管理员临时禁用 API Key。

#### Scenario: 禁用激活的 API Key

- **GIVEN** 用户有一个状态为 `active` 的 API Key
- **WHEN** 用户发送 PATCH 请求到 `/api/v1/users/:id/api-keys/:key_id/disable`
- **THEN** 系统验证请求者是 Key 的所有者或管理员
- **AND** 系统将 API Key 状态更新为 `disabled`
- **AND** 返回 200 OK 和更新后的 API Key 信息
- **AND** 该 Key 无法通过鉴权

#### Scenario: 禁用已禁用的 API Key

- **GIVEN** 用户有一个状态为 `disabled` 的 API Key
- **WHEN** 用户尝试禁用该 API Key
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明 API Key 已是禁用状态

#### Scenario: 禁用已撤销的 API Key

- **GIVEN** 用户有一个状态为 `revoked` 的 API Key
- **WHEN** 用户尝试禁用该 API Key
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明已撤销的 API Key 无法操作

### Requirement: API Key 删除接口

系统 SHALL 提供 API Key 删除接口，允许用户或管理员永久删除 API Key 记录。

#### Scenario: 删除 API Key

- **GIVEN** 用户有一个任意状态的 API Key
- **WHEN** 用户发送 DELETE 请求到 `/api/v1/users/:id/api-keys/:key_id`
- **THEN** 系统验证请求者是 Key 的所有者或管理员
- **AND** 系统从数据库中删除该 API Key 记录
- **AND** 返回 204 No Content
- **AND** 该 API Key 永久无法恢复

#### Scenario: 删除不存在的 API Key

- **GIVEN** 用户有一个不存在的 API Key ID
- **WHEN** 用户尝试删除该 API Key
- **THEN** 返回 404 Not Found
- **AND** 错误信息说明 API Key 不存在

#### Scenario: 无权限删除 API Key

- **GIVEN** 用户 A 有一个 API Key
- **WHEN** 用户 B 尝试删除用户 A 的 API Key（用户 B 不是管理员）
- **THEN** 返回 403 Forbidden
- **AND** 错误信息说明权限不足

#### Scenario: 删除操作与撤销操作的区别

- **GIVEN** 系统提供两种删除方式
- **WHEN** 用户调用现有接口 `DELETE /api/v1/users/:id/api-keys/:key_id`
- **THEN** 系统执行硬删除（从数据库彻底删除记录）
- **AND** 不同于原有的 RevokeAPIKey 接口（将状态设置为 `revoked`）
- **AND** API 文档明确说明此操作不可恢复

