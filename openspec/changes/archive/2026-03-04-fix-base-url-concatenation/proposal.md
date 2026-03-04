# fix-base-url-concatenation Proposal

## 概述

修复 OpenAI Adapter 中 `buildChatURL` 函数的 base URL 拼接逻辑，使其能够正确处理非标准路径结构的 AI Provider。

## 问题背景

当前 `buildChatURL` 函数假设所有 OpenAI 兼容的 API 都遵循 `{base_url}/v1/chat/completions` 的路径结构。然而，部分 AI Provider（如智谱 GLM）使用不同的路径结构：

- 智谱 GLM：`https://open.bigmodel.cn/api/paas/v4/chat/completions`
- 对应的 base_url 应为：`https://open.bigmodel.cn/api/paas/v4`

当使用当前的 `buildChatURL` 逻辑时，会错误地拼接为：
- `https://open.bigmodel.cn/api/paas/v4` + `/v1/chat/completions`
- 结果：`https://open.bigmodel.cn/api/paas/v4/v1/chat/completions`（错误）

## 目标

1. 修改 `buildChatURL` 函数，仅在 base_url 后添加 `/chat/completions`
2. 确保兼容现有的标准 OpenAI 路径结构
3. 更新相关文档

## 影响范围

- `internal/adapter/openai/client.go`：修改 `buildChatURL` 函数
- `docs/api.md`：更新文档中关于 base_url 的说明
- `openspec/specs/provider-adapter/spec.md`：更新相关需求

## 向后兼容性

此变更对标准 OpenAI API 兼容服务的影响：

| base_url 输入 | 当前行为 | 修改后行为 | 兼容性 |
|--------------|---------|-----------|-------|
| `https://api.openai.com` | → `/v1/chat/completions` | → `/chat/completions` | ⚠️ 需用户更新配置 |
| `https://api.openai.com/v1` | → `/v1/chat/completions` | → `/v1/chat/completions` | ✅ 完全兼容 |
| `https://api.openai.com/v1/chat/completions` | → `/v1/chat/completions` | → `/v1/chat/completions` | ✅ 完全兼容 |

## 解决方案

### 1. 修改 `buildChatURL` 逻辑

```go
func buildChatURL(baseURL string) string {
    baseURL = strings.TrimSuffix(baseURL, "/")

    // 如果已经包含完整路径，直接返回
    if strings.HasSuffix(baseURL, "/chat/completions") {
        return baseURL
    }

    // 只添加 /chat/completions
    return baseURL + "/chat/completions"
}
```

### 2. 文档更新

明确说明 base_url 应包含完整的 API 路径前缀：

- 标准配置：`base_url` 设为 `https://api.openai.com/v1`
- 智谱 GLM：`base_url` 设为 `https://open.bigmodel.cn/api/paas/v4`

## 替代方案（已否决）

### 方案 A：添加路径检测逻辑

尝试自动检测 base_url 是否已包含版本号，动态决定是否添加 `/v1`。

**否决原因**：
- 增加复杂度，难以覆盖所有情况
- 可能产生意外的行为
- 让用户明确指定完整路径更清晰

### 方案 B：添加配置选项

在 Provider 配置中添加 `chat_path` 字段，让用户自定义完整路径。

**否决原因**：
- 增加配置复杂度
- 大多数 OpenAI 兼容服务遵循相同模式
- 简化 `buildChatURL` 逻辑更优雅

## 风险评估

| 风险 | 影响 | 缓解措施 |
|-----|-----|---------|
| 现有配置失效 | 中 | 发布说明中明确标注迁移指南 |
| 用户配置错误 | 低 | 文档中提供清晰的配置示例 |

## 相关规格

- `provider-adapter`：Provider 适配器相关需求
- `openai-adapter-chat-method`：OpenAI Adapter Chat 方法实现
- `openai-adapter-chat-stream-method`：OpenAI Adapter ChatStream 方法实现
