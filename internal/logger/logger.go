package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

// Init 初始化全局 logger
// level: debug, info, warn, error
// env: development, production
func Init(level string, env string) {
	// 解析日志级别
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 配置 encoder
	var config zap.Config
	if strings.ToLower(env) == "development" {
		// 开发环境：console 格式，debug 级别
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level = zap.NewAtomicLevelAt(zapLevel)
	} else {
		// 生产环境：JSON 格式
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapLevel)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	L = logger
}

// InitFromEnv 从环境变量初始化 logger
// LOG_LEVEL: 日志级别（默认 info）
// ENV: 环境类型（默认 production）
func InitFromEnv() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "production"
	}

	Init(level, env)
}

// Sync 刷新缓冲区
func Sync() {
	if L != nil {
		_ = L.Sync()
	}
}
