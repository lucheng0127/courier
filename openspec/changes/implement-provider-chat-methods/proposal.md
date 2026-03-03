# Change: 实现 Provider Adapter 的 Chat 和 ChatStream 方法

## Why

当前 OpenAI 和 vLLM Adapter 的 `Chat()` 和 `ChatStream()` 方法仅返回 "not implemented yet" 错误，系统无法实际调用 LLM 服务。这导致整个聊天 API 功能无法正常工作，也无法进行真实的服务测试。

## What Changes

- 实现 OpenAI Adapter 的 `Chat()` 方法，支持非流式对话调用
- 实现 OpenAI Adapter 的 `ChatStream()` 方法，支持 SSE 流式对话调用
- 实现 vLLM Adapter 的 `Chat()` 方法（复用 OpenAI 兼容协议）
- 实现 vLLM Adapter 的 `ChatStream()` 方法（复用 OpenAI 兼容协议）
- 添加 HTTP 客户端封装，支持认证头、超时控制、上下文取消
- 添加请求和响应的格式转换（内部格式 ↔ OpenAI API 格式）
- 支持从 extra_config 读取默认模型参数（如 max_tokens、temperature）

## Impact

- **受影响的规范**: `provider-adapter`
- **受影响的代码**:
  - `internal/adapter/openai/adapter.go` - 实现 Chat 和 ChatStream 方法
  - `internal/adapter/vllm/adapter.go` - 实现 Chat 和 ChatStream 方法
  - `internal/adapter/openai/client.go` (新建) - HTTP 客户端封装
  - `internal/adapter/vllm/client.go` (新建) - HTTP 客户端封装（可复用 openai）
