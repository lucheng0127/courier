# fallback-retry Specification

## Purpose
TBD - created by archiving change add-fallback-retry. Update Purpose after archive.
## Requirements
### Requirement: Fallback 配置

系统 SHALL 支持 Provider 配置 Fallback 模型列表。

#### Scenario: 配置 Fallback 模型列表

- **WHEN** 创建或更新 Provider 时
- **THEN** 可提供 `fallback_models` 字段（JSON 数组）
- **AND** 数组第一个元素为主模型，后续为备用模型
- **AND** 至少包含一个模型

#### Scenario: Fallback 模型列表格式

- **WHEN** 配置 `fallback_models` 时
- **THEN** 格式为 `["model-1", "model-2", "model-3"]`
- **AND** 所有模型必须在同一 Provider 内有效

#### Scenario: 默认行为

- **WHEN** 未配置 `fallback_models` 时
- **THEN** 使用请求中指定的模型
- **AND** 不进行 Fallback

### Requirement: Fallback 触发条件

系统 SHALL 在特定错误条件下触发 Fallback。

#### Scenario: 可重试错误触发 Fallback

- **WHEN** 主模型调用返回超时错误时
- **THEN** 触发 Fallback，尝试下一个备用模型
- **AND** **WHEN** 返回网络错误（连接失败、DNS 解析失败）时
- **THEN** 触发 Fallback
- **AND** **WHEN** 返回 5xx 服务器错误时
- **THEN** 触发 Fallback
- **AND** **WHEN** 返回连接拒绝时
- **THEN** 触发 Fallback

#### Scenario: 不可重试错误不触发 Fallback

- **WHEN** 返回 4xx 客户端错误时（除 429 外）
- **THEN** 不触发 Fallback，直接返回错误
- **AND** **WHEN** 返回认证失败时
- **THEN** 不触发 Fallback
- **AND** **WHEN** 模型不存在时
- **THEN** 不触发 Fallback

### Requirement: Fallback 执行流程

系统 SHALL 按顺序尝试 Fallback 模型列表。

#### Scenario: Fallback 成功

- **WHEN** 主模型失败，备用模型成功时
- **THEN** 返回备用模型的响应
- **AND** 记录 Fallback 发生
- **AND** 日志包含最终使用的模型名称

#### Scenario: Fallback 耗尽

- **WHEN** 所有模型（主模型 + 所有备用模型）都失败时
- **THEN** 返回最后一个错误
- **AND** 响应包含所有尝试的信息
- **AND** 记录 Fallback 耗尽日志

#### Scenario: 同 Provider 内 Fallback

- **WHEN** 请求模型为 `provider/model-a`
- **THEN** Fallback 仅在 `provider` 内进行
- **AND** 尝试顺序为：`provider/model-a` → `provider/model-b` → ...
- **AND** 不会切换到其他 Provider

### Requirement: 超时控制

系统 SHALL 为所有请求设置超时控制。

#### Scenario: 全局超时

- **WHEN** 请求处理时间超过配置的超时时间时
- **THEN** 取消正在进行的请求
- **AND** 返回超时错误
- **AND** 记录超时日志

#### Scenario: Provider 超时

- **WHEN** 单次 Provider 调用超过超时时间时
- **THEN** 取消该次调用
- **AND** 触发 Fallback（如果有备用模型）
- **OR** 返回超时错误

#### Scenario: 超时层级

- **WHEN** 配置了多层超时时
- **THEN** Fallback 总时间不超过全局超时
- **AND** 单次调用时间不超过 Provider 超时

### Requirement: TraceID 生成

系统 SHALL 为每个请求生成唯一的 TraceID。

#### Scenario: TraceID 生成

- **WHEN** 收到 HTTP 请求时
- **THEN** 生成唯一的 TraceID
- **AND** 格式为 `trace-<UUID>`
- **AND** 存储到 context.Context 中

#### Scenario: TraceID 响应 Header

- **WHEN** 返回 HTTP 响应时
- **THEN** 在响应 Header 中包含 `X-Trace-ID`
- **AND** 值为生成的 TraceID

#### Scenario: TraceID 唯一性

- **WHEN** 生成多个 TraceID 时
- **THEN** 每个 TraceID 唯一
- **AND** 使用 UUID 保证唯一性

### Requirement: TraceID 透传

系统 SHALL 在整个请求链路中透传 TraceID。

#### Scenario: TraceID 传递到 Provider

- **WHEN** 调用 Provider API 时
- **THEN** 在 HTTP Header 中传递 `X-Trace-ID`
- **AND** 值为当前请求的 TraceID

#### Scenario: TraceID 在日志中记录

- **WHEN** 记录日志时
- **THEN** 包含 `trace_id` 字段
- **AND** 所有相关日志共享同一个 `trace_id`

### Requirement: 统一结构化日志

系统 SHALL 输出结构化 JSON 格式的日志。

#### Scenario: 日志格式

- **WHEN** 输出日志时
- **THEN** 使用 JSON 格式
- **AND** 包含 `timestamp`、`level`、`trace_id`、`message` 等字段

#### Scenario: 日志级别

- **WHEN** 记录不同类型的日志时
- **THEN** 根据严重程度使用 `debug`、`info`、`warn`、`error` 级别
- **AND** Fallback 发生时使用 `warn` 级别

### Requirement: Fallback 状态跟踪

系统 SHALL 跟踪 Fallback 重试状态。

#### Scenario: Fallback 计数

- **WHEN** 发生 Fallback 时
- **THEN** 记录 `fallback_count`
- **AND** `fallback_count` 表示发生的 Fallback 次数

#### Scenario: 最终模型记录

- **WHEN** 请求完成后
- **THEN** 记录 `final_model`
- **AND** 表示最终成功响应的模型

#### Scenario: 尝试详情记录

- **WHEN** 发生 Fallback 时
- **THEN** 记录每次尝试的详细信息
- **AND** 包含模型名称、错误类型、错误消息

