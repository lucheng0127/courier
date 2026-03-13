## 1. 数据库层实现

- [ ] 1.1 创建 `packages` 表迁移脚本
  - 包含字段：id, name, description, price, status, validity_days, created_at, updated_at
  - 添加唯一索引：idx_packages_name
  - 添加状态枚举约束：draft, online, offline

- [ ] 1.2 创建 `package_quotas` 表迁移脚本
  - 包含字段：id, package_id, provider_name, token_limit
  - 添加外键约束：package_id → packages.id
  - 添加复合索引：idx_package_quotas_package_provider

- [ ] 1.3 创建 `user_packages` 表迁移脚本
  - 包含字段：id, user_id, package_id, status, expires_at, activated_at, created_at
  - 添加外键约束：user_id → users.id, package_id → packages.id
  - 添加索引：idx_user_packages_user_status, idx_user_packages_expires_at

- [ ] 1.4 创建 `package_usage` 表迁移脚本
  - 包含字段：id, user_package_id, user_id, package_id, provider_name, tokens_used, request_id, usage_record_id, created_at
  - 添加外键约束：user_package_id → user_packages.id, usage_record_id → usage_records.id
  - 添加索引：idx_package_usage_user_package, idx_package_usage_user

- [ ] 1.5 扩展 `usage_records` 表
  - 添加可选字段：user_package_id (bigint), package_id (bigint)
  - 添加外键约束（可选）

- [ ] 1.6 运行数据库迁移

## 2. 模型层实现

- [ ] 2.1 创建 `internal/model/package.go`
  - 定义 Package 模型
  - 定义 PackageQuota 模型
  - 定义 PackageStatus 常量
  - 定义请求/响应 DTO（CreatePackageRequest, UpdatePackageRequest, PackageResponse）

- [ ] 2.2 创建 `internal/model/user_package.go`
  - 定义 UserPackage 模型
  - 定义 UserPackageStatus 常量
  - 定义请求/响应 DTO（PurchasePackageRequest, UserPackageResponse）

- [ ] 2.3 创建 `internal/model/package_usage.go`
  - 定义 PackageUsage 模型
  - 定义统计相关 DTO（PackageUsageStats, QuotaInfo）

- [ ] 2.4 扩展 `internal/model/usage.go`
  - 在 UsageRecord 中添加 UserPackageID 和 PackageID 字段

## 3. 数据访问层实现

- [ ] 3.1 创建 `internal/repository/package.go`
  - 实现 CreatePackage 方法
  - 实现 UpdatePackage 方法（仅草稿状态）
  - 实现 GetPackageByID 方法
  - 实现 ListPackages 方法（支持状态筛选）
  - 实现 DeletePackage 方法（仅草稿状态）
  - 实现 UpdatePackageStatus 方法

- [ ] 3.2 创建 `internal/repository/package_quota.go`
  - 实现 CreatePackageQuota 方法
  - 实现 GetPackageQuotas 方法
  - 实现 DeletePackageQuotas 方法（套餐删除时）

- [ ] 3.3 创建 `internal/repository/user_package.go`
  - 实现 CreateUserPackage 方法
  - 实现 GetUserPackageByID 方法
  - 实现 ListUserPackages 方法
  - 实现 GetUserActivePackages 方法（用于配额检测）
  - 实现 UpdateUserPackageStatus 方法（过期处理）
  - 实现 CheckPackageExists 方法（防止删除已购买套餐）

- [ ] 3.4 创建 `internal/repository/package_usage.go`
  - 实现 CreatePackageUsage 方法
  - 实现 GetPackageUsageStats 方法
  - 实现 GetUserTotalQuota 方法
  - 实现 GetPackageRemainingQuota 方法
  - 实现 ListPackageUsageRecords 方法

## 4. 服务层实现

- [ ] 4.1 创建 `internal/service/package.go`
  - 实现 PackageService 结构体
  - 实现创建套餐方法：CreatePackage
  - 实现更新套餐方法：UpdatePackage
  - 实现上架套餐方法：OnlinePackage
  - 实现下架套餐方法：OfflinePackage
  - 实现删除套餐方法：DeletePackage
  - 实现查询套餐方法：GetPackage, ListPackages
  - 实现套餐购买统计方法：GetPackagePurchaseStats

