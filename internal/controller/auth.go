package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// AuthController 认证控制器
type AuthController struct {
	authSvc *service.AuthService
}

// NewAuthController 创建 Auth Controller
func NewAuthController(authSvc *service.AuthService) *AuthController {
	return &AuthController{
		authSvc: authSvc,
	}
}

// RegisterRoutes 注册路由
func (c *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/refresh", c.RefreshToken)
	}
}

// Login 用户登录
// POST /api/v1/auth/login
func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"type":    "invalid_request_error",
		})
		return
	}

	resp, err := c.authSvc.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
			"type":    "authentication_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// RefreshToken 刷新 Token
// POST /api/v1/auth/refresh
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req model.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"type":    "invalid_request_error",
		})
		return
	}

	resp, err := c.authSvc.RefreshToken(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
			"type":    "authentication_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
