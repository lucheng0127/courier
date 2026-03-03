# Provider 配置指南

本文档介绍如何配置各种 LLM Provider，包括 OpenAI 兼容服务、本地 vLLM 等。

## Provider 类型

系统支持以下 Provider 类型：

| 类型 | 说明 | 适用场景 |
|------|------|----------|
| `openai` | OpenAI API 或兼容服务 | OpenAI、通义千问等 |
| `vllm` | vLLM 本地部署服务 | 私有化部署 |

## 配置参数

### 通用参数

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `name` | string | ✓ | Provider 实例名称（全局唯一） |
| `type` | string | ✓ | Provider 类型（`openai` 或 `vllm`） |
| `base_url` | string | ✓ | API 地址 |
| `timeout` | int | ✓ | 超时时间（秒），默认 300 |
| `enabled` | boolean | ✓ | 是否启用，默认 true |
| `api_key` | string | - | API Key（vLLM 可选） |
| `extra_config` | object | - | 扩展配置 |
| `fallback_models` | array | - | Fallback 模型列表 |

### 扩展配置 (extra_config)

支持在 `extra_config` 中设置默认模型参数：

| 参数 | 类型 | 说明 |
|------|------|------|
| `temperature` | float64 | 温度参数（0-2） |
| `max_tokens` | int | 最大生成 token 数 |
| `top_p` | float64 | 核采样参数（0-1） |

> **注意**：请求级参数优先于 `extra_config` 中的默认参数。

## 配置示例

### OpenAI 官方服务

```bash
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "openai-main",
    "type": "openai",
    "base_url": "https://api.openai.com/v1",
    "timeout": 60,
    "api_key": "sk-your-openai-api-key",
    "enabled": true,
    "extra_config": {
      "temperature": 0.7,
      "max_tokens": 2000
    },
    "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
  }'
```

### 通义千问 (Qwen)

通义千问提供 OpenAI 兼容的 API 接口。

**Base URL**: `https://dashscope.aliyuncs.com/compatible-mode/v1`

```bash
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "qwen-main",
    "type": "openai",
    "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1",
    "timeout": 60,
    "api_key": "sk-your-qwen-api-key",
    "enabled": true,
    "extra_config": {
      "temperature": 0.8,
      "max_tokens": 1500
    },
    "fallback_models": ["qwen-max", "qwen-plus", "qwen-turbo"]
  }'
```

**使用示例**：

```bash
# 调用通义千问
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_API_KEY" \
  -d '{
    "model": "qwen-main/qwen-turbo",
    "messages": [{"role": "user", "content": "你好，请介绍一下你自己"}]
  }'
```

**常用模型**：

| 模型 | 说明 |
|------|------|
| `qwen-max` | 最强模型，适合复杂任务 |
| `qwen-plus` | 性能均衡，适合大多数场景 |
| `qwen-turbo` | 快速响应，适合简单任务 |

### 本地 vLLM 部署

vLLM 是一个高性能的 LLM 推理引擎，支持本地部署各种开源模型。

**不需要 API Key**（内网部署）：

```bash
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "local-vllm",
    "type": "vllm",
    "base_url": "http://localhost:8000/v1",
    "timeout": 120,
    "enabled": true,
    "extra_config": {
      "temperature": 0.7,
      "max_tokens": 2048
    }
  }'
```

**带 API Key 的 vLLM**（如果需要认证）：

```bash
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "vllm-protected",
    "type": "vllm",
    "base_url": "http://vllm.internal:8000/v1",
    "timeout": 120,
    "api_key": "your-internal-api-key",
    "enabled": true
  }'
```

**使用示例**：

```bash
# 调用本地 vLLM
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_API_KEY" \
  -d '{
    "model": "local-vllm/Qwen/Qwen2-7B-Instruct",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### 其他 OpenAI 兼容服务

任何提供 OpenAI 兼容 API 的服务都可以配置为 `openai` 类型：

```bash
curl -X POST http://localhost:8080/api/v1/providers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "name": "custom-provider",
    "type": "openai",
    "base_url": "https://your-provider.com/v1",
    "timeout": 60,
    "api_key": "your-api-key",
    "enabled": true,
    "fallback_models": ["model-1", "model-2"]
  }'
```

## Fallback 配置

### 为什么需要 Fallback

- **提高可用性**：当主模型不可用时自动切换
- **成本优化**：优先使用低成本模型，需要时升级到高成本模型
- **负载均衡**：分散请求到多个模型

### 配置 Fallback

在 `fallback_models` 数组中按优先级排列模型：

```json
{
  "name": "smart-fallback",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "api_key": "sk-xxx",
  "fallback_models": [
    "gpt-4o-mini",    // 主模型（快速、便宜）
    "gpt-4o",         // 备用模型（更强但更贵）
    "gpt-3.5-turbo"   // 最后备用
  ]
}
```

### Fallback 触发条件

以下情况会触发 Fallback：
- 超时
- 网络错误（连接失败、DNS 解析失败）
- 5xx 服务器错误
- 连接拒绝

以下情况**不会**触发 Fallback：
- 4xx 客户端错误（除 429 外）
- 认证失败
- 模型不存在

## 管理 Provider

### 查询 Provider 列表

```bash
curl http://localhost:8080/api/v1/providers \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### 更新 Provider

```bash
# 更新 Provider 配置
curl -X PUT http://localhost:8080/api/v1/providers/openai-main \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "timeout": 120,
    "enabled": false
  }'
```

**说明**：
- 只需提供要更新的字段，未提供的字段保持原值
- 如果 Provider 已启用，更新后会自动重载
- 如果从禁用变为启用，会自动初始化并注册
- 如果从启用变为禁用，会自动注销并清理

### 删除 Provider

```bash
# 删除 Provider
curl -X DELETE http://localhost:8080/api/v1/providers/old-provider \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**说明**：
- 删除前会自动注销运行中的 Provider
- 删除操作不可逆，请谨慎操作
- 如果 Provider 正在被使用，建议先禁用再删除

### 重载 Provider

修改 Provider 配置后，也可以手动重载使其生效：

```bash
# 重载所有 Provider
curl -X POST http://localhost:8080/api/v1/providers/reload \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 重载指定 Provider
curl -X POST http://localhost:8080/api/v1/admin/providers/openai-main/reload \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

### 启用/禁用 Provider

```bash
# 禁用
curl -X POST http://localhost:8080/api/v1/admin/providers/openai-main/disable \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 启用
curl -X POST http://localhost:8080/api/v1/admin/providers/openai-main/enable \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**注意**：也可以通过更新 API 的 `enabled` 字段来启用/禁用 Provider，更新操作会自动处理状态转换。

## 最佳实践

1. **配置 Fallback**：为每个 Provider 配置至少 2 个 Fallback 模型
2. **合理设置超时**：根据模型响应时间调整 `timeout` 参数
3. **使用默认参数**：在 `extra_config` 中设置合理的默认值
4. **监控状态**：定期检查 Provider 状态和调用日志
5. **安全存储 API Key**：生产环境应使用加密存储（MVP 阶段使用明文）

## 故障排查

### Provider 不可用

```bash
# 检查 Provider 状态
curl http://localhost:8080/api/v1/providers \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq
```

检查响应中的 `enabled` 和运行状态。

### 认证失败

- 检查 `api_key` 是否正确
- 确认 API Key 未过期
- 验证 `base_url` 是否正确

### 超时错误

- 增加 `timeout` 参数
- 检查网络连接
- 查看 Provider 服务状态

### 模型不存在

- 确认模型名称拼写正确
- 检查 Provider 支持的模型列表
- 验证 `base_url` 对应的服务
