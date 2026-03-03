# usage-tracking Specification

## Purpose
TBD - created by archiving change add-apikey-auth-and-usage. Update Purpose after archive.
## Requirements
### Requirement: 使用记录数据模型

系统 SHALL 提供使用记录数据模型用于追踪和统计 API 使用情况。

#### Scenario: 记录基本信息

- **WHEN** API 请求处理完成后
- **THEN** 记录请求的唯一 `request_id`
- **AND** 记录 `trace_id` 用于链路追踪
- **AND** 记录请求时间戳

#### Scenario: 关联用户和 API Key

- **WHEN** 记录使用信息时
- **THEN** 从请求上下文获取 `user_id`
- **AND** 从请求上下文获取 `api_key_id`
- **AND** 关联到对应的用户和 API Key

#### Scenario: 记录模型和 Provider 信息

- **WHEN** 记录使用信息时
- **THEN** 记录请求的 `model` 参数（如 `openai-main/gpt-4o`）
- **AND** 记录实际调用的 `provider_name`
- **AND** 便于按模型维度统计

#### Scenario: 记录 Token 使用量

- **WHEN** API 请求成功时
- **THEN** 记录 `prompt_tokens`（输入 Token 数）
- **AND** 记录 `completion_tokens`（输出 Token 数）
- **AND** 记录 `total_tokens`（总 Token 数）
- **AND** Token 数为 0 时也记录

#### Scenario: 记录性能指标

- **WHEN** 记录使用信息时
- **THEN** 记录请求耗时 `latency_ms`（毫秒）
- **AND** 记录请求状态 `status`（`success` 或 `error`）
- **AND** 失败时记录 `error_type`

#### Scenario: 级联删除行为

- **WHEN** 用户被删除时
- **THEN** 级联删除该用户的所有使用记录
- **AND** 当 API Key 被删除时
- **AND** 级联删除该 Key 关联的使用记录

### Requirement: 使用记录写入

系统 SHALL 在每次 API 请求处理完成后写入使用记录。

#### Scenario: 成功请求的记录

- **WHEN** Chat API 请求成功完成时
- **THEN** 从响应中提取 Token 使用量
- **AND** 结合上下文中的用户信息构造 UsageRecord
- **AND** 异步写入数据库
- **AND** 写入失败不影响主请求响应

#### Scenario: 失败请求的记录

- **WHEN** Chat API 请求失败时
- **THEN** 仍然创建使用记录
- **AND** Token 使用量记录为 0
- **AND** `status` 设置为 `error`
- **AND** 记录错误类型（如 `timeout`、`provider_error`）

#### Scenario: 流式响应的记录

- **WHEN** 处理流式响应请求时
- **THEN** 在流结束后收集完整的使用量
- **AND** 一次性写入使用记录
- **AND** 记录总耗时（包含流传输时间）

#### Scenario: 异步写入机制

- **WHEN** 写入使用记录时
- **THEN** 使用后台 goroutine 或 channel 批量写入
- **AND** 避免阻塞主请求流程
- **AND** 写入失败记录错误日志

### Requirement: 使用统计查询 API

系统 SHALL 提供按用户和时间范围查询使用统计的 API。

#### Scenario: 查询参数验证

- **WHEN** GET 请求到 `/v1/usage` 时
- **THEN** 必须提供 `user_id` 查询参数
- **AND** 可选提供 `start_date`（默认 30 天前）
- **AND** 可选提供 `end_date`（默认今天）
- **AND** 可选提供 `group_by`（`day` 或 `model`，默认 `day`）

#### Scenario: 按天聚合统计

- **WHEN** `group_by=day` 或未指定时
- **THEN** 返回每天的使用汇总
- **AND** 包含每天的总请求数、Token 数、平均耗时
- **AND** 返回整体统计汇总

#### Scenario: 按模型聚合统计

- **WHEN** `group_by=model` 时
- **THEN** 返回每个模型的使用汇总
- **AND** 包含每个模型的请求数、Token 数、平均耗时
- **AND** 按使用量降序排列

#### Scenario: 统计响应格式

- **WHEN** 返回统计结果时
- **THEN** 包含 `user_id` 和查询的时间范围
- **AND** 包含 `summary` 汇总信息
- **AND** 包含 `daily_breakdown` 或 `model_breakdown` 详细数据

#### Scenario: 空结果处理

- **WHEN** 指定时间范围内没有使用记录时
- **THEN** 返回空数组
- **AND** 汇总数据为 0

### Requirement: Chat 控制器集成

系统 SHALL 在 Chat 控制器中集成使用量记录功能。

#### Scenario: 从上下文获取用户信息

- **WHEN** Chat 请求处理时
- **THEN** 从 Gin Context 获取 `user_id`
- **AND** 从 Gin Context 获取 `api_key_id`
- **AND** 用于关联使用记录

#### Scenario: 记录使用量

- **WHEN** 请求处理完成时（成功或失败）
- **THEN** 调用 UsageService.RecordUsage 方法
- **AND** 传递完整的上下文和使用量信息
- **AND** 使用独立的 context（避免请求取消影响写入）

### Requirement: 管理员鉴权

使用统计查询接口 SHALL 验证管理员权限。

#### Scenario: 管理员 API Key 验证

- **WHEN** 调用 `/v1/usage` 接口时
- **THEN** 验证 `X-Admin-API-Key` Header
- **AND** 验证失败返回 401 Unauthorized

#### Scenario: 跨用户查询权限

- **WHEN** 管理员查询任意用户的使用统计时
- **THEN** 允许查询任何 `user_id` 的数据
- **AND** 普通用户（未来实现）只能查询自己的数据

