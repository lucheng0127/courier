# web-ui Specification

## Purpose

定义 Courier 项目的 Web 用户界面功能，提供用户认证和 API Key 管理的可视化操作界面。

## ADDED Requirements

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

### Requirement: 用户登录界面

系统 SHALL 提供用户登录界面，支持邮箱和密码登录。

#### Scenario: 显示登录页面

- **WHEN** 用户访问 `http://localhost` 且未登录
- **THEN** 显示登录页面
- **AND** 页面包含邮箱输入框
- **AND** 页面包含密码输入框
- **AND** 页面包含登录按钮
- **AND** 页面风格简洁，使用中性色配色

#### Scenario: 登录成功

- **GIVEN** 用户在登录页面输入了正确的邮箱和密码
- **WHEN** 点击登录按钮
- **THEN** 前端调用 `POST /api/v1/auth/login` 接口
- **AND** 登录成功后保存 access_token 和 refresh_token 到 localStorage
- **AND** 跳转到 API Key 管理页面

#### Scenario: 登录失败

- **GIVEN** 用户在登录页面输入了错误的邮箱或密码
- **WHEN** 点击登录按钮
- **THEN** 前端调用 `POST /api/v1/auth/login` 接口
- **AND** 收到 401 错误响应
- **AND** 显示错误提示信息
- **AND** 保持在登录页面

#### Scenario: 表单验证

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

系统 SHALL 提供主布局界面，包含导航栏和内容区域。

#### Scenario: 显示主布局

- **GIVEN** 用户已登录
- **WHEN** 访问受保护的路由
- **THEN** 显示左侧垂直导航栏
- **AND** 显示顶部 Header（包含用户信息和注销按钮）
- **AND** 显示主内容区域

#### Scenario: 导航栏菜单

- **GIVEN** 用户已登录
- **WHEN** 查看左侧导航栏
- **THEN** 显示"API Key"菜单项
- **AND** 点击后跳转到 API Key 管理页面
- **AND** 当前激活的菜单项高亮显示

#### Scenario: 顶部用户信息

- **GIVEN** 用户已登录
- **WHEN** 查看顶部 Header
- **THEN** 显示当前用户邮箱
- **AND** 显示注销按钮

### Requirement: 用户注销

系统 SHALL 提供用户注销功能。

#### Scenario: 点击注销按钮

- **GIVEN** 用户已登录
- **WHEN** 点击顶部 Header 的注销按钮
- **THEN** 清除本地存储的 Token 和用户信息
- **AND** 跳转到登录页面

### Requirement: API Key 管理页面

系统 SHALL 提供 API Key 管理页面，允许用户查看、创建和撤销 API Key。

#### Scenario: 显示 API Key 列表

- **GIVEN** 用户已登录
- **WHEN** 访问 API Key 管理页面
- **THEN** 前端调用 `GET /api/v1/users/:id/api-keys` 接口
- **AND** 以表格形式显示 API Key 列表
- **AND** 每行显示：Key 前缀、名称、状态、创建时间、最后使用时间
- **AND** 提供撤销按钮

#### Scenario: 创建新 API Key

- **GIVEN** 用户在 API Key 管理页面
- **WHEN** 点击"创建 API Key"按钮
- **THEN** 显示创建对话框
- **AND** 对话框包含 Key 名称输入框
- **AND** 对话框包含确认和取消按钮

#### Scenario: 创建 API Key 成功

- **GIVEN** 用户在创建对话框中输入了 Key 名称
- **WHEN** 点击确认按钮
- **THEN** 前端调用 `POST /api/v1/users/:id/api-keys` 接口
- **AND** 成功后显示完整 API Key
- **AND** 提示用户"请妥善保存，此 Key 只显示一次"
- **AND** 提供"复制"和"关闭"按钮
- **AND** 刷新 API Key 列表

#### Scenario: 撤销 API Key

- **GIVEN** 用户在 API Key 管理页面
- **WHEN** 点击某个 API Key 的撤销按钮
- **THEN** 显示确认对话框
- **AND** 对话框提示"确认撤销此 API Key？此操作不可恢复"
- **AND** 点击确认后调用 `DELETE /api/v1/users/:id/api-keys/:key_id` 接口
- **AND** 成功后刷新列表并显示撤销成功提示

#### Scenario: 复制 API Key

- **GIVEN** 创建 API Key 成功后显示了完整 Key
- **WHEN** 点击"复制"按钮
- **THEN** 将 API Key 复制到剪贴板
- **AND** 显示"已复制到剪贴板"提示

### Requirement: 路由控制

系统 SHALL 实现前端路由控制，保护需要认证的页面。

#### Scenario: 未登录访问受保护路由

- **GIVEN** 用户未登录
- **WHEN** 直接访问 `/api-keys` 路由
- **THEN** 自动跳转到登录页面
- **AND** 登录成功后跳转回原目标页面

#### Scenario: 已登录访问登录页

- **GIVEN** 用户已登录
- **WHEN** 访问 `/login` 路由
- **THEN** 自动跳转到 API Key 管理页面

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

#### Scenario: 权限错误

- **GIVEN** 用户已登录但权限不足
- **WHEN** API 返回 403 错误
- **THEN** 显示"权限不足"提示

### Requirement: 响应式设计

系统 SHALL 支持基本的响应式设计，适配不同屏幕尺寸。

#### Scenario: 桌面端显示

- **GIVEN** 用户使用桌面浏览器（宽度 ≥ 768px）
- **WHEN** 访问 Web 界面
- **THEN** 显示完整布局（左侧导航栏 + 主内容区）

#### Scenario: 移动端显示

- **GIVEN** 用户使用移动浏览器（宽度 < 768px）
- **WHEN** 访问 Web 界面
- **THEN** 导航栏折叠为汉堡菜单
- **AND** 点击后展开导航选项

### Requirement: 加载状态

系统 SHALL 在加载数据时显示加载状态。

#### Scenario: API 请求加载中

- **GIVEN** 用户触发需要调用 API 的操作
- **WHEN** API 请求正在处理
- **THEN** 相关按钮显示加载状态（如禁用并显示 loading 图标）
- **AND** 或显示全屏加载遮罩
