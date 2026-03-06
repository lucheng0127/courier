package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/middleware"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// UserController 用户管理控制器
type UserController struct {
	authSvc *service.AuthService
}

// NewUserController 创建 User Controller
func NewUserController(authSvc *service.AuthService) *UserController {
	return &UserController{
		authSvc: authSvc,
	}
}

// RegisterRoutes 注册路由
func (c *UserController) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		// 用户管理（仅管理员）
		users.GET("", c.ListUsers)
		users.PUT("/:id", c.UpdateUser)
		users.DELETE("/:id", c.DeleteUser)
		users.PATCH("/:id/status", c.UpdateUserStatus)

		// 获取用户信息（普通用户可获取自己的，管理员可获取任何人的）
		users.GET("/:id", c.GetUser)

		// API Key 管理（普通用户可管理自己的，管理员可管理任何人的）
		users.POST("/:id/api-keys", c.CreateAPIKey)
		users.GET("/:id/api-keys", c.ListAPIKeys)
		users.PATCH("/:id/api-keys/:key_id/enable", c.EnableAPIKey)
		users.PATCH("/:id/api-keys/:key_id/disable", c.DisableAPIKey)
		users.DELETE("/:id/api-keys/:key_id", c.DeleteAPIKey)
		users.DELETE("/:id/api-keys/:key_id/revoke", c.RevokeAPIKey)
	}
}

// GetUser 获取用户信息
// GET /api/v1/users/:id
// 权限：管理员可获取任何用户，普通用户只能获取自己
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"type":    "invalid_request_error",
		})
		return
	}

	// 权限检查：普通用户只能获取自己的信息
	userID, hasAuth := middleware.GetUserID(ctx)
	userRole, _ := middleware.GetUserRole(ctx)
	if hasAuth && userRole != "admin" && userID != targetID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
			"type":    "permission_error",
		})
		return
	}

	user, err := c.authSvc.GetUserByID(ctx, targetID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
			"type":    "invalid_request_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// CreateAPIKey 为用户创建 API Key
// POST /api/v1/users/:id/api-keys
// 权限：管理员可创建任何用户的，普通用户只能创建自己的
func (c *UserController) CreateAPIKey(ctx *gin.Context) {
	idStr := ctx.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"type":    "invalid_request_error",
		})
		return
	}

	// 权限检查：普通用户只能为自己创建 API Key
	userID, hasAuth := middleware.GetUserID(ctx)
	userRole, _ := middleware.GetUserRole(ctx)
	if hasAuth && userRole != "admin" && userID != targetID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
			"type":    "permission_error",
		})
		return
	}

	var req model.CreateAPIKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"type":    "invalid_request_error",
		})
		return
	}

	response, err := c.authSvc.CreateAPIKey(ctx, targetID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
				"type":    "invalid_request_error",
			})
			return
		}
		if err.Error() == "user is not active" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "User is not active",
				"type":    "invalid_request_error",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create API key",
			"type":    "api_error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// ListAPIKeys 获取用户的 API Key 列表
// GET /api/v1/users/:id/api-keys
// 权限：管理员可查看任何用户的，普通用户只能查看自己的
func (c *UserController) ListAPIKeys(ctx *gin.Context) {
	idStr := ctx.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"type":    "invalid_request_error",
		})
		return
	}

	// 权限检查：普通用户只能查看自己的 API Key
	userID, hasAuth := middleware.GetUserID(ctx)
	userRole, _ := middleware.GetUserRole(ctx)
	if hasAuth && userRole != "admin" && userID != targetID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
			"type":    "permission_error",
		})
		return
	}

	keys, err := c.authSvc.ListAPIKeys(ctx, targetID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
				"type":    "invalid_request_error",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to list API keys",
			"type":    "api_error",
		})
		return
	}

	// 转换为列表项格式（不包含完整 Key）
	items := make([]*model.APIKeyListItem, len(keys))
	for i, key := range keys {
		items[i] = &model.APIKeyListItem{
			ID:         key.ID,
			KeyPrefix:  key.KeyPrefix,
			Name:       key.Name,
			Status:     key.Status,
			LastUsedAt: key.LastUsedAt,
			ExpiresAt:  key.ExpiresAt,
			CreatedAt:  key.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"api_keys": items,
	})
}

