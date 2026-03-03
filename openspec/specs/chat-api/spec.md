# chat-api Specification

## Purpose
TBD - created by archiving change add-chat-api. Update Purpose after archive.
## Requirements
### Requirement: Chat Completions API

系统 SHALL 提供兼容 OpenAI 格式的 `/v1/chat/completions` 端点。

#### Scenario: 非流式请求

- **WHEN** 客户端发送 POST 请求到 `/v1/chat/completions` 且 `stream=false`
- **THEN** 系统返回完整的聊天响应
- **AND** 响应格式符合 OpenAI 规范
- **AND** 包含 `id`、`object`、`created`、`model`、`choices`、`usage` 字段

#### Scenario: 流式请求

- **WHEN** 客户端发送 POST 请求到 `/v1/chat/completions` 且 `stream=true`
- **THEN** 系统返回 SSE (Server-Sent Events) 流式响应
- **AND** 每个数据块以 `data:` 开头
- **AND** 流结束时发送 `data: [DONE]`
- **AND** 响应 Header 包含 `Content-Type: text/event-stream`

#### Scenario: 请求验证

- **WHEN** 请求缺少必填字段（`model` 或 `messages`）
- **THEN** 返回 400 错误
- **AND** 错误信息说明缺少的字段

#### Scenario: 模型未找到

- **WHEN** 请求中的 `model` 参数没有对应的 Provider
- **THEN** 返回 404 错误
- **AND** 错误信息说明模型不可用

### Requirement: API Key 鉴权

系统 SHALL 要求所有 API 请求通过 API Key 鉴权。

#### Scenario: 鉴权 Header 格式

- **WHEN** 客户端发送 API 请求时
- **THEN** 必须在 Header 中提供 `Authorization: Bearer <api_key>`
- **AND** API Key 格式为 `sk-` 开头

#### Scenario: 鉴权失败

- **WHEN** 请求未提供 API Key 或 API Key 无效
- **THEN** 返回 401 Unauthorized
- **AND** 响应体包含 `{"error": {"message": "Invalid API key", "type": "invalid_request_error"}}`

#### Scenario: 鉴权成功

- **WHEN** API Key 验证通过
- **THEN** 请求正常处理
- **AND** 在日志中记录 API Key（脱敏后）

### Requirement: 模型路由

系统 SHALL 根据 `model` 参数将请求路由到对应的 Provider。

#### Scenario: 模型格式

- **WHEN** 请求指定 `model` 参数时
- **THEN** 模型格式必须为 `provider/model_name`
- **AND** `provider` 为 Provider 实例名称
- **AND** `model_name` 为 Provider 端定义的模型名称

#### Scenario: 模型解析和路由

- **WHEN** 请求中的 `model` 为 `openai-main/gpt-4o` 时
- **THEN** 系统解析 `provider` 为 `openai-main`、`model_name` 为 `gpt-4o`
- **AND** 查找名为 `openai-main` 的 Provider 实例
- **AND** 验证 Provider 是否启用
- **AND** 调用该 Provider 的 Chat 方法，传入 `gpt-4o` 作为模型

#### Scenario: 模型格式错误

- **WHEN** 请求中的 `model` 不包含 `/` 分隔符
- **THEN** 返回 400 错误
- **AND** 错误信息说明正确的格式为 `provider/model_name`

#### Scenario: Provider 不存在

- **WHEN** 解析的 `provider` 名称对应的 Provider 实例不存在
- **THEN** 返回 404 错误
- **AND** 错误信息说明 Provider 不存在

#### Scenario: Provider 未启用

- **WHEN** 解析的 `provider` 对应的 Provider 未启用
- **THEN** 返回 403 错误
- **AND** 错误信息说明 Provider 未启用

### Requirement: 请求日志

系统 SHALL 记录所有 Chat API 请求的日志。

#### Scenario: 日志内容

- **WHEN** 处理 Chat API 请求时
- **THEN** 记录请求 ID、API Key（脱敏）、模型名称、Provider、Token 使用量、耗时、状态
- **AND** 日志格式为 JSON

#### Scenario: 错误日志

- **WHEN** 请求失败时
- **THEN** 记录错误信息和堆栈
- **AND** 不记录敏感用户数据

### Requirement: 流式响应取消

系统 SHALL 支持客户端取消流式响应。

#### Scenario: 客户端断开连接

- **WHEN** 客户端在流式响应过程中断开连接
- **THEN** 系统检测到断开
- **AND** 取消上游 Provider 调用
- **AND** 清理相关资源

#### Scenario: Context 超时

- **WHEN** 请求超过配置的超时时间
- **THEN** 取消流式响应
- **AND** 返回超时错误

### Requirement: 响应格式转换

系统 SHALL 将 Provider 响应转换为 OpenAI 格式。

#### Scenario: 非流式响应转换

- **WHEN** Provider 返回非流式响应时
- **THEN** 转换为 OpenAI 格式
- **AND** 生成唯一的 `id`（格式 `chatcmpl-<UUID>`）
- **AND** `object` 字段为 `chat.completion`
- **AND** `created` 字段为 Unix 时间戳

#### Scenario: 流式响应转换

- **WHEN** Provider 返回流式数据块时
- **THEN** 转换为 OpenAI SSE 格式
- **AND** 每个 chunk 包含 `id`、`object`、`created`、`model`、`choices`

### Requirement: 错误处理

系统 SHALL 统一处理和返回错误。

#### Scenario: Provider 错误

- **WHEN** Provider 调用失败时
- **THEN** 返回适当的 HTTP 状态码
- **AND** 错误响应符合 OpenAI 格式
- **AND** 包含错误类型和消息

#### Scenario: 超时错误

- **WHEN** Provider 调用超时时
- **THEN** 返回 408 超时错误
- **AND** 记录超时日志

