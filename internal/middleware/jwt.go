package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/service"
)

const (
	userIDKey    = "user_id"
	userEmailKey = "user_email"
	userRoleKey  = "user_role"
)

// JWTAuth JWT 鉴权中间件
func JWTAuth(jwtSvc service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization Header 提取 Bearer token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Missing authorization header",
				"type":    "authentication_error",
			})
			c.Abort()
			return
		}

		// 验证 Bearer 格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization header format",
				"type":    "authentication_error",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证 Access Token
		claims, err := jwtSvc.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired access token",
				"type":    "authentication_error",
			})
			c.Abort()
			return
		}

		// 将用户信息注入到上下文
		c.Set(userIDKey, claims.UserID)
		c.Set(userEmailKey, claims.UserEmail)
		c.Set(userRoleKey, claims.UserRole)

		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(userIDKey)
	if !exists {
		return 0, false
	}
	return userID.(int64), true
}

// GetUserEmail 从上下文获取用户邮箱
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(userEmailKey)
	if !exists {
		return "", false
	}
	return email.(string), true
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(userRoleKey)
	if !exists {
		return "", false
	}
	return role.(string), true
}
