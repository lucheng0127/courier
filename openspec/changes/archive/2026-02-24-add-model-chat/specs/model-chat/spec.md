## ADDED Requirements

### Requirement: 发起模型对话

系统 MUST 支持通过 API 发起与上游模型的对话。

#### Scenario: 成功发起非流式对话

- **WHEN** 客户端发送 POST 请求到 `/api/v1/models/:model/chat`
- **AND** 请求头包含有效的 API Key
- **AND** 请求体包含 `messages` 数组，`stream` 为 `false` 或不设置
- **AND** `:model` 为已配置的模型名称
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 返回包含 AI 回复的 JSON 响应
- **AND** 响应格式兼容 OpenAI Chat Completions API
- **AND** 记录请求日志

#### Scenario: 成功发起流式对话

- **WHEN** 客户端发送 POST 请求到 `/api/v1/models/:model/chat`
- **AND** 请求头包含有效的 API Key
- **AND** 请求体包含 `messages` 数组，`stream` 为 `true`
- **AND** `:model` 为已配置的模型名称
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 响应 Content-Type 为 `text/event-stream`
- **AND** 以 Server-Sent Events 格式流式返回 AI 回复
- **AND** 每个数据块格式为 `data: <json>\n\n`
- **AND** 流结束时发送 `data: [DONE]\n\n`
- **AND** 记录请求日志

#### Scenario: 模型不存在

- **WHEN** 客户端请求的 `:model` 不在配置的模型列表中
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示模型不存在

#### Scenario: 未认证

- **WHEN** 客户端发送请求时没有提供 API Key 或 API Key 无效
- **THEN** 系统返回 HTTP 401 状态码
- **AND** 不调用上游模型 API

#### Scenario: 请求体无效

- **WHEN** 客户端发送请求时 `messages` 字段缺失或为空数组
- **THEN** 系统返回 HTTP 400 状态码
- **AND** 返回错误信息提示请求参数无效

#### Scenario: 上游模型错误

- **WHEN** 上游模型返回错误响应（如 API Key 无效、配额用尽等）
- **THEN** 系统返回 HTTP 502 或 503 状态码
- **AND** 记录请求日志，状态为 `error`
- **AND** 返回错误信息（不暴露上游敏感信息）

### Requirement: 多轮对话支持

系统 MUST 支持多轮对话上下文。

#### Scenario: 多轮对话

- **WHEN** 客户端在 `messages` 数组中提供历史对话
- **AND** `messages` 包含多条消息，交替使用 `user` 和 `assistant` 角色
- **THEN** 系统将完整的 `messages` 数组发送给上游模型
- **AND** 上游模型基于完整上下文生成回复

### Requirement: 流式响应中断处理

系统 MUST 正确处理客户端断开连接的情况。

#### Scenario: 客户端断开连接

- **WHEN** 流式响应进行中客户端断开连接
- **THEN** 系统取消向上游模型的请求
- **AND** 不再继续转发上游响应
- **AND** 记录请求日志，状态标记为中断

### Requirement: 请求超时

系统 MUST 对上游模型请求设置超时时间。

#### Scenario: 请求超时

- **WHEN** 上游模型在 30 秒内未返回响应
- **THEN** 系统取消请求
- **AND** 返回 HTTP 504 状态码
- **AND** 记录请求日志，状态为 `error`

### Requirement: Token 统计

系统 MUST 记录请求和响应的 Token 使用量。

#### Scenario: 上游返回 Token 信息

- **WHEN** 上游模型响应包含 `usage` 字段（`prompt_tokens`、`completion_tokens`、`total_tokens`）
- **THEN** 系统将这些信息记录到请求日志中

#### Scenario: 上游不返回 Token 信息

- **WHEN** 上游模型响应不包含 `usage` 字段
- **THEN** 系统在请求日志中将 Token 字段记录为 0
- **AND** 请求仍然成功处理
