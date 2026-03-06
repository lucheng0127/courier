# Dashboard UI 设计文档

## Context

Courier 项目当前是一个纯后端的 LLM Gateway 服务，提供 RESTful API 接口。为了提升用户体验，需要添加一个 Web Dashboard 前端项目，让用户可以通过浏览器进行注册、登录和查看系统概览。

**技术约束**：
- 必须通过 Docker Compose 部署
- 使用 nginx 部署编译后的静态文件
- 前端设计遵循 `frontend/design/` 中的设计文档

**利益相关者**：
- 最终用户：需要通过界面注册账户、登录系统、查看使用统计
- 运维人员：需要通过 Docker Compose 一键部署

## Goals / Non-Goals

### Goals
- 提供用户注册功能
- 提供用户登录和注销功能
- 提供 Dashboard 概览页面（使用统计、Provider 状态）
- 通过 Docker Compose 与后端服务一起部署
- 现代化的 UI 设计（参考 Vercel、Stripe Console）

### Non-Goals
- 不实现 API Key 管理界面（后续阶段）
- 不实现 Provider 管理界面（后续阶段）
- 不实现详细的使用统计页面（后续阶段）
- 不实现聊天界面（后续阶段）

## Decisions

### 1. 前端技术栈

**决策**：使用 Vue 3 + TypeScript + Ant Design Vue + Vite

**理由**：
- Vue 3 是成熟的渐进式框架，Composition API 提供更好的类型推导
- Ant Design Vue 是企业级 UI 组件库，设计风格现代、组件丰富
- Vite 提供快速的开发体验和构建速度
- TypeScript 提供类型安全，降低维护成本

**替代方案**：
- View UI Plus：同样可行，但 Ant Design Vue 生态更活跃
- React + Ant Design：同样可行，但团队更熟悉 Vue

### 2. 部署架构

**决策**：nginx 作为静态文件服务器，独立于后端 API 服务

**架构图**：
```
┌─────────────────┐
│   Browser       │
└────────┬────────┘
         │ HTTP
         ▼
┌─────────────────┐
│   nginx:80      │ ← 静态文件服务
│   (frontend)    │
└────────┬────────┘
         │ Proxy
         │ /api/v1/*
         ▼
┌─────────────────┐
│   courier:8080  │ ← API 服务
│   (backend)     │
└─────────────────┘
```

**理由**：
- 前后端分离，职责清晰
- nginx 高效处理静态文件
- nginx 可以作为 API 反向代理，解决 CORS 问题
- 便于独立扩展和部署

### 3. API 通信

**决策**：使用 axios 进行 HTTP 请求，实现请求/响应拦截器

**关键点**：
- 请求拦截器：自动添加 JWT Token 到 Authorization Header
- 响应拦截器：处理 401 错误，自动尝试刷新 Token
- Token 刷新失败时，自动跳转登录页

**理由**：
- 统一的错误处理
- 自动 Token 管理，用户无感知

### 4. 状态管理

**决策**：使用 Pinia 进行状态管理

**理由**：
- Vue 3 官方推荐的状态管理库
- 比 Vuex 更简洁、类型友好
- 适合管理用户认证状态、使用统计数据等全局状态

**状态结构**：
```typescript
{
  // auth store
  user: {
    id: number,
    name: string,
    email: string,
    role: string
  },
  token: {
    accessToken: string,
    refreshToken: string,
    expiresAt: number
  },

  // dashboard store
  stats: {
    totalRequests: number,
    totalTokens: number,
    successRate: number,
    avgLatency: number
  },
  providers: Array<{
    name: string,
    type: string,
    enabled: boolean,
    is_running: boolean
  }>,
  recentActivity: Array<{
    timestamp: string,
    model: string,
    status: string,
    tokens: number
  }>
}
```

### 5. 路由设计

**决策**：使用 Vue Router，基于认证状态控制路由访问

**路由表**：
```typescript
{
  path: '/login',
  name: 'Login',
  component: LoginView,
  meta: { requiresAuth: false }
}
{
  path: '/register',
  name: 'Register',
  component: RegisterView,
  meta: { requiresAuth: false }
}
{
  path: '/',
  name: 'Dashboard',
  component: DashboardView,
  meta: { requiresAuth: true }
}
```

### 6. UI 设计规范

遵循 `frontend/design/ui-style.md` 中的设计规范：

**颜色方案**：
- 主色：`#10A37F` (OpenAI 绿)
- 背景：`#F9FAFB` (浅灰背景)
- 卡片背景：`#FFFFFF` (白色)
- 边框：`#E5E7EB` (浅灰边框)
- 文字：`#111827` (深灰)

**布局结构**：
- 顶部导航栏（60px）：Logo + 用户菜单
- 主内容区域：Dashboard 卡片

**组件风格**：
- 使用 Ant Design Vue 的组件
- 自定义样式覆盖以匹配设计规范
- 卡片式布局，圆角 12px

## Risks / Trade-offs

### 风险 1：CORS 问题
**风险**：前端和后端在不同端口/域名时，可能遇到 CORS 问题

**缓解措施**：
- nginx 配置 `/api/v1/*` 反向代理到后端服务
- 前端统一使用相对路径调用 API

### 风险 2：Token 过期处理
**风险**：Token 过期后用户操作可能失败

**缓解措施**：
- 实现自动 Token 刷新机制
- 刷新失败时提示用户重新登录

### 风险 3：构建产物体积
**风险**：前端打包后体积可能较大（Ant Design Vue）

**缓解措施**：
- 使用 Vite 的代码分割
- 按需导入 Ant Design Vue 组件
- 生产环境启用 gzip 压缩

### 风险 4：API 依赖
**风险**：Dashboard 首页依赖使用统计 API，如果数据量大可能影响性能

**缓解措施**：
- 使用分页和聚合查询
- 添加加载状态和错误处理
- 考虑缓存策略

## Migration Plan

### 实施步骤

1. **初始化前端项目**
   ```bash
   cd frontend
   npm create vite@latest . -- --template vue-ts
   npm install ant-design-vue axios pinia vue-router
   npm install @ant-design/icons-vue
   ```

2. **开发页面组件**
   - 注册页面
   - 登录页面
   - Dashboard 首页
   - 主布局组件

3. **配置 nginx**
   - 静态文件服务
   - API 反向代理

4. **更新 docker-compose.yml**
   - 添加 nginx 服务

5. **测试部署**
   ```bash
   docker-compose up --build
   ```

### 回滚计划

如果出现严重问题：
1. 移除 docker-compose.yml 中的 nginx 服务
2. 删除 frontend 目录
3. 系统恢复到纯 API 服务模式

## Open Questions

1. **Dashboard 数据刷新频率？**
   - 建议页面加载时获取一次，用户可手动刷新

2. **是否需要实时更新？**
   - 当前方案不支持，后续可考虑 WebSocket

3. **是否需要支持多语言？**
   - 当前方案不支持，后续可添加 vue-i18n
