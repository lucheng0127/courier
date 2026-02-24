package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/handler"
	"github.com/lucheng0127/courier/internal/middleware"
)

// SetupRouter 设置路由
func SetupRouter(
	userHandler *handler.UserHandler,
	apiKeyHandler *handler.APIKeyHandler,
	modelHandler *handler.ModelHandler,
	chatHandler *handler.ChatHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
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

		// 模型管理（公开接口）
		models := v1.Group("/models")
		{
			models.GET("", modelHandler.ListModels)

			// 模型对话（需要认证）
			models.POST("/:model/chat", authMiddleware.RequireAuth(), chatHandler.Chat)
		}
	}

	return r
}
