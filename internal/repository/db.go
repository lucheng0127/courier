package repository

import (
	"github.com/lucheng0127/courier/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	err = db.AutoMigrate(&model.User{}, &model.APIKey{}, &model.RequestLog{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
