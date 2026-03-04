# provider-adapter Specification

## Purpose
TBD - created by archiving change add-provider-adapter. Update Purpose after archive.
## Requirements
### Requirement: Provider Interface

系统 SHALL 定义统一的 Provider 接口，用于抽象不同 LLM 供应商的调用方式。

#### Scenario: 接口定义

- **WHEN** 定义 Provider 接口时
- **THEN** 接口必须包含 `Chat()` 方法用于非流式调用
- **AND** 接口必须包含 `ChatStream()` 方法用于流式调用
- **AND** 接口必须包含 `Type()` 方法返回 Provider 类型
- **AND** 接口必须包含 `Name()` 方法返回 Provider 实例名称

#### Scenario: 流式响应

- **WHEN** 调用 `ChatStream()` 方法时
- **THEN** 返回只读 channel 用于接收流式数据块
- **AND** 支持 context 取消机制

### Requirement: Adapter Registry

系统 SHALL 提供 Adapter 注册机制，支持动态扩展新的 Provider 类型。

#### Scenario: Adapter 类型注册

- **WHEN** 新的 Provider 类型实现后
- **THEN** 可通过 `RegisterAdapterType()` 注册工厂函数
- **AND** 注册时指定 Provider 类型标识（如 "openai", "vllm"）

#### Scenario: Adapter 实例化

- **WHEN** 给定 Provider 配置时
- **THEN** 系统根据配置中的 Type 查找对应的工厂函数
- **AND** 调用工厂函数创建 Provider 实例
- **AND** 如果 Type 未注册则返回错误

#### Scenario: 自动注册

- **WHEN** 程序启动时
- **THEN** 各 Adapter 通过 `init()` 函数自动注册到 Registry
- **AND** 无需手动调用注册代码

### Requirement: Provider 配置存储

系统 SHALL 将 Provider 配置存储在数据库中，支持持久化和查询。

#### Scenario: 配置必填字段

- **WHEN** 创建 Provider 配置时
- **THEN** 必须提供 `name` 作为唯一标识
- **AND** 必须提供 `type` 指定 Provider 类型
- **AND** 必须提供 `base_url` 指定 API 地址
- **AND** 必须提供 `timeout` 指定超时时间（秒）
- **AND** 必须提供 `enabled` 指定启用状态（布尔值，默认 true）

#### Scenario: 配置可选字段

- **WHEN** 创建 SaaS 类型 Provider 时
- **THEN** 应提供 `api_key` 用于身份认证
- **WHEN** 创建本地模型 Provider 时
- **THEN** `api_key` 可为空
- **WHEN** Provider 需要额外配置时
- **THEN** 可提供 `extra_config` 存储扩展参数（JSON 格式）

#### Scenario: 配置唯一性

- **WHEN** 创建 Provider 配置时
- **THEN** `name` 字段必须在全局唯一
- **AND** 如果重复则返回唯一性约束错误

### Requirement: Provider 初始化

系统启动时 SHALL 加载所有 Provider 配置并完成 Adapter 初始化。

#### Scenario: 启动加载

- **WHEN** 系统启动时
- **THEN** 从数据库加载所有 Provider 配置
- **AND** 对于 `enabled` 为 true 的 Provider，根据配置创建对应的 Provider 实例
- **AND** 对于 `enabled` 为 false 的 Provider，跳过初始化
- **AND** 将已初始化的实例注册到全局 Provider Registry

#### Scenario: 初始化失败处理

- **WHEN** 单个 Provider 初始化失败时
- **THEN** 记录错误日志
- **AND** 不中断系统启动
- **AND** 继续初始化其他 Provider
- **AND** 该 Provider 不可用，修复配置后可通过重载 API 使其生效

#### Scenario: 未知 Provider 类型

- **WHEN** 配置中的 `type` 未有对应 Adapter 注册时
- **THEN** 记录错误日志
- **AND** 跳过该 Provider 配置

### Requirement: Provider 状态查询

系统 SHALL 提供 API 端点查询当前可用的 Provider 列表和状态。

#### Scenario: 查询 Provider 列表

- **WHEN** 调用 Provider 列表 API 时
- **THEN** 返回所有已配置的 Provider
- **AND** 包含每个 Provider 的名称、类型、启用状态（enabled/disabled）、运行状态（可用/不可用）

#### Scenario: 查询单个 Provider

- **WHEN** 按 name 查询 Provider 时
- **THEN** 返回该 Provider 的详细配置
- **AND** 如果不存在则返回 404

### Requirement: Provider 运行时重载

系统 SHALL 支持运行时重载 Provider 配置，无需重启服务。

#### Scenario: 重载单个 Provider

- **WHEN** 管理员调用重载 API 并指定 Provider name 时
- **THEN** 从数据库重新加载该 Provider 的最新配置
- **AND** 销毁旧的 Adapter 实例
- **AND** 使用新配置创建新的 Adapter 实例
- **AND** 更新 Registry 中的实例引用

#### Scenario: 重载所有 Provider

- **WHEN** 管理员调用重载所有 Provider API 时
- **THEN** 从数据库重新加载所有 Provider 配置
- **AND** 销毁所有旧的 Adapter 实例
- **AND** 使用新配置创建所有新的 Adapter 实例
- **AND** 更新 Registry 中的所有实例引用

#### Scenario: 重载失败处理

- **WHEN** Provider 重载失败时
- **THEN** 记录错误日志
- **AND** 保持旧实例继续运行
- **AND** 返回错误信息给调用方

