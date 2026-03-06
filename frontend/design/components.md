# Components

## 布局组件

### ConsoleLayout

**类型**: 布局组件

**描述**: Console 主布局，包含侧边栏和顶部导航栏

**子组件**:

- `Topbar`
- `Sidebar`
- `MainContent`

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| children | ReactNode | 主内容区子组件 |

---

### Topbar

**类型**: 布局组件

**描述**: 顶部导航栏

**子组件**:

- Logo
- 用户菜单

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| user | User | 当前用户信息 |
| onLogout | () => void | 登出回调 |

---

### Sidebar

**类型**: 布局组件

**描述**: 侧边导航栏

**子组件**:

- 导航菜单项

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| currentPath | string | 当前路径 |
| userRole | 'admin' \| 'user' | 用户角色（控制菜单显示） |

**菜单项**:

| 路径 | 图标 | 标签 | 权限 |
|------|------|------|------|
| / | chart-bar | 概览 | 所有用户 |
| /models | cube | 模型列表 | 所有用户 |
| /api-keys | key | API Keys | 所有用户 |
| /usage | chart-line | 使用统计 | 所有用户 |
| /providers | server | Providers | 仅管理员 |
| /users | users | 用户管理 | 仅管理员 |
| /settings | cog | 设置 | 所有用户 |

---

## 表单组件

### LoginForm

**类型**: 表单组件

**描述**: 登录表单

**字段**:

| 名称 | 类型 | 必填 | 描述 |
|------|------|------|------|
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码 |

**API**:

- `POST /api/v1/auth/login`

**数据结构**:

```typescript
interface LoginFormValues {
  email: string;
  password: string;
}

interface LoginResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
}
```

---

### RegisterForm

**类型**: 表单组件

**描述**: 注册表单

**字段**:

| 名称 | 类型 | 必填 | 描述 |
|------|------|------|------|
| name | string | 是 | 用户姓名 |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码（至少 8 位） |
| confirmPassword | string | 是 | 确认密码 |

**验证规则**:

- 密码至少 8 位
- 确认密码必须与密码一致

**API**:

- `POST /api/v1/auth/register`

**数据结构**:

```typescript
interface RegisterFormValues {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

interface RegisterResponse {
  id: number;
  name: string;
  email: string;
  role: string;
  status: string;
  created_at: string;
}
```

---

### ProviderForm

**类型**: 表单组件

**描述**: Provider 创建/编辑表单

**字段**:

| 名称 | 类型 | 必填 | 描述 |
|------|------|------|------|
| name | string | 是 | Provider 名称 |
| type | string | 是 | Provider 类型 |
| base_url | string | 是 | API Base URL |
| timeout | number | 是 | 超时时间（秒） |
| api_key | string | 否 | API Key |
| enabled | boolean | 是 | 是否启用 |
| fallback_models | string[] | 是 | Fallback 模型列表 |
| extra_config | object | 否 | 额外配置 |

**API**:

- `POST /api/v1/providers` - 创建
- `PUT /api/v1/providers/:name` - 更新

**数据结构**:

```typescript
interface ProviderFormValues {
  name: string;
  type: 'openai' | 'azure' | 'anthropic' | 'custom';
  base_url: string;
  timeout: number;
  api_key?: string;
  enabled: boolean;
  fallback_models: string[];
  extra_config?: Record<string, any>;
}

interface Provider {
  id: number;
  name: string;
  type: string;
  base_url: string;
  timeout: number;
  enabled: boolean;
  fallback_models: string[];
  created_at: string;
  updated_at?: string;
}
```

---

### CreateApiKeyForm

**类型**: 表单组件

**描述**: API Key 创建表单

**字段**:

| 名称 | 类型 | 必填 | 描述 |
|------|------|------|------|
| name | string | 是 | API Key 名称 |

**API**:

- `POST /api/v1/users/:id/api-keys`

**数据结构**:

```typescript
interface CreateApiKeyFormValues {
  name: string;
}

interface ApiKey {
  id: number;
  key: string; // 仅创建时返回
  key_prefix: string;
  name: string;
  status: string;
  created_at: string;
}
```

