package repository

import (
	"github.com/lucheng0127/courier/internal/model"
	"gorm.io/gorm"
)

// RequestLogRepository 请求日志仓储接口
type RequestLogRepository interface {
	Create(log *model.RequestLog) error
}

// requestLogRepository 请求日志仓储实现
type requestLogRepository struct {
	db *gorm.DB
}

// NewRequestLogRepository 创建请求日志仓储
func NewRequestLogRepository(db *gorm.DB) RequestLogRepository {
	return &requestLogRepository{db: db}
}

func (r *requestLogRepository) Create(log *model.RequestLog) error {
	return r.db.Create(log).Error
}
