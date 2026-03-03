package migrate

import (
	"time"
)

// SchemaMigration 追踪 schema 版本
type SchemaMigration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"size:50;not null;index"` // 版本号，如 "v1.0.0"
	Hash      string    `gorm:"size:64;not null;index"` // struct 定义的 hash
	AppliedAt time.Time `gorm:"not null"`
}

// TableName 指定表名
func (SchemaMigration) TableName() string {
	return "schema_migrations"
}
