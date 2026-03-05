# dashboard-ui Specification

## Purpose
TBD - created by archiving change implement-dashboard-ui. Update Purpose after archive.
## Requirements
### Requirement: Web 界面部署

系统 SHALL 提供 Web 界面，通过 Docker Compose 与后端服务一起部署。

#### Scenario: 通过 Docker Compose 启动 Web 界面

- **GIVEN** docker-compose.yml 配置了 nginx 服务
- **WHEN** 执行 `docker-compose up` 命令
- **THEN** nginx 服务启动在 80 端口
- **AND** 可以通过浏览器访问 `http://localhost`
- **AND** nginx 正确提供静态文件服务

#### Scenario: nginx 反向代理 API 请求

- **GIVEN** 前端代码调用 `/api/v1/*` 路径的 API
- **WHEN** 浏览器发送 API 请求
- **THEN** nginx 将请求代理到后端 courier:8080 服务
- **AND** 返回后端 API 响应给前端

### Requirement: 用户注册界面

系统 SHALL 提供用户注册界面，支持用户创建账户。

#### Scenario: 显示注册页面

- **WHEN** 用户访问 `/register` 路由
- **THEN** 显示注册页面
- **AND** 页面包含姓名输入框
- **AND** 页面包含邮箱输入框
- **AND** 页面包含密码输入框
- **AND** 页面包含确认密码输入框
- **AND** 页面包含注册按钮
- **AND** 页面包含登录链接

#### Scenario: 注册成功

- **GIVEN** 用户在注册页面输入了有效的注册信息
- **WHEN** 点击注册按钮
- **THEN** 前端调用 `POST /api/v1/auth/register` 接口
- **AND** 注册成功后跳转到登录页面
- **AND** 显示"注册成功，请登录"提示

#### Scenario: 注册失败

- **GIVEN** 用户在注册页面输入了已存在的邮箱
- **WHEN** 点击注册按钮
- **THEN** 前端调用 `POST /api/v1/auth/register` 接口
- **AND** 收到错误响应
- **AND** 显示错误提示信息（如"邮箱已存在"）
- **AND** 保持在注册页面

#### Scenario: 注册表单验证

- **GIVEN** 用户在注册页面
- **WHEN** 邮箱格式不正确、密码少于 8 位或确认密码不匹配时点击注册按钮
- **THEN** 显示表单验证错误提示
- **AND** 不发送 API 请求

### Requirement: 用户登录界面

系统 SHALL 提供用户登录界面，支持邮箱和密码登录。

#### Scenario: 显示登录页面

- **WHEN** 用户访问 `/login` 路由且未登录
- **THEN** 显示登录页面
- **AND** 页面包含邮箱输入框
- **AND** 页面包含密码输入框
- **AND** 页面包含登录按钮
- **AND** 页面包含注册链接

#### Scenario: 登录成功

- **GIVEN** 用户在登录页面输入了正确的邮箱和密码
- **WHEN** 点击登录按钮
- **THEN** 前端调用 `POST /api/v1/auth/login` 接口
- **AND** 登录成功后保存 access_token 和 refresh_token 到 localStorage
- **AND** 跳转到 Dashboard 首页

#### Scenario: 登录失败

- **GIVEN** 用户在登录页面输入了错误的邮箱或密码
- **WHEN** 点击登录按钮
- **THEN** 前端调用 `POST /api/v1/auth/login` 接口
- **AND** 收到 401 错误响应
- **AND** 显示错误提示信息
- **AND** 保持在登录页面

#### Scenario: 登录表单验证

- **GIVEN** 用户在登录页面
- **WHEN** 邮箱格式不正确或密码为空时点击登录按钮
- **THEN** 显示表单验证错误提示
- **AND** 不发送 API 请求

### Requirement: 认证状态管理

系统 SHALL 管理用户认证状态，自动处理 Token 刷新。

#### Scenario: 请求携带 Token

- **GIVEN** 用户已登录
- **WHEN** 前端发送 API 请求
- **THEN** 请求头自动包含 `Authorization: Bearer <access_token>`

#### Scenario: Token 过期自动刷新

- **GIVEN** 用户已登录
- **WHEN** API 请求返回 401 未授权错误
- **THEN** 前端自动调用 `POST /api/v1/auth/refresh` 接口
- **AND** 使用 refresh_token 获取新的 access_token
- **AND** 重新发送原请求
- **AND** 用户无感知

#### Scenario: Token 刷新失败

- **GIVEN** 用户已登录
- **WHEN** refresh_token 也已过期
- **THEN** 调用 `POST /api/v1/auth/refresh` 返回 401
- **AND** 清除本地存储的 Token
- **AND** 跳转到登录页面
- **AND** 显示"会话已过期，请重新登录"提示

### Requirement: 主布局和导航

系统 SHALL 提供主布局界面，包含顶部导航栏和内容区域。

#### Scenario: 显示主布局

- **GIVEN** 用户已登录
- **WHEN** 访问受保护的路由
- **THEN** 显示顶部导航栏（高度 60px）
- **AND** 顶部导航栏左侧显示 Logo 和产品名称
- **AND** 顶部导航栏右侧显示用户信息和注销按钮
- **AND** 显示主内容区域

#### Scenario: 顶部用户信息

