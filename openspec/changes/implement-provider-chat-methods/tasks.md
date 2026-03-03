## 1. OpenAI 适配器实现

- [ ] 1.1 创建 `internal/adapter/openai/client.go`，封装 HTTP 客户端
  - 支持 Bearer Token 认证
  - 支持配置的超时时间
  - 支持 context 取消机制
  - 支持 TraceID 透传
- [ ] 1.2 在 client.go 中实现 `doChatRequest()` 方法，发送非流式请求
- [ ] 1.3 在 client.go 中实现 `doChatStreamRequest()` 方法，发送流式请求
- [ ] 1.4 实现 `Chat()` 方法
  - 构建 OpenAI API 请求体
  - 调用 client 发送请求
  - 解析响应并转换为内部格式
  - 错误处理和日志记录
- [ ] 1.5 实现 `ChatStream()` 方法
  - 构建 OpenAI API 请求体（stream=true）
  - 调用 client 发送流式请求
  - 逐行解析 SSE 响应
  - 通过 channel 返回数据块
  - 处理 context 取消和资源清理

## 2. vLLM 适配器实现

- [ ] 2.1 创建 `internal/adapter/vllm/client.go`（复用 openai 客户端逻辑）
  - 支持 API Key 可选（vLLM 本地部署可能不需要认证）
  - 其他特性与 openai 客户端一致
- [ ] 2.2 实现 vLLM Adapter 的 `Chat()` 方法
- [ ] 2.3 实现 vLLM Adapter 的 `ChatStream()` 方法

## 3. 请求/响应格式处理

- [ ] 3.1 定义 OpenAI API 请求结构体
- [ ] 3.2 定义 OpenAI API 响应结构体
- [ ] 3.3 定义 OpenAI SSE 流式响应结构体
- [ ] 3.4 实现内部请求格式 → OpenAI 格式的转换函数
- [ ] 3.5 实现 OpenAI 格式 → 内部响应格式的转换函数
- [ ] 3.6 实现 SSE 解析和流式数据块转换

## 4. 扩展配置支持

- [ ] 4.1 支持从 extra_config 读取默认模型参数
  - `max_tokens`
  - `temperature`
  - `top_p`
  - 等其他 OpenAI 支持的参数
- [ ] 4.2 请求级参数优先于默认参数

## 5. 测试

- [ ] 5.1 为 openai client 添加单元测试（mock HTTP 服务器）
- [ ] 5.2 为 vllm client 添加单元测试
- [ ] 5.3 为 `Chat()` 方法添加单元测试
- [ ] 5.4 为 `ChatStream()` 方法添加单元测试
- [ ] 5.5 添加真实服务集成测试（使用环境变量配置测试端点）
  - 支持 OpenAI 兼容服务（如 Qwen/通义千问）
  - 支持本地 vLLM 服务

## 6. 文档

- [ ] 6.1 更新 `docs/` 目录下的 Provider 配置文档
- [ ] 6.2 添加 OpenAI 兼容服务配置示例（如 Qwen、通义千问）
- [ ] 6.3 添加本地 vLLM 配置示例
