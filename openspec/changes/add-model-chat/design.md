## Context

这是 Courier AI API 网关的核心功能实现。用户需要通过 API Key 认证后，调用上游 AI 模型进行对话，系统需要记录完整的请求日志。

约束条件：
- 必须使用现有技术栈（Go、Gin、GORM、SQLite、Zap）
- 模型配置通过 YAML 文件管理，不存储到数据库
- 必须支持流式响应（Server-Sent Events）
- 必须记录完整请求和响应日志
- API Key 认证通过中间件实现
- 禁止全局变量，所有依赖通过构造函数注入
- 分层架构：Handler → Service → Repository / Client

## Goals / Non-Goals

- **Goals**:
  - 实现可用的模型对话功能
  - 提供清晰的 API 接口
  - 支持流式和非流式两种响应模式
  - 记录完整的请求日志用于审计和计费
  - 支持配置多个上游模型

- **Non-Goals**:
  - 用户权限管理（所有认证用户都可调用所有模型）
  - Token 计费功能
  - 速率限制
  - 模型负载均衡
  - 失败重试机制
  - 对话历史存储（仅记录日志，不提供查询）

## Decisions

### 模型配置设计

**YAML 配置结构**：
```yaml
models:
  - name: qwen-turbo
    provider: qwen
    base_url: https://dashscope.aliyuncs.com/compatible-mode/v1
    api_key: ${QWEN_API_KEY}
  - name: deepseek-chat
    provider: deepseek
    base_url: https://api.deepseek.com/v1
    api_key: ${DEEPSEEK_API_KEY}
```

支持环境变量替换，避免敏感信息明文存储。

### API Key 认证中间件

**认证逻辑**：
1. 从请求头 `Authorization: Bearer <api_key>` 提取 API Key
2. 查询数据库验证 API Key 是否存在且状态为 `active`
3. 更新 API Key 的 `last_used_at` 时间戳
4. 将用户信息存入上下文（`gin.Context`）

**错误处理**：
- 缺少 Authorization 头：401 Unauthorized
- API Key 不存在：401 Unauthorized
- API Key 已禁用：401 Unauthorized

### 模型对话 API 设计

**请求格式**（OpenAI 兼容）：
```json
{
  "messages": [
    {"role": "user", "content": "你好"}
  ],
  "stream": true
}
```

**响应格式**：
- 非流式：标准 JSON 响应
- 流式：Server-Sent Events (SSE)，格式为 `data: <json>\n\n`

### 上游模型客户端

**接口设计**：
```go
type ModelClient interface {
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error)
}
```

使用标准库 `net/http` 实现，不依赖第三方 SDK。

### 请求日志设计

**数据模型**：
- `ID` (uint, 主键)
- `UserID` (uint, 外键)
- `ModelName` (string, 模型名称)
- `RequestMessages` (text, JSON 格式的请求消息)
- `ResponseContent` (text, 响应内容)
- `PromptTokens` (int, 请求 Token 数)
- `CompletionTokens` (int, 响应 Token 数)
- `TotalTokens` (int, 总 Token 数)
- `LatencyMs` (int, 响应延迟毫秒)
- `Status` (string, success/error)
- `ErrorMessage` (string, 错误信息)
- `CreatedAt` (time.Time)

**写入策略**：
- 异步写入，不阻塞主请求
- 使用 goroutine + channel 实现

## Risks / Trade-offs

- **流式响应下的日志记录**：流式响应是分段返回的，需要缓存完整响应才能记录
  - **缓解措施**：在内存中缓存流式响应，流结束后写入日志

- **上游 API Key 泄露风险**：配置文件中包含上游 API Key
  - **缓解措施**：支持环境变量，生产环境使用密钥管理服务

- **并发写入日志**：高并发下可能产生大量写操作
  - **缓解措施**：使用 buffered channel 批量写入，当前 MVP 阶段单 goroutine 处理

- **流式响应中断**：客户端断开连接后上游请求仍在继续
  - **缓解措施**：使用 context.Context 传播取消信号

## Migration Plan

无迁移计划，这是新增功能。

## Open Questions

无。
