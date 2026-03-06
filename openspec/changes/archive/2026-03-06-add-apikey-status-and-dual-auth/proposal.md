# 提案：API Key 启用/禁用/删除和 Chat 接口双重认证

## 概述

本提案旨在为 Courier LLM Gateway 添加两个功能：
1. API Key 管理增加启用、禁用和删除接口
2. Chat 接口同时支持 API Key 和 JWT 两种认证方式

## 背景

### 当前状态

**API Key 管理**：
- 支持创建 API Key（`POST /api/v1/users/:id/api-keys`）
- 支持列出 API Key（`GET /api/v1/users/:id/api-keys`）
- 支持撤销 API Key（`DELETE /api/v1/users/:id/api-keys/:key_id`），将状态设置为 `revoked`
- API Key 模型支持 `active`、`disabled`、`revoked` 三种状态
- 但**没有**单独的启用/禁用接口

**Chat 接口认证**：
- Chat 接口（`POST /v1/chat/completions`）当前**仅支持** API Key 认证
- 使用 `middleware.APIKeyAuth` 中间件
- 无法通过 JWT Token 访问

### 问题

1. **API Key 状态管理不完整**：
   - 用户无法临时禁用某个 API Key，只能永久撤销
   - 缺少启用/禁用切换的能力
   - 删除操作实际上是撤销（软删除），而非真正的删除

2. **Chat 接口认证方式单一**：
   - Web UI 用户在 Dashboard 中测试对话时，必须使用 API Key
   - 增加了用户体验复杂度（需要先创建 API Key 再测试对话）
   - 无法直接利用已登录的 JWT Session

## 目标

### 目标 1：API Key 启用/禁用/删除接口

1. **启用 API Key 接口**
   - 将状态设置为 `active`
   - 允许通过 API 重新启用被禁用的 Key

2. **禁用 API Key 接口**
   - 将状态设置为 `inactive`（规范中当前使用 `disabled`，需要统一）
   - 临时禁用 Key，保留重新启用的能力

3. **删除 API Key 接口**
   - 直接从数据库删除 API Key 记录（硬删除）
   - 不同于现有的撤销接口（软删除）

### 目标 2：Chat 接口双重认证

1. **支持 API Key 认证**
   - 保持现有 API Key 认证方式
   - 用于外部 API 调用

2. **支持 JWT 认证**
   - 允许使用 JWT Access Token 访问
   - 用于 Web UI 和已登录用户
   - 使用量统计只记录到用户级别，不关联 API Key

## 影响范围

### 修改的规范

- **apikey-auth** - 添加启用/禁用/删除接口需求
- **chat-api** - 添加 JWT 认证支持

### 相关规范

- **jwt-auth** - JWT 认证已实现，需要在 Chat 接口中应用

## 用户故事

### 作为 API 用户
- 我希望能够临时禁用某个 API Key，而不是永久删除
- 我希望能够重新启用之前禁用的 API Key
- 我希望能够彻底删除不再使用的 API Key

### 作为 Web UI 用户
- 我希望在 Dashboard 中测试对话时，不需要单独创建 API Key
- 我可以直接使用已登录的会话进行对话测试
- 我可以查看我的 API Key 列表，并启用/禁用它们

## 设计决策

### 决策 1：状态值统一

当前规范和代码中存在不一致：
- 规范（`apikey-auth`）中定义：`active`、`disabled`、`revoked`
- 用户需求：使用 `inactive` 表示禁用状态

**决策**：遵循现有规范，使用 `disabled` 表示禁用状态。

### 决策 2：删除接口的实现方式

**选项 A**：硬删除（直接从数据库删除记录）
- 优点：彻底清理数据
- 缺点：无法恢复，可能影响审计

**选项 B**：软删除（标记为已删除）
- 优点：保留审计记录
- 缺点：需要额外的字段和逻辑

**决策**：按照用户需求，实现硬删除（直接从数据库删除）。

### 决策 3：Chat 接口双重认证的实现方式

**选项 A**：创建组合中间件（尝试两种认证方式）
- 优点：单一中间件，逻辑集中
- 缺点：中间件逻辑复杂

**选项 B**：两个独立的中间件，按顺序尝试
- 优点：职责分离，易于测试
- 缺点：需要协调两个中间件

**选项 C**：新的统一认证中间件
- 优点：灵活支持多种认证方式
- 缺点：可能过度设计

**决策**：使用选项 B - 创建一个组合中间件，先尝试 JWT 认证，失败后尝试 API Key 认证。

### 决策 4：JWT 认证时的使用量记录

当使用 JWT 访问 Chat 接口时，如何记录使用量？

**选项 A**：只记录用户 ID，不关联 API Key
- 优点：简单直接，符合用户级别统计的需求
- 缺点：无法区分具体使用哪个 API Key

**选项 B**：使用用户的第一个 `active` API Key
- 优点：保持使用量记录格式一致
- 缺点：增加了复杂度，不准确（用户可能有多个 API Key）

**决策**：使用选项 A - 只记录用户 ID，不关联 API Key。用量统计粒度为用户级别。

## 非目标

以下功能不在本次提案范围内：
- API Key 的批量操作
- API Key 的使用历史查询
- 对话历史的管理
- API Key 权限的细粒度控制（如限制某些模型）

## 依赖关系

- 后端 API Key 管理基础已就绪
- JWT 认证已实现
- Chat 接口已实现

## 验收标准

### API Key 启用/禁用/删除

- [ ] PATCH `/api/v1/users/:id/api-keys/:key_id/enable` 启用 API Key
- [ ] PATCH `/api/v1/users/:id/api-keys/:key_id/disable` 禁用 API Key
- [ ] DELETE `/api/v1/users/:id/api-keys/:key_id` 删除 API Key（硬删除）
- [ ] 只有 Key 的所有者或管理员可以操作
- [ ] 操作后返回更新后的 API Key 信息

### Chat 接口双重认证

- [ ] Chat 接口接受 API Key 认证（`Authorization: Bearer <api_key>`）
- [ ] Chat 接口接受 JWT 认证（`Authorization: Bearer <jwt_token>`）
- [ ] 认证失败返回 401 错误
- [ ] JWT 认证成功时正确记录使用量
- [ ] API Key 认证成功时正确记录使用量（现有行为）

## 风险和缓解

| 风险 | 缓解措施 |
|------|----------|
| 状态值混淆（`disabled` vs `inactive`） | 在 API 文档中明确说明状态值 |
| 硬删除导致数据丢失 | 在 API 文档中警告不可恢复 |
| 双重认证导致的安全问题 | 确保 JWT 和 API Key 认证都有相同的用户验证逻辑 |
| 使用量记录不准确 | 在日志中标注认证类型（JWT/API Key） |
