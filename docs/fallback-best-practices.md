# Fallback 重试最佳实践

## 概述

Courier LLM Gateway 提供了同 Provider 内的模型 Fallback 能力，当主模型失败时自动切换到备用模型，提高服务可用性。

## 配置原则

### 1. 模型选择策略

```
主模型 → 备用模型1 → 备用模型2 → ...
(高能力)  (中等能力)    (低能力/快速)
```

**推荐顺序**：
1. **主模型**：最强能力模型（如 GPT-4o）
2. **备用模型1**：平衡模型（如 GPT-4o-mini）
3. **备用模型2**：快速/低成本模型（如 GPT-3.5-turbo）

### 2. 同质化模型

推荐使用同一系列的模型进行 Fallback：

```json
{
  "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
}
```

**避免**跨系列 Fallback（如 Anthropic → OpenAI），因为输出格式和风格差异较大。

### 3. 超时配置

根据模型响应速度设置合理的超时：

| 模型类型 | 推荐超时 |
|----------|----------|
| GPT-4o | 60 秒 |
| GPT-4o-mini | 30 秒 |
| GPT-3.5-turbo | 20 秒 |
| Claude Opus | 90 秒 |
| 本地 vLLM | 30 秒 |

```json
{
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "timeout": 60,
  "fallback_models": ["gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"]
}
```

## 配置示例

### OpenAI Provider

```json
{
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "timeout": 60,
  "api_key": "sk-xxx",
  "enabled": true,
  "fallback_models": [
    "gpt-4o",
    "gpt-4o-mini",
    "gpt-3.5-turbo"
  ]
}
```

### Anthropic Provider

```json
{
  "name": "anthropic-main",
  "type": "anthropic",
  "base_url": "https://api.anthropic.com/v1",
  "timeout": 90,
  "api_key": "sk-ant-xxx",
  "enabled": true,
  "fallback_models": [
    "claude-3-opus-20240229",
    "claude-3-sonnet-20240229",
    "claude-3-haiku-20240307"
  ]
}
```

### vLLM Provider（本地部署）

```json
{
  "name": "vllm-local",
  "type": "vllm",
  "base_url": "http://localhost:8000/v1",
  "timeout": 30,
  "enabled": true,
  "fallback_models": [
    "llama-2-70b",
    "llama-2-13b",
    "llama-2-7b"
  ]
}
```

## 监控和告警

### 关键指标

1. **Fallback 频率**
   - 正常：< 5%
   - 警告：5-20%
   - 严重：> 20%

2. **Fallback 分布**
   - 监控每个模型的成功率
   - 识别经常失败的模型

3. **延迟影响**
   - Fallback 增加的额外延迟
   - 超时对用户体验的影响

### 日志分析

通过结构化日志分析 Fallback 行为：

```bash
# 查看有 Fallback 的请求
jq 'select(.fallback_count > 0)' logs.jsonl

# 统计 Fallback 频率
jq -r '.fallback_count' logs.jsonl | awk '{count[$1]++} END {for (c in count) print c, count[c]}'
```

## 故障排查

### 常见问题

#### 1. Fallback 频繁触发

**可能原因**：
- 主模型超时配置过短
- Provider 服务不稳定
- 网络连接问题

**解决方案**：
- 增加超时时间
- 检查 Provider 状态
- 检查网络连接

#### 2. 所有模型都失败

**可能原因**：
- Provider API Key 失效
- Provider 服务完全不可用
- 网络完全中断

**排查步骤**：
1. 检查日志中的 `error_type`
2. 验证 API Key 有效性
3. 测试网络连通性

#### 3. 成本增加

**可能原因**：
- 主模型频繁失败，频繁使用备用模型

**解决方案**：
- 监控 Fallback 频率
- 设置告警阈值
- 优化主模型配置

## 最佳实践总结

1. **配置合理的 Fallback 列表**：按能力从高到低排列
2. **设置适当的超时时间**：根据模型特性调整
3. **监控 Fallback 频率**：及时发现异常
4. **定期审查配置**：根据实际情况调整
5. **设置告警**：Fallback 频率异常时及时通知
6. **记录详细日志**：便于问题排查

## 限制说明

当前版本的 Fallback 机制有以下限制：

- 仅支持 **同 Provider 内** 模型 Fallback
- 不支持 **跨 Provider Fallback**
- 不支持 **指数退避**重试策略
- 不支持 **Fallback 次数限制**配置

以上限制可能在后续版本中改进。
