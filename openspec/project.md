# Project Context

## Purpose

本项目旨在实现一个面向企业内部或小规模商业验证的 LLM Gateway（类 OpenRouter 架构）的 MVP 版本。

核心目标：

1. 提供兼容 OpenAI 风格的统一 API 入口（/v1/chat/completions）。
2. 支持多模型接入与基础路由能力（固定路由 + fallback）。
3. 提供基础的 token 统计与调用日志能力。
4. 提供 API Key 鉴权与限流能力。
5. 可通过 Docker Compose 一键部署，适合私有化或小规模上线。

本项目定位为“AI 网关中台基础设施”，不承担模型训练职责，仅负责模型请求转发、路由与治理。

## Tech Stack

- Golang（后端核心服务）
- PostgreSQL（主数据存储）
- Redis（限流、缓存、临时状态存储）
- Docker + Docker Compose（部署）
- RESTful API（对外接口规范）
- SSE（流式响应支持）

## Project Conventions

### Code Style

#### 1. 语言规范

- 使用 Go 1.22+
- 启用 go mod
- 强制使用 gofmt
- go mod 名为 github.com/lucheng0127/courier

#### 2. 命名规范
- 包名：小写、单数名词（如 router, adapter, service）
- 文件名：snake_case
- 结构体：PascalCase
- 方法/函数：PascalCase（对外） / camelCase（内部）
- 数据库字段：snake_case

#### 3. 分层原则
- controller：HTTP 接口层
- service：业务逻辑层
- repository：数据库访问层
- adapter：外部模型适配层
- middleware：鉴权、限流、中间件

禁止：
- controller 直接访问数据库
- adapter 中写业务逻辑

### Architecture Patterns

本项目采用经典分层架构 + Provider Adapter 模式。

整体架构：

Client
   ↓
API Layer (Controller)
   ↓
Service Layer
   ↓
Router Engine
   ↓
Provider Adapter
   ↓
LLM Provider

#### 1. Adapter 模式

每个模型供应商实现统一接口：

```
type Provider interface {
Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
...
}
```

不同厂商实现不同 adapter：

- OpenAIAdapter
- AnthropicAdapter
- LocalVLLMAdapter

禁止在 service 层出现厂商差异判断。

#### 2. Router 模块
Router 负责：

- 模型查找
- fallback 逻辑
- 超时控制
- 重试机制

MVP 阶段仅支持：
- 固定模型路由
- 单次 fallback

#### 3. 中间件模式
统一中间件处理：

- API Key 校验
- 限流
- 请求日志记录
- TraceID 注入

### Testing Strategy

#### 1. 单元测试
- service 层必须有单元测试
- router 逻辑必须覆盖
- adapter 需要 mock 测试

覆盖率目标：≥ 60%

### Git Workflow

采用简化 Git Flow。

主分支：

- main：生产可发布版本
- develop：开发主分支

功能分支：

- feature/xxx
- fix/xxx
- refactor/xxx

Commit 规范：

采用 Conventional Commits：

- feat: 新功能
- fix: 修复
- refactor: 重构
- test: 测试
- chore: 维护

## Domain Context

本项目属于 LLM Gateway（大模型网关）领域。

核心概念：

1. Model Registry  
   存储模型配置，包括：
   - model_name
   - provider
   - endpoint
   - api_key
   - timeout
   - fallback_model

2. Provider  
   指外部模型供应商，如：
   - OpenAI
   - 本地 vLLM
   - 私有模型服务

3. Token Usage  
   记录：
   - prompt_tokens
   - completion_tokens
   - total_tokens

4. Fallback  
   当主模型：
   - 超时
   - 5xx 错误
   - 网络异常

自动切换备用模型。

5. API Key  
   系统内部签发，用于：
   - 用户身份识别
   - 限流
   - 统计

本项目不是训练平台，也不是推理引擎，只是“请求调度与治理层”。

## Important Constraints

### 技术约束

- 必须支持流式输出（SSE）
- 必须支持超时控制（默认 60s）
- 单实例可支撑 300 QPS
- 必须无状态设计（可横向扩展）

### 部署约束

- 必须支持 Docker Compose 一键启动
- 不依赖云厂商特定服务
- 可私有化部署

### 数据约束

- 不长期存储用户 Prompt 内容（默认可配置）
- 日志必须支持脱敏

### 安全约束

- 所有 API 必须通过 API Key

## External Dependencies

1. PostgreSQL
   - 存储：
     - 用户
     - API Key
     - 模型配置
     - 使用日志

2. Redis
   - 限流
   - 短期缓存
   - 分布式锁（可选）

3. 外部模型 API
   - OpenAI 兼容接口
   - 私有模型 HTTP 服务
   - 其他第三方模型 API

本项目当前阶段为 MVP，仅实现：

- 统一 API
- 多模型接入
- 固定路由 + fallback
- 基础限流
- 使用统计
- 管理接口（基础）

不实现：

- 复杂计费系统
- 多租户体系
- 模型灰度发布
- 策略引擎 DSL
- 企业级权限系统
