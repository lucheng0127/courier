# jwt-auth Specification Delta

## Purpose

定义基于 JWT 的身份认证机制，用于除 `/chat/completions` 外的所有管理接口。

## ADDED Requirements

### Requirement: 管理接口 JWT 鉴权

系统 SHALL 使用 JWT Token 对所有管理接口进行身份认证。

#### Scenario: 用户登录获取 Token

- **GIVEN** 数据库中存在状态为 `active` 的用户
- **AND** 用户角色为 `admin` 或 `user`
- **WHEN** 客户端发送 POST 请求到 `/api/v1/auth/login`
- **AND** 请求体包含有效的 `email` 和 `password`
- **THEN** 系统验证密码正确性
- **AND** 验证用户状态为 `active`
- **AND** 生成 Access Token（有效期 15 分钟）
- **AND** 生成 Refresh Token（有效期 7 天）
- **AND** 返回 200 状态码
- **AND** 响应体包含：
  - `access_token`: JWT 访问令牌
  - `refresh_token`: 刷新令牌
  - `token_type`: "Bearer"
  - `expires_in`: Access Token 过期秒数

#### Scenario: 登录失败 - 用户不存在

- **GIVEN** 数据库中不存在该邮箱的用户
- **WHEN** 客户端发送登录请求
- **THEN** 返回 401 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Invalid email or password"
  - `type`: "authentication_error"

#### Scenario: 登录失败 - 密码错误

- **GIVEN** 数据库中存在该用户
- **AND** 密码验证失败
- **WHEN** 客户端发送登录请求
- **THEN** 返回 401 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Invalid email or password"
  - `type`: "authentication_error"

#### Scenario: 登录失败 - 用户状态非活跃

- **GIVEN** 数据库中存在该用户
- **AND** 用户状态为 `disabled`
- **WHEN** 客户端发送登录请求
- **THEN** 返回 403 状态码
- **AND** 响应体包含错误信息：
  - `message`: "User account is disabled"
  - `type`: "permission_error"

#### Scenario: 使用 Access Token 访问管理接口

- **GIVEN** 用户已成功登录并获得有效的 Access Token
- **WHEN** 客户端发送请求到任意管理接口
- **AND** 请求头包含 `Authorization: Bearer <access_token>`
- **THEN** JWT 中间件验证 Token 签名
- **AND** 验证 Token 未过期
- **AND** 从 JWT Claims 中提取用户信息
- **AND** 将 `user_id`、`user_email`、`user_role` 注入到请求上下文
- **AND** 请求继续传递到后续中间件和处理函数

#### Scenario: Access Token 过期

- **GIVEN** 客户端持有已过期的 Access Token
- **WHEN** 客户端使用该 Token 访问管理接口
- **THEN** 返回 401 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Access token is expired"
  - `type`: "authentication_error"

#### Scenario: Access Token 无效

- **GIVEN** 客户端持有伪造的 Access Token
- **WHEN** 客户端使用该 Token 访问管理接口
- **THEN** 返回 401 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Invalid access token"
  - `type`: "authentication_error"

#### Scenario: 刷新 Access Token

- **GIVEN** 用户持有有效的 Refresh Token
- **WHEN** 客户端发送 POST 请求到 `/api/v1/auth/refresh`
- **AND** 请求体包含 `refresh_token`
- **THEN** 系统验证 Refresh Token 的有效性
- **AND** 验证 Refresh Token 未过期
- **AND** 生成新的 Access Token
- **AND** 生成新的 Refresh Token
- **AND** 返回 200 状态码
- **AND** 响应体包含新的 `access_token` 和 `refresh_token`

#### Scenario: Refresh Token 无效

- **GIVEN** 客户端持有无效或过期的 Refresh Token
- **WHEN** 客户端发送刷新请求
- **THEN** 返回 401 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Invalid or expired refresh token"
  - `type`: "authentication_error"

#### Scenario: 缺少 Authorization Header

- **WHEN** 客户端访问需要 JWT 鉴权的接口
- **AND** 请求头中缺少 `Authorization` 字段
- **THEN** 返回 401 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Missing authorization header"
  - `type`: "authentication_error"

#### Scenario: JWT Claims 结构

- **WHEN** 系统生成 JWT Access Token
- **THEN** Token Payload（Claims）包含以下字段：
  - `user_id`: 用户 ID（数字）
  - `user_email`: 用户邮箱（字符串）
  - `user_role`: 用户角色（"admin" 或 "user"）
  - `iss`: Token 发行者（"courier-gateway"）
  - `exp`: 过期时间（Unix 时间戳）
  - `iat`: 签发时间（Unix 时间戳）

### Requirement: JWT 配置

系统 SHALL 支持通过环境变量配置 JWT 相关参数。

#### Scenario: 基本配置

- **GIVEN** 系统启动时
- **THEN** 从环境变量 `JWT_SECRET` 读取 JWT 签名密钥
- **AND** 如果 `JWT_SECRET` 未设置，系统应拒绝启动并返回错误
- **AND** 从环境变量 `JWT_ACCESS_TOKEN_EXPIRES_IN` 读取 Access Token 有效期（默认 15m）
- **AND** 从环境变量 `JWT_REFRESH_TOKEN_EXPIRES_IN` 读取 Refresh Token 有效期（默认 168h）
- **AND** 从环境变量 `JWT_ISSUER` 读取 Token 发行者（默认 "courier-gateway"）

### Requirement: Chat API 保持 API Key 鉴权

`/v1/chat/completions` 接口 SHALL 继续使用 API Key 鉴权，不受 JWT 系统影响。

#### Scenario: Chat API 使用 API Key

- **GIVEN** 用户拥有有效的 API Key
- **WHEN** 客户端发送 POST 请求到 `/v1/chat/completions`
- **AND** 请求头包含 `Authorization: Bearer <api_key>`
- **THEN** 使用现有的 API Key 鉴权中间件
- **AND** 不进行 JWT 验证
- **AND** 请求正常处理

## ADDED Requirements

### Requirement: 密码哈希存储

系统 SHALL 使用安全的哈希算法存储用户密码。

#### Scenario: 密码哈希

- **WHEN** 创建或更新用户密码时
- **THEN** 使用 bcrypt 算法进行哈希
- **AND** 使用 cost factor 12
- **AND** 只存储哈希值，不存储明文密码

### Requirement: 登录速率限制

系统 SHALL 对登录接口实施速率限制，防止暴力破解攻击。

#### Scenario: 速率限制配置

- **GIVEN** 系统配置了登录速率限制
- **WHEN** 同一 IP 地址在 1 分钟内尝试登录超过 5 次
- **THEN** 返回 429 状态码
- **AND** 响应体包含错误信息：
  - `message`: "Too many login attempts, please try again later"
  - `type`: "rate_limit_error"
