package logger

import (
	"go.uber.org/zap"
)

// NewLogger 创建新的日志器
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
