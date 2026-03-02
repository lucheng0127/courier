package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// TraceIDHeader 是响应头中 TraceID 的字段名
	TraceIDHeader = "X-Trace-ID"
	// TraceIDKey 是 Context 中 TraceID 的键名
	TraceIDKey = "trace_id"
)

// TraceID 生成和透传中间件
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成 TraceID：格式为 trace-<UUID>
		traceID := "trace-" + uuid.New().String()

		// 存储到 Context
		c.Set(TraceIDKey, traceID)

		// 设置响应 Header
		c.Header(TraceIDHeader, traceID)

		c.Next()
	}
}

// GetTraceID 从 Context 获取 TraceID
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get(TraceIDKey); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}
