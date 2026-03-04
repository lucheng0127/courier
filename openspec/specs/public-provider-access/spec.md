# public-provider-access Specification

## Purpose

定义 Provider 公开访问接口的规范，允许所有认证用户查询 Provider 信息和模型列表，同时支持按启用状态过滤。

## Requirements
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

