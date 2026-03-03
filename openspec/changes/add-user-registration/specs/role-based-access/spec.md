# role-based-access 规范修改

## Purpose

移除管理员创建用户的功能，用户现在通过自主注册创建。

## MODIFIED Requirements

### Requirement: 用户管理权限

系统 SHALL 只允许管理员用户查看和管理用户账户状态，但不允许创建新用户。

#### Scenario: 查看用户列表

- **GIVEN** 用户已登录
- **AND** 用户角色为 `admin`
- **WHEN** 发送 GET 请求到 `/api/v1/users`
- **THEN** 返回所有用户的列表

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 GET 请求到 `/api/v1/users`
- **THEN** 返回 403 状态码

#### Scenario: 查看任意用户信息

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 GET 请求到 `/api/v1/users/:id`
- **THEN** 可以查看任何用户的信息

#### Scenario: 普通用户查看自己的信息

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/123`
- **THEN** 返回该用户的信息

#### Scenario: 普通用户尝试查看他人信息

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/456`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

#### Scenario: 更新用户信息

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 PUT 请求到 `/api/v1/users/:id`
- **THEN** 可以更新用户信息

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 PUT 请求到 `/api/v1/users/:id`
- **THEN** 返回 403 状态码

#### Scenario: 删除用户

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 DELETE 请求到 `/api/v1/users/:id`
- **THEN** 可以删除用户

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 DELETE 请求到 `/api/v1/users/:id`
- **THEN** 返回 403 状态码

#### Scenario: 更新用户状态

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 PATCH 请求到 `/api/v1/users/:id/status`
- **THEN** 可以更新用户状态（active/disabled）

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 PATCH 请求到 `/api/v1/users/:id/status`
- **THEN** 返回 403 状态码

## REMOVED Requirements

### Requirement: 创建用户（管理员）

管理员不再能够通过 API 创建用户。

#### Scenario: 移除的创建用户接口

- **GIVEN** 管理员用户已登录
- **WHEN** 尝试发送 POST 请求到 `/api/v1/users`
- **THEN** 返回 404 或 405 状态码
- **AND** 不创建新用户

#### Scenario: 用户通过自主注册创建

- **GIVEN** 系统正常运行
- **AND** 提供有效的邮箱和密码
- **WHEN** 发送 POST 请求到 `/api/v1/auth/register`
- **THEN** 创建新用户
- **AND** 新用户角色为 "user"
- **AND** 新用户状态为 "active"
