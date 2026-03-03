# unify-api-and-jwt-auth 设计文档

## 架构设计

### 系统分层

```
┌─────────────────────────────────────────────────────────────┐
│                        API Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  /v1/chat/      │  │  /api/v1/*      │  │/api/v1/auth │ │
│  │  completions    │  │  (管理接口)      │  │   (认证)    │ │
│  └────────┬────────┘  └────────┬────────┘  └──────┬──────┘ │
│           │                    │                   │        │
│  APIKeyAuth│              JWTAuth│           No Auth│        │
└───────────┼────────────────────┼───────────────────┼────────┘
            │                    │                   │
├───────────┼────────────────────┼───────────────────┼────────┤
│           ▼                    ▼                   ▼        │
│  ┌─────────────┐    ┌─────────────────┐   ┌──────────────┐  │
│  │ChatController│    │Other Controllers│   │AuthController│  │
│  └─────────────┘    └────────┬────────┘   └──────────────┘  │
│                               │                               │
│  ┌────────────────────────────┴────────────────────────────┐  │
│  │                    Middleware Layer                      │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │  │
│  │  │APIKeyAuth│  │ JWTAuth  │  │   RequireAdmin       │  │  │
│  │  └──────────┘  └──────────┘  │   (检查 role=admin)   │  │  │
│  │                               └──────────────────────┘  │  │
│  └────────────────────────────────────────────────────────┘  │
│                               │                               │
│  ┌────────────────────────────┴────────────────────────────┐  │
│  │                    Service Layer                         │  │
│  │  ┌──────────────┐  ┌──────────┐  ┌──────────────────┐  │  │
│  │  │RouterService │  │JWTService│  │  AuthService     │  │  │
│  │  └──────────────┘  └──────────┘  └──────────────────┘  │  │
│  └────────────────────────────────────────────────────────┘  │
│                               │                               │
│  ┌────────────────────────────┴────────────────────────────┐  │
│  │                    Repository Layer                      │  │
│  │  ┌──────────────┐  ┌──────────┐  ┌──────────────────┐  │  │
│  │  │UserRepo      │  │KeyRepo   │  │  ProviderRepo    │  │  │
│  │  └──────────────┘  └──────────┘  └──────────────────┘  │  │
│  └────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 鉴权流程设计

### 1. Chat Completions API（API Key 鉴权）

```
┌─────────┐                    ┌─────────┐
│ Client  │                    │Provider │
└────┬────┘                    └────┬────┘
     │                              │
     │ POST /v1/chat/completions    │
     │ Authorization: Bearer sk-xxx │
     │                              │
     ▼                              │
┌────────────────────────────────┐  │
│ APIKeyAuth Middleware          │  │
│ - 验证 API Key 格式            │  │
│ - 从数据库查询 Key             │  │
│ - 检查 Key 状态                │  │
│ - 检查用户状态                 │  │
│ - 注入 user_id, api_key_id     │  │
└────────┬───────────────────────┘  │
         │                          │
         ▼                          │
┌────────────────────────────────┐  │
│ ChatController                 │  │
│ - 验证请求参数                 │  │
│ - 解析 model 参数              │  │
└────────┬───────────────────────┘  │
         │                          │
         ▼                          │
┌────────────────────────────────┐  │
│ RouterService                  │  │
│ - 解析 provider/model_name     │  │
│ - 查找 Provider 实例           │  │
└────────┬───────────────────────┘  │
         │                          │
         ▼                          │
┌────────────────────────────────┐  │
│ Provider Adapter               │  │
│ - 调用上游 API                 │──┼──►
└────────────────────────────────┘  │
                                    │
```

### 2. 管理 API（JWT 鉴权）

```
┌─────────┐
│ Client  │
└────┬────┘
     │
     │ ① POST /api/v1/auth/login
     │    { "email": "...", "password": "..." }
     │
     ▼
┌────────────────────────────────┐
│ AuthController                 │
│ - 验证邮箱密码                 │
│ - 检查用户状态（active）        │
│ - 支持用户登录（user 和 admin） │
└────────┬───────────────────────┘
         │
         ▼
┌────────────────────────────────┐
│ JWTService                     │
│ - 生成 access_token (15min)    │
│ - 生成 refresh_token (7days)   │
│ - Token 包含 user_role Claims  │
│ - 返回登录响应                 │
└────────┬───────────────────────┘
         │
         │ ② { "access_token": "...", "refresh_token": "..." }
         │
         │
