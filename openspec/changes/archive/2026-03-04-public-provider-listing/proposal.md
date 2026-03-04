# public-provider-listing 变更提案

## Why

当前 `GET /api/v1/providers` 接口只允许管理员访问，普通用户无法查看系统中有哪些 Provider 可用，也不清楚各个 Provider 支持哪些模型。这限制了用户对系统能力的了解。

通过将 Provider 查询接口开放给普通用户，用户可以：
1. 查看所有可用的 Provider（支持按启用状态过滤）
2. 通过单独的接口查看每个 Provider 支持的模型列表
3. 更好地选择合适的模型进行 API 调用

## What Changes

### 涉及的能力 (Capabilities)

1. **public-provider-access** - 公开 Provider 访问权限
   - 修改 `GET /api/v1/providers` 接口权限
   - 根据用户角色返回不同的响应字段（管理员看到全部，普通用户只看到非敏感信息）
   - 添加启用状态过滤功能
   - 新增 `GET /api/v1/providers/:name/models` 接口获取模型列表
   - `GET /api/v1/providers/:name` 保持仅管理员访问

### 不涉及的内容

- 不修改 Provider 的创建、更新、删除权限（仍为管理员专属）
- 不修改 `GET /api/v1/providers/:name` 的权限（保持仅管理员访问）
- 不修改 Provider 的启用、禁用、重载操作（仍为管理员专属）
- 不新增除模型列表外的额外接口

## 设计方案

### 1. 新增模型列表接口

**接口**: `GET /api/v1/providers/:name/models`

**权限**: 所有认证用户（JWT）

**路径参数**:
- `name`: Provider 名称

**响应**:
```json
{
  "name": "openai-main",
  "type": "openai",
  "models": [
    "gpt-4o",
    "gpt-4o-mini",
    "gpt-3.5-turbo"
  ]
}
```

**模型数据来源**: 从 Provider 的 `fallback_models` 字段获取

### 2. Provider 列表接口增强

**接口**: `GET /api/v1/providers`

**权限变更**: 从"仅管理员"改为"所有认证用户"

**新增查询参数**:
- `enabled`: 可选，过滤条件
  - `true`: 只返回已启用的 Provider
  - `false`: 只返回已禁用的 Provider
  - 不传参数: 返回所有 Provider

**响应格式（根据用户角色不同）**:

管理员看到的完整响应：
```json
{
  "providers": [
    {
      "provider": {
        "name": "openai-main",
        "type": "openai",
        "base_url": "https://api.openai.com/v1",
        "timeout": 60,
        "enabled": true,
        "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
      },
      "is_running": true
    }
  ]
}
```

普通用户看到的简化响应（不包含敏感信息）：
```json
{
  "providers": [
    {
      "name": "openai-main",
      "type": "openai",
      "base_url": "https://api.openai.com/v1",
      "enabled": true,
      "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
    }
  ]
}
```

**普通用户响应字段说明**：
- `name`: Provider 名称
- `type`: Provider 类型
- `base_url`: API 地址
- `enabled`: 启用状态
- `fallback_models`: 支持的模型列表
- 不包含：`api_key`、`timeout`、`is_running`、`extra_config` 等敏感或运维信息

### 3. 路由注册调整

将 Provider 路由按权限分组：

**修改前**:
```go
adminOnly := jwtAuth.Group("")
adminOnly.Use(middleware.RequireAdmin())
providerCtrl.RegisterRoutes(adminOnly)  // 所有接口都在 adminOnly 下
```

**修改后**:
```go
adminOnly := jwtAuth.Group("")
adminOnly.Use(middleware.RequireAdmin())

// 管理操作 - 需要管理员权限
adminOnly.POST("/providers", providerCtrl.CreateProvider)
adminOnly.PUT("/providers/:name", providerCtrl.UpdateProvider)
adminOnly.DELETE("/providers/:name", providerCtrl.DeleteProvider)
adminOnly.GET("/providers/:name", providerCtrl.GetProvider)  // 保持仅管理员访问

// 查询操作 - 所有认证用户可访问
jwtAuth.GET("/providers", providerCtrl.ListProviders)
jwtAuth.GET("/providers/:name/models", providerCtrl.ListProviderModels)
```

### 4. Controller 实现

**新增方法**: `ListProviderModels`
- 从数据库获取 Provider 配置
- 解析 `fallback_models` JSONB 字段
- 返回 Provider 基本信息 + 模型名称列表

**修改方法**: `ListProviders`
- 添加 `enabled` 查询参数处理
- 支持按启用状态过滤
- **根据用户角色返回不同的响应格式**：
  - 管理员：返回完整的 `ProviderInfo` 结构
  - 普通用户：返回简化的 `PublicProviderInfo` 结构（只包含 name、type、base_url、fallback_models）

**修改方法**: `RegisterRoutes`
- 拆分为管理员路由和普通用户路由的注册

**新增结构体**: `PublicProviderInfo`
```go
type PublicProviderInfo struct {
    Name           string   `json:"name"`
    Type           string   `json:"type"`
    BaseURL        string   `json:"base_url"`
    Enabled        bool     `json:"enabled"`
    FallbackModels []string `json:"fallback_models"`
}
```

## 对现有规范的影响

### 修改的规范

1. **role-based-access** - 添加新的需求说明 Provider 查询接口的权限变更
2. **unified-api-routes** - 更新 Provider API 路由权限和新增模型列表接口

## 向后兼容性

- **向后兼容**: 现有的管理员调用不受影响
- **新功能**: 普通用户现在可以调用查询接口
- **新增参数**: `enabled` 参数为可选，不传时行为与之前一致
- **API 格式**: 响应格式保持不变

## 安全考虑

1. **只读访问**：普通用户只能查询，不能修改
2. **需要认证**：仍需要 JWT Token，不允许匿名访问
3. **敏感信息过滤**：
   - 普通用户响应中不包含 `api_key`、`timeout`、`is_running`、`extra_config` 等敏感或运维信息
   - 普通用户可以看到 `name`、`type`、`base_url`、`enabled`、`fallback_models`
   - 管理员仍然可以看到完整信息
4. **接口权限区分**：
   - `GET /api/v1/providers` - 所有认证用户
   - `GET /api/v1/providers/:name` - 仅管理员（保持不变）
   - `GET /api/v1/providers/:name/models` - 所有认证用户
5. **模型列表公开**：模型列表来自配置，不涉及敏感信息

## 测试计划

1. 普通用户可以成功调用 `GET /api/v1/providers`
2. 普通用户响应只包含 `name`、`type`、`base_url`、`enabled`、`fallback_models` 字段
3. 普通用户响应不包含 `api_key`、`timeout`、`is_running` 字段
4. 管理员用户调用 `GET /api/v1/providers` 返回完整信息
5. 普通用户调用 `GET /api/v1/providers/:name` 返回 403 Forbidden
6. 普通用户可以成功调用 `GET /api/v1/providers/:name/models`
7. 管理员用户可以成功调用 `GET /api/v1/providers/:name/models`
8. 未认证用户调用返回 401
9. 验证 `enabled=true` 参数只返回已启用的 Provider
10. 验证 `enabled=false` 参数只返回已禁用的 Provider
11. 验证模型列表接口正确返回 `fallback_models` 中的模型名称
12. 验证模型列表接口响应只包含 `name`、`type`、`models` 字段
13. 管理操作（POST/PUT/DELETE）仍然需要管理员权限
