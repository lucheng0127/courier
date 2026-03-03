## ADDED Requirements

### Requirement: Provider 更新

系统 SHALL 支持通过 API 更新 Provider 配置。

#### Scenario: 更新 Provider 配置

- **WHEN** 管理员调用 PUT `/api/v1/providers/:name` 更新 Provider 时
- **THEN** 系统验证 Provider 存在
- **AND** 更新数据库中的 Provider 配置
- **AND** 如果 Provider 已启用，重载使其生效
- **AND** 返回 200 OK 和更新后的配置

#### Scenario: 更新时重载 Provider

- **WHEN** Provider 已启用且配置被更新时
- **THEN** 从数据库重新加载配置
- **AND** 销毁旧的 Adapter 实例
- **AND** 使用新配置创建新的 Adapter 实例
- **AND** 更新 Registry 中的实例引用

#### Scenario: 禁用状态更新为启用

- **WHEN** Provider 从 `enabled=false` 更新为 `enabled=true` 时
- **THEN** 从数据库重新加载配置
- **AND** 创建新的 Adapter 实例并注册到 Registry
- **AND** 后续请求可正常使用该 Provider

#### Scenario: 启用状态更新为禁用

- **WHEN** Provider 从 `enabled=true` 更新为 `enabled=false` 时
- **THEN** 从 Registry 中移除该 Provider 实例
- **AND** 销毁 Adapter 实例释放资源
- **AND** 后续请求无法使用该 Provider

#### Scenario: 更新不存在的 Provider

- **WHEN** 尝试更新不存在的 Provider 时
- **THEN** 返回 404 Not Found
- **AND** 错误信息说明 Provider 不存在

#### Scenario: 部分更新参数

- **WHEN** 请求体仅包含部分参数时
- **THEN** 仅更新提供的字段
- **AND** 未提供的字段保持原值

### Requirement: Provider 删除

系统 SHALL 支持通过 API 删除 Provider。

#### Scenario: 删除 Provider

- **WHEN** 管理员调用 DELETE `/api/v1/providers/:name` 删除 Provider 时
- **THEN** 系统验证 Provider 存在
- **AND** 如果 Provider 正在运行，先从 Registry 注销
- **AND** 从数据库删除 Provider 配置
- **AND** 返回 204 No Content

#### Scenario: 删除运行中的 Provider

- **WHEN** 删除正在运行的 Provider 时
- **THEN** 先从 Registry 中注销 Provider
- **AND** 销毁 Adapter 实例释放资源
- **AND** 然后从数据库删除配置

#### Scenario: 删除不存在的 Provider

- **WHEN** 尝试删除不存在的 Provider 时
- **THEN** 返回 404 Not Found
- **AND** 错误信息说明 Provider 不存在

#### Scenario: 删除失败的 Provider

- **WHEN** Provider 处于错误状态时
- **THEN** 仍然允许删除
- **AND** 清理所有相关资源

### Requirement: Provider 更新参数验证

系统 SHALL 验证 Provider 更新请求的参数。

#### Scenario: 必填字段验证

- **WHEN** 更新请求包含 `type`、`base_url` 或 `timeout` 时
- **THEN** 这些字段必须同时提供（binding "required_with" 规则）
- **AND** 如果验证失败返回 400 Bad Request

#### Scenario: 超时参数验证

- **WHEN** 提供 `timeout` 参数时
- **THEN** 值必须大于等于 1
- **AND** 如果验证失败返回 400 Bad Request

#### Scenario: 可选字段处理

- **WHEN** 请求体不包含某个字段时
- **THEN** 该字段保持原值不变
- **AND** `enabled` 字段通过指针类型支持显式设置 false
