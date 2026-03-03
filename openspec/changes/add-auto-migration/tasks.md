## 1. Logger 模块

- [ ] 1.1 添加 `go.uber.org/zap` 依赖到 `go.mod`
- [ ] 1.2 创建 `internal/logger/logger.go` 实现全局 logger
- [ ] 1.3 实现 `Init(level string, env string)` 函数初始化 logger
- [ ] 1.4 实现环境变量日志级别解析（debug/info/warn/error）
- [ ] 1.5 支持开发环境 console 格式和生产环境 JSON 格式
- [ ] 1.6 导出全局 `L *zap.Logger` 变量供其他模块使用

## 2. 数据库迁移模块

- [ ] 2.1 创建 `internal/migrate/migrator.go` 实现 `Migrator` 结构体
- [ ] 2.2 创建 `internal/migrate/schema_migration.go` 定义 `SchemaMigration` struct
- [ ] 2.3 实现 `ensureSchemaTable()` 方法创建 `schema_migrations` 表（使用 GORM AutoMigrate）
- [ ] 2.4 实现 `getCurrentHash()` 方法计算 struct 定义 hash
- [ ] 2.5 实现 `recordVersion()` 方法记录版本和 hash
- [ ] 2.6 实现 `checkSchemaChange()` 方法检测 schema 变化
- [ ] 2.7 实现 `Run()` 方法调用 GORM AutoMigrate
- [ ] 2.8 添加所有需要迁移的 model 到 AutoMigrate 调用中

## 3. Model 定义完善

- [ ] 3.1 确保 `internal/model/provider.go` 的 struct 定义完整
- [ ] 3.2 确保 `internal/model/user.go` 的 struct 定义完整
- [ ] 3.3 确保 `internal/model/api_key.go` 的 struct 定义完整
- [ ] 3.4 确保 `internal/model/usage_record.go` 的 struct 定义完整
- [ ] 3.5 为所有 model 添加 GORM 标签（table、column、index 等）

## 4. 启动流程集成

- [ ] 4.1 在 `cmd/server/main.go` 开头调用 `logger.Init()`
- [ ] 4.2 添加 `AUTO_MIGRATE` 和 `LOG_LEVEL` 环境变量读取
- [ ] 4.3 在数据库连接后调用 `migrator.Run()`
- [ ] 4.4 确保迁移失败时使用 `log.Fatal()` 停止启动
- [ ] 4.5 更新所有 `log.Printf` 为 `logger.L.Info()/Error()` 等

## 5. 全局日志替换

- [ ] 5.1 替换 `cmd/server/main.go` 中的所有 `log.Printf`
- [ ] 5.2 替换 `internal/controller/` 中所有 `log.Printf`
- [ ] 5.3 替换 `internal/service/` 中所有 `log.Printf`
- [ ] 5.4 替换 `internal/middleware/` 中所有 `log.Printf`
- [ ] 5.5 替换 `internal/bootstrap/` 中所有 `log.Printf`
- [ ] 5.6 确保日志包含有意义的上下文字段

## 6. 迁移日志输出

- [ ] 6.1 添加迁移开始日志 "Starting database auto-migration..."
- [ ] 6.2 添加表同步日志 "Syncing table: table_name"（Debug 级别）
- [ ] 6.3 添加 schema 变化检测日志
- [ ] 6.4 添加迁移完成日志 "Database auto-migration completed successfully"
- [ ] 6.5 添加 schema 无变化日志 "Database schema is up to date"

## 7. 清理工作

- [ ] 7.1 删除 `migrations/` 目录下所有 `.sql` 文件
- [ ] 7.2 更新 `docker-compose.yml` 移除无效的 migrate 命令配置
- [ ] 7.3 更新项目文档说明新的迁移方式

## 8. 测试验证

- [ ] 8.1 测试空数据库首次启动（创建所有表）
- [ ] 8.2 测试现有数据库启动（AutoMigrate 仅添加缺失列）
- [ ] 8.3 测试新增 struct 字段后启动（自动添加新列）
- [ ] 8.4 测试 `AUTO_MIGRATE=false` 环境变量禁用迁移
- [ ] 8.5 测试 `LOG_LEVEL=debug` 输出详细日志
- [ ] 8.6 测试 `LOG_LEVEL=info` 仅输出基础日志
- [ ] 8.7 验证结构化日志 JSON 格式输出（生产环境）
- [ ] 8.8 测试迁移失败场景（数据库连接失败等）
