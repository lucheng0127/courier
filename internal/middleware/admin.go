package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const adminRole = "admin"

var adminAPIKey string

func init() {
	// 从环境变量读取管理员 API Key（用于过渡期兼容）
	adminAPIKey = os.Getenv("ADMIN_API_KEY")
}

// RequireAdmin 要求管理员角色的中间件
// 需要在 JWTAuth 中间件之后使用
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetUserRole(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authentication required",
				"type":    "authentication_error",
			})
			c.Abort()
			return
		}

		if role != adminRole {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Admin privileges required",
				"type":    "permission_error",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminAuth 管理员 API Key 鉴权中间件（用于过渡期兼容）
// 保留以支持从 Admin API Key 到 JWT 的平滑迁移
// 使用 AdminAPIKey 作为降级方案
func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 如果没有配置管理员 Key，则跳过验证（开发模式）
		if adminAPIKey == "" {
			ctx.Next()
			return
		}

		apiKey := ctx.GetHeader("X-Admin-API-Key")
		if apiKey == "" {
			// 也支持 Authorization: Bearer 格式
			auth := ctx.GetHeader("Authorization")
			if len(auth) > 7 && auth[:7] == "Bearer " {
				apiKey = auth[7:]
			}
		}

		if apiKey != adminAPIKey {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
				"type":    "authentication_error",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

