# Provider 管理和 Model 列表实施任务清单

## 1. 布局更新

- [x] 1.1 修改主布局结构
  - [x] 1.1.1 将顶部导航栏改为侧边栏 + 顶部导航栏布局
  - [x] 1.1.2 创建 Sidebar 组件（240px 宽）
  - [x] 1.1.3 更新 Topbar 样式
  - [x] 1.1.4 调整主内容区布局

- [x] 1.2 添加导航菜单项
  - [x] 1.2.1 Dashboard 菜单项
  - [x] 1.2.2 Models 菜单项（所有用户）
  - [x] 1.2.3 Providers 菜单项（仅管理员）
  - [x] 1.2.4 API Keys 菜单项（后续）
  - [x] 1.2.5 Settings 菜单项（后续）

- [x] 1.3 实现权限控制
  - [x] 1.3.1 在 Sidebar 中根据角色显示菜单
  - [x] 1.3.2 更新路由守卫权限检查

## 2. Provider API 和状态管理

- [x] 2.1 创建 Provider API 模块
  - [x] 2.1.1 实现 getProviders 接口
  - [x] 2.1.2 实现 createProvider 接口
  - [x] 2.1.3 实现 updateProvider 接口
  - [x] 2.1.4 实现 deleteProvider 接口
  - [x] 2.1.5 实现 reloadProvider 接口
  - [x] 2.1.6 实现 enableProvider 接口
  - [x] 2.1.7 实现 disableProvider 接口
  - [x] 2.1.8 实现 getProviderModels 接口

- [x] 2.2 创建 Provider Store
  - [x] 2.2.1 定义 state：providers 列表, loading
  - [x] 2.2.2 实现 fetchProviders action
  - [x] 2.2.3 实现创建/更新/删除 action
  - [x] 2.2.4 实现重载/启用/禁用 action
  - [x] 2.2.5 实现获取模型列表 action

- [x] 2.3 扩展类型定义
  - [x] 2.3.1 更新 Provider 类型（添加更多字段）
  - [x] 2.3.2 定义 ProviderForm 类型
  - [x] 2.3.3 定义 ModelInfo 类型

## 3. Provider 管理页面

- [x] 3.1 创建 ProvidersView.vue 组件
- [x] 3.2 实现 Provider 列表表格
  - [x] 3.2.1 Provider 名称列
  - [x] 3.2.2 类型列
  - [x] 3.2.3 Base URL 列
  - [x] 3.2.4 启用状态列（Switch/Tag）
  - [x] 3.2.5 运行状态列
  - [x] 3.2.6 Fallback 模型列
  - [x] 3.2.7 创建时间列
  - [x] 3.2.8 操作列

- [x] 3.3 实现操作功能
  - [x] 3.3.1 创建 Provider 按钮
  - [x] 3.3.2 编辑按钮
  - [x] 3.3.3 删除按钮（带确认）
  - [x] 3.3.4 重载按钮
  - [x] 3.3.5 启用/禁用按钮
  - [x] 3.3.6 查看模型按钮（打开抽屉）

- [x] 3.4 创建 Provider 表单模态框
  - [x] 3.4.1 名称输入框（必填）
  - [x] 3.4.2 类型选择器（下拉）
  - [x] 3.4.3 Base URL 输入框（必填）
  - [x] 3.4.4 超时时间输入框（数字）
  - [x] 3.4.5 API Key 输入框（密码）
  - [x] 3.4.6 启用状态开关
  - [x] 3.4.7 Fallback 模型标签输入
  - [x] 3.4.8 提交和取消按钮

- [x] 3.5 创建 Provider 模型列表抽屉
  - [x] 3.5.1 模型列表展示
  - [x] 3.5.2 复制模型标识按钮

- [x] 3.6 实现表单验证
  - [x] 3.6.1 必填字段验证
  - [x] 3.6.2 Base URL 格式验证
  - [x] 3.6.3 Fallback 模型至少一个

- [x] 3.7 实现加载和错误状态
  - [x] 3.7.1 列表加载动画
  - [x] 3.7.2 操作成功/失败提示
  - [x] 3.7.3 错误处理和重试

## 4. Model 列表页面

- [x] 4.1 创建 ModelsView.vue 组件
- [x] 4.2 实现 Model 列表展示
  - [x] 4.2.1 按 Provider 分组（Collapse 组件）
  - [x] 4.2.2 每个分组显示 Provider 信息
  - [x] 4.2.3 模型列表展示
  - [x] 4.2.4 模型标识显示（provider/model）
  - [x] 4.2.5 状态标签

- [x] 4.3 实现交互功能
  - [x] 4.3.1 展开/折叠分组
  - [x] 4.3.2 复制模型标识按钮
  - [x] 4.3.3 刷新按钮

- [x] 4.4 实现加载和空状态
  - [x] 4.4.1 加载动画
  - [x] 4.4.2 空状态提示
  - [x] 4.4.3 错误状态处理

## 5. 路由配置

- [x] 5.1 添加新路由
  - [x] 5.1.1 `/providers` - Provider 管理（管理员）
  - [x] 5.1.2 `/models` - Model 列表（所有用户）
- [x] 5.2 更新路由守卫
  - [x] 5.2.1 添加管理员权限检查
  - [x] 5.2.2 权限不足跳转到首页

## 6. 组件更新

- [x] 6.1 更新 DashboardView.vue
  - [x] 6.1.1 添加 Sidebar 组件
  - [x] 6.1.2 调整布局结构
  - [x] 6.1.3 移除旧的顶部导航样式

## 7. 样式和响应式

- [x] 7.1 实现 Sidebar 样式
  - [x] 7.1.1 宽度 240px
  - [x] 7.1.2 背景色和边框
  - [x] 7.1.3 菜单项样式
  - [x] 7.1.4 激活状态高亮
  - [x] 7.1.5 悬停效果

- [x] 7.2 实现移动端适配
  - [x] 7.2.1 Sidebar 折叠为抽屉
  - [x] 7.2.2 汉堡菜单按钮

## 8. 测试和验证

- [x] 8.1 功能测试
  - [x] 8.1.1 Provider CRUD 操作测试
  - [x] 8.1.2 权限控制测试
  - [x] 8.1.3 Model 列表展示测试

- [x] 8.2 浏览器兼容性测试
  - [x] 8.2.1 Chrome 测试
  - [x] 8.2.2 Firefox 测试
  - [x] 8.2.3 Safari 测试
