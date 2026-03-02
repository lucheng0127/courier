# Change: 添加 API Key 鉴权和用户使用统计功能

## Why

当前系统的 API Key 鉴权是通过环境变量配置简单白名单实现的，缺乏以下核心能力：

1. **用户与 API Key 关联管理**：无法支持一个用户拥有多个 API Key 的场景
2. **动态 API Key 管理**：API Key 只能通过环境变量配置，无法动态增删改查
3. **使用统计追踪**：虽然请求日志中包含 Token 使用量，但没有按用户聚合统计
4. **中间件用户注入**：鉴权中间件只验证 API Key 有效性，未将用户信息注入到请求上下文中

本变更旨在实现完整的 API Key 鉴权体系和用户使用统计功能。

## What Changes

### 新增功能

1. **apikey-auth 能力**：
   - 实现用户（User）与 API Key 的数据模型和数据库表
   - API Key 动态管理（创建、删除、禁用、查询）
   - 中间件从 API Key 解析用户信息并注入到 Gin Context
   - 支持一个用户拥有多个 API Key

2. **usage-tracking 能力**：
   - 实现使用记录（Usage）的数据模型和数据库表
   - 记录每次请求的 Token 使用量并与用户关联
   - 提供使用统计查询接口（按用户、时间范围聚合）
   - 支持按 API Key、模型、时间维度统计

### 架构调整

- 新增 `internal/model/user.go` - 用户和 API Key 模型
- 新增 `internal/model/usage.go` - 使用记录模型
- 新增 `internal/repository/user.go` - 用户数据访问层
- 新增 `internal/repository/usage.go` - 使用记录数据访问层
- 新增 `internal/service/user.go` - 用户管理服务
- 新增 `internal/service/usage.go` - 使用统计服务
- 新增 `internal/controller/user.go` - 用户管理控制器
- 新增 `internal/controller/usage.go` - 使用统计控制器
- 修改 `internal/middleware/apikey.go` - 从数据库验证 API Key 并注入用户信息

### 数据库变更

- 新增 `users` 表：存储用户信息
- 新增 `api_keys` 表：存储 API Key，与 users 关联
- 新增 `usage_records` 表：存储使用记录，与 users 和 api_keys 关联

## Impact

- **Affected specs**:
  - 新增 `apikey-auth` 规范 - API Key 鉴权和用户管理
  - 新增 `usage-tracking` 规范 - 使用统计功能
- **Affected code**:
  - 修改 `internal/middleware/apikey.go` - 增强鉴权中间件
  - 修改 `internal/controller/chat.go` - 使用注入的用户信息记录使用量
- **Database**: 新增 3 张表（users, api_keys, usage_records）
- **API Endpoints**:
  - `POST /v1/users` - 创建用户
  - `GET /v1/users/:id` - 获取用户信息
  - `POST /v1/users/:id/api-keys` - 为用户创建 API Key
  - `GET /v1/users/:id/api-keys` - 获取用户的 API Key 列表
  - `DELETE /v1/users/:id/api-keys/:key_id` - 删除 API Key
  - `GET /v1/usage` - 查询使用统计
- **Dependencies**: 无新增外部依赖
