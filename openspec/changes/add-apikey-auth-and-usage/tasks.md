# Implementation Tasks

本文档按依赖顺序列出实现 add-apikey-auth-and-usage 变更的任务。

## Phase 1: 数据库和基础模型

### 1.1 创建数据库迁移文件

- [ ] 创建 `migrations/000003_create_users.up.sql` - users 表定义
- [ ] 创建 `migrations/000003_create_users.down.sql` - users 表回滚
- [ ] 创建 `migrations/000004_create_api_keys.up.sql` - api_keys 表定义
- [ ] 创建 `migrations/000004_create_api_keys.down.sql` - api_keys 表回滚
- [ ] 创建 `migrations/000005_create_usage_records.up.sql` - usage_records 表定义
- [ ] 创建 `migrations/000005_create_usage_records.down.sql` - usage_records 表回滚

### 1.2 创建数据模型

- [ ] 创建 `internal/model/user.go` - User 和 APIKey 结构体
- [ ] 创建 `internal/model/usage.go` - UsageRecord 结构体

**验证**：运行 `go build` 确保模型定义无语法错误

## Phase 2: 数据访问层 (Repository)

### 2.1 实现 User Repository

- [ ] 创建 `internal/repository/user.go`
- [ ] 实现 `CreateUser` 方法
- [ ] 实现 `GetUserByID` 方法
- [ ] 实现 `GetUserByEmail` 方法
- [ ] 实现 `ListUsers` 方法（分页、状态过滤）
- [ ] 实现 `UpdateUserStatus` 方法

### 2.2 实现 API Key Repository

- [ ] 在 `internal/repository/user.go` 中添加 API Key 相关方法
- [ ] 实现 `CreateAPIKey` 方法
- [ ] 实现 `GetAPIKeyByHash` 方法
- [ ] 实现 `ListAPIKeysByUserID` 方法
- [ ] 实现 `UpdateAPIKeyStatus` 方法
- [ ] 实现 `UpdateKeyLastUsed` 方法

### 2.3 实现 Usage Repository

- [ ] 创建 `internal/repository/usage.go`
- [ ] 实现 `CreateUsageRecord` 方法
- [ ] 实现 `QueryUsageByUserAndTimeRange` 方法
- [ ] 实现 `AggregateUsageByDay` 方法
- [ ] 实现 `AggregateUsageByModel` 方法

**验证**：编写单元测试验证 Repository 方法

## Phase 3: 业务逻辑层 (Service)

### 3.1 实现 Auth Service

- [ ] 创建 `internal/service/auth.go`
- [ ] 实现 `ValidateAPIKey` 方法（哈希验证 + 状态检查 + 过期检查）
- [ ] 实现 `GetUserByID` 方法
- [ ] 实现 `CreateUser` 方法（生成用户 ID）
- [ ] 实现 `CreateAPIKey` 方法（生成 Key + 哈希 + 返回完整 Key）
- [ ] 实现 `ListAPIKeys` 方法
- [ ] 实现 `RevokeAPIKey` 方法
- [ ] 实现 `UpdateKeyLastUsed` 异步方法

### 3.2 实现 Usage Service

- [ ] 创建 `internal/service/usage.go`
- [ ] 实现 `RecordUsage` 方法（异步写入）
- [ ] 实现 `GetUsageStats` 方法（按用户和时间范围查询）
- [ ] 实现 `aggregateByDay` 辅助方法
- [ ] 实现 `aggregateByModel` 辅助方法
- [ ] 实现批量写入 channel 机制

**验证**：编写单元测试验证 Service 逻辑

## Phase 4: 中间件

### 4.1 重构 API Key 鉴权中间件

- [ ] 修改 `internal/middleware/apikey.go`
- [ ] 移除环境变量白名单逻辑
- [ ] 注入 AuthService 依赖
- [ ] 实现从数据库验证 API Key
- [ ] 实现用户信息注入到 Context（`user_id`, `user_email`, `api_key_id`）
- [ ] 实现异步更新 `last_used_at`
- [ ] 处理各种鉴权失败场景（无效、禁用、过期）

### 4.2 实现管理员鉴权中间件

- [ ] 创建 `internal/middleware/admin.go`
- [ ] 实现 `AdminAuth` 中间件
- [ ] 从 `X-Admin-API-Key` Header 验证
- [ ] 支持环境变量配置管理员 Key

**验证**：编写中间件单元测试

## Phase 5: 控制器 (Controller)

### 5.1 实现 User Controller

- [ ] 创建 `internal/controller/user.go`
- [ ] 实现 `CreateUser` handler
- [ ] 实现 `GetUser` handler
- [ ] 实现 `ListAPIKeys` handler
- [ ] 实现 `CreateAPIKey` handler
- [ ] 实现 `RevokeAPIKey` handler
- [ ] 添加请求验证和错误处理

### 5.2 实现 Usage Controller

- [ ] 创建 `internal/controller/usage.go`
- [ ] 实现 `GetUsageStats` handler
- [ ] 实现查询参数解析和验证
- [ ] 实现按天/模型聚合逻辑
- [ ] 添加时间范围默认值处理

### 5.3 集成到 Chat Controller

- [ ] 修改 `internal/controller/chat.go`
- [ ] 在 `logRequestWithRetry` 中添加使用量记录
- [ ] 从 Context 获取用户和 API Key 信息
- [ ] 调用 UsageService.RecordUsage

**验证**：使用 curl/Postman 测试 API 端点

## Phase 6: 路由和依赖注入

### 6.1 注册新路由

- [ ] 修改 `cmd/server/main.go`
- [ ] 注册用户管理路由组（需要管理员中间件）
- [ ] 注册 API Key 管理路由
- [ ] 注册使用统计路由

### 6.2 初始化依赖

- [ ] 在 main.go 中初始化 Repository 层
- [ ] 初始化 Service 层并注入 Repository
- [ ] 初始化 Controller 层并注入 Service
- [ ] 更新中间件初始化（注入 AuthService）

**验证**：启动服务确保所有依赖正确注入

## Phase 7: 测试

### 7.1 单元测试

- [ ] Repository 层单元测试（使用 mock DB）
- [ ] Service 层单元测试
- [ ] 中间件单元测试
- [ ] Controller 集成测试

### 7.2 端到端测试

- [ ] 测试完整的用户创建和 API Key 生成流程
- [ ] 测试使用新 API Key 调用 Chat API
- [ ] 测试使用统计查询
- [ ] 测试各种鉴权失败场景

**验证**：测试覆盖率 ≥ 60%

## Phase 8: 文档

- [ ] 更新 `docs/deployment.md` - 添加数据库迁移说明
- [ ] 创建 API 使用文档（用户和 API Key 管理）
- [ ] 更新 README 添加新功能说明
