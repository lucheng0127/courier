## ADDED Requirements

### Requirement: 套餐使用记录数据模型

系统 SHALL 提供套餐使用记录（PackageUsage）数据模型用于跟踪套餐配额的使用情况。

#### Scenario: 创建套餐使用记录

- **WHEN** Chat API 调用完成后
- **THEN** 记录关联的用户套餐 ID `user_package_id`
- **AND** 记录用户 ID `user_id`
- **AND** 记录套餐 ID `package_id`
- **AND** 记录使用的 Provider `provider_name`
- **AND** 记录实际使用的 Token 数 `tokens_used`
- **AND** 关联原始使用记录 ID `usage_record_id`
- **AND** 关联请求 ID `request_id`
- **AND** 记录创建时间

### Requirement: 配额检测中间件

系统 SHALL 在 Chat API 中集成配额检测中间件，确保用户有足够的套餐配额才能发起请求。

#### Scenario: 中间件执行时机

- **WHEN** 请求到达 `/v1/chat/completions` 端点
- **THEN** 在路由处理函数之前执行配额检测
- **AND** 检测通过后继续处理请求
- **AND** 检测失败时返回错误并终止请求

#### Scenario: 从上下文获取用户信息

- **WHEN** 配额检测中间件执行时
- **THEN** 从 Gin Context 获取 `user_id`
- **AND** 从请求参数获取 `model` 参数
- **AND** 解析 model 获取 `provider_name`

#### Scenario: 查询用户有效套餐

- **WHEN** 检测用户配额时
- **THEN** 查询用户所有状态为 `active` 的套餐
- **AND** 筛选过期时间晚于当前时间的套餐
- **AND** 筛选包含请求 Provider 配额的套餐
- **AND** 筛选包含全局配额（`*`）的套餐

#### Scenario: 套餐配额筛选逻辑

- **GIVEN** 用户有多个套餐
- **AND** 请求的 Provider 为 `openai`
- **WHEN** 筛选可用套餐时
- **THEN** 包含明确配置了 `openai` Provider 配额的套餐
- **AND** 包含配置了全局配额（`*`）的套餐
- **AND** 不包含仅配置其他 Provider 配额的套餐

#### Scenario: 套餐排序规则

- **WHEN** 有多个可用套餐时
- **THEN** 按过期时间升序排列（即将过期的优先）
- **AND** 优先使用特定 Provider 的配额
- **AND** 最后使用全局配额（`*`）

#### Scenario: 计算剩余配额

- **WHEN** 检查套餐剩余配额时
- **THEN** 查询该套餐该 Provider 的已使用 Token 总数
- **AND** 剩余配额 = 套餐配额限制 - 已使用配额
- **AND** 如果套餐配额限制为 0，表示不限量
- **AND** 考虑所有可用套餐的总剩余配额

#### Scenario: 配额充足

- **GIVEN** 用户有足够的总剩余配额
- **WHEN** Chat API 请求通过配额检测
- **THEN** 将可用的套餐列表注入到 Context（用于后续扣减）
- **AND** 继续处理请求

#### Scenario: 配额不足

- **GIVEN** 用户所有套餐的剩余配额不足
- **WHEN** Chat API 请求通过配额检测
- **THEN** 返回 429 Too Many Requests
- **AND** 响应体包含错误信息：
  - `message`: "套餐配额已用完，请购买新套餐"
  - `type`: "quota_exceeded"
  - `provider`: 请求的 Provider 名称
- **AND** 终止请求处理

#### Scenario: 无可用套餐

- **GIVEN** 用户没有购买任何套餐
- **WHEN** Chat API 请求通过配额检测
- **THEN** 返回 429 Too Many Requests
- **AND** 响应体包含错误信息：
  - `message`: "您还没有购买套餐，请先购买套餐"
  - `type`: "no_package"
- **AND** 终止请求处理

#### Scenario: JWT 认证用户的配额检测

- **GIVEN** 用户通过 JWT 认证
- **WHEN** 发起 Chat API 请求
- **THEN** 从 Context 获取 `user_id`
- **AND** 执行配额检测
- **AND** 不检查 API Key 相关的配额

#### Scenario: API Key 认证用户的配额检测

