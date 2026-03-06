# Change: 添加 Provider 管理和 Model 列表功能

## Why

当前 Dashboard 已完成基础的登录、注册和首页功能，但缺少以下核心功能：

1. **管理员无法管理 Provider**：无法创建、编辑、删除、启用/禁用 Provider
2. **用户无法查看可用的模型列表**：普通用户需要知道有哪些 Provider 和模型可以调用

这些功能对于 LLM Gateway 的日常使用和管理至关重要。

## What Changes

- **新增 Provider 管理页面**（仅管理员）
  - Provider 列表展示
  - 创建新 Provider
  - 编辑现有 Provider
  - 删除 Provider
  - 启用/禁用 Provider
  - 重载 Provider
  - 查看 Provider 模型列表

- **新增 Model 列表页面**（所有用户）
  - 显示所有 Provider 的模型列表
  - 显示完整的模型标识（provider/model）
  - 按 Provider 分组展示
  - 模型状态显示

- **侧边栏导航**：添加导航菜单项

- **权限控制**：Provider 管理页面仅管理员可访问

## Impact

- **修改规格**：`dashboard-ui` - 新增 Provider 管理和 Model 列表需求
- **受影响代码**：
  - 新增 `frontend/src/views/ProvidersView.vue` - Provider 管理页面
  - 新增 `frontend/src/views/ModelsView.vue` - Model 列表页面
  - 修改 `frontend/src/router/index.ts` - 添加新路由
  - 修改 `frontend/src/views/DashboardView.vue` - 添加侧边栏导航
  - 新增 `frontend/src/stores/providers.ts` - Provider 状态管理
  - 新增 `frontend/src/api/providers.ts` - Provider API

- **依赖 API**：
  - `GET /api/v1/providers` - 获取 Provider 列表
  - `POST /api/v1/providers` - 创建 Provider
  - `PUT /api/v1/providers/:name` - 更新 Provider
  - `DELETE /api/v1/providers/:name` - 删除 Provider
  - `POST /api/v1/admin/providers/:name/reload` - 重载 Provider
  - `POST /api/v1/admin/providers/:name/enable` - 启用 Provider
  - `POST /api/v1/admin/providers/:name/disable` - 禁用 Provider
  - `GET /api/v1/providers/:name/models` - 获取 Provider 模型列表
