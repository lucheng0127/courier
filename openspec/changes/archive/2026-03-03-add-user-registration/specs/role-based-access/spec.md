# role-based-access 规范修改

## Purpose

移除管理员创建用户的功能，用户现在通过自主注册创建。

## MODIFIED Requirements

### Requirement: 用户角色

系统 SHALL 支持用户角色区分，包含 `admin` 和 `user` 两种角色。

#### Scenario: Admin 角色定义

- **GIVEN** 用户角色为 `admin`
- **THEN** 该用户可以访问所有管理接口
- **AND** 可以读取、更新、删除其他用户
- **AND** 可以管理 Providers（创建、读取、更新、删除、重载、启用、禁用）
- **AND** 可以查询所有用户的使用统计
- **AND** 可以为任何用户创建 API Key

#### Scenario: User 角色定义

- **GIVEN** 用户角色为 `user`
- **THEN** 该用户只能访问自己的资源
- **AND** 可以使用自己的 API Key 调用 Chat API
- **AND** 可以查询自己的使用统计
- **AND** 可以查看自己的用户信息
- **AND** 不能访问管理接口
- **AND** 不能管理其他用户

### Requirement: 用户管理权限

系统 SHALL 只允许管理员用户查看和管理用户账户状态，用户通过自主注册创建。

#### Scenario: 查看用户列表

- **GIVEN** 用户已登录
- **AND** 用户角色为 `admin`
- **WHEN** 发送 GET 请求到 `/api/v1/users`
- **THEN** 返回所有用户的列表

- **GIVEN** 用户已登录
- **AND** 用户角色为 `user`
- **WHEN** 发送 GET 请求到 `/api/v1/users`
- **THEN** 返回 403 状态码

#### Scenario: 管理员查看任意用户信息

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

#### Scenario: 管理员为用户创建 API Key

- **GIVEN** 管理员用户已登录
- **WHEN** 发送 POST 请求到 `/api/v1/users/:id/api-keys`
- **THEN** 可以为任何用户创建 API Key

#### Scenario: 普通用户为自己创建 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 POST 请求到 `/api/v1/users/123/api-keys`
- **THEN** 可以创建 API Key

#### Scenario: 普通用户尝试为他人创建 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 POST 请求到 `/api/v1/users/456/api-keys`
- **THEN** 返回 403 状态码
- **AND** 响应体包含权限错误信息

#### Scenario: 普通用户查询自己的 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/123/api-keys`
- **THEN** 返回该用户的 API Key 列表

#### Scenario: 普通用户尝试查询他人 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 GET 请求到 `/api/v1/users/456/api-keys`
- **THEN** 返回 403 状态码

#### Scenario: 普通用户撤销自己的 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 DELETE 请求到 `/api/v1/users/123/api-keys/:key_id`
- **THEN** 可以撤销该 API Key

#### Scenario: 普通用户尝试撤销他人 API Key

- **GIVEN** 普通用户已登录
- **AND** 用户 ID 为 123
- **WHEN** 发送 DELETE 请求到 `/api/v1/users/456/api-keys/:key_id`
- **THEN** 返回 403 状态码
