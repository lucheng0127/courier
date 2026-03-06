# API Key 管理和对话页面 - 实现任务

## 阶段 1：API Key 管理

### 1.1 创建 API Key 相关文件和类型
- [ ] 创建 `frontend/src/api/api-keys.ts` - API Key 相关的 API 调用
- [ ] 创建 `frontend/src/stores/api-keys.ts` - API Key 状态管理
- [ ] 在 `frontend/src/types/index.ts` 中添加 ApiKey 类型定义（如尚未定义）

### 1.2 实现 API Key 列表页面
- [ ] 创建 `frontend/src/views/ApiKeysView.vue`
- [ ] 实现页面布局：顶部创建按钮 + 表格
- [ ] 实现表格列：名称、Key 前缀、状态、创建时间、最后使用时间、操作
- [ ] 添加删除功能（带确认对话框）
- [ ] 添加加载状态和错误处理

### 1.3 实现创建 API Key 模态框
- [ ] 添加"创建 API Key"按钮到页面顶部
- [ ] 创建表单模态框：名称输入字段
- [ ] 实现表单验证
- [ ] 调用 `POST /api/v1/users/:id/api-keys` API
- [ ] 创建成功后显示完整 Key（模态框）
- [ ] 添加复制完整 Key 的功能
- [ ] 关闭后无法再次查看完整 Key

### 1.4 添加路由和导航
- [ ] 在 `frontend/src/router/index.ts` 中添加 `/api-keys` 路由
- [ ] 在 Sidebar 中添加"API Keys"菜单项

### 1.5 测试和优化
- [ ] 测试创建 API Key 流程
- [ ] 测试删除 API Key 流程
- [ ] 测试错误处理（网络错误、权限错误）
- [ ] 确保响应式设计正常工作

## 阶段 2：对话页面

### 2.1 创建对话相关文件和类型
- [ ] 创建 `frontend/src/api/chat.ts` - Chat API 调用
- [ ] 创建 `frontend/src/stores/chat.ts` - 对话状态管理
- [ ] 在 `frontend/src/types/index.ts` 中添加对话相关类型

### 2.2 实现对话页面基础结构
- [ ] 创建 `frontend/src/views/ChatView.vue`
- [ ] 实现三栏布局：选择区域、对话区域、输入区域
- [ ] 或实现两栏布局：顶部选择 + 中间对话 + 底部输入

### 2.3 实现 Provider 和 Model 选择器
- [ ] 添加 Provider 下拉选择器
- [ ] 添加 Model 下拉选择器
- [ ] 实现联动：选择 Provider 后动态更新 Model 列表
- [ ] 显示完整的模型标识（provider/model）

### 2.4 实现 API Key 检查逻辑
- [ ] 页面加载时检查用户是否有 API Key
- [ ] 如果没有 API Key，显示提示页面
- [ ] 提供"创建 API Key"按钮跳转到 `/api-keys`
- [ ] 有 API Key 时显示对话界面

### 2.5 实现流式响应
- [ ] 实现 `fetch` 调用 `POST /v1/chat/completions`
- [ ] 设置 `stream: true` 参数
- [ ] 实现 `ReadableStream` 读取
- [ ] 解析 SSE 格式（`data: {...}\n\n`）
- [ ] 实时追加消息内容到界面
- [ ] 处理 `data: [DONE]` 结束标记

### 2.6 实现对话历史管理
- [ ] 定义消息类型（role, content）
- [ ] 实现消息列表渲染
- [ ] 区分用户消息和 AI 消息样式
- [ ] 支持多轮对话（传递历史 messages）

### 2.7 实现输入和发送功能
- [ ] 创建文本输入框（支持多行）
- [ ] 创建发送按钮
- [ ] 实现 Enter 发送（Shift+Enter 换行）
- [ ] 发送时禁用输入和按钮
- [ ] 响应完成后重新启用

### 2.8 优化用户体验
- [ ] 添加 Markdown 渲染（可选，使用 marked 或类似库）
- [ ] 添加自动滚动到最新消息
- [ ] 添加错误处理和重试机制
- [ ] 添加流式响应的加载动画

### 2.9 添加路由和导航
- [ ] 在 `frontend/src/router/index.ts` 中添加 `/chat` 路由
- [ ] 在 Sidebar 中添加"Chat"菜单项

### 2.10 测试和优化
- [ ] 测试无 API Key 时的提示流程
- [ ] 测试 Provider 和 Model 选择
- [ ] 测试流式响应
- [ ] 测试多轮对话
- [ ] 测试错误处理
- [ ] 确保响应式设计正常工作

## 阶段 3：集成和优化

### 3.1 样式统一
- [ ] 确保新页面与现有页面样式一致
- [ ] 使用相同的颜色主题
- [ ] 使用相同的组件样式（Card, Table, Button 等）

### 3.2 性能优化
- [ ] 对话消息列表使用虚拟滚动（如消息很多）
- [ ] 优化流式响应的渲染性能

### 3.3 安全性检查
- [ ] 确保用户只能管理自己的 API Key
- [ ] 确保 API Key 正确传递给 Chat API
- [ ] 确保错误信息不泄露敏感数据

### 3.4 文档和测试
- [ ] 更新用户使用文档
- [ ] 进行端到端测试
- [ ] 进行跨浏览器测试

## 依赖关系

- 任务 1.x 可以并行开始
- 任务 2.x 依赖 1.x 完成（因为需要 API Key 功能）
- 任务 3.x 在所有功能完成后进行
