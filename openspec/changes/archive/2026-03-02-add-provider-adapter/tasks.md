# Tasks: Provider Adapter 实现

## 1. 项目初始化

- [x] 1.1 初始化 Go module (`go mod init github.com/lucheng0127/courier`)
- [x] 1.2 创建项目目录结构（`internal/adapter`, `internal/repository`, `internal/service`, `internal/bootstrap`）
- [ ] 1.3 配置 gofmt 和 lint 工具

## 2. 数据库层实现

- [x] 2.1 设计并创建 `providers` 表迁移文件
  - 字段：id, name, type, base_url, timeout (默认 300), api_key (nullable), extra_config (JSON), created_at, updated_at, enabled
  - 约束：name 唯一索引
- [x] 2.2 实现 Provider 数据模型（`internal/model/provider.go`）
- [x] 2.3 实现 Provider Repository（`internal/repository/provider.go`）
  - `Create()` - 创建 Provider
  - `GetByID()` - 按 ID 查询
  - `GetByName()` - 按 name 查询
  - `List()` - 列出所有 Provider
  - `Update()` - 更新 Provider
  - `Delete()` - 删除 Provider
- [ ] 2.4 编写 Repository 单元测试（使用 mock 数据库）

## 3. Adapter 接口层实现

- [x] 3.1 定义 Provider 接口（`internal/adapter/provider.go`）
  - `Chat()` 方法签名
  - `ChatStream()` 方法签名
  - `Type()` 和 `Name()` 方法
- [x] 3.2 定义请求/响应模型
  - `ChatRequest` - 通用请求结构
  - `ChatResponse` - 非流式响应结构
  - `ChatStreamChunk` - 流式响应块
- [x] 3.3 实现 Adapter Registry（`internal/adapter/registry.go`）
  - `RegisterAdapterType()` 注册函数
  - `NewAdapter()` 工厂函数
  - `GetProvider()` 查询函数
  - 并发安全（使用 sync.RWMutex）
- [x] 3.4 实现 ProviderConfig 结构（`internal/adapter/config.go`）
  - 映射数据库字段到 Go 结构
  - ExtraConfig JSON 解析

## 4. 具体 Adapter 实现（示例）

- [x] 4.1 创建 OpenAI Adapter 框架（`internal/adapter/openai/adapter.go`）
  - 实现 Provider 接口
  - 支持 HTTP 认证头
  - 暂不实现实际调用逻辑（后续变更）
- [x] 4.2 创建 vLLM Adapter 框架（`internal/adapter/vllm/adapter.go`）
  - 实现 Provider 接口
  - 支持无 API Key 调用
- [x] 4.3 在 `init()` 中注册 Adapter 类型

## 5. Provider 管理服务

- [x] 5.1 实现 Provider Service（`internal/service/provider.go`）
  - `CreateProvider()` - 创建并初始化
  - `GetProvider()` - 获取 Provider 实例
  - `ListProviders()` - 列出所有 Provider 及状态
  - `ReloadProvider(name)` - 重新加载指定 Provider
  - `ReloadAllProviders()` - 重新加载所有 Provider
  - `EnableProvider()` - 启用 Provider
  - `DisableProvider()` - 禁用 Provider
- [x] 5.2 实现 Registry Replace 方法（原子替换实例）
- [ ] 5.3 编写 Service 单元测试

## 6. 系统启动初始化

- [x] 6.1 实现 Provider Bootstrap（`internal/bootstrap/provider.go`）
  - `InitProviders()` - 启动时加载所有 Provider
  - 错误处理：单个失败不影响其他
  - 日志记录：记录初始化结果
- [x] 6.2 集成到 main 函数

## 7. API 层实现

- [x] 7.1 实现 Provider 管理 API（`internal/controller/provider.go`）
  - `POST /api/v1/providers` - 创建 Provider
  - `GET /api/v1/providers` - 列出所有 Provider
  - `GET /api/v1/providers/:name` - 获取单个 Provider
  - `PUT /api/v1/providers/:name` - 更新 Provider
  - `DELETE /api/v1/providers/:name` - 删除 Provider
- [x] 7.2 实现 Provider 重载 API（`internal/controller/provider_reload.go`）
  - `POST /api/v1/admin/providers/reload` - 重载所有 Provider
  - `POST /api/v1/admin/providers/:name/reload` - 重载指定 Provider
  - `POST /api/v1/admin/providers/:name/enable` - 启用 Provider
  - `POST /api/v1/admin/providers/:name/disable` - 禁用 Provider
- [x] 7.3 添加 API 认证中间件
- [ ] 7.4 编写 API 集成测试

## 8. 配置和部署

- [x] 8.1 创建 Docker Compose 配置（包含 PostgreSQL）
- [x] 8.2 创建数据库迁移脚本
- [x] 8.3 编写部署文档

## 9. 测试和验证

- [ ] 9.1 端到端测试：创建 OpenAI Provider 并查询
- [ ] 9.2 端到端测试：创建 vLLM Provider 并查询
- [ ] 9.3 测试初始化失败场景（无效 type）
- [ ] 9.4 测试并发访问 Provider Registry
- [ ] 9.5 测试超时控制（默认 300 秒）
- [ ] 9.6 测试重载单个 Provider
- [ ] 9.7 测试重载所有 Provider
- [ ] 9.8 测试重载失败场景（保持旧实例运行）
- [ ] 9.9 测试并发重载场景

## 依赖关系

- 任务 2 和 3 可以并行进行
- 任务 4 依赖任务 3 完成
- 任务 5 依赖任务 2 和 3 完成
- 任务 5.2（Registry Replace）需要与任务 3 同步设计
- 任务 6 依赖任务 5 完成
- 任务 7 依赖任务 5 完成
- 任务 8 和 9 可以在实现完成后并行进行
