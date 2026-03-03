package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/lucheng0127/courier/internal/logger"
	"github.com/lucheng0127/courier/internal/model"
)

// Migrator 数据库迁移器
type Migrator struct {
	db              *gorm.DB
	currentVersion  string
	schemaHash      string
	models          []interface{} // 需要迁移的 models
}

// NewMigrator 创建迁移器
func NewMigrator(db *gorm.DB, version string) *Migrator {
	return &Migrator{
		db:             db,
		currentVersion: version,
		models:         []interface{}{},
	}
}

// RegisterModels 注册需要迁移的 model
func (m *Migrator) RegisterModels(models ...interface{}) {
	m.models = append(m.models, models...)
}

// Run 执行自动迁移
func (m *Migrator) Run() error {
	logger.L.Info("Starting database auto-migration...")

	// 1. 确保 schema_migrations 表存在
	if err := m.ensureSchemaTable(); err != nil {
		logger.L.Error("Failed to create schema_migrations table",
			zap.String("error", err.Error()))
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// 2. 计算当前 schema hash
	m.schemaHash = m.calculateSchemaHash()

	// 3. 检查 schema 变化
	lastHash, err := m.getLastHash()
	if err != nil {
		logger.L.Error("Failed to get last schema hash",
			zap.String("error", err.Error()))
		return fmt.Errorf("failed to get last schema hash: %w", err)
	}

	if lastHash != "" && lastHash != m.schemaHash {
		// Schema 发生变化，记录警告但继续执行
		logger.L.Warn("Schema has changed",
			zap.String("old_hash", lastHash),
			zap.String("new_hash", m.schemaHash))
	}

	// 4. 执行 AutoMigrate
	if err := m.autoMigrate(); err != nil {
		logger.L.Error("Failed to auto migrate",
			zap.String("error", err.Error()))
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	logger.L.Info("Database auto-migration completed successfully",
		zap.String("version", m.currentVersion),
		zap.String("hash", m.schemaHash))

	// 5. 记录版本
	if err := m.recordVersion(); err != nil {
		logger.L.Error("Failed to record version",
			zap.String("error", err.Error()))
		return fmt.Errorf("failed to record version: %w", err)
	}

	return nil
}

// ensureSchemaTable 确保 schema_migrations 表存在
func (m *Migrator) ensureSchemaTable() error {
	return m.db.AutoMigrate(&SchemaMigration{})
}

// calculateSchemaHash 计算 schema 定义 hash
// MVP 阶段简化：基于版本号生成 hash
// 生产环境应基于 struct 定义生成 hash
func (m *Migrator) calculateSchemaHash() string {
	data := m.currentVersion + ":courier-schema"
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// getLastHash 获取上次记录的 schema hash
func (m *Migrator) getLastHash() (string, error) {
	var migration SchemaMigration
	result := m.db.Order("applied_at DESC").First(&migration)
	if result.Error != nil {
		if result.Error == gormlogger.ErrRecordNotFound {
			return "", nil // 表为空，首次运行
		}
		return "", result.Error
	}
	return migration.Hash, nil
}

// autoMigrate 执行 GORM AutoMigrate
func (m *Migrator) autoMigrate() error {
	// 禁用 GORM 默认 logger，使用我们自己的
	m.db.Logger = gormlogger.Default.LogMode(gormlogger.Silent)

	// 基础表
	models := []interface{}{
		&SchemaMigration{},
		&model.Provider{},
		&model.User{},
		&model.APIKey{},
		&model.UsageRecord{},
	}

	// 添加注册的额外 models
	models = append(models, m.models...)

	// 记录同步的表
	for _, m := range models {
		logger.L.Debug("Syncing table",
			zap.String("table", getTableName(m)))
	}

	return m.db.AutoMigrate(models...)
}

// recordVersion 记录当前版本
func (m *Migrator) recordVersion() error {
	// 检查是否已存在当前版本
	var existing SchemaMigration
	result := m.db.Where("version = ?", m.currentVersion).First(&existing)

	if result.Error != nil {
		if result.Error == gormlogger.ErrRecordNotFound {
			// 不存在，创建新记录
			migration := SchemaMigration{
				Version:   m.currentVersion,
				Hash:      m.schemaHash,
				AppliedAt: m.db.NowFunc(),
			}
			return m.db.Create(&migration).Error
		}
		return result.Error
	}

	// 已存在，更新 hash 和时间
	existing.Hash = m.schemaHash
	existing.AppliedAt = m.db.NowFunc()
	return m.db.Save(&existing).Error
}

// getTableName 获取 model 对应的表名
func getTableName(model interface{}) string {
	if tn, ok := model.(interface{ TableName() string }); ok {
		return tn.TableName()
	}
	return "unknown"
}
