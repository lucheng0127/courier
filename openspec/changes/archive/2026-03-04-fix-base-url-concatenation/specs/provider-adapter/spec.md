# provider-adapter Spec Delta

## MODIFIED Requirements

### Requirement: OpenAI Adapter Chat 方法实现

OpenAI Adapter SHALL 修改 URL 构建逻辑，直接在 base_url 后添加 `/chat/completions`，不再自动添加 `/v1` 路径。

#### Scenario: 非 SSE 流式调用

- **WHEN** 调用 OpenAI Adapter 的 `Chat()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/chat/completions`
- **AND** 请求体格式符合 OpenAI API 规范
- **AND** 携带 `Authorization: Bearer {api_key}` 请求头（如果配置了 API Key）
- **AND** 使用 Provider 配置的超时时间
- **AND** 返回解析后的 `ChatResponse`

#### Scenario: base_url 配置要求

- **WHEN** 配置 OpenAI 类型的 Provider 时
- **THEN** `base_url` MUST 包含完整的 API 路径前缀
- **AND** 对于标准 OpenAI API，`base_url` 应设为 `https://api.openai.com/v1`
- **AND** 对于智谱 GLM，`base_url` 应设为 `https://open.bigmodel.cn/api/paas/v4`
- **AND** 系统在 base_url 后追加 `/chat/completions` 构建完整请求路径

#### Scenario: URL 路径构建

- **WHEN** base_url 为 `https://api.openai.com/v1`
- **THEN** 完整请求 URL SHALL 为 `https://api.openai.com/v1/chat/completions`
- **WHEN** base_url 为 `https://open.bigmodel.cn/api/paas/v4`
- **THEN** 完整请求 URL SHALL 为 `https://open.bigmodel.cn/api/paas/v4/chat/completions`
- **WHEN** base_url 已包含完整路径（如 `https://api.openai.com/v1/chat/completions`）
- **THEN** 系统不重复添加路径，直接使用该 URL

---

### Requirement: OpenAI Adapter ChatStream 方法实现

OpenAI Adapter SHALL 修改流式调用的 URL 构建逻辑，直接在 base_url 后添加 `/chat/completions`。

#### Scenario: SSE 流式调用

- **WHEN** 调用 OpenAI Adapter 的 `ChatStream()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/chat/completions`
- **AND** 请求体中 `stream` 字段设置为 `true`
- **AND** 携带 `Authorization: Bearer {api_key}` 请求头（如果配置了 API Key）
- **AND** 使用 Provider 配置的超时时间
- **AND** 返回只读 channel 用于接收流式数据

---

### Requirement: OpenAI 兼容服务支持

系统 SHALL 支持配置 OpenAI 兼容的第三方服务，要求 base_url 包含完整的 API 路径前缀。

#### Scenario: 配置第三方服务

- **WHEN** 配置 OpenAI 类型的 Provider 时
- **THEN** `base_url` 可指向任何 OpenAI 兼容的服务端点
- **AND** `base_url` MUST 包含完整的 API 路径前缀（如 `/v1`、`/api/paas/v4` 等）
- **AND** `api_key` 使用该服务提供的密钥
- **AND** Adapter 正常工作，无需额外配置

#### Scenario: 通义千问配置示例

- **WHEN** 配置通义千问服务时
- **THEN** `type` 设置为 `openai`
- **AND** `base_url` 设置为 `https://dashscope.aliyuncs.com/compatible-mode/v1`
- **AND** `api_key` 设置为通义千问的 API Key
- **AND** 可在 `extra_config` 中设置默认模型参数

#### Scenario: 智谱 GLM 配置示例

- **WHEN** 配置智谱 GLM 服务时
- **THEN** `type` 设置为 `openai`
- **AND** `base_url` 设置为 `https://open.bigmodel.cn/api/paas/v4`
- **AND** `api_key` 设置为智谱 GLM 的 API Key
- **AND** 可在 `extra_config` 中设置默认模型参数

---

### Requirement: vLLM Adapter Chat 方法实现

vLLM Adapter SHALL 修改 URL 构建逻辑，与 OpenAI Adapter 保持一致，直接在 base_url 后添加 `/chat/completions`。

#### Scenario: 非 SSE 流式调用

- **WHEN** 调用 vLLM Adapter 的 `Chat()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/chat/completions`
- **AND** 请求体格式符合 OpenAI API 规范
- **AND** 如果配置了 API Key，携带 `Authorization: Bearer {api_key}` 请求头
- **AND** 如果未配置 API Key，不发送认证头
- **AND** 返回解析后的 `ChatResponse`

---

### Requirement: vLLM Adapter ChatStream 方法实现

vLLM Adapter SHALL 修改流式调用的 URL 构建逻辑，与 OpenAI Adapter 保持一致。

#### Scenario: SSE 流式调用

- **WHEN** 调用 vLLM Adapter 的 `ChatStream()` 方法时
- **THEN** 发送 HTTP POST 请求到 `{base_url}/chat/completions`
- **AND** 请求体中 `stream` 字段设置为 `true`
- **AND** 如果配置了 API Key，携带认证头；否则不发送
- **AND** 返回只读 channel 用于接收流式数据
