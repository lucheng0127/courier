# Change: 实现上游模型调用功能

## Why

当前系统已完成用户管理和 API Key 管理功能，但尚未实现核心的 AI 模型网关能力。用户需要通过 API Key 进行身份认证后，调用上游 AI 模型进行对话。

## What Changes

### 新增功能

1. **模型配置管理**
   - 通过 YAML 配置上游模型列表
   - 每个模型配置包含：base_url（HTTPS 地址）、api_key、model_name
   - 支持多个上游模型提供商

2. **模型列表查询**
   - 提供 API 接口获取当前可用的模型列表
   - 返回模型名称和提供商信息

3. **API Key 认证中间件**
   - 实现认证中间件，验证请求头中的 API Key
   - 解析 API Key 对应的用户信息
   - 仅对需要认证的接口生效

4. **模型对话功能**
   - RESTful API：`POST /api/v1/models/:model/chat`
   - 支持流式响应（Server-Sent Events）
   - 支持多轮对话上下文

5. **请求日志记录**
   - 记录完整的请求内容（用户消息、模型名称）
   - 记录完整的响应内容（AI 回复）
   - 记录 Token 使用量（如果上游提供）
   - 记录请求时间、响应时间、状态

### 数据模型

- **ModelConfig**：模型配置（从 YAML 加载，不存储到数据库）
- **RequestLog**：请求日志（存储到 SQLite）

### API 接口

```
# 公开接口
GET    /api/v1/models                  查询可用模型列表

# 需要认证的接口
POST   /api/v1/models/:model/chat      发起模型对话（支持流式）
```

## Impact

- 涉及规范：新增 `model-management`、`model-chat`、`request-logging`、`api-auth`
- 影响代码：
  - 新增 `internal/client/` 目录，存放上游模型客户端
  - 新增 `internal/middleware/auth.go`，API Key 认证中间件
  - 新增 `internal/model/request_log.go`，请求日志模型
  - 新增 `internal/repository/request_log.go`，请求日志仓储
  - 扩展 `pkg/config/config.go`，添加模型配置结构
  - 扩展 `internal/router/router.go`，添加新路由
