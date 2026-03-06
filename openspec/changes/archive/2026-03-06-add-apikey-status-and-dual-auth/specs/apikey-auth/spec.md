# apikey-auth 规范变更

## ADDED Requirements

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

## MODIFIED Requirements

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
