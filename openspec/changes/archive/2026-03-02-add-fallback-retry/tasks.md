# Tasks: Fallback 重试和可观测性实现

## 1. 数据库层变更

- [x] 1.1 创建数据库迁移文件
  - 添加 `fallback_models` JSONB 字段
  - 添加 `timeout` INTEGER 字段（默认 30）
- [x] 1.2 更新 Provider 数据模型
- [x] 1.3 更新 Repository

## 2. TraceID 中间件

- [x] 2.1 实现 TraceID 生成中间件（`internal/middleware/trace.go`）
  - 生成格式：`trace-<UUID>`
  - 存储到 context.Context
  - 设置响应 Header：`X-Trace-ID`
- [x] 2.2 在 apikey 中间件后注册 trace 中间件

## 3. 重试服务

- [x] 3.1 实现重试服务（`internal/service/retry.go`）
  - `RetryWithFallback()` - 带 Fallback 的重试逻辑
  - 支持超时控制
  - 记录每次尝试的结果
- [x] 3.2 实现错误分类
  - 判断是否可重试的函数
  - 区分可重试错误和不可重试错误

## 4. 增强日志

- [x] 4.1 更新日志模型（`internal/model/log.go`）
  - 添加 `trace_id` 字段
  - 添加 `fallback_count` 字段
  - 添加 `final_model` 字段
  - 添加 `attempt_details` 字段
- [x] 4.2 实现结构化日志记录器
  - JSON 格式输出
  - 支持日志级别
- [x] 4.3 集成到 Chat 控制器

## 5. 修改 Chat 控制器

- [x] 5.1 集成 TraceID
  - 从 Context 读取 TraceID
  - 传递给日志和 Provider
- [x] 5.2 集成超时控制
  - 使用 context.WithTimeout
  - 处理超时取消
- [x] 5.3 集成 Fallback 重试
  - 调用重试服务
  - 处理 Fallback 结果
- [x] 5.4 更新响应格式（Fallback 耗尽时）

## 6. 修改路由服务

- [x] 6.1 添加 Fallback 模型列表获取
  - 从 Provider 配置读取 `fallback_models`
  - 验证模型列表

## 7. 修改 Adapter

- [x] 7.1 更新 Adapter 接口（可选）
  - 考虑是否需要传递 TraceID 给 Provider
- [x] 7.2 在 HTTP 请求中添加 TraceID Header

## 8. 配置和部署

- [x] 8.1 更新 Provider 配置示例
  - 添加 `fallback_models` 示例
  - 添加 `timeout` 示例
- [x] 8.2 更新部署文档

## 9. 测试和验证

- [x] 9.1 测试 TraceID 生成和透传
- [x] 9.2 测试超时控制
- [x] 9.3 测试 Fallback 功能（单次失败）
- [x] 9.4 测试 Fallback 功能（多次失败）
- [x] 9.5 测试 Fallback 耗尽
- [x] 9.6 测试不可重试错误（不触发 Fallback）
- [x] 9.7 测试日志包含 TraceID
- [ ] 9.8 测试跨 Provider 调用不触发 Fallback（需集成测试环境）
- [ ] 9.9 性能测试：超时取消（已实现单元测试）

## 10. 文档

- [x] 10.1 更新 API 文档（Fallback 配置）
- [x] 10.2 编写 Fallback 最佳实践
- [x] 10.3 更新运维文档

## 依赖关系

- 任务 1 可独立进行
- 任务 2、3、4 可以并行进行
- 任务 5 依赖任务 2、3、4
- 任务 6 依赖任务 1
- 任务 7 可以与任务 5 并行
- 任务 8、9、10 依赖任务 1-7
