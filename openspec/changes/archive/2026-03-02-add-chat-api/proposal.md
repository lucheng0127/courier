# Change: 添加 Chat API 层

## Why

当前系统已完成 Provider Adapter 层实现，可以对接多个 LLM 供应商。但缺少对外统一的 API 接口，无法接收客户端请求并路由到合适的 Provider。

本变更旨在实现兼容 OpenAI API 风格的 Chat Completions 接口，作为 LLM Gateway 的核心入口能力。

## What Changes

- 新增 `chat-api` 能力：实现 `/v1/chat/completions` 端点
- 实现 API Key 鉴权中间件
- 实现模型路由引擎（根据 `model` 参数选择 Provider）
- 支持非流式和流式（SSE）响应
- 统一请求/响应格式（OpenAI 兼容）
- 请求日志记录

## Impact

- **Affected specs**: 新增 `chat-api` 规范
- **Affected code**:
  - 新增 `internal/controller/chat.go` - Chat API 控制器
  - 新增 `internal/service/router.go` - 模型路由服务
  - 新增 `internal/middleware/apikey.go` - API Key 鉴权中间件
  - 新增 `internal/model/chat.go` - Chat 请求/响应模型
  - 修改 `cmd/server/main.go` - 注册新路由
- **Database**: 无新增表（API Key 管理后续实现，MVP 阶段使用配置或固定值）
- **Dependencies**: 无新增外部依赖
