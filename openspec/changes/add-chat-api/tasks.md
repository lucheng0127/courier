# Tasks: Chat API 实现

## 1. 数据模型层

- [x] 1.1 定义 Chat 请求模型（`internal/model/chat.go`）
  - `ChatRequest` - 兼容 OpenAI 格式
  - `ChatMessage` - 消息结构
  - `ChatResponse` - 非流式响应
  - `ChatStreamChunk` - 流式响应块
- [x] 1.2 定义日志模型（`internal/model/log.go`）
  - `ChatLog` - 请求日志结构

## 2. 路由服务层

- [x] 2.1 实现模型路由服务（`internal/service/router.go`）
  - `ParseModel(model)` - 解析 `provider/model_name` 格式
  - `ResolveProvider(providerName)` - 根据 provider 名称获取 Provider 实例
  - `GetAvailableModels()` - 获取所有可用模型列表（格式：`provider/model_name`）
  - 验证 Provider 是否启用
- [ ] 2.2 编写路由服务单元测试
  - 测试正确格式解析
  - 测试缺少 `/` 分隔符的错误
  - 测试 Provider 不存在
  - 测试 Provider 未启用

## 3. API Key 鉴权中间件

- [x] 3.1 实现 API Key 验证中间件（`internal/middleware/apikey.go`）
  - 从 Authorization Header 提取 API Key
  - 验证 API Key 有效性
  - 返回 401 错误（验证失败时）
- [x] 3.2 配置 API Key 白名单（环境变量或配置文件）
- [ ] 3.3 编写中间件单元测试

## 4. Chat API 控制器

- [x] 4.1 实现 Chat Completions 控制器（`internal/controller/chat.go`）
  - `ChatCompletions()` - 主处理函数
  - 支持非流式和流式响应
  - 调用路由服务获取 Provider
  - 调用 Provider 的 Chat/ChatStream 方法
  - 转换响应格式
- [x] 4.2 实现流式响应处理
  - 设置正确的 SSE Header
  - 逐块发送数据
  - 处理客户端断开连接
- [x] 4.3 实现响应格式转换
  - Provider 格式 → OpenAI 格式
  - 生成请求 ID 和时间戳
- [ ] 4.4 编写控制器单元测试

## 5. 请求日志

- [x] 5.1 实现日志记录器（集成在控制器中）
  - `LogRequest()` - 记录请求
  - API Key 脱敏
  - JSON 格式输出
- [x] 5.2 集成到 Chat 控制器

## 6. 路由注册

- [x] 6.1 在 main.go 中注册 Chat API 路由
  - `POST /v1/chat/completions`
- [x] 6.2 应用 API Key 鉴权中间件
- [ ] 6.3 配置 CORS（如果需要）

## 7. 配置管理

- [ ] 7.1 定义 Provider 模型配置格式
  - `extra_config.models` - 模型列表
  - `extra_config.default_model` - 默认模型
- [ ] 7.2 更新现有 Provider 配置示例

## 8. 测试和验证

- [x] 8.1 单元测试：非流式请求成功
- [x] 8.2 单元测试：流式请求成功
- [x] 8.3 单元测试：模型未找到错误
- [x] 8.4 单元测试：API Key 鉴权失败
- [x] 8.5 单元测试：请求日志正确记录
- [x] 8.6 集成测试：完整请求流程
- [x] 8.7 集成测试：多 Provider 模型路由
- [x] 8.8 性能测试：流式响应性能

## 9. 文档

- [x] 9.1 编写 API 使用文档
- [x] 9.2 编写配置示例
- [x] 9.3 更新部署文档

## 依赖关系

- 任务 1 和 2 可以并行进行
- 任务 3 可以独立进行
- 任务 4 依赖任务 1、2、3
- 任务 5 可以与任务 4 并行
- 任务 6 依赖任务 4、5
- 任务 7 可以与任务 4 并行
- 任务 8 依赖任务 1-7
- 任务 9 可以在实现完成后进行
