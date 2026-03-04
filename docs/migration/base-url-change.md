# Base URL 配置迁移指南

## 变更说明

从版本 v0.2.0 开始，Provider 的 `base_url` 配置要求已更新。**`base_url` 必须包含完整的 API 路径前缀**（如 `/v1`、`/api/paas/v4` 等）。

系统会在 `base_url` 后自动追加 `/chat/completions` 来构建完整的请求路径。

## 变更原因

此前，系统会自动在 `base_url` 后添加 `/v1/chat/completions`。但部分 AI Provider（如智谱 GLM）使用不同的路径结构（如 `/api/paas/v4/chat/completions`），导致无法正确适配。

为支持更多 Provider，我们简化了 URL 构建逻辑：**仅在 base_url 后追加 `/chat/completions`**。

## 迁移步骤

### 1. 检查现有 Provider 配置

运行以下命令查看所有 Provider 配置：

```bash
curl -X GET http://your-server/api/v1/providers \
  -H "Authorization: Bearer <admin-jwt-token>"
```

### 2. 更新 base_url 配置

对于每个 Provider，按照下表更新 `base_url`：

#### OpenAI

**旧配置**（如果使用了不含路径的 base_url）：
```json
{
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com"
}
```

**新配置**：
```json
{
  "name": "openai-main",
  "type": "openai",
  "base_url": "https://api.openai.com/v1"
}
```

#### 智谱 GLM

```json
{
  "name": "zhipu-glm",
  "type": "openai",
  "base_url": "https://open.bigmodel.cn/api/paas/v4",
  "api_key": "your-api-key"
}
```

#### 通义千问

```json
{
  "name": "qwen-main",
  "type": "openai",
  "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1",
  "api_key": "your-api-key"
}
```

#### vLLM（本地部署）

```json
{
  "name": "vllm-local",
  "type": "vllm",
  "base_url": "http://localhost:8000/v1"
}
```

### 3. 更新 Provider 配置

使用 PUT 请求更新每个 Provider：

```bash
curl -X PUT http://your-server/api/v1/providers/openai-main \
  -H "Authorization: Bearer <admin-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "base_url": "https://api.openai.com/v1",
    "timeout": 60
  }'
```

### 4. 重载 Provider

更新配置后，重载 Provider 使其生效：

```bash
# 重载所有 Provider
curl -X POST http://your-server/api/v1/admin/providers/reload \
  -H "Authorization: Bearer <admin-jwt-token>"

# 或重载单个 Provider
curl -X POST http://your-server/api/v1/admin/providers/openai-main/reload \
  -H "Authorization: Bearer <admin-jwt-token>"
```

## 常见 Provider 配置参考

| Provider | base_url |
|----------|----------|
| OpenAI | `https://api.openai.com/v1` |
| 智谱 GLM | `https://open.bigmodel.cn/api/paas/v4` |
| 通义千问 | `https://dashscope.aliyuncs.com/compatible-mode/v1` |
| DeepSeek | `https://api.deepseek.com/v1` |
| Moonshot（月之暗面） | `https://api.moonshot.cn/v1` |
| 百度文心 | `https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop` |
| 本地 vLLM | `http://localhost:8000/v1` |
| 本地 Ollama | `http://localhost:11434/v1` |

## 验证配置

更新后，可以通过发送测试请求验证配置是否正确：

```bash
curl -X POST http://your-server/v1/chat/completions \
  -H "Authorization: Bearer <api-key>" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "provider-name/model-name",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## 回滚

如果更新后出现问题，可以回滚到旧配置并使用兼容模式：

1. 将 `base_url` 改回旧值（不含 `/v1`）
2. 确保完整路径包含 `/chat/completions`

例如：
```json
{
  "base_url": "https://api.openai.com/v1/chat/completions"
}
```

这样系统会检测到已包含完整路径，不会重复添加。

## 问题排查

### 问题：请求失败，提示路径不存在

**原因**：`base_url` 配置不正确，缺少路径前缀。

**解决方案**：按照上述表格更新 `base_url`，确保包含完整路径前缀。

### 问题：某个特定 Provider 无法使用

**原因**：该 Provider 可能使用非标准的路径结构。

**解决方案**：
1. 查看该 Provider 的官方文档，确定完整的 API 路径
2. 将 `base_url` 设置为完整路径（包含 `/chat/completions` 之前的部分）

例如，如果 Provider 的完整路径是 `https://api.example.com/v2/api/chat/completions`，则：
```json
{
  "base_url": "https://api.example.com/v2/api"
}
```

## 联系支持

如有任何问题，请查看项目文档或提交 Issue。
