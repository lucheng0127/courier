# unify-api-and-jwt-auth 实施任务

## 阶段 1：基础设施准备

### 1.1 数据库迁移

- [x] 创建数据库迁移文件 `000006_add_user_role.up.sql`
  - 添加 `role` 字段到 `users` 表
  - 设置默认值为 `'user'`
  - 添加 CHECK 约束
- [x] 创建数据库迁移文件 `000006_add_user_role.down.sql`
  - 移除 `role` 字段
- [x] 创建数据库迁移文件 `000007_add_password_hash.up.sql`
  - 添加 `password_hash` 字段到 `users` 表
- [x] 创建数据库迁移文件 `000007_add_password_hash.down.sql`
  - 移除 `password_hash` 字段
- [x] 运行数据库迁移（待部署时执行）

### 1.2 模型层更新

- [x] 更新 `internal/model/user.go`
  - 添加 `Role` 字段到 `User` 结构体
  - 添加 `PasswordHash` 字段到 `User` 结构体
- [x] 在 `internal/model/user.go` 中添加新的请求/响应模型
  - `LoginRequest`
  - `LoginResponse`
  - `RefreshTokenRequest`
  - `RefreshTokenResponse`
- [x] 创建 `internal/model/jwt.go`
  - 添加 `JWTClaims` 模型
  - 添加 `TokenType` 常量
  - 添加 `TokenInfo` 模型

### 1.3 Repository 层更新

- [x] 更新 `internal/repository/user.go`
  - 在 `CreateUser` 查询中包含 `role` 和 `password_hash` 字段
  - 在 `GetUserByID` 查询中包含 `role` 字段
  - 在 `GetUserByEmail` 查询中包含 `role` 字段
  - 在 `ListUsers` 查询中包含 `role` 字段
  - 添加 `UpdateUser` 方法
  - 添加 `UpdatePassword` 方法
  - 添加 `GetUserByEmailWithPassword` 方法（用于登录验证）

## 阶段 2：JWT 服务实现

### 2.1 JWT Service

- [x] 添加依赖 `github.com/golang-jwt/jwt/v5`
- [x] 创建 `internal/service/jwt.go`
  - 实现 `JWTService` 接口定义
  - 实现 `GenerateAccessToken` 方法
  - 实现 `GenerateRefreshToken` 方法
  - 实现 `ValidateAccessToken` 方法
  - 实现 `ValidateRefreshToken` 方法
  - 实现 `GetAccessTokenExpiration` 方法
- [x] 添加 JWT 配置结构
  - 从环境变量读取配置
  - 验证必需的配置项（JWT_SECRET）

### 2.2 密码哈希工具

- [x] 创建 `internal/pkg/password/password.go`
  - 实现 `HashPassword` 函数（使用 bcrypt）
  - 实现 `VerifyPassword` 函数
  - 设置 bcrypt cost factor 为 12

## 阶段 3：中间件实现

### 3.1 JWT 鉴权中间件

- [x] 创建 `internal/middleware/jwt.go`
  - 实现 `JWTAuth` 中间件函数
  - 验证 Authorization Header 格式
  - 验证 JWT Token
  - 提取用户信息并注入上下文
  - 处理各种错误情况
  - 添加辅助函数 `GetUserID`、`GetUserEmail`、`GetUserRole`

### 3.2 角色验证中间件

- [x] 更新 `internal/middleware/admin.go`
  - 新增 `RequireAdmin` 中间件
  - 从上下文获取用户角色
  - 验证角色是否为 `admin`
  - 返回适当的错误响应
- [x] 保留原有的 Admin API Key 支持作为过渡
  - 保留 `AdminAuth` 中间件用于降级兼容

### 3.3 速率限制中间件（可选）

- [ ] 创建 `internal/middleware/rate_limit.go`
  - 实现基于 IP 的速率限制
  - 专门用于登录接口
  - 配置：1 分钟内最多 5 次尝试

## 阶段 4：认证服务

### 4.1 Auth Service 更新

- [x] 更新 `internal/service/auth.go`
  - 添加 `Login` 方法
  - 添加 `RefreshToken` 方法
  - 添加 `CreateAdminUser` 方法（用于初始化）
  - 添加 `EnsureInitialAdmin` 方法

### 4.2 初始管理员创建

- [x] 在 `main.go` 中调用 `EnsureInitialAdmin`
  - 检查是否存在管理员用户
  - 从环境变量读取初始管理员信息
  - 创建默认管理员账户

## 阶段 5：认证控制器

### 5.1 AuthController 实现

- [x] 创建 `internal/controller/auth.go`
  - 创建 `AuthController` 结构体
  - 实现 `NewAuthController` 构造函数
  - 实现 `Login` 处理函数
  - 实现 `RefreshToken` 处理函数
  - 实现 `RegisterRoutes` 方法

### 5.2 路由注册