- **GIVEN** 用户通过 API Key 认证
- **WHEN** 发起 Chat API 请求
- **THEN** 从 Context 获取 `user_id`（关联的用户 ID）
- **AND** 执行配额检测
- **AND** 检测该用户的套餐配额

#### Scenario: 管理员用户配额检测

- **GIVEN** 用户角色为 `admin`
- **WHEN** 发起 Chat API 请求
- **THEN** 跳过配额检测
- **AND** 继续处理请求

### Requirement: 配额扣减逻辑

系统 SHALL 在 Chat API 请求完成后扣减套餐配额。

#### Scenario: 请求成功后扣减配额

- **WHEN** Chat API 请求成功完成
- **THEN** 从响应中获取实际使用的 Token 数
- **AND** 从 Context 获取之前检测的可用套餐列表
- **AND** 按优先级顺序扣减配额
- **AND** 创建套餐使用记录
- **AND** 使用独立 Context 异步执行，避免影响请求响应

#### Scenario: 配额扣减顺序

- **GIVEN** 用户有多个可用套餐
- **AND** 总共使用了 1000 Token
- **WHEN** 扣减配额时
- **THEN** 优先扣减即将过期的套餐配额
- **AND** 如果第一个套餐配额不足，扣减第二个套餐
- **AND** 依此类推，直到扣减完所有使用的 Token

#### Scenario: 扣减到不限量套餐

- **GIVEN** 用户套餐配额为 0（不限量）
- **WHEN** 扣减配额时
- **THEN** 记录使用量但不检查剩余配额
- **AND** 不限量套餐的配额永远不会不足

#### Scenario: 请求失败不扣减配额

- **WHEN** Chat API 请求失败
- **THEN** 不扣减套餐配额
- **AND** 不创建套餐使用记录

### Requirement: 配额缓存机制

系统 SHALL 使用缓存提高配额检测性能。

#### Scenario: 缓存用户套餐配额

- **WHEN** 查询用户套餐配额时
- **THEN** 首先检查缓存是否存在
- **AND** 缓存存在时直接使用缓存数据
- **AND** 缓存不存在时从数据库查询并写入缓存
- **AND** 缓存 TTL 为 5 分钟

#### Scenario: 缓存键格式

- **WHEN** 生成缓存键时
- **THEN** 格式为 `package_quota:{user_id}:{provider_name}`
- **AND** 存储总剩余配额和可用套餐列表

#### Scenario: 配额扣减后更新缓存

- **WHEN** 扣减套餐配额后
- **THEN** 更新缓存中的剩余配额
- **AND** 如果缓存不存在，重新计算并写入缓存

#### Scenario: 缓存失效策略

- **WHEN** 套餐过期或状态变更时
- **THEN** 清除相关缓存
- **AND** 下次查询时重新加载数据

### Requirement: 配额预警机制

系统 SHALL 提供配额即将用尽的预警。

#### Scenario: 配额不足警告

- **GIVEN** 用户套餐剩余配额小于 20%
- **WHEN** 用户发起 Chat API 请求
- **THEN** 在响应头添加警告信息
- **AND** 响应头 `X-Quota-Warning` 包含：
  - `package_id`: 套餐 ID
  - `remaining_tokens`: 剩余配额
  - `message`: "套餐配额即将用完"

#### Scenario: 配额已用完警告

- **GIVEN** 用户套餐配额已用完
- **WHEN** 用户发起 Chat API 请求
- **THEN** 在错误响应中提示配额已用完
- **AND** 提示用户购买新套餐或等待其他套餐生效

### Requirement: 配额统计查询

系统 SHALL 提供配额统计查询 API。

#### Scenario: 查询用户总配额

- **WHEN** GET 请求到 `/api/v1/user/quota`
- **THEN** 返回用户所有激活套餐的总配额
- **AND** 包含总配额数、已使用数、剩余数
- **AND** 按 Provider 分组显示配额

#### Scenario: 查询指定 Provider 配额

- **WHEN** GET 请求到 `/api/v1/user/quota?provider=openai`
- **THEN** 返回指定 Provider 的配额信息
- **AND** 包含可用套餐列表
- **AND** 包含每个套餐的剩余配额
