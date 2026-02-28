# Change: 添加 Provider Adapter 支持多模型接入

## Why

当前系统需要支持对接多种 AI 模型供应商（如 OpenAI、Anthropic、本地 vLLM 等），以实现统一的 LLM Gateway 能力。需要一个可扩展的适配器模式来抽象不同厂商的差异，支持流式输出，并允许动态配置和扩展新的 Provider 类型。

## What Changes

- 新增 `provider-adapter` 能力：定义 Provider 接口和 Adapter 模式
- 实现 Provider 配置的数据库存储（表：`providers`）
- 实现系统启动时的 Provider 配置加载和 Adapter 初始化机制
- 支持多种 Provider 类型的自动注册和实例化
- 支持第三方 SaaS 模型（需要 API Key）和私有本地模型（API Key 可选）
- 支持流式响应（SSE）和非流式响应
- 提供扩展配置（extra_config）支持 Provider 特定参数
- **支持运行时重载 Provider**：管理员新增/修改 Provider 后可通过 API 立即生效，无需重启服务
- **默认超时时间 300 秒**

## Impact

- **Affected specs**: 新增 `provider-adapter` 规范
- **Affected code**:
  - 新增 `internal/adapter/` - Adapter 层实现
  - 新增 `internal/repository/provider.go` - Provider 数据访问
  - 新增 `internal/service/provider.go` - Provider 管理服务
  - 修改 `internal/bootstrap/` - 添加 Provider 初始化逻辑
- **Database**: 新增 `providers` 表
- **Dependencies**: 无新增外部依赖（使用 Go 标准库和现有 HTTP 客户端）
