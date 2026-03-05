# Dashboard UI 实施任务清单

## 1. 项目初始化

- [ ] 1.1 清理或初始化 `frontend/` 目录
- [ ] 1.2 使用 Vite 初始化 Vue 3 + TypeScript 项目
  ```bash
  cd frontend
  npm create vite@latest . -- --template vue-ts
  ```
- [ ] 1.3 安装依赖：
  - [ ] 1.3.1 ant-design-vue（UI 组件库）
  - [ ] 1.3.2 @ant-design/icons-vue（图标库）
  - [ ] 1.3.3 axios（HTTP 请求）
  - [ ] 1.3.4 pinia（状态管理）
  - [ ] 1.3.5 vue-router（路由）
- [ ] 1.4 配置 Vite 开发服务器代理（开发环境）
- [ ] 1.5 配置 TypeScript 路径别名（@/ 指向 src/）

## 2. 项目基础配置

- [ ] 2.1 创建项目目录结构
  ```
  src/
  ├── api/           # API 服务层
  ├── assets/        # 静态资源
  ├── components/    # 通用组件
  ├── layouts/       # 布局组件
  ├── router/        # 路由配置
  ├── stores/        # Pinia stores
  ├── styles/        # 全局样式
  ├── types/         # TypeScript 类型
  ├── utils/         # 工具函数
  └── views/         # 页面组件
  ```
- [ ] 2.2 配置 Ant Design Vue 全局引入
- [ ] 2.3 自定义 Ant Design 主题色（#10A37F）
- [ ] 2.4 配置 Vue Router
- [ ] 2.5 配置 Pinia Store
- [ ] 2.6 配置 axios 实例和拦截器
- [ ] 2.7 定义全局样式（参考 `frontend/design/ui-style.md`）

## 3. 类型定义

- [ ] 3.1 定义 User 类型
- [ ] 3.2 定义 Token 类型
- [ ] 3.3 定义 ApiKey 类型
- [ ] 3.4 定义 Provider 类型
- [ ] 3.5 定义 UsageRecord 类型
- [ ] 3.6 定义 DashboardStats 类型

## 4. API 服务层

- [ ] 4.1 创建 axios 实例配置
  - [ ] 4.1.1 配置基础 URL
  - [ ] 4.1.2 实现请求拦截器（添加 Token）
  - [ ] 4.1.3 实现响应拦截器（处理错误和 Token 刷新）
- [ ] 4.2 创建 auth API 模块
  - [ ] 4.2.1 实现 register 接口
  - [ ] 4.2.2 实现 login 接口
  - [ ] 4.2.3 实现 refreshToken 接口
- [ ] 4.3 创建 usage API 模块
  - [ ] 4.3.1 实现 getUsageStats 接口
  - [ ] 4.3.2 实现 getRecentActivity 接口
- [ ] 4.4 创建 providers API 模块
  - [ ] 4.4.1 实现 getProviders 接口

## 5. 状态管理（Pinia）

- [ ] 5.1 创建 auth store（用户认证状态）
  - [ ] 5.1.1 定义 state：user, token, isAuthenticated
  - [ ] 5.1.2 实现 register action
  - [ ] 5.1.3 实现 login action
  - [ ] 5.1.4 实现 logout action
  - [ ] 5.1.5 实现 refreshToken action
  - [ ] 5.1.6 实现从 localStorage 恢复状态
- [ ] 5.2 创建 dashboard store（Dashboard 数据）
  - [ ] 5.2.1 定义 state：stats, providers, recentActivity, loading
  - [ ] 5.2.2 实现 fetchDashboardData action
  - [ ] 5.2.3 实现 refreshDashboardData action

## 6. 路由配置

- [ ] 6.1 定义路由表
  - [ ] 6.1.1 `/login` - 登录页
  - [ ] 6.1.2 `/register` - 注册页
  - [ ] 6.1.3 `/` - Dashboard 首页
- [ ] 6.2 实现路由守卫（认证检查）
- [ ] 6.3 实现未认证跳转逻辑
- [ ] 6.4 实现已登录访问登录/注册页的跳转逻辑

## 7. 页面组件

### 7.1 注册页面

- [ ] 7.1.1 创建 RegisterView.vue 组件
- [ ] 7.1.2 实现注册表单（姓名、邮箱、密码、确认密码）
- [ ] 7.1.3 实现表单验证规则
  - [ ] 7.1.3.1 邮箱格式验证
  - [ ] 7.1.3.2 密码至少 8 位
  - [ ] 7.1.3.3 确认密码匹配
- [ ] 7.1.4 连接 auth store 的 register action
- [ ] 7.1.5 处理注册成功/失败状态
- [ ] 7.1.6 添加错误提示显示
- [ ] 7.1.7 添加登录链接

### 7.2 登录页面

