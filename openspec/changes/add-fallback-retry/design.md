# Design: Fallback 重试和可观测性架构设计

## Context

当前系统已实现基本的 Chat API，支持多 Provider 接入和模型路由。但在生产环境中，需要更高的可靠性和可观测性：

1. **可靠性**：模型调用可能因各种原因失败（网络超时、服务不可用、限流等），需要自动重试和 Fallback
2. **超时控制**：防止请求无限期等待，影响用户体验
3. **可观测性**：需要 TraceID 追踪请求链路，结构化日志便于监控和诊断

## Goals / Non-Goals

### Goals

- 实现 Fallback 机制：同 Provider 内模型失败时自动切换备用模型
- 实现超时控制：为所有请求设置合理的超时时间
- 实现 TraceID 生成和透传：在 HTTP 层生成，传递给 Provider
- 增强统一日志：结构化 JSON 输出，包含 TraceID 和完整上下文
- 实现请求重试计数：记录重试次数和最终使用的模型

### Non-Goals

- 跨 Provider Fallback（后续实现）
- 复杂的重试策略（如指数退避）
- 分布式追踪（如 OpenTelemetry）
- 日志聚合到外部系统（如 ELK）

## Decisions

### 1. Fallback 策略

**Scope**：仅在 **同一 Provider 内** 进行模型 Fallback

**触发条件**（满足任一即触发 Fallback）：
1. 超时错误
2. 网络错误（连接失败、DNS 解析失败）
3. 5xx 服务器错误
4. 连接拒绝
5. 服务不可用

**不触发 Fallback**：
- 4xx 客户端错误（如参数错误）
- 认证失败
- 模型不存在

**Fallback 流程**：
```
1. 尝试主模型：provider/model-primary
2. 失败 → 尝试备用模型1：provider/model-fallback-1
3. 失败 → 尝试备用模型2：provider/model-fallback-2
4. ... 直到成功或列表耗尽
```

**配置格式**：
```go
// Provider 配置
{
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "fallback_models": [
    "gpt-4o",           // 主模型
    "gpt-4o-mini",     // 备用模型1
    "gpt-3.5-turbo"    // 备用模型2
  ]
}
```

### 2. 超时控制

**默认超时**：30 秒（可配置）

**超时层级**：
```
1. 全局超时（最外层）- 控制整个请求的最大时间
2. Provider 超时（单次调用）- 控制单次 Provider 调用的最大时间
3. Fallback 超时 - 控制所有 Fallback 尝试的总时间
```

**实现方式**：
- 使用 `context.WithTimeout` 创建带超时的上下文
- 超时时取消正在进行的请求
- 记录超时日志

### 3. TraceID 生成和透传

**生成位置**：HTTP 中间件层（第一个处理请求的中间件）

**格式**：`trace-<UUID>`（例如：`trace-550e8400-e29b-41d4-a716-446655440000`）

**透传方式**：
1. HTTP 中间件生成 TraceID
2. 存储到 `context.Context` 中
3. 在整个请求处理链路中从 Context 读取
4. 传递给 Provider（通过 HTTP Header 或自定义字段）

**Header 名称**：`X-Trace-ID`

**日志集成**：所有日志记录包含 TraceID

### 4. 统一日志

**格式**：结构化 JSON

**日志级别**：
- `debug` - 详细调试信息
- `info` - 正常业务流程
- `warn` - 警告（如重试）
- `error` - 错误（如最终失败）

**日志字段**：
```json
{
  "timestamp": "2026-03-02T12:00:00Z",
  "level": "info",
  "trace_id": "trace-550e8400-e29b-41d4-a716-446655440000",
  "request_id": "chatcmpl-123",
  "api_key": "sk-...key",
  "provider": "openai-main",
  "model": "gpt-4o",
  "fallback_count": 1,
  "final_model": "gpt-4o-mini",
  "status": "success",
  "latency_ms": 1250,
  "error": ""
}
```

### 5. 数据库模型变更

**providers 表新增字段**：
```sql
ALTER TABLE providers ADD COLUMN fallback_models JSONB;
ALTER TABLE providers ADD COLUMN timeout INTEGER DEFAULT 30;
```

**fallback_models 格式**：
```json
["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
```

### 6. API 变更

**请求体变更**（可选，用于覆盖默认配置）：
```json
{
  "model": "openai-main/gpt-4o",
  "messages": [...],
  "timeout": 60,           // 可选：覆盖默认超时
  "max_fallbacks": 2      // 可选：限制最大 fallback 次数
}
```

### 7. 错误处理

**Fallback 耗尽**：返回最后一个错误，包含所有尝试的信息

```json
{
  "error": {
    "message": "All models failed. Tried: gpt-4o (timeout), gpt-4o-mini (500), gpt-3.5-turbo (timeout)",
    "type": "service_unavailable",
    "details": [
      {"model": "gpt-4o", "error": "timeout"},
      {"model": "gpt-4o-mini", "error": "500 Internal Server Error"},
      {"model": "gpt-3.5-turbo", "error": "timeout"}
    ]
  }
}
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Fallback 导致成本增加 | 记录 Fallback 使用情况，支持配置限制 |
| Fallback 掩盖真实问题 | 详细记录每次尝试的错误信息 |
| 超时时间设置不当 | 提供合理的默认值，支持配置覆盖 |
| TraceID 冲突/重复 | 使用 UUID 保证唯一性 |
| 日志量过大 | 支持日志级别配置 |

## Migration Plan

**数据库迁移**：
1. 添加 `fallback_models` 字段
2. 添加 `timeout` 字段
3. 为现有 Provider 设置默认值

**部署步骤**：
1. 执行数据库迁移
2. 部署新代码
3. 更新 Provider 配置（添加 fallback_models）
4. 重载 Provider
5. 验证 Fallback 功能

**回滚策略**：
- 回滚代码版本
- 恢复数据库 schema

## Open Questions

1. **Fallback 次数限制？**
   - 建议：MVP 不限制，后续可配置 `max_fallbacks`

2. **是否支持跨 Provider Fallback？**
   - 建议：MVP 不支持，后续可作为独立功能

3. **TraceID 是否需要支持外部传入？**
   - 建议：MVP 自动生成，后续支持从 Header 读取

4. **超时时间的层级关系？**
   - 建议：全局超时 > Fallback 总超时 > 单次调用超时
