# request-logging Specification

## Purpose
TBD - created by archiving change add-model-chat. Update Purpose after archive.
## Requirements
### Requirement: 请求日志记录

系统 MUST 记录所有模型对话请求的完整信息。

#### Scenario: 成功请求的日志记录

- **WHEN** 模型对话请求成功完成
- **THEN** 系统创建一条请求日志记录
- **AND** 日志包含：`user_id`、`model_name`、`request_messages`（JSON 格式）、`response_content`、`prompt_tokens`、`completion_tokens`、`total_tokens`、`latency_ms`、`status`（值为 `success`）、`created_at`

#### Scenario: 失败请求的日志记录

- **WHEN** 模型对话请求失败（上游错误、超时等）
- **THEN** 系统创建一条请求日志记录
- **AND** 日志包含：`user_id`、`model_name`、`request_messages`、`error_message`、`status`（值为 `error`）、`created_at`

#### Scenario: 流式请求的日志记录

- **WHEN** 模型对话请求为流式响应
- **THEN** 系统缓存完整的响应内容
- **AND** 流结束后创建一条请求日志记录
- **AND** 日志包含完整的 `response_content`

#### Scenario: 中断请求的日志记录

- **WHEN** 流式响应被客户端中断
- **THEN** 系统创建一条请求日志记录
- **AND** 日志包含已接收的部分响应内容
- **AND** `status` 记录为 `interrupted`
- **AND** `error_message` 记录中断原因

### Requirement: 异步日志写入

系统 MUST 使用异步方式写入请求日志，不阻塞主请求。

#### Scenario: 异步写入

- **WHEN** 请求完成需要记录日志
- **THEN** 日志写入操作在单独的 goroutine 中执行
- **AND** 主请求不等待日志写入完成即可返回

#### Scenario: 写入失败不影响请求

- **WHEN** 日志写入操作失败
- **THEN** 不影响已返回给客户端的响应
- **AND** 错误被记录到系统日志中

### Requirement: 延迟统计

系统 MUST 记录请求的响应延迟。

#### Scenario: 记录延迟

- **WHEN** 模型对话请求完成
- **THEN** 系统计算从请求开始到响应结束的时间差
- **AND** 将延迟（毫秒）记录到 `latency_ms` 字段

### Requirement: 敏感信息保护

系统 MUST 在日志中保护敏感信息。

#### Scenario: 不记录上游 API Key

- **WHEN** 记录请求日志
- **THEN** 日志中不包含上游模型的 API Key
- **AND** 日志中不包含用户的 API Key

#### Scenario: 不记录请求头

- **WHEN** 记录请求日志
- **THEN** 日志中不包含原始 HTTP 请求头
- **AND** 仅记录业务相关的 `messages` 内容

### Requirement: 日志数据持久化

系统 MUST 将请求日志持久化到 SQLite 数据库。

#### Scenario: 数据库写入

- **WHEN** 异步日志写入执行时
- **THEN** 日志数据被写入 `request_logs` 表
- **AND** 写入失败时错误被记录到系统日志

### Requirement: Token 统计

系统 MUST 记录请求的 Token 使用量。

#### Scenario: 记录 Token 使用量

- **WHEN** 上游模型响应包含 Token 使用信息
- **THEN** `prompt_tokens`、`completion_tokens`、`total_tokens` 被记录到日志
- **AND** 当上游不提供这些信息时，记录为 0

