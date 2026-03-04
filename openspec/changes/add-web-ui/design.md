# Web UI 设计文档

## Context

Courier 项目当前是一个纯后端的 LLM Gateway 服务，提供 RESTful API 接口。为了提升用户体验，需要添加一个 Web 管理界面，让用户可以通过浏览器进行登录和 API Key 管理操作。

**技术约束**：
- 必须通过 Docker Compose 部署
- 使用 nginx 部署编译后的静态文件
- 界面风格简洁、模块化、卡片式组件、中性色配色

**利益相关者**：
- 最终用户：需要通过界面管理 API Key
- 运维人员：需要通过 Docker Compose 一键部署

## Goals / Non-Goals

### Goals
- 提供用户登录和注销功能
- 提供 API Key 管理功能（列表、创建、撤销）
- 通过 Docker Compose 与后端服务一起部署
- 简洁、现代的 UI 设计

### Non-Goals
- 不实现 Provider 管理界面（后续阶段）
- 不实现使用统计界面（后续阶段）
- 不实现聊天界面（后续阶段）
- 不实现用户注册界面（用户通过 API 注册）

## Decisions

### 1. 前端技术栈

**决策**：使用 Vue 3 + View UI (iView) + Vite

**理由**：
- Vue 3 是成熟的渐进式框架，学习曲线平缓
- View UI (iView) 是成熟的 Vue UI 组件库，提供丰富的组件
- Vite 提供快速的开发体验和构建速度
- 用户明确要求使用 iView

**替代方案**：
- React + Ant Design：同样可行，但用户要求使用 iView
- 纯原生 JavaScript：开发效率低，不利于维护

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
- 适合管理用户认证状态、API Key 列表等全局状态

**状态结构**：
```typescript
{
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
  apiKeys: Array<{
    id: number,
    key_prefix: string,
    name: string,
    status: string,
    created_at: string,
    last_used_at: string
  }>
}
```

### 5. 路由设计

**决策**：使用 Vue Router，基于角色控制路由访问

**路由表**：
```typescript
{
  path: '/login',
  name: 'Login',
  component: LoginView
}
{
  path: '/',
  redirect: '/api-keys'
}
{
  path: '/api-keys',
  name: 'APIKeys',
  component: APIKeysView,
  meta: { requiresAuth: true }
}
```

### 6. UI 设计规范

**颜色方案**（中性色）：
- 主色：`#515a6e` (View UI 默认主色)
- 背景：`#f8f8f9` (浅灰背景)
- 卡片背景：`#ffffff` (白色)
- 边框：`#dcdee2` (浅灰边框)
- 文字：`#515a6e` (深灰)

**布局结构**：
- 左侧导航栏（宽度 ~200px）
- 顶部 Header（包含用户信息和注销按钮）
- 主内容区域（卡片式内容）

**组件风格**：
- 使用 View UI 的 Card 组件作为主要容器
- 表格使用 Table 组件展示 API Key 列表
- 表单使用 Form 组件
- 按钮使用 Button 组件

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
**风险**：前端打包后体积可能较大

**缓解措施**：
- 使用 Vite 的代码分割
- 按需加载 View UI 组件
- 生产环境启用 gzip 压缩

## Migration Plan

### 实施步骤

1. **初始化前端项目**
   ```bash
   cd frontend
   npm create vite@latest . -- --template vue-ts
   npm install view-ui-plus axios pinia vue-router
   ```

2. **开发页面组件**
   - 登录页面
   - 主布局（导航栏 + 内容区）
   - API Key 管理页面

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

1. **前端是否需要支持国际化？**
   - 当前方案不支持，后续可添加 vue-i18n

2. **是否需要主题切换功能？**
   - 当前方案只支持浅色主题

3. **API Key 创建后是否需要复制到剪贴板提示？**
   - 建议添加，提升用户体验
