# Chat API 使用文档

## 概述

Courier LLM Gateway 提供 OpenAI 兼容的 Chat Completions API，支持对接多个 LLM Provider。

## 端点

```
POST /v1/chat/completions
```

## 认证

所有请求需要通过 Bearer Token 进行 API Key 认证：

```http
Authorization: Bearer sk-your-api-key
```

### 配置 API Key

通过环境变量 `API_KEYS` 配置白名单（逗号分隔）：

```bash
export API_KEYS="sk-key-1,sk-key-2,sk-key-3"
```

开发模式下，任何以 `sk-` 开头的 Key 都会被接受。

## 请求格式

### 非流式请求

```json
{
  "model": "provider-name/model-name",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello!"}
  ],
  "temperature": 0.7,
  "max_tokens": 1000,
  "stream": false
}
```

### 流式请求

```json
{
  "model": "provider-name/model-name",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ],
  "stream": true
}
```

## 模型格式

采用 **OpenRouter 风格**：`provider/model_name`

- `provider` - Provider 实例名称（在系统配置中定义）
- `model_name` - 模型名称（由 Provider 定义）

### 示例

```json
{
  "model": "openai-main/gpt-4o"
}
```

```json
{
  "model": "vllm-local/llama-2-7b"
}
```

## 响应格式

### 非流式响应

```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "openai-main/gpt-4o",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! How can I help you today?"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 9,
    "total_tokens": 19
  }
}
```

### 流式响应

```
data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"openai-main/gpt-4o","choices":[{"index":0,"delta":{"content":"Hello"}}]}

data: [DONE]
```

## 错误响应

### 模型格式错误

```json
{
  "error": {
    "message": "invalid model format: gpt-4 (expected format: provider/model_name)",
    "type": "invalid_request_error"
  }
}
```

### Provider 不存在

```json
{
  "error": {
    "message": "provider not found: unknown-provider",
    "type": "invalid_request_error"
  }
}
```

### API Key 无效

```json
{
  "error": {
    "message": "Invalid API key",
    "type": "invalid_request_error"
  }
}
```

## 使用示例

### cURL

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-key" \
  -d '{
    "model": "openai-main/gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### Python

```python
import requests

response = requests.post(
    "http://localhost:8080/v1/chat/completions",
    headers={
        "Authorization": "Bearer sk-your-key",
        "Content-Type": "application/json"
    },
    json={
        "model": "openai-main/gpt-4o",
        "messages": [{"role": "user", "content": "Hello!"}]
    }
)

print(response.json())
```

### JavaScript/Node.js

```javascript
const response = await fetch('http://localhost:8080/v1/chat/completions', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer sk-your-key',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    model: 'openai-main/gpt-4o',
    messages: [{ role: 'user', content: 'Hello!' }]
  })
});

const data = await response.json();
console.log(data);
```

## 配置 Provider

在创建 Provider 时，系统会自动将其注册为可用的模型前缀。

例如，创建名为 `openai-main` 的 Provider：

```bash
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "openai-main",
    "type": "openai",
    "base_url": "https://api.openai.com/v1",
    "timeout": 300,
    "api_key": "sk-xxx",
    "enabled": true
  }'
```

然后可以使用 `openai-main/<model>` 格式调用模型。

## 日志

每个请求都会记录日志（JSON 格式），包含：

- `trace_id` - 链路追踪 ID
- `request_id` - 请求 ID
- `api_key` - API Key（脱敏）
- `model` - 完整模型名称
- `provider_name` - Provider 名称
- `provider_type` - Provider 类型
- `model_name` - 模型名称
- `fallback_count` - Fallback 次数
- `final_model` - 最终使用的模型
- `attempt_details` - 每次尝试的详情
- `prompt_tokens` - 输入 token 数
- `completion_tokens` - 输出 token 数
- `total_tokens` - 总 token 数
- `latency_ms` - 请求耗时
- `status` - 状态（success/error）
- `error` - 错误信息（如果有）
- `timestamp` - 时间戳

### 日志示例

```json
{
  "timestamp": "2026-03-02T12:00:00Z",
  "level": "info",
  "message": "Chat request completed",
  "trace_id": "trace-550e8400-e29b-41d4-a716-446655440000",
  "request_id": "chatcmpl-123",
  "api_key": "sk-...key",
  "provider": "openai-main",
  "model": "openai-main/gpt-4o",
  "fallback_count": 1,
  "final_model": "gpt-4o-mini",
  "latency_ms": 1250,
  "status": "success"
}
```

## Fallback 重试机制

当模型调用失败时，系统会自动尝试 Fallback 列表中的下一个模型。

### 触发条件

- 超时错误
- 网络错误（连接失败、DNS 解析失败）
- 5xx 服务器错误
- 连接拒绝

### 不触发 Fallback

- 4xx 客户端错误（参数错误）
- 认证失败
- 模型不存在

### Fallback 耗尽响应

当所有 Fallback 模型都失败时：

```json
{
  "error": {
    "message": "All models failed after 3 attempts. Last error: timeout",
    "type": "service_unavailable",
    "details": [
      {
        "model": "gpt-4o",
        "error_type": "timeout",
        "duration_ms": 30000
      },
      {
        "model": "gpt-4o-mini",
        "error_type": "server_error",
        "duration_ms": 2500
      },
      {
        "model": "gpt-3.5-turbo",
        "error_type": "timeout",
        "duration_ms": 30000
      }
    ]
  }
}
```

### TraceID

每个请求都会生成唯一的 TraceID，可通过响应头获取：

```http
X-Trace-ID: trace-550e8400-e29b-41d4-a716-446655440000
```

TraceID 用于追踪请求链路，方便排查问题。