- **GIVEN** 用户已登录
- **WHEN** 查看顶部导航栏
- **THEN** 显示当前用户邮箱
- **AND** 显示注销按钮

### Requirement: 用户注销

系统 SHALL 提供用户注销功能。

#### Scenario: 点击注销按钮

- **GIVEN** 用户已登录
- **WHEN** 点击顶部导航栏的注销按钮
- **THEN** 清除本地存储的 Token 和用户信息
- **AND** 跳转到登录页面

### Requirement: Dashboard 首页

系统 SHALL 提供 Dashboard 首页，展示系统使用统计和状态概览。

#### Scenario: 显示 Dashboard 首页

- **GIVEN** 用户已登录
- **WHEN** 访问 `/` 路由
- **THEN** 前端调用 `GET /api/v1/usage` 接口获取使用统计
- **AND** 前端调用 `GET /api/v1/providers` 接口获取 Provider 状态
- **AND** 显示统计卡片（总请求数、总 Token 数、成功率、平均延迟）
- **AND** 显示 Provider 状态卡片
- **AND** 显示最近活动记录

#### Scenario: 使用统计卡片

- **GIVEN** 用户在 Dashboard 首页
- **WHEN** 查看使用统计卡片
- **THEN** 显示总请求数（数字格式化显示）
- **AND** 显示总 Token 数（数字格式化显示）
- **AND** 显示成功率（百分比格式）
- **AND** 显示平均延迟（毫秒格式）

#### Scenario: Provider 状态卡片

- **GIVEN** 用户在 Dashboard 首页
- **WHEN** 查看 Provider 状态卡片
- **THEN** 显示所有 Provider 的名称
- **AND** 显示每个 Provider 的运行状态（运行中/已停止）
- **AND** 使用绿色标签表示运行中
- **AND** 使用红色标签表示已停止

#### Scenario: 最近活动记录

- **GIVEN** 用户在 Dashboard 首页
- **WHEN** 查看最近活动记录
- **THEN** 显示最近的 API 调用记录（最多 10 条）
- **AND** 每条记录显示时间、模型名称、状态、Token 数量
- **AND** 使用绿色标签表示成功
- **AND** 使用红色标签表示失败

#### Scenario: Dashboard 加载状态

- **GIVEN** 用户访问 Dashboard 首页
- **WHEN** 数据正在加载
- **THEN** 显示加载动画或骨架屏
- **AND** 加载完成后显示实际数据

#### Scenario: Dashboard 错误处理

- **GIVEN** 用户访问 Dashboard 首页
- **WHEN** API 请求失败
- **THEN** 显示错误提示信息
- **AND** 提供重试按钮

### Requirement: 路由控制

系统 SHALL 实现前端路由控制，保护需要认证的页面。

#### Scenario: 未登录访问受保护路由

- **GIVEN** 用户未登录
- **WHEN** 直接访问 `/` 路由
- **THEN** 自动跳转到登录页面
- **AND** 登录成功后跳转回原目标页面（Dashboard）

#### Scenario: 已登录访问登录页

- **GIVEN** 用户已登录
- **WHEN** 访问 `/login` 路由
- **THEN** 自动跳转到 Dashboard 首页

#### Scenario: 已登录访问注册页

- **GIVEN** 用户已登录
- **WHEN** 访问 `/register` 路由
- **THEN** 自动跳转到 Dashboard 首页

### Requirement: 错误处理

系统 SHALL 提供友好的错误提示和处理机制。

#### Scenario: 网络错误

- **GIVEN** 用户正在进行操作
- **WHEN** 网络请求失败（如网络断开）
- **THEN** 显示"网络错误，请检查网络连接"提示
- **AND** 提供重试按钮

#### Scenario: 服务器错误

- **GIVEN** 用户正在进行操作
- **WHEN** API 返回 500 错误
- **THEN** 显示"服务器错误，请稍后重试"提示

#### Scenario: 注册速率限制

- **GIVEN** 用户尝试注册
- **WHEN** API 返回 429 速率限制错误
- **THEN** 显示"注册请求过于频繁，请稍后再试"提示

### Requirement: 响应式设计

系统 SHALL 支持基本的响应式设计，适配不同屏幕尺寸。

#### Scenario: 桌面端显示

- **GIVEN** 用户使用桌面浏览器（宽度 ≥ 1024px）
- **WHEN** 访问 Web 界面
- **THEN** 显示完整布局（顶部导航栏 + 主内容区）
- **AND** Dashboard 统计卡片按网格布局显示

#### Scenario: 移动端显示

- **GIVEN** 用户使用移动浏览器（宽度 < 640px）
- **WHEN** 访问 Web 界面
- **THEN** 顶部导航栏保持显示
- **AND** Dashboard 统计卡片单列布局显示

### Requirement: 加载状态

系统 SHALL 在加载数据时显示加载状态。

#### Scenario: API 请求加载中

- **GIVEN** 用户触发需要调用 API 的操作
- **WHEN** API 请求正在处理
- **THEN** 相关按钮显示加载状态（禁用并显示 loading 图标）
- **AND** 或显示全屏加载遮罩

#### Scenario: 页面初始加载

- **GIVEN** 用户访问页面
- **WHEN** 页面数据正在加载
- **THEN** 显示骨架屏或加载动画
- **AND** 加载完成后显示实际内容

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