#### Scenario: 新增 Provider 后重载

- **WHEN** 管理员通过 API 新增 Provider 配置后
- **THEN** 可立即调用重载 API 使新 Provider 生效
- **AND** 无需重启服务

### Requirement: Provider 启用/禁用

系统 SHALL 支持运行时启用或禁用 Provider。

#### Scenario: 禁用 Provider

- **WHEN** 管理员调用禁用 API 并指定 Provider name 时
- **THEN** 更新数据库中该 Provider 的 `enabled` 字段为 false
- **AND** 从 Registry 中移除该 Provider 实例
- **AND** 销毁 Adapter 实例释放资源
- **AND** 后续请求无法使用该 Provider

#### Scenario: 启用 Provider

- **WHEN** 管理员调用启用 API 并指定 Provider name 时
- **THEN** 更新数据库中该 Provider 的 `enabled` 字段为 true
- **AND** 从数据库重新加载该 Provider 的配置
- **AND** 创建新的 Adapter 实例并注册到 Registry
- **AND** 后续请求可正常使用该 Provider

#### Scenario: 启用/禁用失败处理

- **WHEN** 启用或禁用操作失败时
- **THEN** 记录错误日志
- **AND** 返回错误信息给调用方
- **AND** 保持原有状态不变

### Requirement: 第三方 SaaS 模型支持

系统 SHALL 支持需要 API Key 的第三方 SaaS 模型服务。

#### Scenario: API Key 认证

- **WHEN** 调用 SaaS Provider 时
- **THEN** 使用配置中的 `api_key` 进行身份认证
- **AND** API Key 通过 HTTP Header 传递（如 `Authorization: Bearer <key>`）

#### Scenario: API Key 安全存储

- **WHEN** 存储 API Key 时
- **THEN** 应考虑加密存储（MVP 阶段可使用明文，后续升级）

### Requirement: 私有本地模型支持

系统 SHALL 支持私有部署的本地模型，API Key 可选。

#### Scenario: 无 API Key 调用

- **WHEN** 配置中未提供 `api_key` 时
- **THEN** Adapter 不发送认证头
- **AND** 正常调用模型 API

#### Scenario: 内网地址支持

- **WHEN** `base_url` 为内网地址时
- **THEN** 系统应正常处理
- **AND** 支持 HTTP 和 HTTPS 协议

### Requirement: 扩展配置支持

系统 SHALL 支持 Provider 特定的扩展配置。

#### Scenario: 额外参数传递

- **WHEN** Provider 需要额外配置时
- **THEN** 可通过 `extra_config` 字段传递 JSON 格式的配置
- **AND** Adapter 负责解析和验证这些配置

#### Scenario: 常见扩展配置

- **WHEN** 配置 OpenAI 兼容接口时
- **THEN** `extra_config` 可支持 `model` 字段指定模型名称
- **AND** 可支持 `organization` 字段指定组织 ID
- **WHEN** 配置 vLLM 时
- **THEN** `extra_config` 可支持 `max_tokens`、`temperature` 等默认参数

### Requirement: 流式输出支持

系统 SHALL 支持 SSE（Server-Sent Events）流式响应。

#### Scenario: 流式响应格式

- **WHEN** Provider 返回流式响应时
- **THEN** 通过 channel 逐块返回数据
- **AND** 每个数据块包含 delta 内容和完成状态

#### Scenario: 流式上下文取消

- **WHEN** 客户端取消请求或超时时
- **THEN** 通过 context.Context 传递取消信号
- **AND** Adapter 应中断流式响应并清理资源

### Requirement: 超时控制

系统 SHALL 支持为每个 Provider 配置独立的超时时间。

#### Scenario: 请求超时

- **WHEN** Provider 调用超过配置的 `timeout` 时间时
- **THEN** 取消请求并返回超时错误
- **AND** 超时时间从 Provider 配置读取

#### Scenario: 默认超时

- **WHEN** 未指定 timeout 时
- **THEN** 使用默认值 300 秒

### Requirement: OpenAI Adapter Chat 方法实现

OpenAI Adapter SHALL 实现完整的 `Chat()` 方法，支持调用 OpenAI 兼容的聊天完成 API。`base_url` 必须包含完整的 API 路径前缀（如 `/v1`），系统会在其后追加 `/chat/completions` 构建完整请求路径。

#### Scenario: 非 SSE 流式调用

- **WHEN** 调用 OpenAI Adapter 的 `Chat()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/chat/completions`
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
- **THEN** 发送 HTTP POST 请求到 `{base_url}/chat/completions`
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

vLLM Adapter SHALL 实现完整的 `Chat()` 方法，复用 OpenAI 兼容协议。`base_url` 必须包含完整的 API 路径前缀。

#### Scenario: 非 SSE 流式调用

- **WHEN** 调用 vLLM Adapter 的 `Chat()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/chat/completions`
- **AND** 请求体格式符合 OpenAI API 规范
- **AND** 如果配置了 API Key，携带 `Authorization: Bearer {api_key}` 请求头
- **AND** 如果未配置 API Key，不发送认证头
- **AND** 返回解析后的 `ChatResponse`

### Requirement: vLLM Adapter ChatStream 方法实现

vLLM Adapter SHALL 实现完整的 `ChatStream()` 方法，复用 OpenAI 兼容协议。

#### Scenario: SSE 流式调用

- **WHEN** 调用 vLLM Adapter 的 `ChatStream()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/chat/completions`
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

