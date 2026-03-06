# dashboard-ui Spec Delta

## ADDED Requirements

### Requirement: API Key 管理页面

系统 SHALL 提供 API Key 管理页面，允许用户创建、查看和删除自己的 API Key。

#### Scenario: 显示 API Key 管理页面

- **GIVEN** 用户已登录
- **WHEN** 访问 `/api-keys` 路由
- **THEN** 显示 API Key 管理页面
- **AND** 前端调用 `GET /api/v1/users/:id/api-keys` 接口（:id 为当前用户 ID）
- **AND** 页面顶部显示"创建 API Key"按钮
- **AND** 以表格形式显示用户的 API Key 列表
- **AND** 每行显示：名称、Key 前缀、状态、创建时间、最后使用时间、操作按钮

#### Scenario: 创建 API Key

- **GIVEN** 用户在 API Key 管理页面
- **WHEN** 点击"创建 API Key"按钮
- **THEN** 显示创建表单模态框
- **AND** 表单包含：名称输入框（必填）
- **AND** 提交后调用 `POST /api/v1/users/:id/api-keys` 接口（:id 为当前用户 ID）
- **AND** 成功后显示新的模态框展示完整的 API Key
- **AND** 显示复制按钮用于复制完整 Key
- **AND** 提示用户"请妥善保存此 Key，关闭后将无法再次查看"
- **AND** 关闭模态框后刷新列表并显示成功提示

#### Scenario: 查看 API Key 列表

- **GIVEN** 用户在 API Key 管理页面
- **WHEN** 查看 API Key 列表表格
- **THEN** 每个显示的字段包括：
  - **名称**: 用户定义的 Key 名称
  - **Key 前缀**: sk-xxx... 格式的前 10 位
  - **状态**: active（已启用）或 disabled（已禁用），使用标签显示
  - **创建时间**: 格式化的时间字符串
  - **最后使用时间**: 格式化的时间字符串，如未使用则显示"-"
  - **操作**: 删除按钮

#### Scenario: 删除 API Key

- **GIVEN** 用户在 API Key 管理页面
- **WHEN** 点击某个 API Key 的删除按钮
- **THEN** 显示确认对话框："确定要删除此 API Key 吗？删除后无法恢复。"
- **AND** 确认后调用 `DELETE /api/v1/users/:id/api-keys/:key_id` 接口
- **AND** 成功后刷新列表并显示"API Key 已删除"提示

#### Scenario: API Key 空状态

- **GIVEN** 用户在 API Key 管理页面
- **WHEN** 用户没有创建任何 API Key
- **THEN** 显示空状态提示："您还没有创建任何 API Key"
- **AND** 显示"创建 API Key"按钮引导用户创建

#### Scenario: API Key 列表加载状态

- **GIVEN** 用户访问 API Key 管理页面
- **WHEN** 数据正在加载
- **THEN** 显示加载动画或骨架屏
- **AND** 加载完成后显示实际数据

### Requirement: 对话页面

系统 SHALL 提供 AI 对话页面，允许用户选择 Provider 和 Model 进行流式对话测试。

#### Scenario: 显示对话页面

- **GIVEN** 用户已登录
- **WHEN** 访问 `/chat` 路由
- **THEN** 显示对话页面
- **AND** 页面包含以下区域：
  - **选择区域**: Provider 下拉选择器、Model 下拉选择器
  - **对话区域**: 消息列表显示
  - **输入区域**: 文本输入框、发送按钮

#### Scenario: Provider 和 Model 选择

- **GIVEN** 用户在对话页面
- **WHEN** 查看 Provider 下拉选择器
- **THEN** 显示所有启用的 Provider
- **AND** Provider 选项显示格式为 "名称 (类型)"

- **GIVEN** 用户在对话页面选择了 Provider
- **WHEN** 查看 Model 下拉选择器
- **THEN** 显示该 Provider 的所有 fallback_models
- **AND** Model 列表通过 `GET /api/v1/providers/:name/models` 接口获取

- **GIVEN** 用户在对话页面选择了 Provider 和 Model
- **WHEN** 查看选择显示区域
- **THEN** 显示完整的模型标识（provider/model 格式）
- **AND** 标识可被复制用于 API 调用

#### Scenario: 无 API Key 时显示提示

