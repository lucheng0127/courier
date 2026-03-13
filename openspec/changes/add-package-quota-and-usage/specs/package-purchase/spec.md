## ADDED Requirements

### Requirement: 用户套餐数据模型

系统 SHALL 提供用户套餐（UserPackage）数据模型用于记录用户购买的套餐。

#### Scenario: 创建用户套餐

- **WHEN** 用户购买套餐时
- **THEN** 系统生成唯一的用户套餐 ID
- **AND** 记录用户 ID `user_id`
- **AND** 记录套餐 ID `package_id`
- **AND** 设置状态为 `active`
- **AND** 计算过期时间 `expires_at` = 当前时间 + 套餐有效期天数
- **AND** 记录激活时间 `activated_at`
- **AND** 记录购买时间 `created_at`

#### Scenario: 用户套餐状态

- **WHEN** 查看用户套餐状态时
- **THEN** 状态可以是 `active`（激活中）
- **AND** 状态可以是 `expired`（已过期）

#### Scenario: 套餐叠加

- **WHEN** 用户购买已拥有的套餐时
- **THEN** 允许重复购买
- **AND** 创建新的用户套餐记录
- **AND** 配额累加计算

### Requirement: 套餐购买 API

系统 SHALL 提供套餐购买 API，允许用户购买套餐（含模拟支付）。

#### Scenario: 购买套餐请求

- **WHEN** POST 请求到 `/api/v1/user/packages`
- **THEN** 验证用户已登录（JWT 或 API Key）
- **AND** 请求体包含 `package_id`（必填）
- **AND** 验证套餐状态为 `online`
- **AND** 模拟支付成功（直接通过）

#### Scenario: 购买套餐成功

- **WHEN** 套餐购买成功时
- **THEN** 创建用户套餐记录
- **AND** 计算过期时间
- **AND** 返回 201 Created
- **AND** 响应体包含用户套餐信息
- **AND** 响应体包含过期时间

#### Scenario: 购买不可用的套餐

- **GIVEN** 套餐状态为 `draft` 或 `offline`
- **WHEN** 尝试购买套餐
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明套餐不可购买

#### Scenario: 购买不存在的套餐

- **WHEN** 尝试购买不存在的套餐
- **THEN** 返回 404 Not Found
- **AND** 错误信息说明套餐不存在

### Requirement: 我的套餐查询 API

系统 SHALL 提供我的套餐查询 API，允许用户查看自己购买的套餐。

#### Scenario: 查询我的套餐列表

- **WHEN** GET 请求到 `/api/v1/user/packages`
- **THEN** 验证用户已登录
- **AND** 仅返回当前用户购买的套餐
- **AND** 包含套餐基本信息
- **AND** 包含配额使用情况
- **AND** 包含过期时间
- **AND** 支持按状态筛选（`?status=active|expired`）

#### Scenario: 套餐列表排序

- **WHEN** 返回套餐列表时
- **THEN** 按创建时间降序排列（最新的在前）
- **AND** 激活的套餐排在过期的套餐前面

### Requirement: 用户套餐详情查询 API

系统 SHALL 提供用户套餐详情查询 API。

#### Scenario: 查询套餐详情

- **WHEN** GET 请求到 `/api/v1/user/packages/:id`
- **THEN** 验证用户已登录
- **AND** 验证套餐属于当前用户
- **AND** 返回套餐详细信息
- **AND** 包含套餐基本信息
- **AND** 包含配额使用详情（按 Provider 分组）
- **AND** 包含过期时间

#### Scenario: 查询他人套餐

- **GIVEN** 用户 A 查询用户 B 的套餐
- **WHEN** GET 请求到 `/api/v1/user/packages/:id`
- **THEN** 返回 403 Forbidden
- **AND** 错误信息说明无权访问

### Requirement: 套餐使用统计查询 API

系统 SHALL 提供套餐使用统计查询 API。

#### Scenario: 查询套餐使用统计

- **WHEN** GET 请求到 `/api/v1/user/packages/:id/usage`
- **THEN** 验证用户已登录
- **AND** 验证套餐属于当前用户
- **AND** 返回套餐使用统计
- **AND** 包含总使用 Token 数
- **AND** 包含剩余配额
- **AND** 包含按 Provider 分组的使用统计
- **AND** 包含按天分组的使用趋势

#### Scenario: 使用统计包含剩余配额

- **WHEN** 返回使用统计时
- **THEN** 包含 `remaining_tokens` 字段
- **AND** 计算方式：`套餐配额 - 已使用配额`
- **AND** 如果配额为 0（不限量），显示为 "不限量"

#### Scenario: 使用统计包含使用记录

- **WHEN** 返回使用统计时
- **THEN** 包含 `usage_records` 数组
- **AND** 每条记录包含请求时间、Provider、使用的 Token 数
- **AND** 支持分页

### Requirement: 套餐过期处理

系统 SHALL 自动处理过期套餐。

#### Scenario: 套餐过期状态更新

- **WHEN** 用户套餐过期时间到达
- **THEN** 系统自动将套餐状态更新为 `expired`
- **AND** 过期套餐无法继续使用配额

#### Scenario: 过期套餐配额不可用

- **GIVEN** 用户套餐状态为 `expired`
- **WHEN** 用户发起 Chat 请求
- **THEN** 配额检测中间件忽略该套餐的配额

### Requirement: 管理员查询用户套餐 API

系统 SHALL 提供管理员查询用户套餐的 API。

#### Scenario: 管理员查询用户套餐

- **GIVEN** 用户具有管理员权限
- **WHEN** GET 请求到 `/api/v1/admin/users/:user_id/packages`
- **THEN** 返回指定用户的所有套餐
- **AND** 包含配额使用情况
- **AND** 支持按状态筛选

#### Scenario: 管理员查询套餐购买记录

- **GIVEN** 用户具有管理员权限
- **WHEN** GET 请求到 `/api/v1/admin/packages/:package_id/purchases`
- **THEN** 返回购买该套餐的所有用户
- **AND** 包含购买时间
- **AND** 包含套餐状态
- **AND** 支持分页
