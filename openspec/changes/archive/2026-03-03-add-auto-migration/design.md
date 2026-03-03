## Context

Courier 项目使用 PostgreSQL 作为主数据存储。当前问题：
- 使用 SQL 文件管理迁移，需要手动执行
- 缺少自动化的 schema 同步机制
- 使用标准 `log` 包，日志结构化程度不足
- 缺少生产环境友好的日志系统

## Goals / Non-Goals

**Goals**:
- 服务启动时自动同步数据库 schema
- 使用 GORM AutoMigrate 实现自动化
- 提供结构化日志输出（使用 zap）
- Schema 变更可通过代码追踪

**Non-Goals**:
- 不支持迁移回滚 - GORM AutoMigrate 不支持删除列/表
- 不支持多数据库锁机制 - 单实例部署
- 不保留 SQL 迁移文件（完全迁移到代码定义）

## Decisions

### 1. 迁移方式：ORM AutoMigrate vs SQL 文件
**决策**: 使用 GORM 的 `AutoMigrate` 功能

**原因**:
- **类型安全**: Schema 定义与 Go struct 绑定，避免不一致
- **自动化**: 启动时自动同步，无需手动执行 SQL
- **简单直接**: 无需维护额外的 SQL 文件
- **已有依赖**: 项目已使用 GORM（sqlx 可逐步替换或共存）

**实现方式**:
```go
db.AutoMigrate(&model.Provider{}, &model.User{}, &model.APIKey{}, &model.UsageRecord{})
```

**替代方案对比**:
| 方案 | 优点 | 缺点 |
|------|------|------|
| SQL 文件 | 精确控制、支持复杂 DDL | 需手动执行、易遗漏 |
| golang-migrate | 工具成熟、支持版本管理 | 需额外工具、SQL 维护 |
| **GORM AutoMigrate** | 自动化、类型安全 | 不支持删除列/表 |

### 2. 日志库选择：zap vs log vs logrus
**决策**: 使用 `uber-go/zap`

**原因**:
- **高性能**: 零分配日志，适合生产环境
- **结构化**: JSON 输出，便于日志采集和分析
- **分级**: 支持 Debug/Info/Warn/Error/Fatal
- **上下文**: 支持添加字段，便于追踪

**替代方案**:
- `log`: 标准库，简单但无结构化
- `logrus`: 流行但性能低于 zap

**配置**:
```go
// 开发环境：console 格式，debug 级别
// 生产环境：JSON 格式，info 级别
logger, _ := zap.NewProduction()
```

### 3. Schema 版本追踪
**决策**: 使用 schema hash 而非版本号

```go
type SchemaMigration struct {
    Version   string    // 版本号（如 "v1.2.0"）
    Hash      string    // struct 定义 hash
    AppliedAt time.Time // 应用时间
}
```

**原因**:
- 版本号便于人类阅读
- Hash 可检测 schema 定义是否变化
- 便于排查 schema 不一致问题

### 4. 迁移执行时机
**决策**: 在 `main.go` 中数据库连接建立后立即执行

```
初始化 Logger → 数据库连接 → 执行 AutoMigrate → 初始化 repositories → 启动服务
```

**原因**:
- 迁移在所有数据库操作前完成
- 失败时立即停止启动，避免不一致状态
- repositories 初始化时 schema 已就绪

### 5. Logger 初始化
**决策**: 在 `main.go` 最开始初始化，作为全局单例

```go
// internal/logger/logger.go
var L *zap.Logger

func Init(level string) {
    L = zap.NewProduction()
}
```

**原因**:
- 所有模块共享同一 logger 实例
- 统一配置（级别、格式）
- 便于测试时可替换

### 6. 环境变量控制
**决策**: 支持 `AUTO_MIGRATE` 和 `LOG_LEVEL` 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `AUTO_MIGRATE` | `true` | 是否启用自动迁移 |
| `LOG_LEVEL` | `info` | 日志级别（debug/info/warn/error） |

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| GORM AutoMigrate 不支持删除列/表 | 需要删除时手动执行 SQL 或使用原始迁移 |
| 并发迁移导致竞争 | 单实例部署，无此风险；未来可添加 advisory lock |
| Schema 定义与 struct 不同步 | 通过代码 review 和测试保证 |
| 大量数据时迁移慢 | AutoMigrate 仅修改结构，不涉及数据迁移 |

## Migration Plan

### 部署步骤
1. 更新代码（添加 struct 定义、AutoMigrate 调用）
2. 添加 zap 依赖 (`go get go.uber.org/zap`)
3. 服务启动
4. AutoMigrate 自动同步 schema
5. 迁移完成后服务正常启动

### 回滚计划
- 代码回滚后，旧版本 struct 定义可能不匹配
- AutoMigrate 仅添加缺失的列，不会破坏现有数据
- 如需删除列：手动执行 SQL `ALTER TABLE ... DROP COLUMN ...`

### 现有数据库处理
- 首次启动时创建 `schema_migrations` 表
- 计算当前 struct 定义的 hash
- 记录版本和 hash，便于后续检测变化

## Open Questions

无 - MVP 范围明确
