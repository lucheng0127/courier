# public-provider-access Specification Delta

## MODIFIED Requirements

### Requirement: Provider 查询权限

系统 SHALL 允许所有认证用户（管理员和普通用户）查询 Provider 列表，但响应内容根据用户角色不同。

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

#### Scenario: 未认证用户查询 Provider 列表

- **GIVEN** 用户未通过 JWT 鉴权
- **WHEN** 发送 GET 请求到 `/api/v1/providers`
- **THEN** 返回 401 状态码
- **AND** 响应体包含认证错误信息

#### Scenario: 普通用户尝试查询单个 Provider

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 发送 GET 请求到 `/api/v1/providers/:name`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

#### Scenario: 管理员用户查询单个 Provider

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `admin`
- **WHEN** 发送 GET 请求到 `/api/v1/providers/:name`
- **THEN** 返回 200 状态码或 404（如果不存在）
- **AND** 响应体包含该 Provider 的完整信息

## ADDED Requirements

### Requirement: Provider 启用状态过滤

系统 SHALL 支持通过 `enabled` 参数过滤 Provider 列表。

#### Scenario: 查询已启用的 Provider

- **GIVEN** 用户已通过 JWT 鉴权
- **WHEN** 发送 GET 请求到 `/api/v1/providers?enabled=true`
- **THEN** 返回 200 状态码
- **AND** 响应体只包含 `enabled=true` 的 Provider

#### Scenario: 查询已禁用的 Provider

- **GIVEN** 用户已通过 JWT 鉴权
- **WHEN** 发送 GET 请求到 `/api/v1/providers?enabled=false`
- **THEN** 返回 200 状态码
- **AND** 响应体只包含 `enabled=false` 的 Provider

#### Scenario: 不传递 enabled 参数

- **GIVEN** 用户已通过 JWT 鉴权
- **WHEN** 发送 GET 请求到 `/api/v1/providers`（不带 enabled 参数）
- **THEN** 返回 200 状态码
- **AND** 响应体包含所有 Provider（无论启用状态）

#### Scenario: 无效的 enabled 参数

- **GIVEN** 用户已通过 JWT 鉴权
- **WHEN** 发送 GET 请求到 `/api/v1/providers?enabled=invalid`
- **THEN** 返回 400 状态码
- **AND** 响应体包含参数错误信息

### Requirement: Provider 模型列表接口

系统 SHALL 提供接口查询 Provider 支持的模型列表。

#### Scenario: 获取 Provider 模型列表

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** Provider 存在且配置了 `fallback_models`
- **WHEN** 发送 GET 请求到 `/api/v1/providers/:name/models`
- **THEN** 返回 200 状态码
- **AND** 响应体包含：
  - `name`: Provider 名称
  - `type`: Provider 类型
  - `models`: 字符串数组，包含模型名称列表
- **AND** `models` 数组来源于 Provider 的 `fallback_models` 配置

#### Scenario: Provider 不存在

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** Provider 不存在
- **WHEN** 发送 GET 请求到 `/api/v1/providers/:name/models`
- **THEN** 返回 404 状态码
- **AND** 响应体包含错误信息

#### Scenario: Provider 未配置 fallback_models

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** Provider 存在但 `fallback_models` 为空
- **WHEN** 发送 GET 请求到 `/api/v1/providers/:name/models`
- **THEN** 返回 200 状态码
- **AND** 响应体中的 `models` 为空数组

#### Scenario: 未认证用户访问模型列表

- **GIVEN** 用户未通过 JWT 鉴权
- **WHEN** 发送 GET 请求到 `/api/v1/providers/:name/models`
- **THEN** 返回 401 状态码
- **AND** 响应体包含认证错误信息

### Requirement: Provider 管理权限保持不变

系统 SHALL 继续只允许管理员用户管理 Provider（创建、更新、删除）。

#### Scenario: 普通用户尝试创建 Provider

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 发送 POST 请求到 `/api/v1/providers`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

#### Scenario: 普通用户尝试更新 Provider

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 发送 PUT 请求到 `/api/v1/providers/:name`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

#### Scenario: 普通用户尝试删除 Provider

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `user`
- **WHEN** 发送 DELETE 请求到 `/api/v1/providers/:name`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

#### Scenario: 管理员用户创建 Provider

- **GIVEN** 用户已通过 JWT 鉴权
- **AND** 用户角色为 `admin`
- **WHEN** 发送 POST 请求到 `/api/v1/providers`
- **THEN** 请求成功处理（返回 201 或其他适当的成功状态码）
