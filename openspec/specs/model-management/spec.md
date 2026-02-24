# model-management Specification

## Purpose
TBD - created by archiving change add-model-chat. Update Purpose after archive.
## Requirements
### Requirement: 模型列表查询

系统 MUST 支持查询当前可用的上游模型列表。

#### Scenario: 成功查询模型列表

- **WHEN** 客户端发送 GET 请求到 `/api/v1/models`
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 返回模型列表数组，每个模型包含 `name`、`provider` 字段
- **AND** 返回的模型列表来自配置文件中配置的有效模型

#### Scenario: 无可用模型

- **WHEN** 配置文件中没有配置任何模型
- **THEN** 系统返回 HTTP 200 状态码
- **AND** 返回空数组

### Requirement: 模型配置验证

系统 MUST 在启动时验证模型配置的有效性。

#### Scenario: 配置有效

- **WHEN** 配置文件中的模型配置包含有效的 `name`、`provider`、`base_url`、`api_key`
- **THEN** 系统正常启动
- **AND** 该模型可用于对话

#### Scenario: 缺少必需字段

- **WHEN** 配置文件中的模型配置缺少 `name`、`base_url` 或 `api_key` 字段
- **THEN** 系统启动失败
- **AND** 返回错误信息指出哪个模型的哪个字段缺失

#### Scenario: API Key 为空

- **WHEN** 配置文件中的模型配置 `api_key` 为空字符串或环境变量未设置
- **THEN** 系统启动失败
- **AND** 返回错误信息指出该模型的 API Key 未配置

### Requirement: 模型名称唯一性

系统 MUST 确保配置的模型名称唯一。

#### Scenario: 模型名称重复

- **WHEN** 配置文件中存在多个相同 `name` 的模型配置
- **THEN** 系统启动失败
- **AND** 返回错误信息指出模型名称重复

### Requirement: 环境变量支持

系统 MUST 支持在模型配置中使用环境变量。

#### Scenario: 使用环境变量

- **WHEN** 配置文件中的 `api_key` 格式为 `${ENV_VAR_NAME}`
- **THEN** 系统从环境变量中读取该值
- **AND** 环境变量存在时使用该值
- **AND** 环境变量不存在时视为配置错误

#### Scenario: 环境变量未设置

- **WHEN** 配置文件中的 `api_key` 引用的环境变量不存在
- **THEN** 系统启动失败
- **AND** 返回错误信息指出该环境变量未设置

### Requirement: 模型存在性验证

系统 MUST 在处理对话请求前验证模型是否存在。

#### Scenario: 模型存在

- **WHEN** 客户端请求的模型名称在配置中存在
- **THEN** 系统继续处理请求

#### Scenario: 模型不存在

- **WHEN** 客户端请求的模型名称在配置中不存在
- **THEN** 系统返回 HTTP 404 状态码
- **AND** 返回错误信息提示模型不存在

