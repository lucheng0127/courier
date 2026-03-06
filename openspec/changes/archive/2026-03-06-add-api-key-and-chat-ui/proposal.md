# 提案：API Key 管理和对话页面

## 概述

本提案旨在为 Courier LLM Gateway 的 Dashboard 添加 API Key 管理页面和对话页面，完善前端功能，使用户能够直接通过 Web 界面管理 API Key 并进行 AI 对话测试。

## 背景

当前 Dashboard 已实现：
- 用户注册和登录
- Provider 管理（管理员）
- Model 列表查看
- 基础的侧边栏导航和布局

缺失的功能：
- API Key 创建和管理界面
- AI 对话测试界面
- 完整的用户自助工作流

## 目标

1. **API Key 管理页面**
   - 用户可以创建、查看、删除自己的 API Key
   - 显示 API Key 的状态、创建时间、最后使用时间
   - 创建时显示完整 Key（仅一次）

2. **对话页面**
   - 用户可以选择 Provider 和 Model 进行对话
   - 支持流式响应（stream mode）
   - 如果用户没有 API Key，提示用户先创建
   - 保存对话历史，支持多轮对话

## 影响范围

### 修改的规范

- **dashboard-ui** - 添加新的页面需求

### 相关规范

- **apikey-auth** - API Key 管理的 API 接口已定义
- **chat-api** - Chat API 的流式响应已实现

## 设计决策

### API Key 页面设计

**位置**：`/api-keys`

**功能**：
1. 顶部显示"创建 API Key"按钮
2. 表格展示所有 API Key：
   - 名称
   - Key 前缀（sk-xxx...）
   - 状态
   - 创建时间
   - 最后使用时间
   - 操作（删除）

**创建流程**：
1. 点击"创建 API Key"按钮
2. 弹出模态框，输入 Key 名称
3. 提交后显示完整的 Key（带复制按钮）
4. 关闭模态框后无法再次查看完整 Key

### 对话页面设计

**位置**：`/chat`

**功能**：
1. 左侧（或顶部）选择区域：
   - Provider 下拉选择
   - Model 下拉选择（根据选中的 Provider 动态更新）
   - 显示完整模型标识（provider/model）

2. 中间对话区域：
   - 消息列表（用户消息 + AI 响应）
   - 流式响应实时显示
   - Markdown 渲染支持

3. 底部输入区域：
   - 文本输入框
   - 发送按钮

**检查逻辑**：
- 进入页面时检查用户是否有 API Key
- 如果没有，显示提示："您还没有创建 API Key，请先创建 API Key 再进行对话"
- 提供"创建 API Key"按钮跳转到 API Key 管理页面

**API Key 选择逻辑**：
- 默认使用用户 API Key 列表中第一个状态为 `active` 的 Key
- 后续可扩展为让用户手动选择特定的 API Key

### 技术实现要点

1. **流式响应处理**
   - 使用 `fetch` API + `ReadableStream`
   - 解析 SSE 格式（`data: {...}\n\n`）
   - 实时追加到消息内容

2. **状态管理**
   - 创建 `useChatStore` 管理对话状态
   - 创建 `useApiKeyStore` 管理 API Key
   - 对话历史保存在内存中（可扩展到 localStorage）

3. **样式一致性**
   - 使用与现有页面相同的卡片、表格、按钮样式
   - 使用 Ant Design Vue 组件
   - 响应式设计适配移动端

## 用户故事

### 作为一个普通用户
- 我希望能够创建自己的 API Key
- 我希望能够查看我的所有 API Key
- 我希望能够删除不需要的 API Key
- 我希望能够测试 AI 对话功能
- 我希望能够选择不同的 Provider 和 Model 进行对话
- 我希望能够看到流式的响应效果

### 作为一个管理员
- 我希望能够看到所有用户的 API Key（后续需求）
- 当前提案仅关注用户自己的 API Key 管理

## 非目标

以下功能不在本次提案范围内：
- 管理员查看其他用户的 API Key
- 对话历史的持久化存储
- 对话导出功能
- 多会话管理
- 系统提示词配置

## 依赖关系

- 后端 API Key 接口已完成（`POST /api/v1/users/:id/api-keys`、`GET /api/v1/users/:id/api-keys`、`DELETE /api/v1/users/:id/api-keys/:key_id`）
- 后端 Chat API 已完成（`POST /v1/chat/completions` 支持 stream 模式）
- 前端路由和认证基础已就绪

## 风险和缓解

| 风险 | 缓解措施 |
|------|----------|
| 流式响应解析复杂 | 参考 Ant Design Vue 的 Message 组件实现 |
| API Key 安全存储 | 仅在创建时显示，提醒用户保存 |
| Provider/Model 联动选择 | 使用 computed 属性实现动态更新 |
| 对话历史管理 | 先实现内存存储，后续扩展到后端 |

## 验收标准

- [ ] 用户可以创建 API Key
- [ ] 创建后显示完整 Key（仅一次）
- [ ] 用户可以查看所有 API Key 列表
- [ ] 用户可以删除 API Key
- [ ] 对话页面正确显示 Provider 和 Model 选择器
- [ ] 无 API Key 时显示提示
- [ ] 对话响应以流式方式显示
- [ ] 支持多轮对话
- [ ] 页面样式与现有页面保持一致