---

### ChangePasswordForm

**类型**: 表单组件

**描述**: 修改密码表单

**字段**:

| 名称 | 类型 | 必填 | 描述 |
|------|------|------|------|
| currentPassword | string | 是 | 当前密码 |
| newPassword | string | 是 | 新密码 |
| confirmPassword | string | 是 | 确认新密码 |

**验证规则**:

- 新密码至少 8 位
- 确认密码必须与新密码一致

---

## 表格组件

### ProviderTable

**类型**: 表格组件

**描述**: Provider 列表表格

**列**:

| 名称 | 描述 | 宽度 |
|------|------|------|
| name | Provider 名称 | 180px |
| type | 类型 | 120px |
| base_url | Base URL | 200px |
| enabled | 启用状态 | 100px |
| is_running | 运行状态 | 100px |
| fallback_models | Fallback 模型 | 200px |
| created_at | 创建时间 | 180px |
| actions | 操作 | 150px |

**操作**:

- 编辑 - 打开编辑表单
- 删除 - 删除 Provider（需确认）
- 重载 - 重载 Provider
- 启用/禁用 - 切换状态

**API**:

- `GET /api/v1/providers` - 获取列表

**数据结构**:

```typescript
interface ProviderWithStatus {
  provider: Provider;
  is_running: boolean;
}
```

---

### ApiKeyTable

**类型**: 表格组件

**描述**: API Key 列表表格

**列**:

| 名称 | 描述 | 宽度 |
|------|------|------|
| key_prefix | Key 前缀 | 150px |
| name | 名称 | 200px |
| status | 状态 | 100px |
| created_at | 创建时间 | 180px |
| last_used_at | 最后使用时间 | 180px |
| actions | 操作 | 100px |

**操作**:

- 删除 - 撤销 Key（需确认）
- 复制 - 复制完整 Key

**API**:

- `GET /api/v1/users/:id/api-keys` - 获取列表

**数据结构**:

```typescript
interface ApiKeyListItem {
  id: number;
  key_prefix: string;
  name: string;
  status: string;
  created_at: string;
  last_used_at?: string;
}
```

---

### ModelListTable

**类型**: 表格组件

**描述**: 模型列表表格

**列**:

| 名称 | 描述 | 宽度 |
|------|------|------|
| provider_name | Provider 名称 | 150px |
| model | 模型名称 | 200px |
| full_name | 完整标识 | 250px |
| enabled | 状态 | 100px |
| actions | 操作 | 100px |

**操作**:

- 测试 - 打开测试模态框

**数据结构**:

```typescript
interface ModelInfo {
  provider_name: string;
  type: string;
  model: string;
  full_name: string;
  enabled: boolean;
}
```

---

### UsageRecordsTable

**类型**: 表格组件

**描述**: 使用记录表格

**列**:

| 名称 | 描述 | 宽度 |
|------|------|------|
| timestamp | 时间 | 180px |
| provider_name | Provider | 150px |
| model | 模型 | 200px |
| prompt_tokens | Prompt Tokens | 120px |
| completion_tokens | Completion Tokens | 140px |
| total_tokens | Total Tokens | 120px |
| latency_ms | 延迟 | 100px |
| status | 状态 | 100px |

**API**:

- `GET /api/v1/usage` - 获取记录

**数据结构**:

```typescript
interface UsageRecord {
  id: number;
  user_id: number;
  model: string;
  provider_name: string;
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  latency_ms: number;
  status: 'success' | 'error';
  timestamp: string;
}
```

---

### UserTable

**类型**: 表格组件

**描述**: 用户列表表格

**列**:

| 名称 | 描述 | 宽度 |
|------|------|------|
| id | 用户 ID | 80px |
| name | 姓名 | 150px |
| email | 邮箱 | 200px |
| role | 角色 | 100px |
| status | 状态 | 100px |
| created_at | 创建时间 | 180px |
| actions | 操作 | 150px |

**操作**:

- 编辑 - 打开编辑表单
- 禁用/启用 - 修改状态
- 删除 - 删除用户（需确认）
- 查看 API Keys - 跳转到 API Key 页面

