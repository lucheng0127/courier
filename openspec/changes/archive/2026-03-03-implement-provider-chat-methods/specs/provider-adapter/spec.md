## ADDED Requirements

### Requirement: OpenAI Adapter Chat 方法实现

OpenAI Adapter SHALL 实现完整的 `Chat()` 方法，支持调用 OpenAI 兼容的聊天完成 API。

#### Scenario: 非 SSE 流式调用

- **WHEN** 调用 OpenAI Adapter 的 `Chat()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/v1/chat/completions`
- **AND** 请求体格式符合 OpenAI API 规范
- **AND** 携带 `Authorization: Bearer {api_key}` 请求头（如果配置了 API Key）
- **AND** 使用 Provider 配置的超时时间
- **AND** 返回解析后的 `ChatResponse`

#### Scenario: 请求格式转换

- **WHEN** 构建请求体时
- **THEN** 将内部 `ChatRequest` 格式转换为 OpenAI API 格式
- **AND** 映射 `Messages` → `messages`
- **AND** 映射 `Model` → `model`
- **AND** 映射 `Temperature` → `temperature`（如果提供）
- **AND** 映射 `MaxTokens` → `max_tokens`（如果提供）

#### Scenario: 响应格式转换

- **WHEN** 收到 OpenAI API 响应时
- **THEN** 将 OpenAI 响应格式转换为内部 `ChatResponse` 格式
- **AND** 提取 `id`、`model`、`choices`、`usage` 字段
- **AND** 映射 `choices[0].message` → `Choices[0].Message`
- **AND** 映射 `usage` → `Usage`

#### Scenario: 默认参数支持

- **WHEN** Provider 配置的 `extra_config` 中包含默认参数时
- **THEN** 使用这些参数作为请求的默认值
- **AND** 请求级参数优先于默认参数
- **AND** 支持的默认参数包括：`max_tokens`、`temperature`、`top_p` 等

#### Scenario: 错误处理

- **WHEN** HTTP 请求失败时
- **THEN** 返回描述性错误
- **AND** 包含 HTTP 状态码和响应体
- **WHEN** 请求超时时
- **THEN** 返回超时错误
- **AND** 错误类型可被 RetryService 识别为可重试错误

### Requirement: OpenAI Adapter ChatStream 方法实现

OpenAI Adapter SHALL 实现完整的 `ChatStream()` 方法，支持 SSE 流式调用。

#### Scenario: SSE 流式调用

- **WHEN** 调用 OpenAI Adapter 的 `ChatStream()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/v1/chat/completions`
- **AND** 请求体中 `stream` 字段设置为 `true`
- **AND** 携带 `Authorization: Bearer {api_key}` 请求头（如果配置了 API Key）
- **AND** 使用 Provider 配置的超时时间
- **AND** 返回只读 channel 用于接收流式数据

#### Scenario: SSE 响应解析

- **WHEN** 收到 SSE 流式响应时
- **THEN** 逐行解析 SSE 数据
- **AND** 忽略以 `:` 开头的注释行
- **AND** 解析 `data:` 行中的 JSON 数据
- **AND** 将每个数据块转换为 `ChatStreamChunk` 格式
- **AND** 通过 channel 发送给调用方

#### Scenario: 流结束检测

- **WHEN** 收到 `data: [DONE]` 标记时
- **THEN** 关闭 channel
- **AND** 清理 HTTP 连接资源

#### Scenario: Context 取消

- **WHEN** context 被取消时
- **THEN** 中断 HTTP 请求
- **AND** 关闭 channel
- **AND** 清理所有相关资源

#### Scenario: 错误处理

- **WHEN** 流式请求失败时
- **THEN** 通过 channel 发送错误信息或直接返回错误
- **AND** 清理 HTTP 连接资源

### Requirement: vLLM Adapter Chat 方法实现

vLLM Adapter SHALL 实现完整的 `Chat()` 方法，复用 OpenAI 兼容协议。

#### Scenario: 非 SSE 流式调用

- **WHEN** 调用 vLLM Adapter 的 `Chat()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/v1/chat/completions`
- **AND** 请求体格式符合 OpenAI API 规范
- **AND** 如果配置了 API Key，携带 `Authorization: Bearer {api_key}` 请求头
- **AND** 如果未配置 API Key，不发送认证头
- **AND** 返回解析后的 `ChatResponse`

### Requirement: vLLM Adapter ChatStream 方法实现

vLLM Adapter SHALL 实现完整的 `ChatStream()` 方法，复用 OpenAI 兼容协议。

#### Scenario: SSE 流式调用

- **WHEN** 调用 vLLM Adapter 的 `ChatStream()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/v1/chat/completions`
- **AND** 请求体中 `stream` 字段设置为 `true`
- **AND** 如果配置了 API Key，携带认证头；否则不发送
- **AND** 返回只读 channel 用于接收流式数据

#### Scenario: 响应处理

- **WHEN** 收到 vLLM 的 SSE 流式响应时
- **THEN** 使用与 OpenAI 相同的解析逻辑
- **AND** 将数据块转换为 `ChatStreamChunk` 格式
- **AND** 通过 channel 发送给调用方

### Requirement: HTTP 客户端封装

Adapter SHALL 提供封装良好的 HTTP 客户端，支持认证、超时、取消等特性。

#### Scenario: 认证支持

- **WHEN** 配置了 API Key 时
- **THEN** 请求头包含 `Authorization: Bearer {api_key}`
- **WHEN** 未配置 API Key 时
- **THEN** 不发送认证头（支持本地模型服务）

#### Scenario: 超时控制

- **WHEN** 发送请求时
- **THEN** 使用 Provider 配置的超时时间
- **AND** 通过 `context.WithTimeout()` 设置请求超时
- **AND** 超时后取消请求并返回超时错误

#### Scenario: TraceID 透传

- **WHEN** 发送请求时
- **THEN** 如果 context 中包含 TraceID
- **AND** 在请求头中添加 `X-Trace-ID` 字段
- **AND** 值为 context 中的 TraceID

#### Scenario: Context 取消

- **WHEN** context 被取消时
- **THEN** 立即中断 HTTP 请求
- **AND** 返回 context 取消错误

### Requirement: OpenAI 兼容服务支持

系统 SHALL 支持配置 OpenAI 兼容的第三方服务（如 Qwen/通义千问）。

#### Scenario: 配置第三方服务

- **WHEN** 配置 OpenAI 类型的 Provider 时
- **THEN** `base_url` 可指向任何 OpenAI 兼容的服务端点
- **AND** `api_key` 使用该服务提供的密钥
- **AND** Adapter 正常工作，无需额外配置

#### Scenario: 示例配置

- **WHEN** 配置通义千问服务时
- **THEN** `type` 设置为 `openai`
- **AND** `base_url` 设置为 `https://dashscope.aliyuncs.com/compatible-mode/v1`
- **AND** `api_key` 设置为通义千问的 API Key
- **AND** 可在 `extra_config` 中设置默认模型参数