- [ ] 7.2.1 创建 LoginView.vue 组件
- [ ] 7.2.2 实现登录表单（邮箱、密码）
- [ ] 7.2.3 实现表单验证规则
- [ ] 7.2.4 连接 auth store 的 login action
- [ ] 7.2.5 处理登录成功/失败状态
- [ ] 7.2.6 添加错误提示显示
- [ ] 7.2.7 添加注册链接

### 7.3 主布局组件

- [ ] 7.3.1 创建 MainLayout.vue 组件
- [ ] 7.3.2 实现顶部导航栏（高度 60px）
  - [ ] 7.3.2.1 左侧 Logo 和产品名称
  - [ ] 7.3.2.2 右侧用户信息和注销按钮
- [ ] 7.3.3 实现主内容区域（router-view）
- [ ] 7.3.4 应用设计规范样式（颜色、间距）

### 7.4 Dashboard 首页

- [ ] 7.4.1 创建 DashboardView.vue 组件
- [ ] 7.4.2 实现页面加载时获取 Dashboard 数据
- [ ] 7.4.3 实现使用统计卡片组件
  - [ ] 7.4.3.1 总请求数卡片
  - [ ] 7.4.3.2 总 Token 数卡片
  - [ ] 7.4.3.3 成功率卡片
  - [ ] 7.4.3.4 平均延迟卡片
- [ ] 7.4.4 实现 Provider 状态卡片组件
  - [ ] 7.4.4.1 Provider 列表展示
  - [ ] 7.4.4.2 运行状态标签（绿色/红色）
- [ ] 7.4.5 实现最近活动表格组件
  - [ ] 7.4.5.1 时间列
  - [ ] 7.4.5.2 模型列
  - [ ] 7.4.5.3 状态列
  - [ ] 7.4.5.4 Token 数列
- [ ] 7.4.6 实现加载状态（骨架屏）
- [ ] 7.4.7 实现错误处理和重试
- [ ] 7.4.8 实现数据刷新功能

## 8. 通用组件

- [ ] 8.1 创建 StatCard 组件（统计卡片）
- [ ] 8.2 创建 ProviderStatusCard 组件（Provider 状态卡片）
- [ ] 8.3 创建 ActivityTable 组件（活动记录表格）
- [ ] 8.4 创建 LoadingSpinner 组件（加载动画）

## 9. 样式和设计规范

- [ ] 9.1 定义全局 CSS 变量（颜色、间距、字体）
- [ ] 9.2 自定义 Ant Design 主题
  - [ ] 9.2.1 主色：#10A37F
  - [ ] 9.2.2 圆角、间距等设计 token
- [ ] 9.3 实现卡片样式（参考设计规范）
- [ ] 9.4 实现响应式布局（桌面端和移动端）
- [ ] 9.5 添加过渡动画效果

## 10. Docker 部署配置

- [ ] 10.1 创建 `frontend/Dockerfile`
  - [ ] 10.1.1 基于 node:alpine 镜像
  - [ ] 10.1.2 多阶段构建（构建 + nginx 部署）
  - [ ] 10.1.3 安装依赖并构建生产版本
  - [ ] 10.1.4 复制到 nginx 静态目录
- [ ] 10.2 创建 `frontend/nginx.conf`
  - [ ] 10.2.1 配置静态文件服务
  - [ ] 10.2.2 配置 API 反向代理（/api/v1/* → courier:8080）
  - [ ] 10.2.3 配置 SPA 路由支持（try_files）
- [ ] 10.3 更新根目录 `docker-compose.yml`
  - [ ] 10.3.1 添加 nginx 服务
  - [ ] 10.3.2 配置端口映射 80:80
- [ ] 10.4 创建 `.dockerignore` 文件

## 11. 测试和验证

- [ ] 11.1 本地开发环境测试
  - [ ] 11.1.1 启动开发服务器
  - [ ] 11.1.2 测试注册功能
  - [ ] 11.1.3 测试登录功能
  - [ ] 11.1.4 测试注销功能
  - [ ] 11.1.5 测试 Dashboard 数据展示
  - [ ] 11.1.6 测试路由控制
- [ ] 11.2 Docker 构建测试
  - [ ] 11.2.1 构建 frontend 镜像
  - [ ] 11.2.2 通过 docker-compose 启动
  - [ ] 11.2.3 验证 nginx 服务正常
  - [ ] 11.2.4 验证 API 反向代理工作
- [ ] 11.3 端到端功能测试
  - [ ] 11.3.1 完整注册流程测试
  - [ ] 11.3.2 完整登录流程测试
  - [ ] 11.3.3 Token 自动刷新测试
  - [ ] 11.3.4 Dashboard 数据加载测试
  - [ ] 11.3.5 错误处理测试
- [ ] 11.4 浏览器兼容性测试
  - [ ] 11.4.1 Chrome 测试
  - [ ] 11.4.2 Firefox 测试
  - [ ] 11.4.3 Safari 测试

## 12. 文档

- [ ] 12.1 更新 README.md 添加 Web UI 访问说明
- [ ] 12.2 创建 frontend/README.md（前端项目说明）
- [ ] 12.3 记录环境变量配置（如有）