**API**:

- `GET /api/v1/users` - 获取列表

**数据结构**:

```typescript
interface User {
  id: number;
  name: string;
  email: string;
  role: 'admin' | 'user';
  status: 'active' | 'disabled';
  created_at: string;
}
```

---

## 数据展示组件

### UsageStatsCard

**类型**: 数据卡片组件

**描述**: 使用统计概览卡片

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| totalRequests | number | 总请求数 |
| totalTokens | number | 总 Token 数 |
| successRate | number | 成功率（百分比） |
| avgLatency | number | 平均延迟（毫秒） |

---

### RequestChart

**类型**: 图表组件

**描述**: 请求趋势图表

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| data | ChartData[] | 图表数据 |
| timeRange | '7d' \| '30d' | 时间范围 |
| onTimeRangeChange | (range) => void | 时间范围变化回调 |

**数据结构**:

```typescript
interface ChartData {
  date: string;
  requests: number;
  tokens: number;
}
```

---

### ProviderStatusCard

**类型**: 数据卡片组件

**描述**: Provider 状态概览卡片

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| providers | ProviderSummary[] | Provider 列表 |

**数据结构**:

```typescript
interface ProviderSummary {
  name: string;
  type: string;
  enabled: boolean;
  is_running: boolean;
}
```

---

## 模态框组件

### ProviderFormModal

**类型**: 模态框组件

**描述**: Provider 创建/编辑模态框

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| provider | Provider \| null | 编辑的 Provider（null 为新建） |
| open | boolean | 是否打开 |
| onClose | () => void | 关闭回调 |
| onSuccess | () => void | 成功回调 |

---

### CreateApiKeyModal

**类型**: 模态框组件

**描述**: 创建 API Key 模态框

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| open | boolean | 是否打开 |
| onClose | () => void | 关闭回调 |
| onSuccess | (key: string) => void | 成功回调，返回完整 Key |

---

### ApiKeyCreatedModal

**类型**: 模态框组件

**描述**: API Key 创建成功展示模态框

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| open | boolean | 是否打开 |
| apiKey | string | 完整的 API Key |
| onClose | () => void | 关闭回调 |

**功能**:

- 显示完整 API Key
- 复制按钮
- 安全提示

---

### ModelTestModal

**类型**: 模态框组件

**描述**: 模型测试模态框

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| open | boolean | 是否打开 |
| model | string | 模型标识 |
| onClose | () => void | 关闭回调 |

**功能**:

- 消息输入框
- 参数配置
- 发送测试请求
- 显示响应结果

---

### UserEditModal

**类型**: 模态框组件

**描述**: 用户编辑模态框

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| user | User \| null | 编辑的用户（null 为新建） |
| open | boolean | 是否打开 |
| onClose | () => void | 关闭回调 |
| onSuccess | () => void | 成功回调 |

---

## 通用组件

### StatusBadge

**类型**: 标签组件

**描述**: 状态标签

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| status | 'active' \| 'disabled' \| 'success' \| 'error' \| 'loading' | 状态值 |
| text | string | 显示文本 |

**样式**:

- active: 绿色背景
- disabled: 红色背景
- success: 绿色背景
- error: 红色背景
- loading: 蓝色背景 + 动画

---

### ConfirmModal

**类型**: 模态框组件

**描述**: 确认对话框

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| open | boolean | 是否打开 |
| title | string | 标题 |
| message | string | 消息内容 |
| onConfirm | () => void | 确认回调 |
| onCancel | () => void | 取消回调 |

---

### DateRangePicker

**类型**: 选择器组件

**描述**: 日期范围选择器

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| value | [Date, Date] \| null | 选中的日期范围 |
| onChange | (range: [Date, Date]) => void | 变化回调 |

---

### CopyButton

**类型**: 按钮组件

**描述**: 复制按钮

**属性**:

| 名称 | 类型 | 描述 |
|------|------|------|
| text | string | 要复制的文本 |
| onCopy | () => void | 复制成功回调 |

**功能**:

- 点击复制文本到剪贴板
- 复制成功后显示提示
