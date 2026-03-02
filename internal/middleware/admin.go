package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var adminAPIKey string

func init() {
	// 从环境变量读取管理员 API Key
	adminAPIKey = os.Getenv("ADMIN_API_KEY")
}

// AdminAuth 管理员鉴权中间件
func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-Admin-API-Key")

		// 如果没有配置管理员 Key，则跳过验证（开发模式）
		if adminAPIKey == "" {
			ctx.Next()
			return
		}

		// 验证管理员 Key
		if apiKey != adminAPIKey {
			// 也支持 Authorization: Bearer 格式
			auth := ctx.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if apiKey != adminAPIKey {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Unauthorized",
					"type":    "authentication_error",
				},
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
