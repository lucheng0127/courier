# fix-base-url-concatenation Tasks

## 实施任务

### 1. 修改 buildChatURL 函数

**文件**: `internal/adapter/openai/client.go`

**任务**:
- 简化 `buildChatURL` 函数逻辑
- 移除自动添加 `/v1` 的逻辑
- 仅在 base_url 后添加 `/chat/completions`
- 保留对已包含完整路径的检测

**验证**:
- 函数输入 `https://api.openai.com/v1` → 输出 `https://api.openai.com/v1/chat/completions`
- 函数输入 `https://open.bigmodel.cn/api/paas/v4` → 输出 `https://open.bigmodel.cn/api/paas/v4/chat/completions`
- 函数输入 `https://api.openai.com/v1/chat/completions` → 输出 `https://api.openai.com/v1/chat/completions`（不变）

---

### 2. 更新 API 文档

**文件**: `docs/api.md`

**任务**:
- 在 Provider 管理章节中，更新 `base_url` 字段的说明
- 添加不同 Provider 的配置示例
- 明确说明 base_url 应包含完整的 API 路径前缀

**示例添加**:
```json
{
  "name": "zhipu-glm",
  "type": "openai",
  "base_url": "https://open.bigmodel.cn/api/paas/v4",
  "api_key": "your-api-key",
  "timeout": 60,
  "enabled": true
}
```

---

### 3. 更新 Provider Adapter 规格

**文件**: `openspec/specs/provider-adapter/spec.md`

**任务**:
- 修改 `OpenAI Adapter Chat 方法实现` 需求中的场景描述
- 移除关于自动添加 `/v1` 的说明
- 明确 base_url 应包含完整路径前缀的要求

---

### 4. 添加单元测试

**文件**: `internal/adapter/openai/client_test.go`

**任务**:
- 为 `buildChatURL` 函数添加单元测试
- 覆盖各种输入情况：
  - 标准路径（包含 `/v1`）
  - 非标准路径（如 `/api/paas/v4`）
  - 已包含完整路径
  - 末尾有斜杠的情况

---

### 5. 验证 vLLM Adapter

**任务**:
- 确认 vLLM Adapter 使用 openai 包的 `buildChatURL` 函数
- 验证修改后对 vLLM Provider 无负面影响

---

### 6. 创建迁移指南

**文件**: `docs/migration/base-url-change.md`（新建）

**任务**:
- 说明配置变更的必要性
- 提供配置迁移示例
- 列出常见 Provider 的正确配置方式

---

## 任务顺序

1. ✏️ 修改 `buildChatURL` 函数
2. 🧪 添加单元测试
3. 📖 更新 API 文档
4. 📋 更新规格文档
5. ✅ 验证 vLLM Adapter
6. 📝 创建迁移指南

任务 1-2 可以并行进行，任务 3-6 依赖任务 1 完成。

---

## 验收标准

- [x] `buildChatURL` 函数能正确处理各种 base_url 格式
- [x] 单元测试通过，覆盖率达到要求
- [x] 文档更新完成，包含清晰的配置示例
- [x] 规格文档同步更新
- [x] vLLM Adapter 功能正常