┌─────────┘
     │
     │ ③ 管理员访问 POST /api/v1/providers
     │    或 普通用户访问 GET /api/v1/usage
     │    Authorization: Bearer <access_token>
     │
     ▼
┌────────────────────────────────┐
│ JWTAuth Middleware             │
│ - 验证 JWT 签名                │
│ - 检查过期时间                 │
│ - 注入 user_id, user_role      │
└────────┬───────────────────────┘
         │
         ├──► 管理员接口 ──► RequireAdmin ──► ProviderController
         │                                  (检查 user_role == "admin")
         │
         └──► 用户接口 ──► UserController/UsageController
                                    (根据 user_id 和 user_role 权限控制)
```

**权限说明**：
- **管理员（admin）**：可以访问所有管理接口，包括 Provider 管理、用户管理、查看所有用户的使用统计
- **普通用户（user）**：可以登录获取 Token，管理自己的 API Key，查看自己的使用统计

### 3. Token 刷新流程

```
┌─────────┐
│ Client  │
└────┬────┘
     │
     │ POST /api/v1/auth/refresh
     │ { "refresh_token": "..." }
     │
     ▼
┌────────────────────────────────┐
│ JWTService                     │
│ - 验证 refresh_token 签名       │
│ - 检查 refresh_token 过期时间   │
│ - 生成新的 access_token        │
│ - 生成新的 refresh_token       │
└────────┬───────────────────────┘
         │
         │ { "access_token": "...", "refresh_token": "..." }
         │
┌─────────┘
```

## 数据模型设计

### Users 表变更

```sql
-- 添加角色字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'user'
CHECK (role IN ('user', 'admin'));

-- 添加密码哈希字段（用于管理员登录）
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);

