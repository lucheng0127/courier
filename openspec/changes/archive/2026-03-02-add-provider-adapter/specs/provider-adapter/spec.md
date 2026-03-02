## ADDED Requirements

### Requirement: Provider Interface

系统 SHALL 定义统一的 Provider 接口，用于抽象不同 LLM 供应商的调用方式。

#### Scenario: 接口定义

- **WHEN** 定义 Provider 接口时
- **THEN** 接口必须包含 `Chat()` 方法用于非流式调用
- **AND** 接口必须包含 `ChatStream()` 方法用于流式调用
- **AND** 接口必须包含 `Type()` 方法返回 Provider 类型
- **AND** 接口必须包含 `Name()` 方法返回 Provider 实例名称

#### Scenario: 流式响应

- **WHEN** 调用 `ChatStream()` 方法时
- **THEN** 返回只读 channel 用于接收流式数据块
- **AND** 支持 context 取消机制

### Requirement: Adapter Registry

系统 SHALL 提供 Adapter 注册机制，支持动态扩展新的 Provider 类型。

#### Scenario: Adapter 类型注册

- **WHEN** 新的 Provider 类型实现后
- **THEN** 可通过 `RegisterAdapterType()` 注册工厂函数
- **AND** 注册时指定 Provider 类型标识（如 "openai", "vllm"）

#### Scenario: Adapter 实例化

- **WHEN** 给定 Provider 配置时
- **THEN** 系统根据配置中的 Type 查找对应的工厂函数
- **AND** 调用工厂函数创建 Provider 实例
- **AND** 如果 Type 未注册则返回错误

#### Scenario: 自动注册

- **WHEN** 程序启动时
- **THEN** 各 Adapter 通过 `init()` 函数自动注册到 Registry
- **AND** 无需手动调用注册代码

### Requirement: Provider 配置存储

系统 SHALL 将 Provider 配置存储在数据库中，支持持久化和查询。

#### Scenario: 配置必填字段

- **WHEN** 创建 Provider 配置时
- **THEN** 必须提供 `name` 作为唯一标识
- **AND** 必须提供 `type` 指定 Provider 类型
- **AND** 必须提供 `base_url` 指定 API 地址
- **AND** 必须提供 `timeout` 指定超时时间（秒）
- **AND** 必须提供 `enabled` 指定启用状态（布尔值，默认 true）

#### Scenario: 配置可选字段

- **WHEN** 创建 SaaS 类型 Provider 时
- **THEN** 应提供 `api_key` 用于身份认证
- **WHEN** 创建本地模型 Provider 时
- **THEN** `api_key` 可为空
- **WHEN** Provider 需要额外配置时
- **THEN** 可提供 `extra_config` 存储扩展参数（JSON 格式）

#### Scenario: 配置唯一性

- **WHEN** 创建 Provider 配置时
- **THEN** `name` 字段必须在全局唯一
- **AND** 如果重复则返回唯一性约束错误

### Requirement: Provider 初始化

系统启动时 SHALL 加载所有 Provider 配置并完成 Adapter 初始化。

#### Scenario: 启动加载

- **WHEN** 系统启动时
- **THEN** 从数据库加载所有 Provider 配置
- **AND** 对于 `enabled` 为 true 的 Provider，根据配置创建对应的 Provider 实例
- **AND** 对于 `enabled` 为 false 的 Provider，跳过初始化
- **AND** 将已初始化的实例注册到全局 Provider Registry

#### Scenario: 初始化失败处理

- **WHEN** 单个 Provider 初始化失败时
- **THEN** 记录错误日志
- **AND** 不中断系统启动
- **AND** 继续初始化其他 Provider
- **AND** 该 Provider 不可用，修复配置后可通过重载 API 使其生效

#### Scenario: 未知 Provider 类型

- **WHEN** 配置中的 `type` 未有对应 Adapter 注册时
- **THEN** 记录错误日志
- **AND** 跳过该 Provider 配置

### Requirement: Provider 状态查询

系统 SHALL 提供 API 端点查询当前可用的 Provider 列表和状态。

#### Scenario: 查询 Provider 列表

- **WHEN** 调用 Provider 列表 API 时
- **THEN** 返回所有已配置的 Provider
- **AND** 包含每个 Provider 的名称、类型、启用状态（enabled/disabled）、运行状态（可用/不可用）

#### Scenario: 查询单个 Provider

- **WHEN** 按 name 查询 Provider 时
- **THEN** 返回该 Provider 的详细配置
- **AND** 如果不存在则返回 404

### Requirement: Provider 运行时重载

系统 SHALL 支持运行时重载 Provider 配置，无需重启服务。

#### Scenario: 重载单个 Provider

