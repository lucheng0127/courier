# dashboard-ui Specification Delta

## ADDED Requirements

### Requirement: 侧边栏导航

系统 SHALL 提供侧边栏导航，包含主要功能菜单项。

#### Scenario: 显示侧边栏

- **GIVEN** 用户已登录
- **WHEN** 访问任意受保护的路由
- **THEN** 显示侧边栏（宽度 240px）
- **AND** 侧边栏包含以下菜单项：
  - Dashboard（图标：dashboard）
  - Models（图标：block）
  - Providers（图标：server，仅管理员可见）
- **AND** 当前页面对应的菜单项高亮显示

#### Scenario: 菜单项导航

- **GIVEN** 用户在侧边栏中
- **WHEN** 点击某个菜单项
- **THEN** 跳转到对应页面
- **AND** 菜单项状态更新为激活

#### Scenario: 管理员权限显示

- **GIVEN** 用户角色为 admin
- **WHEN** 查看侧边栏
- **THEN** 显示所有菜单项（包括 Providers）

#### Scenario: 普通用户权限显示

- **GIVEN** 用户角色为 user
- **WHEN** 查看侧边栏
- **THEN** 不显示 Providers 菜单项

### Requirement: Provider 管理页面

系统 SHALL 提供 Provider 管理页面，允许管理员创建、编辑、删除和管理 Provider。

#### Scenario: 显示 Provider 管理页面

- **GIVEN** 用户已登录且角色为 admin
- **WHEN** 访问 `/providers` 路由
- **THEN** 显示 Provider 管理页面
- **AND** 前端调用 `GET /api/v1/providers` 接口
- **AND** 以表格形式显示 Provider 列表
- **AND** 每行显示：名称、类型、Base URL、启用状态、运行状态、Fallback 模型、创建时间、操作按钮

#### Scenario: 创建 Provider

- **GIVEN** 用户在 Provider 管理页面
- **WHEN** 点击"创建 Provider"按钮
- **THEN** 显示创建表单模态框
- **AND** 表单包含：名称、类型、Base URL、超时时间、API Key、启用状态、Fallback 模型
- **AND** 提交后调用 `POST /api/v1/providers` 接口
- **AND** 成功后刷新列表并显示成功提示

#### Scenario: 编辑 Provider

- **GIVEN** 用户在 Provider 管理页面
- **WHEN** 点击某个 Provider 的编辑按钮
- **THEN** 显示编辑表单模态框，预填充当前 Provider 数据
- **AND** 修改后调用 `PUT /api/v1/providers/:name` 接口
- **AND** 成功后刷新列表并显示成功提示

#### Scenario: 删除 Provider

- **GIVEN** 用户在 Provider 管理页面
- **WHEN** 点击某个 Provider 的删除按钮
- **THEN** 显示确认对话框
- **AND** 确认后调用 `DELETE /api/v1/providers/:name` 接口
- **AND** 成功后刷新列表并显示成功提示

#### Scenario: 启用/禁用 Provider

- **GIVEN** 用户在 Provider 管理页面
- **WHEN** 点击启用/禁用按钮
- **THEN** 调用 `POST /api/v1/admin/providers/:name/enable` 或 `disable` 接口
- **AND** 成功后刷新列表

#### Scenario: 重载 Provider

- **GIVEN** 用户在 Provider 管理页面
- **WHEN** 点击重载按钮
- **THEN** 调用 `POST /api/v1/admin/providers/:name/reload` 接口
- **AND** 显示加载状态
- **AND** 成功后刷新 Provider 运行状态

#### Scenario: 查看 Provider 模型列表

- **GIVEN** 用户在 Provider 管理页面
- **WHEN** 点击"查看模型"按钮
- **THEN** 打开抽屉组件
- **AND** 调用 `GET /api/v1/providers/:name/models` 接口
- **AND** 在抽屉中显示该 Provider 的所有模型

#### Scenario: Provider 权限控制

- **GIVEN** 用户角色为 user
- **WHEN** 直接访问 `/providers` 路由
- **THEN** 自动跳转到 Dashboard 首页
- **AND** 显示"权限不足"提示

#### Scenario: Provider 表单验证

- **GIVEN** 用户在创建/编辑 Provider 表单中
- **WHEN** 必填字段为空或格式不正确时提交
- **THEN** 显示表单验证错误
- **AND** 不提交 API 请求

### Requirement: Model 列表页面

系统 SHALL 提供 Model 列表页面，展示所有可用的模型。

#### Scenario: 显示 Model 列表页面

- **GIVEN** 用户已登录
- **WHEN** 访问 `/models` 路由
- **THEN** 显示 Model 列表页面
- **AND** 前端调用 `GET /api/v1/providers` 接口获取所有 Provider
- **AND** 对每个 Provider 调用 `GET /api/v1/providers/:name/models` 接口获取模型列表
- **AND** 按 Provider 分组展示所有模型

#### Scenario: Provider 分组展示

- **GIVEN** 用户在 Model 列表页面
- **WHEN** 查看模型列表
- **THEN** 每个 Provider 显示为一个可折叠的面板
- **AND** 面板标题显示 Provider 名称和类型
- **AND** 展开时显示该 Provider 的所有模型

#### Scenario: 模型信息展示

- **GIVEN** 用户展开某个 Provider 分组
- **THEN** 显示该 Provider 的所有模型
- **AND** 每个模型显示：
  - 模型名称
  - 完整标识（provider/model 格式）
  - 状态标签（enabled/disabled）
- **AND** 提供"复制标识"按钮

#### Scenario: 复制模型标识

- **GIVEN** 用户在 Model 列表页面
- **WHEN** 点击某个模型的"复制"按钮
- **THEN** 将完整的模型标识（provider/model）复制到剪贴板
- **AND** 显示"已复制"提示

#### Scenario: Model 列表加载状态

- **GIVEN** 用户访问 Model 列表页面
- **WHEN** 数据正在加载
- **THEN** 显示加载动画或骨架屏
- **AND** 加载完成后显示实际数据

#### Scenario: Model 列表空状态

- **GIVEN** 系统中没有可用的 Provider 或模型
- **WHEN** 用户访问 Model 列表页面
- **THEN** 显示空状态提示
- **AND** 提示用户联系管理员添加 Provider

### Requirement: 权限控制扩展

系统 SHALL 实现基于角色的页面访问控制。

#### Scenario: 管理员访问所有页面

- **GIVEN** 用户角色为 admin
- **WHEN** 访问任意页面
- **THEN** 允许访问

#### Scenario: 普通用户访问受限页面

- **GIVEN** 用户角色为 user
- **WHEN** 访问仅管理员可访问的页面（如 /providers）
- **THEN** 自动跳转到 Dashboard 首页
- **AND** 显示"权限不足"提示
