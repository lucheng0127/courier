# 任务清单

## 概述

本文档列出了实现"API Key 启用/禁用/删除和 Chat 接口双重认证"功能的所有任务。

## 任务列表

### 阶段 1：数据层实现

- [ ] **1.1 Repository 层 - 添加 API Key 删除方法**
  - 在 `UserRepository` 接口中添加 `DeleteAPIKey` 方法
  - 在 `userRepository` 实现中添加硬删除逻辑
  - 添加单元测试验证删除功能

### 阶段 2：Service 层实现

- [ ] **2.1 Service 层 - 启用 API Key**
  - 在 `AuthService` 中添加 `EnableAPIKey` 方法
  - 验证用户权限和 API Key 所有权
  - 调用 Repository 更新状态为 `active`
  - 返回更新后的 API Key 信息

- [ ] **2.2 Service 层 - 禁用 API Key**
  - 在 `AuthService` 中添加 `DisableAPIKey` 方法
  - 验证用户权限和 API Key 所有权
  - 调用 Repository 更新状态为 `disabled`
  - 返回更新后的 API Key 信息

- [ ] **2.3 Service 层 - 删除 API Key**
  - 在 `AuthService` 中添加 `DeleteAPIKey` 方法
  - 验证用户权限和 API Key 所有权
  - 调用 Repository 硬删除记录
  - 返回成功状态

### 阶段 3：Controller 层实现

- [ ] **3.1 Controller 层 - 启用 API Key 接口**
  - 在 `UserController` 中添加 `EnableAPIKey` 方法
  - 路由：`PATCH /api/v1/users/:id/api-keys/:key_id/enable`
  - 实现权限验证（所有者或管理员）
  - 返回 200 OK 和更新后的 API Key 信息
  - 处理各种错误场景

- [ ] **3.2 Controller 层 - 禁用 API Key 接口**
  - 在 `UserController` 中添加 `DisableAPIKey` 方法
  - 路由：`PATCH /api/v1/users/:id/api-keys/:key_id/disable`
  - 实现权限验证（所有者或管理员）
  - 返回 200 OK 和更新后的 API Key 信息
  - 处理各种错误场景

- [ ] **3.3 Controller 层 - 删除 API Key 接口**
  - 修改或新增删除 API Key 的方法
  - 路由：`DELETE /api/v1/users/:id/api-keys/:key_id`（硬删除）
  - 或新增路由：`DELETE /api/v1/users/:id/api-keys/:key_id/permanent`
  - 实现权限验证（所有者或管理员）
  - 返回 204 No Content
  - 处理各种错误场景

- [ ] **3.4 路由注册**
  - 在 `RegisterRoutes` 中添加新路由
  - 确保路由在正确的鉴权组下（JWT 认证）

### 阶段 4：双重认证中间件

- [ ] **4.1 创建双重认证中间件**
  - 在 `middleware` 包中创建 `dual_auth.go`
  - 实现 `DualAuth` 中间件函数
  - 先尝试 JWT 认证，失败后尝试 API Key 认证
  - 两者都失败时返回 401 错误

- [ ] **4.2 注入认证类型**
  - 在 JWT 认证成功时注入 `auth_type=jwt`
  - 在 API Key 认证成功时注入 `auth_type=apikey`

- [ ] **4.3 中间件测试**
  - 测试 JWT 认证成功场景
  - 测试 API Key 认证成功场景
  - 测试两者都失败场景
  - 测试认证类型注入

### 阶段 5：Chat 接口集成

- [ ] **5.1 修改 Chat 接口认证中间件**
  - 将 `middleware.APIKeyAuth` 替换为 `middleware.DualAuth`
  - 更新中间件调用参数

- [ ] **5.2 更新使用量记录**
  - 修改使用量记录逻辑，包含 `auth_type` 字段
  - 确保 JWT 认证时只记录用户 ID，api_key_id 为 NULL
  - 确保 API Key 认证时记录用户 ID 和 API Key ID

- [ ] **5.3 更新日志记录**
  - 在日志中添加 `auth_type` 字段
  - 区分 JWT 和 API Key 认证的请求

### 阶段 6：前端集成

- [ ] **6.1 Dashboard - API Key 管理页面**
  - 添加启用/禁用按钮
  - 添加删除按钮（硬删除）
  - 显示 API Key 状态
  - 处理操作确认

- [ ] **6.2 Dashboard - 对话页面**
  - 移除必须创建 API Key 的限制
  - 支持直接使用 JWT Token 进行对话
  - 显示使用的认证方式

### 阶段 7：测试

- [ ] **7.1 单元测试**
  - Repository 层测试
  - Service 层测试
  - Controller 层测试
  - 中间件测试

- [ ] **7.2 集成测试**
  - API Key 启用/禁用流程
  - API Key 删除流程
  - JWT 认证访问 Chat 接口
  - API Key 认证访问 Chat 接口
  - 使用量记录验证

- [ ] **7.3 端到端测试**
  - 用户登录后直接测试对话（无需 API Key）
  - 用户创建 API Key 并测试对话
  - 用户禁用 API Key 后无法使用
  - 用户重新启用 API Key 后可以使用

### 阶段 8：文档和部署

- [ ] **8.1 API 文档更新**
  - 更新 API Key 管理接口文档
  - 添加启用/禁用/删除接口说明
  - 更新 Chat 接口认证说明

- [ ] **8.2 数据库迁移**
  - 确认数据库 schema 无需变更
  - 如需添加索引，创建迁移脚本

- [ ] **8.3 部署准备**
  - 准备发布说明
  - 准备升级指南
  - 准备回滚方案

## 依赖关系

1. **阶段 1 必须最先完成**：数据层是所有其他层的基础
2. **阶段 2 依赖阶段 1**：Service 层需要 Repository 层的方法
3. **阶段 3 依赖阶段 2**：Controller 层需要 Service 层的方法
4. **阶段 4 可以与阶段 2、3 并行**：中间件相对独立
5. **阶段 5 依赖阶段 4**：Chat 接口需要新的中间件
6. **阶段 6 依赖阶段 5**：前端需要后端接口完成
7. **阶段 7 在所有开发任务完成后进行**

## 可并行执行的任务

以下任务可以并行执行：
- 任务 2.1、2.2、2.3（在 1.1 完成后）
- 任务 3.1、3.2、3.3（在对应 Service 方法完成后）
- 任务 4.1 和阶段 3 的 Controller 任务

## 验收标准

### API Key 启用/禁用/删除

- [ ] 可以成功启用被禁用的 API Key
- [ ] 可以成功禁用激活的 API Key
- [ ] 可以成功删除任意状态的 API Key
- [ ] 权限验证正确（所有者或管理员）
- [ ] 错误处理正确（Key 不存在、无权限等）

### Chat 接口双重认证

- [ ] 可以使用 JWT Token 访问 Chat 接口
- [ ] 可以使用 API Key 访问 Chat 接口
- [ ] 两种认证方式都失败时返回正确错误
- [ ] JWT 认证时正确记录使用量（只记录用户 ID，不关联 API Key）
- [ ] API Key 认证时正确记录使用量（记录用户 ID 和 API Key ID）
- [ ] 日志包含认证类型信息
