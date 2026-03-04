# Change: 添加 Web 管理界面

## Why

当前 Courier 项目仅提供 RESTful API，用户需要通过命令行或自定义客户端与系统交互，缺乏直观的可视化管理界面。添加 Web 管理界面可以：

1. 降低用户使用门槛，提供直观的操作界面
2. 简化用户登录、API Key 管理等日常操作
3. 为后续添加更多管理功能（Provider 管理、使用统计等）奠定基础
4. 提升项目的专业性和用户体验

## What Changes

- **新增 `frontend/` 目录**：存放前端项目代码
- **添加 Docker Compose 配置**：新增 nginx 服务，部署编译后的静态文件
- **实现用户认证功能**：登录页面、登录状态管理、自动 Token 刷新、注销功能
- **实现 API Key 管理页面**：查看 API Key 列表、创建新 Key、撤销 Key
- **界面风格**：简洁、模块化、卡片式组件、中性色配色
- **UI 框架**：使用 View UI (iView)

**阶段范围**：第一阶段仅实现用户登录、注销和 API Key 管理功能。

## Impact

- **新增规格**：`web-ui` - Web 用户界面能力
- **受影响代码**：
  - 新增 `frontend/` 目录（Vue 3 + View UI 项目）
  - 修改 `docker-compose.yml`（添加 nginx 服务）
  - 新增 `frontend/Dockerfile` 和 `frontend/nginx.conf`
- **依赖 API**：
  - `POST /api/v1/auth/login` - 登录
  - `POST /api/v1/auth/refresh` - 刷新 Token
  - `GET /api/v1/users/:id` - 获取用户信息
  - `GET /api/v1/users/:id/api-keys` - 获取 API Key 列表
  - `POST /api/v1/users/:id/api-keys` - 创建 API Key
  - `DELETE /api/v1/users/:id/api-keys/:key_id` - 撤销 API Key