// RevokeAPIKey 撤销 API Key
// DELETE /api/v1/users/:id/api-keys/:key_id
// 权限：管理员可撤销任何用户的，普通用户只能撤销自己的
func (c *UserController) RevokeAPIKey(ctx *gin.Context) {
	idStr := ctx.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"type":    "invalid_request_error",
		})
		return
	}

	// 权限检查：普通用户只能撤销自己的 API Key
	userID, hasAuth := middleware.GetUserID(ctx)
	userRole, _ := middleware.GetUserRole(ctx)
	if hasAuth && userRole != "admin" && userID != targetID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
			"type":    "permission_error",
		})
		return
	}

	keyIDStr := ctx.Param("key_id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid API key ID",
			"type":    "invalid_request_error",
		})
		return
	}

	if err := c.authSvc.RevokeAPIKey(ctx, targetID, keyID); err != nil {
		if err.Error() == "api key not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "API key not found",
				"type":    "invalid_request_error",
			})
			return
		}
		if err.Error() == "api key does not belong to user" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "API key does not belong to user",
				"type":    "permission_error",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to revoke API key",
			"type":    "api_error",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// EnableAPIKey 启用 API Key
// PATCH /api/v1/users/:id/api-keys/:key_id/enable
// 权限：管理员可启用任何用户的，普通用户只能启用自己的
func (c *UserController) EnableAPIKey(ctx *gin.Context) {
	idStr := ctx.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"type":    "invalid_request_error",
		})
		return
	}

	// 权限检查：普通用户只能启用自己的 API Key
	userID, hasAuth := middleware.GetUserID(ctx)
	userRole, _ := middleware.GetUserRole(ctx)
	if hasAuth && userRole != "admin" && userID != targetID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
			"type":    "permission_error",
		})
		return
	}

	keyIDStr := ctx.Param("key_id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid API key ID",
			"type":    "invalid_request_error",
		})
		return
	}

	key, err := c.authSvc.EnableAPIKey(ctx, targetID, keyID)
	if err != nil {
		if err.Error() == "api key not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "API key not found",
				"type":    "invalid_request_error",
			})
			return
		}
		if err.Error() == "api key does not belong to user" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "API key does not belong to user",
				"type":    "permission_error",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to enable API key",
			"type":    "api_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, key)
}

// DisableAPIKey 禁用 API Key
// PATCH /api/v1/users/:id/api-keys/:key_id/disable
// 权限：管理员可禁用任何用户的，普通用户只能禁用自己的
func (c *UserController) DisableAPIKey(ctx *gin.Context) {
	idStr := ctx.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"type":    "invalid_request_error",
		})
		return
	}

	// 权限检查：普通用户只能禁用自己的 API Key
	userID, hasAuth := middleware.GetUserID(ctx)
	userRole, _ := middleware.GetUserRole(ctx)
	if hasAuth && userRole != "admin" && userID != targetID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
			"type":    "permission_error",
		})
		return
	}

	keyIDStr := ctx.Param("key_id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid API key ID",
			"type":    "invalid_request_error",
		})
		return
	}

	key, err := c.authSvc.DisableAPIKey(ctx, targetID, keyID)
	if err != nil {
		if err.Error() == "api key not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "API key not found",
				"type":    "invalid_request_error",
			})
			return
		}
		if err.Error() == "api key does not belong to user" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "API key does not belong to user",
				"type":    "permission_error",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to disable API key",
			"type":    "api_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, key)
}

// DeleteAPIKey 删除 API Key（硬删除）
// DELETE /api/v1/users/:id/api-keys/:key_id
// 权限：管理员可删除任何用户的，普通用户只能删除自己的
func (c *UserController) DeleteAPIKey(ctx *gin.Context) {
	idStr := ctx.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"type":    "invalid_request_error",
		})
		return
	}

	// 权限检查：普通用户只能删除自己的 API Key
	userID, hasAuth := middleware.GetUserID(ctx)
	userRole, _ := middleware.GetUserRole(ctx)
	if hasAuth && userRole != "admin" && userID != targetID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
			"type":    "permission_error",
		})
		return
	}

	keyIDStr := ctx.Param("key_id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid API key ID",
			"type":    "invalid_request_error",
		})
		return
	}

	if err := c.authSvc.DeleteAPIKey(ctx, targetID, keyID); err != nil {
		if err.Error() == "api key not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "API key not found",
				"type":    "invalid_request_error",
			})
			return
		}
		if err.Error() == "api key does not belong to user" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "API key does not belong to user",
				"type":    "permission_error",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete API key",
			"type":    "api_error",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ListUsers 列出所有用户（仅管理员）
// GET /api/v1/users
func (c *UserController) ListUsers(ctx *gin.Context) {
	// 这个接口应该在 adminOnly 路由组中，所以不需要额外权限检查
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Not implemented yet",
	})
}

// UpdateUser 更新用户信息（仅管理员）
// PUT /api/v1/users/:id
func (c *UserController) UpdateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Not implemented yet",
	})
}

// DeleteUser 删除用户（仅管理员）
// DELETE /api/v1/users/:id
func (c *UserController) DeleteUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Not implemented yet",
	})
}

// UpdateUserStatus 更新用户状态（仅管理员）
// PATCH /api/v1/users/:id/status
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Not implemented yet",
	})
}