- **WHEN** 管理员调用重载 API 并指定 Provider name 时
- **THEN** 从数据库重新加载该 Provider 的最新配置
- **AND** 销毁旧的 Adapter 实例
- **AND** 使用新配置创建新的 Adapter 实例
- **AND** 更新 Registry 中的实例引用

#### Scenario: 重载所有 Provider

- **WHEN** 管理员调用重载所有 Provider API 时
- **THEN** 从数据库重新加载所有 Provider 配置
- **AND** 销毁所有旧的 Adapter 实例
- **AND** 使用新配置创建所有新的 Adapter 实例
- **AND** 更新 Registry 中的所有实例引用

#### Scenario: 重载失败处理

- **WHEN** Provider 重载失败时
- **THEN** 记录错误日志
- **AND** 保持旧实例继续运行
- **AND** 返回错误信息给调用方

#### Scenario: 新增 Provider 后重载

- **WHEN** 管理员通过 API 新增 Provider 配置后
- **THEN** 可立即调用重载 API 使新 Provider 生效
- **AND** 无需重启服务

### Requirement: Provider 启用/禁用

系统 SHALL 支持运行时启用或禁用 Provider。

#### Scenario: 禁用 Provider

- **WHEN** 管理员调用禁用 API 并指定 Provider name 时
- **THEN** 更新数据库中该 Provider 的 `enabled` 字段为 false
- **AND** 从 Registry 中移除该 Provider 实例
- **AND** 销毁 Adapter 实例释放资源
- **AND** 后续请求无法使用该 Provider

#### Scenario: 启用 Provider

- **WHEN** 管理员调用启用 API 并指定 Provider name 时
- **THEN** 更新数据库中该 Provider 的 `enabled` 字段为 true
- **AND** 从数据库重新加载该 Provider 的配置
- **AND** 创建新的 Adapter 实例并注册到 Registry
- **AND** 后续请求可正常使用该 Provider

#### Scenario: 启用/禁用失败处理

- **WHEN** 启用或禁用操作失败时
- **THEN** 记录错误日志
- **AND** 返回错误信息给调用方
- **AND** 保持原有状态不变

### Requirement: 第三方 SaaS 模型支持

系统 SHALL 支持需要 API Key 的第三方 SaaS 模型服务。

#### Scenario: API Key 认证

- **WHEN** 调用 SaaS Provider 时
- **THEN** 使用配置中的 `api_key` 进行身份认证
- **AND** API Key 通过 HTTP Header 传递（如 `Authorization: Bearer <key>`）

#### Scenario: API Key 安全存储

- **WHEN** 存储 API Key 时
- **THEN** 应考虑加密存储（MVP 阶段可使用明文，后续升级）

### Requirement: 私有本地模型支持

系统 SHALL 支持私有部署的本地模型，API Key 可选。

#### Scenario: 无 API Key 调用

- **WHEN** 配置中未提供 `api_key` 时
- **THEN** Adapter 不发送认证头
- **AND** 正常调用模型 API

#### Scenario: 内网地址支持

- **WHEN** `base_url` 为内网地址时
- **THEN** 系统应正常处理
- **AND** 支持 HTTP 和 HTTPS 协议

### Requirement: 扩展配置支持

系统 SHALL 支持 Provider 特定的扩展配置。

#### Scenario: 额外参数传递

- **WHEN** Provider 需要额外配置时
- **THEN** 可通过 `extra_config` 字段传递 JSON 格式的配置
- **AND** Adapter 负责解析和验证这些配置

#### Scenario: 常见扩展配置

- **WHEN** 配置 OpenAI 兼容接口时
- **THEN** `extra_config` 可支持 `model` 字段指定模型名称
- **AND** 可支持 `organization` 字段指定组织 ID
- **WHEN** 配置 vLLM 时
- **THEN** `extra_config` 可支持 `max_tokens`、`temperature` 等默认参数

### Requirement: 流式输出支持

系统 SHALL 支持 SSE（Server-Sent Events）流式响应。

#### Scenario: 流式响应格式

- **WHEN** Provider 返回流式响应时
- **THEN** 通过 channel 逐块返回数据
- **AND** 每个数据块包含 delta 内容和完成状态

#### Scenario: 流式上下文取消

- **WHEN** 客户端取消请求或超时时
- **THEN** 通过 context.Context 传递取消信号
- **AND** Adapter 应中断流式响应并清理资源

### Requirement: 超时控制

系统 SHALL 支持为每个 Provider 配置独立的超时时间。

#### Scenario: 请求超时

- **WHEN** Provider 调用超过配置的 `timeout` 时间时
- **THEN** 取消请求并返回超时错误
- **AND** 超时时间从 Provider 配置读取

#### Scenario: 默认超时

- **WHEN** 未指定 timeout 时
- **THEN** 使用默认值 300 秒