-- 为管理员创建索引
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
```

### 用户模型更新

```go
// User 用户模型
type User struct {
    ID          int64     `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Email       string    `json:"email" db:"email"`
    Password    string    `json:"-" db:"password_hash"` // 登录密码
    Role        string    `json:"role" db:"role"`       // user, admin
    Status      string    `json:"status" db:"status"`   // active, disabled
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    TokenType    string `json:"token_type"` // Bearer
    ExpiresIn    int    `json:"expires_in"` // 秒
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// JWTClaims JWT 声明
type JWTClaims struct {
    UserID    int64  `json:"user_id"`
    UserEmail string `json:"user_email"`
    UserRole  string `json:"user_role"`
    jwt.RegisteredClaims
}
```

## 路由注册统一设计

### Controller 接口

所有 Controller 应实现统一的 `RegisterRoutes` 方法：

```go
// RegisterRoutes 定义路由注册接口
type RouteController interface {
    RegisterRoutes(r *gin.RouterGroup)
}
```

### 各 Controller 的路由分组

```go
// main.go 中的路由组织
func setupRoutes(router *gin.Engine, services *Services) {
    // API v1 组（管理接口）
    api := router.Group("/api/v1")

    // 认证接口（无需鉴权）
    authCtrl := controller.NewAuthController(services.Auth, services.JWT)
    authCtrl.RegisterRoutes(api)

    // 需要 JWT 鉴权的组
    jwtAuth := api.Group("")
    jwtAuth.Use(middleware.JWTAuth(services.JWT))

    // 需要 Admin 角色的组
    adminOnly := jwtAuth.Group("")
    adminOnly.Use(middleware.RequireAdmin())

    // Provider 管理
    providerCtrl := controller.NewProviderController(services.Provider)
    providerCtrl.RegisterRoutes(adminOnly)

    // Provider 运维
    reloadCtrl := controller.NewProviderReloadController(services.Provider)
    reloadCtrl.RegisterRoutes(adminOnly)

    // 用户管理（管理员可管理所有用户，普通用户可查看自己）
    userCtrl := controller.NewUserController(services.Auth)
    userCtrl.RegisterRoutes(jwtAuth)

    // 使用统计（普通用户可查看自己的，管理员可查看所有）
    usageCtrl := controller.NewUsageController(services.Usage)
    usageCtrl.RegisterRoutes(jwtAuth)

    // Chat API（使用 API Key 鉴权）
    v1 := router.Group("/v1")
    chatGroup := v1.Group("")
    chatGroup.Use(middleware.APIKeyAuth(services.Auth), middleware.TraceID())

    chatCtrl := controller.NewChatController(services.Router, services.Usage)
    chatCtrl.RegisterRoutes(chatGroup)
}
```

### UserController 路由注册

```go
func (c *UserController) RegisterRoutes(r *gin.RouterGroup) {
    users := r.Group("/users")
    {
        // 用户管理（仅管理员）
        users.POST("", c.CreateUser)
        users.GET("", c.ListUsers)
        users.PUT("/:id", c.UpdateUser)
        users.DELETE("/:id", c.DeleteUser)
        users.PATCH("/:id/status", c.UpdateUserStatus)

        // 获取用户信息（普通用户可获取自己的，管理员可获取任何人的）
        users.GET("/:id", c.GetUser)

        // API Key 管理（普通用户可管理自己的，管理员可管理任何人的）
        users.POST("/:id/api-keys", c.CreateAPIKey)
        users.GET("/:id/api-keys", c.ListAPIKeys)
        users.DELETE("/:id/api-keys/:key_id", c.RevokeAPIKey)
    }
}
```

**权限说明**：
- **管理员**：可创建、列出、更新、删除任何用户；可管理任何用户的 API Key
- **普通用户**：
  - 可获取自己的用户信息（`GET /api/v1/users/:自己的ID`）
  - 可管理自己的 API Key
  - 访问其他用户资源时返回 403

### UsageController 路由注册

```go
func (c *UsageController) RegisterRoutes(r *gin.RouterGroup) {
    r.GET("/usage", c.GetUsageStats)
    // 权限说明：
    // - 管理员：可查询任意用户或所有用户的统计（通过 user_id 参数过滤）
    // - 普通用户：只能查询自己的统计（自动过滤 user_id 参数）
}
```

### ChatController 路由注册

```go
func (c *ChatController) RegisterRoutes(r *gin.RouterGroup) {
    r.POST("/chat/completions", c.ChatCompletions)
    // 未来可能的扩展：
    // r.POST("/completions", c.Completions)
    // r.POST("/embeddings", c.Embeddings)
}
```

## JWT 配置设计

### 环境变量

```bash
# JWT 密钥（必须设置，生产环境使用强随机密钥）
JWT_SECRET=your-secret-key-here

# Access Token 过期时间（默认 15 分钟）
JWT_ACCESS_TOKEN_EXPIRES_IN=15m

# Refresh Token 过期时间（默认 7 天）
JWT_REFRESH_TOKEN_EXPIRES_IN=168h

# Token 发行者
JWT_ISSUER=courier-gateway
```

### JWTService 实现

```go
type JWTService interface {
    // GenerateAccessToken 生成访问令牌
    GenerateAccessToken(user *model.User) (string, error)

    // GenerateRefreshToken 生成刷新令牌
    GenerateRefreshToken(userID int64) (string, error)

    // ValidateAccessToken 验证访问令牌
    ValidateAccessToken(token string) (*model.JWTClaims, error)

    // ValidateRefreshToken 验证刷新令牌
    ValidateRefreshToken(token string) (*model.JWTClaims, error)
}
```

## 迁移策略

### 向后兼容方案

在迁移期间，可以同时支持两种鉴权方式：

```go
// AdminAuth 支持两种方式
func FlexibleAdminAuth(jwtSvc service.JWTService) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // 先尝试 JWT 鉴权
        authHeader := ctx.GetHeader("Authorization")
        if strings.HasPrefix(authHeader, "Bearer ") {
            token := strings.TrimPrefix(authHeader, "Bearer ")
            if claims, err := jwtSvc.ValidateAccessToken(token); err == nil {
                ctx.Set("user_id", claims.UserID)
                ctx.Set("user_role", claims.UserRole)
                ctx.Next()
                return
            }
        }

        // 降级到 Admin API Key
        adminKey := ctx.GetHeader("X-Admin-API-Key")
        if adminKey == os.Getenv("ADMIN_API_KEY") {
            ctx.Next()
            return
        }

        ctx.JSON(401, gin.H{"error": "unauthorized"})
        ctx.Abort()
    }
}
```

## 安全考虑

1. **密码存储**：使用 bcrypt 进行密码哈希
2. **JWT 密钥**：从环境变量读取，不硬编码
3. **Token 过期**：Access Token 短期有效，Refresh Token 长期可撤销
4. **HTTPS**：生产环境强制 HTTPS
5. **速率限制**：登录接口添加速率限制防止暴力破解
