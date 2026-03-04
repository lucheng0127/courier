# role-based-access Specification Delta

## MODIFIED Requirements

### Requirement: Provider 管理权限

系统 SHALL 允许管理员用户管理 Provider，所有认证用户可查询 Provider 列表和模型列表。

#### Scenario: 普通用户查询 Provider 列表

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 发送 GET 请求到 `/api/v1/providers`
- **THEN** 返回 200 状态码
- **AND** 响应体只包含非敏感信息：
  - `name`: Provider 名称
  - `type`: Provider 类型
  - `base_url`: API 地址
  - `enabled`: 启用状态
  - `fallback_models`: 支持的模型列表
- **AND** 响应体不包含敏感信息：
  - `api_key`: API 密钥
  - `timeout`: 超时配置
  - `is_running`: 运行状态
  - `extra_config`: 额外配置

#### Scenario: 管理员用户查询 Provider 列表

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `admin`
- **WHEN** 发送 GET 请求到 `/api/v1/providers`
- **THEN** 返回 200 状态码
- **AND** 响应体包含完整的 Provider 信息（包括敏感信息）

#### Scenario: 普通用户查询 Provider 模型列表

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 发送 GET 请求到 `/api/v1/providers/:name/models`
- **THEN** 返回 200 状态码
- **AND** 响应体包含 Provider 名称、类型和模型列表
