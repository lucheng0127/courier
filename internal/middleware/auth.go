package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/repository"
)

const (
	// UserIDKey 上下文中用户 ID 的 key
	UserIDKey = "user_id"
	// UserInfoKey 上下文中用户信息的 key
	UserInfoKey = "user_info"
)

// UserInfo 用户信息
type UserInfo struct {
	ID     uint
	Name   string
	Email  string
	UserID uint
}

// AuthMiddleware API Key 认证中间件
type AuthMiddleware struct {
	apiKeyRepo repository.APIKeyRepository
	userRepo   repository.UserRepository
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(apiKeyRepo repository.APIKeyRepository, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		apiKeyRepo: apiKeyRepo,
		userRepo:   userRepo,
	}
}

// RequireAuth 要求认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式无效"})
			c.Abort()
			return
		}

		apiKeyStr := parts[1]

		// 查询 API Key
		apiKey, err := m.apiKeyRepo.FindByKey(apiKeyStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key 无效"})
			c.Abort()
			return
		}

		// 检查 API Key 状态
		if apiKey.Status != "active" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key 已禁用"})
			c.Abort()
			return
		}

		// 更新 last_used_at
		now := time.Now()
		apiKey.LastUsedAt = &now
		if err := m.apiKeyRepo.UpdateLastUsedAt(apiKey.ID, now); err != nil {
			// 记录错误但不阻断请求
			// 可以考虑使用日志记录
		}

		// 获取用户信息
		user, err := m.userRepo.FindByID(apiKey.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(UserIDKey, user.ID)
		c.Set(UserInfoKey, &UserInfo{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			UserID: user.ID,
		})

		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUserInfo 从上下文获取用户信息
func GetUserInfo(c *gin.Context) (*UserInfo, bool) {
	userInfo, exists := c.Get(UserInfoKey)
	if !exists {
		return nil, false
	}
	info, ok := userInfo.(*UserInfo)
	return info, ok
}
