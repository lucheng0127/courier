# public-provider-listing 实施任务

本文档按依赖顺序列出实现 public-provider-listing 变更的任务。

## 阶段 1: 规范更新

### 1.1 创建 spec delta - public-provider-access

**文件**: `openspec/changes/public-provider-listing/specs/public-provider-access/spec.md`

**状态**: ✅ 已完成

**内容**:
- 定义 Provider 查询接口的新权限要求
- 定义模型列表接口的规范
- 定义 `enabled` 过滤参数的使用方式

## 阶段 2: Controller 层实现

### 2.1 新增 PublicProviderInfo 结构体

**文件**: `internal/controller/provider.go`

**状态**: ✅ 已完成

**新增结构体**:
```go
// PublicProviderInfo 普通用户可见的 Provider 信息（不包含敏感信息）
type PublicProviderInfo struct {
    Name           string   `json:"name"`
    Type           string   `json:"type"`
    BaseURL        string   `json:"base_url"`
    Enabled        bool     `json:"enabled"`
    FallbackModels []string `json:"fallback_models"`
}
```

### 2.2 移除 ProviderController.RegisterRoutes

**文件**: `internal/controller/provider.go`

**状态**: ✅ 已完成

**变更**:
- 移除 `RegisterRoutes` 方法（路由在 main.go 中直接注册）

### 2.3 实现 ListProviderModels 方法

**文件**: `internal/controller/provider.go`

**状态**: ✅ 已完成

**新增方法**:
```go
// ListProviderModels 获取 Provider 支持的模型列表
// GET /api/v1/providers/:name/models
func (c *ProviderController) ListProviderModels(ctx *gin.Context)
```

**实现逻辑**:
1. 从路径参数获取 Provider 名称
2. 调用 Service 层获取 Provider 配置
3. 解析 `fallback_models` JSONB 字段为字符串数组
4. 返回响应：`{"name": "...", "type": "...", "models": ["...", "..."]}`
5. 处理 Provider 不存在的情况（返回 404）

### 2.4 修改 ListProviders 方法

**文件**: `internal/controller/provider.go`

**状态**: ✅ 已完成

**变更**:
- 添加 `enabled` 查询参数支持
- 参数验证：如果传入 `enabled`，必须是 "true" 或 "false"
- 将过滤逻辑传递给 Service 层
- 根据用户角色返回不同的响应格式

### 2.5 更新 main.go 路由注册

**文件**: `cmd/server/main.go`

**状态**: ✅ 已完成

**变更**:
- 移除 `providerCtrl.RegisterRoutes(adminOnly)` 调用
- 分别注册管理员路由和普通用户路由：
  ```go
  // 管理操作
  adminOnly.POST("/providers", providerCtrl.CreateProvider)
  adminOnly.PUT("/providers/:name", providerCtrl.UpdateProvider)
  adminOnly.DELETE("/providers/:name", providerCtrl.DeleteProvider)
  adminOnly.GET("/providers/:name", providerCtrl.GetProvider)  // 保持仅管理员

  // 查询操作 - 所有认证用户
  jwtAuth.GET("/providers", providerCtrl.ListProviders)
  jwtAuth.GET("/providers/:name/models", providerCtrl.ListProviderModels)
  ```

## 阶段 3: Service 层实现

### 3.1 修改 ProviderService.ListProviders

**文件**: `internal/service/provider.go`

**状态**: ✅ 已完成

**变更**:
- 添加 `enabledFilter` 参数（可为空、true、false）
- 实现过滤逻辑

### 3.2 添加 ProviderService.GetProviderByName

**文件**: `internal/service/provider.go`

**状态**: ✅ 已完成

**变更**:
- 添加 `GetProviderByName` 方法用于从数据库获取 Provider 配置
- 更新 Controller 接口定义

### 3.3 修改 ProviderRepository 接口

**文件**: `internal/repository/provider.go`

**状态**: ✅ 已完成（无需修改，使用现有 List 方法）

**变更**:
- Service 层实现过滤逻辑，无需修改 Repository

## 阶段 4: 测试

### 4.1 单元测试 - ListProviderModels

**文件**: `internal/controller/provider_public_test.go`

**状态**: ✅ 已完成

**测试用例**:
1. ✅ 正常获取模型列表
2. ✅ Provider 不存在返回 404

### 4.2 单元测试 - ListProviders with enabled filter

**文件**: `internal/controller/provider_public_test.go`

**状态**: ✅ 已完成

**测试用例**:
1. ✅ 普通用户获取简化响应
2. ✅ 管理员获取完整响应
3. ✅ `enabled=true` 只返回已启用的
4. ✅ `enabled=false` 只返回已禁用的
5. ✅ 无效的 `enabled` 参数返回 400

### 4.3 更新现有测试

**文件**: `internal/controller/provider_update_test.go`

**状态**: ✅ 已完成

**变更**:
- 更新 Mock 服务以支持新接口

## 验证清单

完成开发后，验证以下功能：

- [x] 普通用户 JWT Token 可以调用 `GET /api/v1/providers`
- [x] 普通用户响应只包含 `name`、`type`、`base_url`、`enabled`、`fallback_models` 字段
- [x] 普通用户响应不包含 `api_key`、`timeout`、`is_running` 字段
- [x] 管理员 JWT Token 调用 `GET /api/v1/providers` 返回完整信息
- [x] 普通用户 JWT Token 调用 `GET /api/v1/providers/:name` 返回 403 Forbidden（由中间件处理）
- [x] 普通用户 JWT Token 可以调用 `GET /api/v1/providers/:name/models`
- [x] 管理员 JWT Token 可以调用 `GET /api/v1/providers/:name`（完整信息）
- [x] 管理员 JWT Token 可以调用 `GET /api/v1/providers/:name/models`
- [x] 无 JWT Token 返回 401 Unauthorized（由中间件处理）
- [x] `enabled=true` 只返回已启用的 Provider
- [x] `enabled=false` 只返回已禁用的 Provider
- [x] `enabled=invalid` 返回 400 Bad Request
- [x] `/providers/:name/models` 返回正确的模型名称列表
- [x] 普通用户调用 POST/PUT/DELETE 返回 403 Forbidden（由中间件处理）
- [x] 管理员调用 POST/PUT/DELETE 正常工作

## 实施总结

所有任务已完成。主要变更包括：

1. **Controller 层**:
   - 新增 `PublicProviderInfo` 和 `ProviderModelsResponse` 结构体
   - 新增 `ListProviderModels` 方法
   - 修改 `ListProviders` 方法支持 enabled 过滤和角色区分
   - 移除 `RegisterRoutes` 方法

2. **Service 层**:
   - 新增 `GetProviderByName` 方法
   - 修改 `ListProviders` 方法支持 enabled 过滤

3. **路由注册**:
   - 管理操作保持在 `adminOnly` 路由组
   - 查询操作移至 `jwtAuth` 路由组

4. **测试**:
   - 新增 `provider_public_test.go` 包含相关测试用例
   - 更新现有 Mock 服务支持新接口
