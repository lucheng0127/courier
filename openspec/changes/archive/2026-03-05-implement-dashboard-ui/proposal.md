# Change: 实现 Dashboard 前端项目

## Why

当前 Courier 项目仅提供 RESTful API，缺乏可视化管理界面。为了提升用户体验和系统可用性，需要实现 Web Dashboard 前端项目。

第一阶段聚焦核心用户流程：
1. **用户认证流程**：注册、登录、注销
2. **Dashboard 概览**：展示系统状态和使用统计

## What Changes

- **新增 `frontend/` 目录**：存放前端项目代码（Vue 3 + Ant Design Vue + Vite）
- **添加 Docker 部署配置**：nginx 服务部署编译后的静态文件
- **实现用户注册功能**：Web 界面用户注册
- **实现用户登录/注销功能**：JWT Token 认证，自动刷新
- **实现 Dashboard 首页**：使用统计概览、Provider 状态、最近活动

**技术选型**：
- 前端框架：Vue 3 + TypeScript
- UI 组件库：Ant Design Vue
- 构建工具：Vite
- 状态管理：Pinia
- 路由：Vue Router
- HTTP 客户端：axios

**阶段范围**：第一阶段仅实现注册、登录、注销和 Dashboard 首页。

## Impact

- **新增规格**：`dashboard-ui` - Dashboard 用户界面能力
- **受影响代码**：
  - 新增 `frontend/` 目录（Vue 3 + Ant Design Vue 项目）
  - 修改 `docker-compose.yml`（添加 nginx 服务）
  - 新增 `frontend/Dockerfile` 和 `frontend/nginx.conf`
- **依赖 API**：
  - `POST /api/v1/auth/register` - 注册
  - `POST /api/v1/auth/login` - 登录
  - `POST /api/v1/auth/refresh` - 刷新 Token
  - `GET /api/v1/users/:id` - 获取用户信息
  - `GET /api/v1/usage` - 获取使用统计
  - `GET /api/v1/providers` - 获取 Provider 状态
