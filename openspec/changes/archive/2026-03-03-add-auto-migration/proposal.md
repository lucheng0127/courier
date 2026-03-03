# Change: 数据库自动迁移功能

## Why

当前项目存在以下问题：
1. 数据库迁移文件存在（`migrations/` 目录），但需要手动执行 SQL
2. `docker-compose.yml` 配置了 `migrate up` 命令，但 server 程序未实现此功能
3. 缺少启动时的数据库 schema 验证，可能导致数据库结构与代码不一致
4. 缺少调试日志，难以排查启动相关问题
5. 使用标准 `log` 包，日志结构化程度不足，不利于生产环境监控

用户在使用时遇到 `fallback_models` 列不存在的错误，就是因为迁移未执行。

## What Changes

1. **废弃 SQL 文件迁移方式**
   - 移除 `migrations/` 目录中的 `.sql` 文件
   - 不再依赖 SQL 脚本管理 schema 变更

2. **新增基于 ORM 的自动迁移模块** (`internal/migrate/`)
   - 使用 GORM 的 `AutoMigrate` 功能自动创建/更新表结构
   - 启动时自动检测 schema 变更并应用
   - 迁移失败时阻止服务启动

3. **新增 schema version 表** (`schema_migrations`)
   - 记录当前 schema 版本（使用 hash 或版本号）
   - 便于追踪 schema 变更历史

4. **集成到启动流程** (`cmd/server/main.go`)
   - 数据库连接后首先执行自动迁移
   - 支持通过环境变量控制是否启用自动迁移

5. **使用 zap 日志库**
   - 替换标准 `log` 包为 `uber-go/zap`
   - 支持结构化日志输出
   - 支持日志级别配置（debug/info/warn/error）

6. **添加调试日志**
   - 迁移开始/结束日志
   - 每个表的创建/更新状态日志
   - 数据库连接状态日志

## Impact

- **影响的功能**:
  - 服务启动流程
  - 数据库初始化
  - 全局日志系统

- **影响的代码**:
  - `cmd/server/main.go` - 添加迁移调用、zap 初始化
  - 新增 `internal/migrate/` 包
  - 新增 `internal/logger/` 包（zap 封装）
  - 移除 `migrations/` 目录下的 SQL 文件
  - 更新所有使用 `log.Printf` 的代码为 zap logger
  - `docker-compose.yml` - 移除无效的 `migrate up` 命令配置

- **依赖变更**:
  - 新增 `go.uber.org/zap`
  - 使用 GORM 的 AutoMigrate 功能

- **兼容性**:
  - 向后兼容：AutoMigrate 仅添加缺失的列/索引，不删除数据
  - 可通过环境变量 `AUTO_MIGRATE=false` 禁用自动迁移

## docker-compose.yml 变更

移除无效的 command 配置：

```yaml
# 之前（无效）
courier:
  command: ["/app/server", "migrate", "up"]

# 之后（正常运行）
courier:
  # 无需 command 配置，启动时自动执行迁移
```
