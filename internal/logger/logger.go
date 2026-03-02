package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel 日志级别
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Logger 结构化日志记录器
type Logger struct {
	mu     sync.Mutex
	out    io.Writer
	level  LogLevel
	fields map[string]any
}

// NewLogger 创建日志记录器
func NewLogger(out io.Writer, level LogLevel) *Logger {
	if out == nil {
		out = os.Stdout
	}
	return &Logger{
		out:    out,
		level:  level,
		fields: make(map[string]any),
	}
}

// With 添加字段
func (l *Logger) With(fields map[string]any) *Logger {
	newLogger := &Logger{
		out:    l.out,
		level:  l.level,
		fields: make(map[string]any),
	}
	// 复制现有字段
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	// 添加新字段
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

// log 内部日志方法
func (l *Logger) log(level LogLevel, msg string, fields map[string]any) {
	// 检查日志级别
	if !l.shouldLog(level) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// 构建日志条目
	entry := make(map[string]any)
	entry["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	entry["level"] = level
	entry["message"] = msg

	// 添加固定字段
	for k, v := range l.fields {
		entry[k] = v
	}

	// 添加临时字段
	for k, v := range fields {
		entry[k] = v
	}

	// JSON 序列化
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(l.out, `{"timestamp":"%s","level":"%s","message":"failed to marshal log: %v"}`+"\n",
			time.Now().UTC().Format(time.RFC3339), level, err)
		return
	}

	l.out.Write(data)
	l.out.Write([]byte("\n"))
}

// shouldLog 检查是否应该记录日志
func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
	}
	return levels[level] >= levels[l.level]
}

// Debug 记录 Debug 级别日志
func (l *Logger) Debug(msg string, fields map[string]any) {
	l.log(LevelDebug, msg, fields)
}

// Info 记录 Info 级别日志
func (l *Logger) Info(msg string, fields map[string]any) {
	l.log(LevelInfo, msg, fields)
}

// Warn 记录 Warn 级别日志
func (l *Logger) Warn(msg string, fields map[string]any) {
	l.log(LevelWarn, msg, fields)
}

// Error 记录 Error 级别日志
func (l *Logger) Error(msg string, fields map[string]any) {
	l.log(LevelError, msg, fields)
}

// 全局日志实例
var defaultLogger = NewLogger(os.Stdout, LevelInfo)

// SetDefaultLevel 设置默认日志级别
func SetDefaultLevel(level LogLevel) {
	defaultLogger.level = level
}

// Debug 全局 Debug 日志
func Debug(msg string, fields map[string]any) {
	defaultLogger.Debug(msg, fields)
}

// Info 全局 Info 日志
func Info(msg string, fields map[string]any) {
	defaultLogger.Info(msg, fields)
}

// Warn 全局 Warn 日志
func Warn(msg string, fields map[string]any) {
	defaultLogger.Warn(msg, fields)
}

// Error 全局 Error 日志
func Error(msg string, fields map[string]any) {
	defaultLogger.Error(msg, fields)
}

// With 全局 With 方法
func With(fields map[string]any) *Logger {
	return defaultLogger.With(fields)
}
