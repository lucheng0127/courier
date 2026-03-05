# Pages

## 登录页

**路径**: `/login`

**描述**: 用户登录页面

**布局**: 居中布局，无侧边栏

**组件**:

- `LoginForm` - 登录表单
  - 邮箱输入框
  - 密码输入框
  - 登录按钮
  - 注册链接

**API**:

- `POST /api/v1/auth/login` - 用户登录

**导航**:

- 登录成功 → 跳转到 Dashboard

---

## 注册页

**路径**: `/register`

**描述**: 用户注册页面

**布局**: 居中布局，无侧边栏

**组件**:

- `RegisterForm` - 注册表单
  - 姓名输入框
  - 邮箱输入框
  - 密码输入框
  - 确认密码输入框
  - 注册按钮
  - 登录链接

**API**:

- `POST /api/v1/auth/register` - 用户注册

**导航**:

- 注册成功 → 跳转到 登录页

---

## Dashboard

**路径**: `/`

**描述**: 概览页面，显示使用统计和系统状态

**布局**: Console 布局（侧边栏 + 主内容区）

**组件**:

- `UsageStatsCard` - 使用统计卡片
  - 总请求数
  - 总 Token 数
  - 成功率
  - 平均延迟

- `RequestChart` - 请求趋势图表
  - 按日期/时间显示请求数量
  - 可切换时间范围（7天/30天）

- `ProviderStatusCard` - Provider 状态卡片
  - 显示各 Provider 运行状态
  - 点击跳转到 Provider 列表

- `RecentActivityTable` - 最近活动记录
  - 最近的 API 调用记录
  - 显示时间、模型、状态、Token 数

**API**:

- `GET /api/v1/usage` - 获取使用统计
- `GET /api/v1/providers` - 获取 Provider 状态

**权限**:

- 所有已认证用户

---

## Provider 管理

**路径**: `/providers`

**描述**: Provider 列表和管理页面（仅管理员）

**布局**: Console 布局

**组件**:

- `ProviderTable` - Provider 列表表格
  - Provider 名称
  - 类型
  - Base URL
  - 状态（已启用/已禁用）
  - 运行状态
  - Fallback 模型数量
  - 创建时间
  - 操作按钮（编辑/删除/重载/启用/禁用）

- `ProviderFormModal` - Provider 创建/编辑表单
  - 名称输入框
  - 类型选择器
  - Base URL 输入框
  - 超时时间输入框
  - API Key 输入框
  - 启用状态开关
  - Fallback 模型列表
  - 额外配置 JSON 编辑器

- `ProviderModelListDrawer` - Provider 模型列表抽屉
  - 显示 Provider 的所有可用模型

**API**:

- `GET /api/v1/providers` - 获取 Provider 列表
- `POST /api/v1/providers` - 创建 Provider
- `PUT /api/v1/providers/:name` - 更新 Provider
- `DELETE /api/v1/providers/:name` - 删除 Provider
- `POST /api/v1/admin/providers/:name/reload` - 重载 Provider
- `POST /api/v1/admin/providers/:name/enable` - 启用 Provider
- `POST /api/v1/admin/providers/:name/disable` - 禁用 Provider
- `GET /api/v1/providers/:name/models` - 获取 Provider 模型列表

**权限**:

- 仅管理员

---

## 模型列表

**路径**: `/models`

**描述**: 所有 Provider 的模型列表

**布局**: Console 布局

**组件**:

- `ModelListTable` - 模型列表表格
  - Provider 名称
  - 模型名称
  - 完整模型标识（provider/model）
  - 状态
  - 操作（测试调用）

- `ModelTestModal` - 模型测试模态框
  - 模型选择器
  - 消息输入框
  - 参数配置
  - 测试结果展示

**API**:

- `GET /api/v1/providers` - 获取所有 Provider
- `GET /api/v1/providers/:name/models` - 获取 Provider 模型列表
- `POST /v1/chat/completions` - 测试模型调用

**权限**:

- 所有已认证用户

---

## API Key 管理

**路径**: `/api-keys`

**描述**: API Key 管理页面

**布局**: Console 布局

**组件**:

- `ApiKeyTable` - API Key 列表表格
  - Key 前缀（sk-xxx...）
  - 名称
  - 状态
  - 创建时间
  - 最后使用时间
  - 操作按钮（删除）

- `CreateApiKeyModal` - 创建 API Key 模态框
  - 名称输入框
  - 创建按钮

- `ApiKeyCreatedModal` - API Key 创建成功展示
  - 完整 Key 显示
  - 复制按钮
  - 确认提示

**API**:

- `GET /api/v1/users/:id/api-keys` - 获取 API Key 列表
- `POST /api/v1/users/:id/api-keys` - 创建 API Key
- `DELETE /api/v1/users/:id/api-keys/:key_id` - 撤销 API Key

**权限**:

- 管理员可查看所有用户的 API Key
- 普通用户只能查看自己的 API Key

---

## 使用统计

**路径**: `/usage`

**描述**: 详细的使用统计页面

**布局**: Console 布局

**组件**:

- `UsageFilterBar` - 筛选工具栏
  - 用户选择器（管理员）
  - Provider 选择器
  - 模型选择器
  - 日期范围选择器
  - 导出按钮

- `UsageStatsCards` - 统计卡片组
  - 总请求数
  - 总 Token 数
  - 成功率
  - 平均延迟
  - 总费用（如支持）

- `UsageChart` - 使用趋势图表
  - 请求量趋势
  - Token 消耗趋势
  - 可切换图表类型

- `UsageRecordsTable` - 使用记录表格
  - 时间
  - Provider
  - 模型
  - Prompt Tokens
  - Completion Tokens
  - Total Tokens
  - 延迟
  - 状态

**API**:

- `GET /api/v1/usage` - 获取使用统计记录

**权限**:

- 管理员可查询所有用户
- 普通用户只能查询自己

---

## 用户管理

**路径**: `/users`

**描述**: 用户管理页面（仅管理员）

**布局**: Console 布局

**组件**:

- `UserTable` - 用户列表表格
  - 用户 ID
  - 姓名
  - 邮箱
  - 角色
  - 状态
  - 创建时间
  - 操作按钮（编辑/禁用/删除/查看 API Keys）

- `UserEditModal` - 用户编辑模态框
  - 姓名输入框
  - 角色选择器
  - 保存按钮

- `UserStatusModal` - 用户状态修改模态框
  - 状态选择器
  - 确认按钮

**API**:

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户
- `PATCH /api/v1/users/:id/status` - 更新用户状态
- `DELETE /api/v1/users/:id` - 删除用户

**权限**:

- 仅管理员

---

## 个人设置

**路径**: `/settings`

**描述**: 个人设置页面

**布局**: Console 布局

**组件**:

- `ProfileCard` - 个人信息卡片
  - 姓名
  - 邮箱
  - 角色
  - 注册时间

- `ChangePasswordForm` - 修改密码表单
  - 当前密码输入框
  - 新密码输入框
  - 确认密码输入框
  - 保存按钮

**API**:

- `GET /api/v1/users/:id` - 获取个人信息
- `PUT /api/v1/users/:id` - 更新个人信息

**权限**:

- 所有已认证用户（仅能查看和修改自己）
