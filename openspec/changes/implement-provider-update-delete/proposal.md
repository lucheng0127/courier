# Change: 实现 Provider 更新和删除功能

## Why

当前 `internal/controller/provider.go` 中的 `UpdateProvider` 和 `DeleteProvider` 方法仅返回 "not implemented yet" 错误。虽然 Repository 层已提供 `Update` 和 `Delete` 方法，但缺少 Service 层的完整实现和 Controller 层的调用逻辑，导致无法通过 API 更新或删除 Provider 配置。

## What Changes

- 在 Service 层添加 `UpdateProvider` 方法，支持更新 Provider 配置
- 在 Service 层添加 `DeleteProvider` 方法，支持删除 Provider
- 更新后自动重载 Provider（如果已启用）
- 删除前先注销并清理运行中的实例
- 在 Controller 层实现 `UpdateProvider` 和 `DeleteProvider` 的完整逻辑
- 添加参数验证和错误处理

## Impact

- **受影响的规范**: `provider-adapter`
- **受影响的代码**:
  - `internal/service/provider.go` - 添加 UpdateProvider 和 DeleteProvider 方法
  - `internal/controller/provider.go` - 实现 UpdateProvider 和 DeleteProvider 接口逻辑
