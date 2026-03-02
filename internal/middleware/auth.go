package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminAuth 管理员认证中间件
// TODO: MVP 阶段简单实现，后续可扩展为 JWT 或 OAuth
func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取 API Key
		apiKey := ctx.GetHeader("X-Admin-API-Key")
		if apiKey == "" {
			// 也支持 Authorization: Bearer <token> 格式
			auth := ctx.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		// TODO: 从环境变量或配置读取管理员 API Key
		expectedKey := ctx.GetHeader("X-Admin-API-Key-Expected")
		if expectedKey != "" && apiKey != expectedKey {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		// MVP 阶段：如果没有配置 API Key，则跳过认证
		// 生产环境必须配置
		ctx.Next()
	}
}