- [x] 在 `main.go` 中添加认证路由
  - 创建 `/api/v1/auth` 路由组（无需鉴权）
  - 注册 `POST /login` 路由
  - 注册 `POST /refresh` 路由

## 阶段 6：现有 Controller 重构

### 6.1 UserController 重构

- [x] 更新 `internal/controller/user.go`
  - 实现 `RegisterRoutes` 方法
  - 添加权限检查（普通用户只能访问自己的资源）
  - 添加 `ListUsers` 处理函数（占位实现）
  - 添加 `UpdateUser` 处理函数（占位实现）
  - 添加 `DeleteUser` 处理函数（占位实现）
  - 添加 `UpdateUserStatus` 处理函数（占位实现）

### 6.2 UsageController 重构

- [x] 更新 `internal/controller/usage.go`
  - 实现 `RegisterRoutes` 方法
  - 添加权限检查（普通用户只能查询自己的统计）

### 6.3 ChatController 重构

- [x] 更新 `internal/controller/chat.go`
  - 实现 `RegisterRoutes` 方法

### 6.4 ProviderReloadController 更新

- [x] 验证 `RegisterRoutes` 方法实现正确
  - 确保路由使用 `/api/v1/admin/providers/*` 路径

## 阶段 7：main.go 路由重构

### 7.1 路由组织

- [x] 重构 `cmd/server/main.go` 中的路由注册
  - 统一使用 `/api/v1` 前缀
  - 创建 JWT 鉴权组
  - 在 JWT 组下创建 Admin 角色组
  - 移除硬编码的路由注册
  - 调用各 Controller 的 `RegisterRoutes` 方法
- [x] 组织 Chat API 路由
  - 保持 `/v1/chat/completions` 路径
  - 使用 API Key 鉴权

### 7.2 中间件应用

- [x] 确保中间件应用顺序正确
  - 全局：日志、恢复、CORS
  - 管理 API：JWTAuth
  - Admin 接口：RequireAdmin
  - Chat API：APIKeyAuth、TraceID

## 阶段 8：测试

### 8.1 单元测试

- [x] JWT Service 单元测试
- [x] 密码工具单元测试
- [ ] Auth Service 单元测试
- [ ] Repository 层测试

### 8.2 集成测试

- [ ] 登录流程测试
  - 成功登录
  - 密码错误
  - 用户不存在
  - 用户被禁用
- [ ] JWT 鉴权测试
  - 有效 Token
  - 过期 Token
  - 无效 Token
  - 缺少 Header
- [ ] 角色权限测试
  - 管理员访问所有接口
  - 普通用户访问自己的资源
  - 普通用户被拒绝访问管理接口
- [ ] Token 刷新测试
  - 成功刷新
  - 无效 Refresh Token
  - 过期 Refresh Token

### 8.3 API 测试

- [ ] 测试所有管理接口的新路径
- [ ] 测试 Chat API 保持原有行为
- [ ] 测试所有使用 `RegisterRoutes` 的接口

## 阶段 9：文档与部署

### 9.1 API 文档

- [x] 更新 API 文档
  - 记录新的认证方式
  - 更新所有接口路径
  - 添加错误响应示例

### 9.2 部署准备

- [x] 添加环境变量文档
  - `JWT_SECRET`
  - `JWT_ACCESS_TOKEN_EXPIRES_IN`
  - `JWT_REFRESH_TOKEN_EXPIRES_IN`
  - `INITIAL_ADMIN_EMAIL`
  - `INITIAL_ADMIN_PASSWORD`
- [ ] 创建迁移指南
  - 如何从 Admin API Key 迁移到 JWT
  - 客户端代码更新示例

### 9.3 清理

- [ ] 移除旧的 Admin API Key 相关代码（可选，在过渡期后）
- [ ] 更新 README
- [ ] 更新 CHANGELOG

## 验收标准

- [x] 所有管理接口使用 `/api/v1` 前缀
- [x] Chat API 保持 `/v1/chat/completions` 路径
- [x] 管理接口使用 JWT 鉴权
- [x] Chat API 使用 API Key 鉴权
- [x] 只有管理员可以访问管理接口
- [x] 普通用户可以登录并访问自己的资源
- [x] 所有 Controller 实现 `RegisterRoutes` 方法
- [x] 核心单元测试通过（JWT、密码工具）
- [ ] 所有集成测试通过
- [x] API 文档已更新
- [x] 环境变量文档已更新

## 依赖关系

```
阶段 1 (数据库) → 阶段 2 (模型/Repository) → 阶段 3 (JWT服务) → 阶段 4 (中间件) → 阶段 5 (认证服务) → 阶段 6 (AuthController) → 阶段 7 (重构Controller) → 阶段 8 (main.go重构) → 阶段 9 (测试) → 阶段 10 (文档)
```

可并行任务：
- 阶段 2 的模型更新和 Repository 更新可以并行
- 阶段 6 的各 Controller 重构可以并行
- 单元测试可以与实现并行编写
