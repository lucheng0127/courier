## 1. Service 层实现

- [x] 1.1 在 `internal/service/provider.go` 中实现 `UpdateProvider` 方法
  - 验证 Provider 存在性
  - 更新数据库中的 Provider 配置
  - 如果 Provider 已启用，重载使其生效
  - 如果 Provider 从禁用变为启用，初始化并注册
  - 如果 Provider 从启用变为禁用，注销并清理
  - 返回更新后的 Provider 信息

- [x] 1.2 在 `internal/service/provider.go` 中实现 `DeleteProvider` 方法
  - 验证 Provider 存在性
  - 从 Registry 中注销 Provider（如果正在运行）
  - 从数据库删除 Provider 配置
  - 返回成功或错误

## 2. Controller 层实现

- [x] 2.1 在 `internal/controller/provider.go` 中实现 `UpdateProvider` 方法
  - 从 URL 参数获取 Provider name
  - 绑定和验证请求体
  - 调用 Service 层的 UpdateProvider 方法
  - 返回 200 OK 和更新后的 Provider（成功）
  - 返回 404 Not Found（Provider 不存在）
  - 返回 400 Bad Request（参数验证失败）
  - 返回 500 Internal Server Error（服务器错误）

- [x] 2.2 在 `internal/controller/provider.go` 中实现 `DeleteProvider` 方法
  - 从 URL 参数获取 Provider name
  - 调用 Service 层的 DeleteProvider 方法
  - 返回 204 No Content（成功）
  - 返回 404 Not Found（Provider 不存在）
  - 返回 500 Internal Server Error（服务器错误）

## 3. 测试

- [x] 3.1 为 Service 层 UpdateProvider 添加单元测试
- [x] 3.2 为 Service 层 DeleteProvider 添加单元测试
- [x] 3.3 为 Controller 层 UpdateProvider 添加集成测试
- [x] 3.4 为 Controller 层 DeleteProvider 添加集成测试

## 4. 文档

- [x] 4.1 更新 `docs/provider-config.md` 添加更新和删除 Provider 的说明
- [x] 4.2 更新 `docs/deployment.md` 中的 API 接口说明和 Qwen Provider 示例
