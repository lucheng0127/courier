## Context

这是 Courier AI API 网关项目的初始实现。项目需要支持用户管理和 API Key 管理功能，作为后续 AI 模型网关功能的基础。

约束条件：
- 单体应用架构，禁止过度抽象
- 必须使用 SQLite、Gin、GORM、Zap
- 禁止全局变量，所有依赖通过构造函数注入
- 分层架构：Handler → Service → Repository

## Goals / Non-Goals

- **Goals**:
  - 实现最小可行的用户和 API Key 管理功能
  - 提供清晰的 RESTful API 接口
  - 确保代码结构清晰，易于维护和扩展

- **Non-Goals**:
  - 用户认证/登录（当前仅支持通过 API 直接创建）
  - API Key 过期机制
  - 权限管理
  - 多租户支持

## Decisions

### 数据模型设计

**User 模型**：
- `ID` (uint, 主键)
- `Name` (string, 用户名，唯一)
- `Email` (string, 邮箱，可选)
- `CreatedAt` (time.Time)
- `UpdatedAt` (time.Time)

**APIKey 模型**：
- `ID` (uint, 主键)
- `UserID` (uint, 外键)
- `Key` (string, API Key 值，唯一索引，格式 `ck_<random>`，长度 32 字符)
- `Status` (string, 状态：active/disabled)
- `LastUsedAt` (time.Time, 可为空)
- `ExpiresAt` (time.Time, 可为空，预留字段)
- `CreatedAt` (time.Time)
- `UpdatedAt` (time.Time)

### API 路由设计

```
# 用户管理
POST   /api/v1/users              创建用户
GET    /api/v1/users              查询用户列表
GET    /api/v1/users/:id          查询用户详情

# API Key 管理（作为 Users 的子资源）
POST   /api/v1/users/:id/apikeys            为用户生成 API Key
GET    /api/v1/users/:id/apikeys            查询用户的 API Key 列表
DELETE /api/v1/users/:id/apikeys/:keyid     删除指定的 API Key
PUT    /api/v1/users/:id/apikeys/:keyid/disable 禁用指定的 API Key
```

### 依赖注入结构

main.go 负责组装所有依赖：
```
DB → Logger → Repository → Service → Handler → Router → Server
```

### API Key 生成方案

使用 `crypto/rand` 生成 24 字节随机数据，进行 hex 编码，添加 `ck_` 前缀，确保唯一性和安全性。

## Risks / Trade-offs

- **API Key 泄露风险**：API Key 以明文存储在数据库中，当前版本不加密
  - **缓解措施**：生产环境应考虑使用哈希存储或密钥管理服务

- **并发创建用户风险**：用户名可能重复
  - **缓解措施**：在数据库层面添加唯一索引约束

- **删除 API Key 为物理删除**：无法恢复
  - **缓解措施**：未来可考虑软删除（添加 deleted_at 字段）

## Migration Plan

无迁移计划，这是全新的初始实现。

## Open Questions

无。
