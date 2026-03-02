# Change: 添加 Fallback 重试和可观测性能力

## Why

当前系统已实现基本的 Chat API 功能，但缺少生产环境必需的可靠性保障和可观测性能力：

1. **无 Fallback 机制**：当模型调用失败时，无法自动切换到备用模型，导致服务中断
2. **无超时控制**：请求可能无限期等待，影响用户体验
3. **日志不完整**：缺少 TraceID，难以追踪请求链路
4. **缺少可观测性**：无法有效监控和诊断问题

本变更旨在添加 Fallback 重试、超时控制、TraceID 透传和统一日志等生产级能力。

## What Changes

- 新增 `fallback-retry` 能力：实现同 Provider 内的模型 Fallback
- 实现 Fallback 配置：每个 Provider 可配置主模型和备用模型列表
- 实现超时控制：为所有请求设置合理的超时时间
- 实现 TraceID 生成和透传：在 HTTP 中间件中生成 TraceID，并在整个请求链路中传递
- 增强统一日志：所有日志包含 TraceID，支持结构化输出
- 实现请求重试计数和状态跟踪

## Impact

- **Affected specs**:
  - 新增 `fallback-retry` 规范
  - 修改 `chat-api` 规范（添加 Fallback 相关需求）
  - 修改 `provider-adapter` 规范（添加 fallback 配置支持）
- **Affected code**:
  - 修改 `internal/controller/chat.go` - 添加 Fallback 逻辑
  - 新增 `internal/middleware/trace.go` - TraceID 中间件
  - 修改 `internal/service/router.go` - 添加 Fallback 路由
  - 新增 `internal/service/retry.go` - 重试服务
  - 修改 `internal/model/provider.go` - 添加 Fallback 配置字段
  - 修改 `internal/model/log.go` - 添加 TraceID 和重试相关字段
  - 修改 `internal/middleware/apikey.go` - 透传 TraceID
- **Database**: 修改 `providers` 表，添加 `fallback_models` 字段
- **Dependencies**: 新增 `github.com/google/uuid`（已添加）