- **GIVEN** 用户已登录但没有 API Key
- **WHEN** 访问 `/chat` 路由
- **THEN** 前端检查用户是否有 API Key（调用 `GET /api/v1/users/:id/api-keys`）
- **AND** 如果没有 API Key，显示提示页面
- **AND** 提示内容："您还没有创建 API Key，请先创建 API Key 再进行对话"
- **AND** 提供"创建 API Key"按钮，点击后跳转到 `/api-keys` 页面

#### Scenario: 发送对话消息

- **GIVEN** 用户在对话页面且已选择 Provider 和 Model
- **AND** 用户有至少一个 active 状态的 API Key
- **WHEN** 在输入框中输入消息并点击发送
- **THEN** 前端自动选择用户 API Key 列表中第一个状态为 `active` 的 Key
- **AND** 使用该 Key 调用 `POST /v1/chat/completions` 接口
- **AND** 请求 Header 包含 `Authorization: Bearer <api_key>`
- **AND** 请求参数包括：
  - `model`: provider/model 格式
  - `messages`: 包含历史对话的数组
  - `stream`: true
- **AND** 发送期间禁用输入框和发送按钮
- **AND** 发送失败显示错误提示并重新启用输入

#### Scenario: 自动选择可用 API Key

- **GIVEN** 用户有多个 API Key
- **AND** 部分状态为 active，部分状态为 disabled
- **WHEN** 在对话页面发送消息
- **THEN** 系统自动使用第一个状态为 `active` 的 API Key
- **AND** 不显示 API Key 选择器（当前版本）
- **AND** 在对话界面不暴露使用的 Key 值

#### Scenario: 流式响应处理

- **GIVEN** 用户发送了对话消息
- **WHEN** 接收到流式响应（SSE 格式）
- **THEN** 前端使用 `ReadableStream` 读取响应
- **AND** 解析每个 `data:` 行
- **AND** 将响应内容实时追加到对话消息中
- **AND** 使用 Markdown 渲染 AI 响应内容
- **AND** 遇到 `data: [DONE]` 时停止读取
- **AND** 重新启用输入框和发送按钮

#### Scenario: 多轮对话支持

- **GIVEN** 用户在对话页面
- **WHEN** 进行多轮对话时
- **THEN** 前端维护对话历史消息数组
- **AND** 每次发送新消息时携带完整的历史 messages
- **AND** 历史消息在当前会话中保存在内存中

#### Scenario: 对话消息显示

- **GIVEN** 用户在对话页面
- **WHEN** 查看对话消息列表
- **THEN** 用户消息显示在右侧，头像标识为"用户"
- **AND** AI 消息显示在左侧，头像标识为"AI"
- **AND** 消息按时间顺序从上到下排列
- **AND** 最新的消息自动滚动到视图内

#### Scenario: 输入区域交互

- **GIVEN** 用户在对话页面
- **WHEN** 在输入框中输入内容
- **THEN** 支持 Enter 键发送消息
- **AND** 支持 Shift+Enter 换行
- **AND** 输入框支持多行文本输入

#### Scenario: 对话页面错误处理

- **GIVEN** 用户在对话页面
- **WHEN** 发送消息时遇到错误（网络错误、API 错误、超时等）
- **THEN** 在对话区域显示错误提示消息
- **AND** 提示消息样式与普通消息不同（如红色背景）
- **AND** 重新启用输入框和发送按钮允许重试

#### Scenario: 对话页面响应式设计

- **GIVEN** 用户使用不同设备访问对话页面
- **WHEN** 在桌面端（宽度 ≥ 1024px）
- **THEN** Provider 和 Model 选择器在页面顶部横向排列
- **AND** 对话区域占据主要空间

- **GIVEN** 用户在移动端（宽度 < 640px）
- **THEN** Provider 和 Model 选择器垂直排列
- **AND** 输入框固定在底部
- **AND** 对话区域占据剩余空间

### Requirement: 侧边栏菜单扩展

系统 SHALL 在侧边栏中添加新的菜单项。

#### Scenario: 显示 API Keys 菜单项

- **GIVEN** 用户已登录
- **WHEN** 查看侧边栏
- **THEN** 显示"API Keys"菜单项（图标：key）
- **AND** 菜单项位于"Models"之后、"Settings"之前

#### Scenario: 显示 Chat 菜单项

- **GIVEN** 用户已登录
- **WHEN** 查看侧边栏
- **THEN** 显示"Chat"菜单项（图标：message）
- **AND** 菜单项位于"API Keys"之后、"Settings"之前
