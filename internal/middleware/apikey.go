package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var validAPIKeys map[string]bool

func init() {
	validAPIKeys = make(map[string]bool)

	// 从环境变量读取 API Key 白名单（逗号分隔）
	if keys := os.Getenv("API_KEYS"); keys != "" {
		keyList := strings.Split(keys, ",")
		for _, key := range keyList {
			if trimmed := strings.TrimSpace(key); trimmed != "" {
				validAPIKeys[trimmed] = true
			}
		}
	}

	// 如果未配置白名单，则允许所有 sk- 开头的 Key（开发模式）
	if len(validAPIKeys) == 0 {
		// 开发模式：任何 sk- 开头的 Key 都通过
	}
}

// APIKeyAuth API Key 鉴权中间件
func APIKeyAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Missing API key",
					"type":    "invalid_request_error",
				},
			})
			ctx.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Invalid authorization format. Use: Bearer <api_key>",
					"type":    "invalid_request_error",
				},
			})
			ctx.Abort()
			return
		}

		apiKey := parts[1]

		// 验证 API Key
		if !isValidAPIKey(apiKey) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Invalid API key",
					"type":    "invalid_request_error",
				},
			})
			ctx.Abort()
			return
		}

		// 将 API Key 存储在上下文中（脱敏后）
		maskedKey := maskAPIKey(apiKey)
		ctx.Set("api_key", apiKey)
		ctx.Set("api_key_masked", maskedKey)

		ctx.Next()
	}
}

// isValidAPIKey 验证 API Key 是否有效
func isValidAPIKey(key string) bool {
	// 如果配置了白名单，严格验证
	if len(validAPIKeys) > 0 {
		return validAPIKeys[key]
	}

	// 开发模式：任何 sk- 开头的 Key 都通过
	return strings.HasPrefix(key, "sk-")
}

// maskAPIKey 脱敏 API Key
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "sk-****"
	}
	return key[:7] + "..." + key[len(key)-4:]
}
