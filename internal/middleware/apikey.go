package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/service"
)

// APIKeyAuth API Key 鉴权中间件
func APIKeyAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从 Authorization Header 提取 Bearer token
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

		// 从数据库验证 API Key
		keyRecord, err := authService.ValidateAPIKey(ctx, apiKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Invalid API key",
					"type":    "invalid_request_error",
				},
			})
			ctx.Abort()
			return
		}

		// 获取关联用户信息
		user, err := authService.GetUserByID(ctx, keyRecord.UserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "User not found",
					"type":    "invalid_request_error",
				},
			})
			ctx.Abort()
			return
		}

		if user.Status != "active" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"message": "User account is disabled",
					"type":    "permission_error",
				},
			})
			ctx.Abort()
			return
		}

		// 注入到 Context
		ctx.Set("user_id", user.ID)
		ctx.Set("user_email", user.Email)
		ctx.Set("api_key_id", keyRecord.ID)
		ctx.Set("api_key_masked", maskAPIKey(apiKey))

		// 异步更新 last_used_at
		go authService.UpdateKeyLastUsed(ctx, keyRecord.ID)

		ctx.Next()
	}
}

// maskAPIKey 脱敏 API Key
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "sk-****"
	}
	return key[:7] + "..." + key[len(key)-4:]
}
