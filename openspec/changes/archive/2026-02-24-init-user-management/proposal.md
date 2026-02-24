# Change: 初始化项目并实现用户与 API Key 管理

## Why

当前项目需要完成基础架构搭建，包括 Go 模块初始化、用户管理功能和 API Key 管理功能。这是 AI API 网关服务的核心基础，用户通过 API Key 进行身份认证和访问控制。

## What Changes

- 初始化 Go 模块 `github.com/lucheng0127/courier`
- 实现用户管理功能：创建用户、查询用户列表、查询用户详情
- 实现 API Key 管理功能：生成 API Key、查询 API Key 列表、删除 API Key、禁用 API Key
- API Keys 作为 Users 的子资源（一个用户可以有多个 API Key）
- 提供 RESTful HTTP API 接口
- 实现 SQLite 数据持久化

## Impact

- 涉及规范：`user-management`、`api-key-management`
- 影响代码：
  - 新增项目基础结构（`cmd/server/main.go`、`internal/` 目录）
  - 新增用户相关模型、仓储、服务和处理器
  - 新增 API Key 相关模型、仓储、服务和处理器
  - 新增数据库迁移脚本
