# jwt-auth Specification

## Purpose
TBD - created by archiving change unify-api-and-jwt-auth. Update Purpose after archive.
## Requirements
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

