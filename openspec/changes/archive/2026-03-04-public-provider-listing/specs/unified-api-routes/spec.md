# unified-api-routes Specification Delta

## MODIFIED Requirements

### Requirement: Provider 管理 API 路径

Provider 查询接口权限 SHALL 已调整，所有认证用户可访问。

#### Scenario: Provider 查询路径权限

- **WHEN** 访问 Provider 查询接口
- **THEN** 接口路径和权限如下：
  - `GET /api/v1/providers` - 列出所有 Provider（所有认证用户）
  - `GET /api/v1/providers/:name/models` - 获取模型列表（所有认证用户）
  - `GET /api/v1/providers/:name` - 获取单个 Provider（仅管理员）
  - `POST /api/v1/providers` - 创建 Provider（仅管理员）
  - `PUT /api/v1/providers/:name` - 更新 Provider（仅管理员）
  - `DELETE /api/v1/providers/:name` - 删除 Provider（仅管理员）

## ADDED Requirements

### Requirement: Provider 模型列表 API 路径

系统 SHALL 提供专门的接口获取 Provider 支持的模型列表。

#### Scenario: 模型列表路径

- **WHEN** 访问 Provider 模型列表接口
- **THEN** 接口路径为 `GET /api/v1/providers/:name/models`
- **AND** 所有认证用户可访问
