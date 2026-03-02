# Design: Chat API 架构设计

## Context

系统已实现 Provider Adapter 层，能够对接多种 LLM 供应商。现在需要实现对外统一的 Chat API 层，接收客户端请求并路由到合适的 Provider。

## Goals / Non-Goals

### Goals

- 兼容 OpenAI `/v1/chat/completions` API 格式
- 支持非流式和流式（SSE）响应
- 根据 `model` 参数路由到对应 Provider
- API Key 鉴权
- 请求日志记录
- 错误处理和转换

### Non-Goals

- 复杂的路由策略（如负载均衡、灰度发布）
- API Key 管理系统（MVP 阶段使用配置或固定值）
- 请求缓存
- 复杂的限流（后续实现）

## Decisions

### 1. API 端点设计

采用 OpenAI 兼容格式：

```
POST /v1/chat/completions
```

**请求格式**（OpenRouter 兼容）：
```json
{
  "model": "openai/gpt-4o",   // 必填：provider/model_name 格式
  "messages": [
    {"role": "user", "content": "Hello"}
  ],
  "stream": false,            // 可选：是否流式响应
  "temperature": 0.7,
  "max_tokens": 1000
}
```

**响应格式**（OpenAI 兼容）：
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-4",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "Hello!"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}
```

**理由**:
- OpenAI 格式是事实标准
- 降低客户端迁移成本
- 便于与现有工具集成

### 2. 模型路由策略

采用 **OpenRouter 风格**的模型命名：`provider/model_name`

**格式说明**：
- `provider` - Provider 实例名称（对应 Provider 配置中的 `name`）
- `model_name` - 模型名称（由 Provider 端定义）

**示例**：
```json
{
  "model": "openai-main/gpt-4o"      // 使用 openai-main Provider 的 gpt-4o 模型
  "model": "openai-main/gpt-3.5-turbo"
  "model": "vllm-local/llama-2-7b"   // 使用 vllm-local Provider 的 llama-2-7b 模型
  "model": "anthropic/claude-3-opus"
}
```

**路由流程**：
```
1. 解析 model 参数：split("/") 分割为 [provider_name, model_name]
2. 根据 provider_name 查找对应的 Provider 实例
3. 验证 Provider 是否启用
4. 调用 Provider 的 Chat() 或 ChatStream() 方法，传入 model_name
5. 转换响应格式为 OpenAI 格式
6. 返回给客户端
```

**错误处理**：
- 格式错误（不含 `/`）：返回 400，提示正确格式
- Provider 不存在：返回 404
- Provider 未启用：返回 403
- 模型不存在：由 Provider 返回错误

**优势**：
- 明确指定 Provider，无歧义
- 避免模型名称冲突
- 客户端可精确控制使用哪个 Provider
- 便于多租户和成本控制

**Provider 配置**：
```go
// Provider 配置示例
{
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  // 支持的模型列表（可选，用于验证）
  "extra_config": {
    "models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
  }
}
```

### 3. API Key 鉴权

MVP 阶段使用**固定白名单**或**简单验证**：

```go
// 方式 1: 环境变量配置白名单
var validAPIKeys = map[string]bool{
    "sk-test-key-1": true,
    "sk-test-key-2": true,
}

// 方式 2: 任何以 "sk-" 开头的 Key 都通过（仅用于开发）
```

**中间件流程**：
```
1. 从 Header 读取 Authorization: Bearer <token>
2. 验证 token 是否有效
3. 有效则放行，无效则返回 401
```

**理由**:
- MVP 阶段简化实现
- 后续可扩展为数据库存储和动态管理

### 4. 流式响应实现

使用 Server-Sent Events (SSE)：

```go
func (c *ChatController) ChatCompletions(ctx *gin.Context) {
    if stream {
        // 流式响应
        chunks := provider.ChatStream(ctx, req)

        ctx.Header("Content-Type", "text/event-stream")
        ctx.Header("Cache-Control", "no-cache")
        ctx.Header("Connection", "keep-alive")

        for chunk := range chunks {
            data := formatSSE(chunk)
            ctx.Writer.Write(data)
            ctx.Writer.Flush()
        }
    } else {
        // 非流式响应
        resp := provider.Chat(ctx, req)
        ctx.JSON(http.StatusOK, resp)
    }
}
```

**SSE 格式**（OpenAI 兼容）：
```
data: {"id":"chatcmpl-123","choices":[{"delta":{"content":"Hello"}}]}

data: [DONE]
```

### 5. 请求日志

记录每次请求的关键信息：

```go
type ChatLog struct {
    RequestID   string    // 请求 ID（UUID）
    APIKey      string    // API Key（脱敏）
    Model       string    // 模型名称
    Provider    string    // 实际调用的 Provider
    PromptTokens int      // 输入 token 数
    CompletionTokens int   // 输出 token 数
    TotalTokens int       // 总 token 数
    Latency     int64     // 请求耗时（毫秒）
    Status      string    // 状态：success/error
    Error       string    // 错误信息（如果有）
    Timestamp   time.Time // 请求时间
}
```

**日志输出**：
- MVP 阶段：写入标准日志（JSON 格式）
- 后续：写入数据库或 Elasticsearch

### 6. 目录结构

```
internal/
├── controller/
│   └── chat.go              # Chat API 控制器
├── service/
│   └── router.go            # 模型路由服务
├── middleware/
│   └── apikey.go            # API Key 鉴权中间件
├── model/
│   └── chat.go              # Chat 请求/响应模型
└── logger/
    └── chat.go              # Chat 请求日志
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| 模型路由冲突（多个 Provider 支持同一模型） | MVP 选择第一个，记录日志，后续支持优先级配置 |
| 流式响应连接中断 | 实现 context 取消传播，确保资源清理 |
| API Key 泄露 | 日志中脱敏处理，响应中不返回 |
| Provider 调用超时影响用户体验 | 设置合理超时时间，返回友好错误信息 |

## Migration Plan

由于是新功能，无迁移需求。

部署步骤：
1. 部署新代码
2. 配置模型到 Provider 的映射
3. 测试 API 端点
4. 客户端切换到新 API

回滚策略：
- 回滚代码版本
- 恢复旧版本 API（如果有）

## Open Questions

1. **是否支持别名模型？**
   - 建议：MVP 不支持，后续可通过配置添加 `alias` 字段
   - 例如：配置 `gpt-4o` 为 `openai-main/gpt-4o` 的别名

2. **是否需要支持 Function Calling？**
   - 建议：MVP 不支持，作为后续增强功能

3. **流式响应的 token 统计如何处理？**
   - 建议：流式响应结束后在 `finish` 消息中返回完整统计

4. **是否需要支持模型版本管理？**
   - 建议：MVP 不支持，后续可通过模型名称后缀实现（如 `gpt-4o@20240101`）
