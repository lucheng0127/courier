# Design: Provider Adapter 架构设计

## Context

系统需要对接多种 LLM Provider，包括：
- 第三方 SaaS 服务（OpenAI、Anthropic、Azure OpenAI 等）
- 私有部署的本地模型（vLLM、Ollama 等）

不同 Provider 的 API 风格类似但不完全一致，需要统一的抽象层来管理这些差异。

## Goals / Non-Goals

### Goals
- 统一的 Provider 接口抽象
- 支持流式和非流式响应
- 可动态扩展新的 Provider 类型
- 配置驱动的 Adapter 初始化
- 支持 Provider 特定的扩展配置
- 支持运行时重载 Provider（无需重启服务）

### Non-Goals
- 实现具体的 Provider 调用逻辑（在后续变更中）
- 复杂的负载均衡和路由策略（在 Router 变更中）
- Provider 健康检查和自动故障转移（后续优化）

## Decisions

### 1. Provider Interface 设计

```go
type Provider interface {
    // Chat 完成对话调用
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

    // ChatStream 流式对话调用
    ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatStreamChunk, error)

    // Type 返回 Provider 类型标识
    Type() string

    // Name 返回 Provider 实例名称
    Name() string
}
```

**理由**:
- `Chat` 和 `ChatStream` 分离，清晰表达流式/非流式语义
- `Type()` 和 `Name()` 用于标识和日志

### 2. Adapter 注册模式

采用工厂注册模式，支持运行时动态扩展：

```go
type AdapterFactory func(config *ProviderConfig) (Provider, error)

var adapterRegistry = map[string]AdapterFactory{}

func RegisterAdapterType(providerType string, factory AdapterFactory) {
    adapterRegistry[providerType] = factory
}

func NewAdapter(config *ProviderConfig) (Provider, error) {
    factory, ok := adapterRegistry[config.Type]
    if !ok {
        return nil, ErrUnknownProviderType
    }
    return factory(config)
}
```

**理由**:
- 支持在 `init()` 中自动注册新 Adapter
- 新增 Provider 类型无需修改核心代码
- 符合开闭原则

**替代方案**:
1. 硬编码 switch-case - 不利于扩展
2. 反射 + 配置文件 - 过于复杂，类型不安全

### 3. Provider 数据模型

```go
type ProviderConfig struct {
    ID          int64             `json:"id" db:"id"`
    Name        string            `json:"name" db:"name"`               // 必填：唯一标识
    Type        string            `json:"type" db:"type"`               // 必填：openai, anthropic, vllm 等
    BaseURL     string            `json:"base_url" db:"base_url"`       // 必填：API 地址
    Timeout     int               `json:"timeout" db:"timeout"`         // 必填：超时时间（秒）
    APIKey      *string           `json:"api_key" db:"api_key"`         // 可选：SaaS 需要
    ExtraConfig JSON              `json:"extra_config" db:"extra_config"` // 可选：扩展配置
    CreatedAt   time.Time         `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}
```

**理由**:
- `APIKey` 使用指针支持 NULL 值（本地模型不需要）
- `ExtraConfig` 使用 JSON 类型存储扩展参数
- `Timeout` 单位为秒，便于人类理解，默认值 300 秒

### 4. 系统启动初始化流程

```
1. 从数据库加载所有 Provider 配置
2. 对每个配置调用 NewAdapter() 创建实例
3. 将实例注册到全局 Provider Registry
4. 记录初始化结果（成功/失败）
```

**失败处理策略**:
- 单个 Provider 初始化失败不影响其他 Provider
- 记录错误日志但不中断启动
- 提供 API 端点查询 Provider 状态

### 5. 运行时重载机制

系统 SHALL 支持通过 API 触发 Provider 重载，无需重启服务。

```go
// 重载单个 Provider
func ReloadProvider(name string) error {
    // 1. 从数据库加载最新配置
    config, err := repo.GetByName(name)
    if err != nil {
        return err
    }

    // 2. 创建新 Adapter 实例
    newProvider, err := NewAdapter(config)
    if err != nil {
        return err
    }

    // 3. 原子替换 Registry 中的实例
    registry.Replace(name, newProvider)

    return nil
}

// 重载所有 Provider
func ReloadAllProviders() error {
    // 1. 从数据库加载所有配置
    configs, err := repo.List()
    if err != nil {
        return err
    }

    // 2. 逐个重载，失败不影响其他
    for _, config := range configs {
        if err := ReloadProvider(config.Name); err != nil {
            log.Errorf("Failed to reload provider %s: %v", config.Name, err)
        }
    }

    return nil
}
```

**API 端点**:
- `POST /api/v1/admin/providers/reload` - 重载所有 Provider
- `POST /api/v1/admin/providers/:name/reload` - 重载指定 Provider

**理由**:
- 管理员新增/修改 Provider 后可立即生效
- 重载失败不影响其他正在运行的 Provider
- 原子替换保证请求处理的一致性

**失败处理**:
- 重载失败时保持旧实例继续运行
- 记录详细错误日志
- 返回错误信息给调用方

### 6. 目录结构

```
internal/
├── adapter/
│   ├── provider.go          # Provider 接口定义
│   ├── registry.go          # Adapter 注册逻辑
│   ├── openai/
│   │   └── adapter.go       # OpenAI Adapter（后续实现）
│   └── vllm/
│       └── adapter.go       # vLLM Adapter（后续实现）
├── repository/
│   └── provider.go          # Provider 数据访问
├── service/
│   └── provider.go          # Provider 管理服务
└── bootstrap/
    └── provider.go          # Provider 初始化
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Adapter 初始化失败导致部分 Provider 不可用 | 记录详细日志，提供状态查询 API |
| ExtraConfig JSON 类型缺少类型检查 | 在 Adapter 层进行校验和转换 |
| 流式响应的 channel 泄漏 | 使用 context 取消机制，确保资源清理 |
| 并发访问 Provider Registry | 使用 sync.RWMutex 保护注册表 |

## Migration Plan

由于是新功能，无迁移需求。

部署步骤：
1. 执行数据库迁移创建 `providers` 表
2. 部署新代码
3. 通过 API 或数据库初始化 Provider 配置
4. 重启服务完成 Adapter 初始化

回滚策略：
- 删除 `providers` 表
- 回滚代码版本

## Open Questions

1. **ExtraConfig 是否需要 JSON Schema 校验？**
   - 建议：MVP 阶段不做强制校验，由各 Adapter 内部处理

2. **Provider 配置是否需要版本控制？**
   - 建议：MVP 阶段不支持，使用 `updated_at` 跟踪变更时间

3. **是否需要自动监听数据库变更实现自动重载？**
   - 建议：MVP 阶段使用 API 手动触发，后续可考虑使用 PostgreSQL NOTIFY/LISTEN 实现自动重载
