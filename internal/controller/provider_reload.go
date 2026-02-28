package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/service"
)

// ProviderReloadController Provider 重载 API 控制器
type ProviderReloadController struct {
	svc *service.ProviderService
}

// NewProviderReloadController 创建 Provider Reload Controller
func NewProviderReloadController(svc *service.ProviderService) *ProviderReloadController {
	return &ProviderReloadController{svc: svc}
}

// ReloadAllProviders 重载所有 Provider
// POST /api/v1/admin/providers/reload
func (c *ProviderReloadController) ReloadAllProviders(ctx *gin.Context) {
	if err := c.svc.ReloadAllProviders(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to reload providers",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "all providers reloaded"})
}

// ReloadProvider 重载指定 Provider
// POST /api/v1/admin/providers/:name/reload
func (c *ProviderReloadController) ReloadProvider(ctx *gin.Context) {
	name := ctx.Param("name")

	if err := c.svc.ReloadProvider(ctx, name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to reload provider",
			"provider": name,
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "provider reloaded",
		"provider": name,
	})
}

// EnableProvider 启用 Provider
// POST /api/v1/admin/providers/:name/enable
func (c *ProviderReloadController) EnableProvider(ctx *gin.Context) {
	name := ctx.Param("name")

	if err := c.svc.EnableProvider(ctx, name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to enable provider",
			"provider": name,
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "provider enabled",
		"provider": name,
	})
}

// DisableProvider 禁用 Provider
// POST /api/v1/admin/providers/:name/disable
func (c *ProviderReloadController) DisableProvider(ctx *gin.Context) {
	name := ctx.Param("name")

	if err := c.svc.DisableProvider(ctx, name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to disable provider",
			"provider": name,
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "provider disabled",
		"provider": name,
	})
}

// RegisterRoutes 注册路由
func (c *ProviderReloadController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/admin/providers/reload", c.ReloadAllProviders)
	r.POST("/admin/providers/:name/reload", c.ReloadProvider)
	r.POST("/admin/providers/:name/enable", c.EnableProvider)
	r.POST("/admin/providers/:name/disable", c.DisableProvider)
}
