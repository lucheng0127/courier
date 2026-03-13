## ADDED Requirements

### Requirement: 套餐数据模型

系统 SHALL 提供套餐（Package）数据模型用于管理可购买的套餐产品。

#### Scenario: 创建套餐基本字段

- **WHEN** 创建新套餐时
- **THEN** 系统生成唯一的套餐 ID
- **AND** 存储套餐名称 `name`（必填）
- **AND** 存储套餐描述 `description`
- **AND** 存储套餐价格 `price`（分为单位）
- **AND** 设置初始状态为 `draft`
- **AND** 存储有效期天数 `validity_days`
- **AND** 记录创建时间和更新时间

#### Scenario: 套餐状态定义

- **WHEN** 查看套餐状态时
- **THEN** 状态可以是 `draft`（草稿）
- **AND** 状态可以是 `online`（已上架）
- **AND** 状态可以是 `offline`（已下架）

#### Scenario: 套餐名称唯一性

- **WHEN** 创建套餐时名称已存在
- **THEN** 返回 409 Conflict 错误
- **AND** 错误信息说明套餐名称已被使用

### Requirement: 套餐配额数据模型

系统 SHALL 提供套餐配额（PackageQuota）数据模型用于定义套餐包含的 Provider Token 配额。

#### Scenario: 创建套餐配额

- **WHEN** 为套餐添加配额时
- **THEN** 系统生成唯一的配额 ID
- **AND** 关联套餐 ID `package_id`
- **AND** 存储 Provider 名称 `provider_name`
- **AND** 存储 Token 配额限制 `token_limit`
- **AND** Provider 名称可以是 `*`（表示所有 Provider）

#### Scenario: 不限量配额

- **WHEN** 套餐配额的 `token_limit` 设置为 0 时
- **THEN** 表示该 Provider 不限量使用

#### Scenario: 多 Provider 配额

- **WHEN** 创建套餐时
- **THEN** 可以为多个 Provider 分别设置配额
- **AND** 可以同时设置特定 Provider 配额和全局配额（`*`）
- **AND** 优先使用特定 Provider 的配额

### Requirement: 套餐创建 API

系统 SHALL 提供套餐创建 API，允许管理员创建新套餐。

#### Scenario: 创建套餐请求

- **WHEN** POST 请求到 `/api/v1/admin/packages`
- **THEN** 验证用户具有管理员权限
- **AND** 请求体包含 `name`（必填）
- **AND** 请求体包含 `description`（可选）
- **AND** 请求体包含 `price`（必填，单位：分）
- **AND** 请求体包含 `validity_days`（必填，天数）
- **AND** 请求体包含 `quotas` 数组（必填），每个元素包含 `provider_name` 和 `token_limit`

#### Scenario: 创建套餐成功

- **WHEN** 套餐创建成功时
- **THEN** 返回 201 Created
- **AND** 响应体包含完整的套餐信息
- **AND** 响应体包含套餐配额列表
- **AND** 套餐状态为 `draft`

#### Scenario: 创建套餐失败

- **WHEN** 请求参数无效时
- **THEN** 返回 400 Bad Request
- **AND** 响应体包含错误详情

### Requirement: 套餐更新 API

系统 SHALL 提供套餐更新 API，仅允许更新草稿状态的套餐。

#### Scenario: 更新草稿套餐

- **GIVEN** 套餐状态为 `draft`
- **WHEN** PUT 请求到 `/api/v1/admin/packages/:id`
- **THEN** 允许更新套餐信息
- **AND** 返回 200 OK 和更新后的套餐信息

#### Scenario: 更新非草稿套餐

- **GIVEN** 套餐状态为 `online` 或 `offline`
- **WHEN** 尝试更新套餐
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明只有草稿状态的套餐可以更新

### Requirement: 套餐上架 API

系统 SHALL 提供套餐上架 API，允许管理员上架套餐使其可被用户购买。

#### Scenario: 上架草稿套餐

- **GIVEN** 套餐状态为 `draft`
- **WHEN** POST 请求到 `/api/v1/admin/packages/:id/online`
- **THEN** 将套餐状态更新为 `online`
- **AND** 返回 200 OK

#### Scenario: 上架已下架套餐

- **GIVEN** 套餐状态为 `offline`
- **WHEN** POST 请求到 `/api/v1/admin/packages/:id/online`
- **THEN** 将套餐状态更新为 `online`
- **AND** 返回 200 OK

#### Scenario: 上架已上架套餐

- **GIVEN** 套餐状态为 `online`
- **WHEN** 尝试上架套餐
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明套餐已是上架状态

### Requirement: 套餐下架 API

系统 SHALL 提供套餐下架 API，允许管理员下架套餐使其无法被购买。

#### Scenario: 下架已上架套餐

- **GIVEN** 套餐状态为 `online`
- **WHEN** POST 请求到 `/api/v1/admin/packages/:id/offline`
- **THEN** 将套餐状态更新为 `offline`
- **AND** 返回 200 OK
- **AND** 已购买该套餐的用户不受影响

#### Scenario: 下架已下架套餐

- **GIVEN** 套餐状态为 `offline`
- **WHEN** 尝试下架套餐
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明套餐已是下架状态

### Requirement: 套餐删除 API

系统 SHALL 提供套餐删除 API，仅允许删除草稿状态的套餐。

#### Scenario: 删除草稿套餐

- **GIVEN** 套餐状态为 `draft`
- **AND** 没有用户购买该套餐
- **WHEN** DELETE 请求到 `/api/v1/admin/packages/:id`
- **THEN** 删除套餐记录
- **AND** 删除关联的配额记录
- **AND** 返回 204 No Content

#### Scenario: 删除非草稿套餐

- **GIVEN** 套餐状态为 `online` 或 `offline`
- **WHEN** 尝试删除套餐
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明只能删除草稿状态的套餐

#### Scenario: 删除已购买套餐

- **GIVEN** 套餐状态为 `draft`
- **AND** 有用户已购买该套餐
- **WHEN** 尝试删除套餐
- **THEN** 返回 400 Bad Request
- **AND** 错误信息说明已有用户购买该套餐，无法删除

### Requirement: 套餐查询 API

系统 SHALL 提供套餐查询 API，允许管理员查询所有套餐。

#### Scenario: 查询套餐列表

- **WHEN** GET 请求到 `/api/v1/admin/packages`
- **THEN** 验证管理员权限
- **AND** 返回所有套餐的列表
- **AND** 每个套餐包含配额信息
- **AND** 支持按状态筛选（`?status=draft|online|offline`）

#### Scenario: 查询套餐详情

- **WHEN** GET 请求到 `/api/v1/admin/packages/:id`
- **THEN** 返回指定套餐的完整信息
- **AND** 包含配额列表
- **AND** 包含购买统计（购买人数、激活人数）

### Requirement: 公开套餐查询 API

系统 SHALL 提供公开套餐查询 API，允许用户查询可购买的套餐列表。

#### Scenario: 查询可购买套餐

- **WHEN** GET 请求到 `/api/v1/packages`
- **THEN** 仅返回状态为 `online` 的套餐
- **AND** 不返回敏感信息（如成本价）
- **AND** 支持分页

#### Scenario: 查询套餐详情

- **WHEN** GET 请求到 `/api/v1/packages/:id`
- **THEN** 返回套餐详情
- **AND** 仅当套餐状态为 `online` 时返回
- **AND** 包含配额信息
- **AND** 不返回敏感信息