- [ ] 4.2 创建 `internal/service/user_package.go`
  - 实现 UserPackageService 结构体
  - 实现购买套餐方法：PurchasePackage（含模拟支付）
  - 实现查询我的套餐方法：ListMyPackages, GetMyPackage
  - 实现套餐使用统计方法：GetPackageUsageStats, ExportPackageUsage
  - 实现过期套餐处理方法：ProcessExpiredPackages

- [ ] 4.3 创建 `internal/service/quota.go`
  - 实现 QuotaService 结构体
  - 实现配额检测方法：CheckQuota
  - 实现配额扣减方法：DeductQuota
  - 实现剩余配额查询方法：GetRemainingQuota
  - 实现配额缓存操作方法

- [ ] 4.4 扩展 `internal/service/usage.go`
  - 在 RecordUsage 方法中添加套餐使用记录逻辑
  - 添加按套餐聚合统计方法

## 5. 控制器层实现

- [ ] 5.1 创建 `internal/controller/package_admin.go`
  - 实现 PackageAdminController 结构体
  - 实现创建套餐接口：POST /api/v1/admin/packages
  - 实现更新套餐接口：PUT /api/v1/admin/packages/:id
  - 实现上架套餐接口：POST /api/v1/admin/packages/:id/online
  - 实现下架套餐接口：POST /api/v1/admin/packages/:id/offline
  - 实现删除套餐接口：DELETE /api/v1/admin/packages/:id
  - 实现查询套餐接口：GET /api/v1/admin/packages, GET /api/v1/admin/packages/:id
  - 实现套餐购买统计接口：GET /api/v1/admin/packages/:id/purchases

- [ ] 5.2 创建 `internal/controller/package_public.go`
  - 实现 PackagePublicController 结构体
  - 实现查询可购买套餐接口：GET /api/v1/packages, GET /api/v1/packages/:id

- [ ] 5.3 创建 `internal/controller/user_package.go`
  - 实现 UserPackageController 结构体
  - 实现购买套餐接口：POST /api/v1/user/packages
  - 实现查询我的套餐接口：GET /api/v1/user/packages, GET /api/v1/user/packages/:id
  - 实现套餐使用统计接口：GET /api/v1/user/packages/:id/usage
  - 实现导出使用记录接口：GET /api/v1/user/packages/:id/usage/export
  - 实现查询配额接口：GET /api/v1/user/quota

- [ ] 5.4 创建 `internal/middleware/quota.go`
  - 实现 QuotaCheckMiddleware 中间件
  - 实现配额检测逻辑
  - 实现套餐筛选和排序逻辑
  - 实现配额缓存逻辑
  - 将可用套餐信息注入 Context

- [ ] 5.5 修改 `internal/controller/chat.go`
  - 在 ChatCompletions 方法中集成配额扣减逻辑
  - 请求完成后从 Context 获取套餐信息
  - 调用 QuotaService.DeductQuota 扣减配额

## 6. 路由注册

- [ ] 6.1 在 `cmd/server/main.go` 中注册新路由
  - 注册管理员套餐管理路由（使用 RequireAdmin 中间件）
  - 注册公开套餐查询路由（使用 Auth 中间件）
  - 注册用户套餐管理路由（使用 Auth 中间件）
  - 在 Chat 路由上添加配额检测中间件

## 7. 缓存实现

- [ ] 7.1 实现配额缓存逻辑
  - 使用 Redis 缓存用户配额信息
  - 实现 TTL 机制（5分钟）
  - 实现缓存更新和失效逻辑

## 8. 定时任务

- [ ] 8.1 实现过期套餐处理定时任务
  - 每小时执行一次
  - 将过期套餐状态更新为 expired
  - 清理相关缓存

## 9. 测试

- [ ] 9.1 编写单元测试
  - PackageService 测试
  - UserPackageService 测试
  - QuotaService 测试
  - 配额检测中间件测试

- [ ] 9.2 编写集成测试
  - 套餐创建、上架、购买流程测试
  - 配额检测和扣减流程测试
  - 过期套餐处理测试

- [ ] 9.3 编写性能测试
  - 配额检测中间件性能测试
  - 缓存命中率和性能对比测试

## 10. 文档和部署

- [ ] 10.1 更新 API 文档
  - 添加套餐管理 API 文档
  - 添加套餐购买 API 文档
  - 添加配额相关错误码说明

- [ ] 10.2 更新数据库迁移文档

- [ ] 10.3 准备部署配置
  - 添加环境变量配置（如缓存 TTL）
  - 更新 Docker Compose 配置（如需要）
