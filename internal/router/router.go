package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/handler"
)

// SetupRouter 设置路由
func SetupRouter(userHandler *handler.UserHandler, apiKeyHandler *handler.APIKeyHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		// 用户管理
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)

			// API Key 管理（作为 Users 的子资源）
			userID := users.Group("/:id/apikeys")
			{
				userID.POST("", apiKeyHandler.GenerateAPIKey)
				userID.GET("", apiKeyHandler.ListAPIKeys)
				userID.DELETE("/:keyid", apiKeyHandler.DeleteAPIKey)
				userID.PUT("/:keyid/disable", apiKeyHandler.DisableAPIKey)
			}
		}
	}

	return r
}
