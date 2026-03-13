## MODIFIED Requirements

### Requirement: 使用记录数据模型

系统 SHALL 提供使用记录数据模型用于追踪和统计 API 使用情况，同时支持按套餐维度统计。

#### Scenario: 关联套餐使用记录

- **WHEN** 记录使用信息时
- **AND** 请求使用了套餐配额
- **THEN** 记录关联的 `user_package_id`
- **AND** 记录关联的 `package_id`
- **AND** 便于按套餐维度统计使用情况

#### Scenario: 原有使用记录字段保持不变

- **WHEN** 记录使用信息时
- **THEN** 保持原有的所有字段不变
- **AND** 新增字段为可选（`user_package_id` 和 `package_id`）

### Requirement: 使用统计查询 API

系统 SHALL 提供按用户和时间范围查询使用统计的 API，同时支持按套餐维度聚合。

#### Scenario: 按套餐聚合统计

- **WHEN** `group_by=package` 时
- **THEN** 返回每个套餐的使用汇总
- **AND** 包含每个套餐的请求数、Token 数、平均耗时
- **AND** 按使用量降序排列

#### Scenario: 套餐使用详情统计

- **WHEN** 查询指定套餐的使用统计时
- **THEN** 请求参数包含 `package_id`
- **AND** 返回该套餐的详细使用情况
- **AND** 包含按 Provider 分组的统计
- **AND** 包含按天分组的趋势

#### Scenario: 统计响应格式扩展

- **WHEN** 返回统计结果时
- **THEN** 原有字段保持不变
- **AND** 新增可选字段 `package_breakdown`（按套餐聚合时）
- **AND** 保持向后兼容性

### Requirement: 管理员套餐使用统计查询

系统 SHALL 提供管理员查询套餐使用统计的 API。

#### Scenario: 管理员查询所有套餐使用统计

- **GIVEN** 用户具有管理员权限
- **WHEN** GET 请求到 `/api/v1/admin/packages/usage`
- **THEN** 返回所有套餐的使用统计汇总
- **AND** 包含每个套餐的购买人数、激活人数、总使用 Token 数
- **AND** 支持按时间范围筛选

#### Scenario: 管理员查询指定套餐使用统计

- **GIVEN** 用户具有管理员权限
- **WHEN** GET 请求到 `/api/v1/admin/packages/:id/usage`
- **THEN** 返回指定套餐的详细使用统计
- **AND** 包含购买用户列表
- **AND** 包含每个用户的使用量
- **AND** 包含按 Provider 分组的使用统计

#### Scenario: 管理员查询套餐使用趋势

- **GIVEN** 用户具有管理员权限
- **WHEN** GET 请求到 `/api/v1/admin/packages/:id/trend`
- **THEN** 返回套餐的使用趋势
- **AND** 包含按天的购买量
- **AND** 包含按天的使用量
- **AND** 支持指定时间范围

### Requirement: 用户套餐使用统计导出

系统 SHALL 支持用户导出自己的套餐使用记录。

#### Scenario: 导出套餐使用记录

- **WHEN** GET 请求到 `/api/v1/user/packages/:id/usage/export`
- **THEN** 验证用户已登录
- **AND** 验证套餐属于当前用户
- **AND** 返回 CSV 格式的使用记录
- **AND** 包含请求时间、Provider、使用的 Token 数、请求 ID

#### Scenario: 导出记录限制

- **WHEN** 导出使用记录时
- **THEN** 最多导出最近 90 天的记录
- **AND** 最多导出 10000 条记录
- **AND** 超出限制时提示用户调整查询范围
