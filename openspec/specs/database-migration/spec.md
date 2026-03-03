# database-migration Specification

## Purpose
TBD - created by archiving change add-auto-migration. Update Purpose after archive.
## Requirements
### Requirement: 数据库自动迁移

系统 SHALL 在服务启动时使用 GORM AutoMigrate 自动同步数据库 schema，确保数据库结构与 Go struct 定义一致。

#### Scenario: 服务启动时自动执行迁移
- **WHEN** 服务启动且数据库连接成功
- **THEN** 系统自动调用 GORM AutoMigrate 同步所有表的 schema
- **AND** 所有迁移执行完成后服务才开始接受请求

#### Scenario: 迁移执行失败时阻止启动
- **WHEN** AutoMigrate 执行失败
- **THEN** 服务立即停止启动并返回错误信息
- **AND** 日志中记录失败的原因和相关表名

#### Scenario: Schema 增量更新
- **WHEN** Go struct 定义新增字段或修改类型
- **THEN** AutoMigrate 自动添加缺失的列或修改列类型
- **AND** 现有数据保持不变

### Requirement: Schema 版本追踪

系统 SHALL 维护 `schema_migrations` 表来追踪当前 schema 版本和定义 hash。

#### Scenario: 创建 schema_migrations 表
- **WHEN** 表不存在时首次执行迁移
- **THEN** 系统自动创建 `schema_migrations` 表
- **AND** 表结构包含 `version` (VARCHAR)、`hash` (VARCHAR) 和 `applied_at` (TIMESTAMP) 字段

#### Scenario: 记录 schema 版本
- **WHEN** AutoMigrate 成功执行
- **THEN** 系统在 `schema_migrations` 表中插入或更新版本记录
- **AND** `version` 字段记录当前版本号（如 "v1.0.0"）
- **AND** `hash` 字段记录 struct 定义的 hash 值
- **AND** `applied_at` 字段记录当前时间戳

#### Scenario: 检测 schema 变化
- **WHEN** 服务启动时检测到 hash 值与记录不同
- **THEN** 系统记录 warning 日志提示 schema 已变化
- **AND** 继续执行 AutoMigrate 更新 schema

### Requirement: 结构化日志系统

系统 SHALL 使用 uber-go/zap 库提供结构化日志输出。

#### Scenario: Logger 初始化
- **WHEN** 服务启动时
- **THEN** 系统初始化全局 zap logger 实例
- **AND** 根据 `LOG_LEVEL` 环境变量设置日志级别
- **AND** 开发环境使用 console 格式，生产环境使用 JSON 格式

#### Scenario: 日志级别配置
- **WHEN** 设置 `LOG_LEVEL=debug`
- **THEN** 输出所有级别的日志（Debug/Info/Warn/Error）
- **WHEN** 设置 `LOG_LEVEL=info`（默认）
- **THEN** 仅输出 Info/Warn/Error 级别的日志

#### Scenario: 结构化日志输出
- **WHEN** 记录日志时
- **THEN** 日志包含时间戳、级别、消息等结构化字段
- **AND** 支持添加自定义字段用于追踪（如 trace_id、user_id）

### Requirement: 迁移日志输出

系统 SHALL 使用 zap 输出详细的迁移执行日志。

#### Scenario: 迁移开始日志
- **WHEN** 开始执行迁移
- **THEN** 输出 "Starting database auto-migration..." 日志（Info 级别）

#### Scenario: 表同步日志
- **WHEN** AutoMigrate 处理每个表
- **THEN** 输出 "Syncing table: table_name" 日志（Debug 级别）
- **AND** 表结构变更时输出 "Table updated: table_name, changes: ..." 日志（Info 级别）

#### Scenario: 迁移完成日志
- **WHEN** 所有迁移执行完成
- **THEN** 输出 "Database auto-migration completed successfully" 日志（Info 级别）
- **AND** 包含当前 schema 版本和 hash

#### Scenario: Schema 无变化日志
- **WHEN** AutoMigrate 检测到 schema 无需变更
- **THEN** 输出 "Database schema is up to date" 日志（Info 级别）

### Requirement: 环境变量控制

系统 SHALL 支持通过环境变量控制自动迁移和日志行为。

#### Scenario: 禁用自动迁移
- **WHEN** 环境变量 `AUTO_MIGRATE=false`
- **THEN** 系统跳过 AutoMigrate 执行
- **AND** 输出 "Auto migration disabled by AUTO_MIGRATE=false" 日志（Warn 级别）

#### Scenario: 默认启用迁移
- **WHEN** 环境变量 `AUTO_MIGRATE` 未设置或为 `true`
- **THEN** 系统正常执行 AutoMigrate

#### Scenario: 日志级别配置
- **WHEN** 设置 `LOG_LEVEL` 环境变量
- **THEN** Logger 使用对应的日志级别
- **AND** 无效值时默认使用 info 级别

### Requirement: 全局 Logger 替换

系统 SHALL 使用 zap logger 替换所有使用标准 `log` 包的代码。

#### Scenario: main.go 中的日志
- **WHEN** `cmd/server/main.go` 需要输出日志
- **THEN** 使用全局 zap logger 实例
- **AND** 包含有意义的上下文字段

#### Scenario: controller 层日志
- **WHEN** controller 处理请求时需要记录日志
- **THEN** 使用 zap logger 记录请求信息
- **AND** 包含 trace_id、user_id 等追踪字段

#### Scenario: service 层日志
- **WHEN** service 层执行业务逻辑时需要记录日志
- **THEN** 使用 zap logger 记录关键操作和错误
- **AND** 错误日志包含完整的错误堆栈信息

