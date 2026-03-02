package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// UserController 用户管理控制器
type UserController struct {
	authService *service.AuthService
}

// NewUserController 创建 User Controller
func NewUserController(authService *service.AuthService) *UserController {
	return &UserController{
		authService: authService,
	}
}

// CreateUser 创建用户
// POST /v1/users
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	user, err := c.authService.CreateUser(ctx, &req)
	if err != nil {
		// 检查是否是邮箱重复错误
		if err.Error() == "email already exists" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"message": "Email already exists",
					"type":    "invalid_request_error",
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to create user",
				"type":    "api_error",
			},
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// GetUser 获取用户信息
// GET /v1/users/:id
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid user ID",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	user, err := c.authService.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"message": "User not found",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// CreateAPIKey 为用户创建 API Key
// POST /v1/users/:id/api-keys
func (c *UserController) CreateAPIKey(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid user ID",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	var req model.CreateAPIKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	response, err := c.authService.CreateAPIKey(ctx, userID, &req)
	if err != nil {
		// 检查错误类型
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": "User not found",
					"type":    "invalid_request_error",
				},
			})
			return
		}
		if err.Error() == "user is not active" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": "User is not active",
					"type":    "invalid_request_error",
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to create API key",
				"type":    "api_error",
			},
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// ListAPIKeys 获取用户的 API Key 列表
// GET /v1/users/:id/api-keys
func (c *UserController) ListAPIKeys(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid user ID",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	keys, err := c.authService.ListAPIKeys(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": "User not found",
					"type":    "invalid_request_error",
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to list API keys",
				"type":    "api_error",
			},
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
// DELETE /v1/users/:id/api-keys/:key_id
func (c *UserController) RevokeAPIKey(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid user ID",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	keyIDStr := ctx.Param("key_id")
	keyID, err := strconv.ParseInt(keyIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid API key ID",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	if err := c.authService.RevokeAPIKey(ctx, userID, keyID); err != nil {
		if err.Error() == "api key not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": "API key not found",
					"type":    "invalid_request_error",
				},
			})
			return
		}
		if err.Error() == "api key does not belong to user" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"message": "API key does not belong to user",
					"type":    "permission_error",
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to revoke API key",
				"type":    "api_error",
			},
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
