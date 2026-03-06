package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/service"
)

const (
	authTypeKey = "auth_type"
	jwtAuthType = "jwt"
	apiKeyAuthType = "apikey"
)

// DualAuth 双重认证中间件（JWT 或 API Key）
// 先尝试 JWT 认证，失败后尝试 API Key 认证
func DualAuth(authService *service.AuthService, jwtSvc service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		// 尝试 JWT 认证
		if tryJWTAuth(ctx, jwtSvc, authHeader) {
			ctx.Next()
			return
		}

		// JWT 失败，尝试 API Key 认证
		if tryAPIKeyAuth(ctx, authService, authHeader) {
			ctx.Next()
			return
		}

		// 两者都失败
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": "Authentication failed. Please provide a valid API key or JWT token.",
				"type":    "authentication_error",
			},
		})
		ctx.Abort()
	}
}

// tryJWTAuth 尝试 JWT 认证
func tryJWTAuth(ctx *gin.Context, jwtSvc service.JWTService, authHeader string) bool {
	if authHeader == "" {
		return false
	}

	// 验证 Bearer 格式
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 验证 Access Token
	claims, err := jwtSvc.ValidateAccessToken(token)
	if err != nil {
		return false
	}

	// 将用户信息注入到上下文
	ctx.Set(userIDKey, claims.UserID)
	ctx.Set(userEmailKey, claims.UserEmail)
	ctx.Set(userRoleKey, claims.UserRole)
	ctx.Set(authTypeKey, jwtAuthType)

	return true
}

// tryAPIKeyAuth 尝试 API Key 认证
func tryAPIKeyAuth(ctx *gin.Context, authService *service.AuthService, authHeader string) bool {
	if authHeader == "" {
		return false
	}

	// 验证 Bearer 格式
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	apiKey := strings.TrimPrefix(authHeader, "Bearer ")

	// 从数据库验证 API Key
	keyRecord, err := authService.ValidateAPIKey(ctx, apiKey)
	if err != nil {
		return false
	}

	// 获取关联用户信息
	user, err := authService.GetUserByID(ctx, keyRecord.UserID)
	if err != nil {
		return false
	}

	if user.Status != "active" {
		return false
	}

	// 注入到 Context
	ctx.Set(userIDKey, user.ID)
	ctx.Set(userEmailKey, user.Email)
	ctx.Set("api_key_id", keyRecord.ID)
	ctx.Set("api_key_masked", maskAPIKey(apiKey))
	ctx.Set(authTypeKey, apiKeyAuthType)

	// 异步更新 last_used_at
	go authService.UpdateKeyLastUsed(ctx, keyRecord.ID)

	return true
}

// GetAuthType 从上下文获取认证类型
func GetAuthType(c *gin.Context) (string, bool) {
	authType, exists := c.Get(authTypeKey)
	if !exists {
		return "", false
	}
	return authType.(string), true
}
