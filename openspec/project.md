# Project Context

## Purpose

构建一个最小可行版本（MVP）的 AI API 网关服务，实现以下功能：

系统职责：

- 为用户生成 API Key
- 用户通过 API Key 调用本服务接口
- 服务验证 API Key 并解析用户信息
- 将请求转发至上游 AI 模型（当前仅支持 Qwen）
- 记录每次请求的完整日志与统计信息
- 系统架构支持未来扩展至多个 AI 模型

该项目为单体应用架构，主要目标是验证业务流程与架构设计合理性。

## Tech Stack

语言与框架：

- Go 1.22+
- Gin (HTTP 框架)
- GORM (ORM)
- SQLite (数据库)
- Zap (结构化日志)
- net/http (上游模型调用)

辅助工具：

- crypto/rand (API Key 生成)
- context (请求控制)
- uuid (可选)

## Project Conventions

### 项目结构

采用清晰分层结构：

```
/cmd/server/main.go
/internal
/handler
/middleware
/service
/repository
/model
/client
/router
/pkg
```

原则：

- handler 不允许直接访问数据库
- repository 只负责数据持久化
- service 负责业务逻辑
- client 负责调用外部 API
- middleware 仅做请求拦截处理
- 所有依赖通过构造函数注入

### 依赖注入原则

- 禁止全局变量保存 DB 或 Logger
- 所有依赖通过结构体注入
- main.go 负责依赖初始化与组装
- 所有配置通过yaml文件进行配置

### 错误处理规范

- 不允许忽略 error
- 不允许 panic
- 所有 error 必须明确处理
- 外部错误不得暴露内部细节

### Code Style

遵循 Go 官方代码规范：

- 使用 gofmt
- 使用 go vet
- 结构体命名使用 PascalCase
- 私有变量使用 camelCase
- 错误变量统一命名为 err
- 结构体构造函数命名为 NewXxx()

日志规范：

- 使用 zap.Logger
- 必须使用结构化日志
- 不允许使用 fmt.Println
- 不允许使用 log.Println


### Architecture Patterns

#### 分层架构（Layered Architecture）

系统分为：

1. Transport Layer（HTTP）
2. Middleware Layer（认证）
3. Service Layer（业务逻辑）
4. Repository Layer（数据访问）
5. Client Layer（外部模型调用）

#### 接口驱动设计（Interface-driven Design）

所有外部依赖必须通过接口抽象，例如：

- ModelClient 接口
- UserRepository 接口
- RequestLogRepository 接口

禁止在业务层直接依赖具体实现。

#### 单体架构（Monolithic）

当前版本为单体应用。

禁止：

- 微服务拆分
- 分布式架构
- 过度抽象

### Testing Strategy
[Explain your testing approach and requirements]

### Git Workflow

采用简化 Git Flow：

- main：稳定分支
- develop：开发分支
- feature/xxx：功能分支

Commit 规范：

- feat: 新功能
- fix: 修复
- refactor: 重构
- docs: 文档
- chore: 其他

## Domain Context

系统本质是一个 AI API 网关。

核心领域概念：

- User（用户）
- API Key（鉴权凭证）
- Model（AI 模型）
- Chat Request（用户请求）
- Upstream Provider（上游模型服务）
- Request Log（请求日志）

当前支持模型：

- Qwen（单模型）

未来扩展：

- DeepSeek
- 多模型路由
- 权重分配
- 失败回退机制

## Important Constraints

- 数据库必须使用 SQLite
- 日志必须使用 zap
- 单请求超时 30 秒
- 必须支持 API Key 鉴权
- 必须记录完整请求与响应
- 不允许阻塞主流程写日志
- 必须支持未来多模型扩展
- 不允许在 Handler 中写复杂逻辑
- 不允许在 Handler 中访问数据库
- API Key 使用 crypto/rand 生成
- 不记录上游 API 密钥
- 错误信息不暴露内部实现
- 不在日志中记录敏感字段（如上游密钥）

## External Dependencies

上游依赖：

- Qwen HTTP API

运行依赖：

- SQLite 文件数据库
- 本地运行

日志依赖：

- go.uber.org/zap

ORM：

- gorm.io/gorm
- gorm.io/driver/sqlite

HTTP 框架：

- github.com/gin-gonic/gin
