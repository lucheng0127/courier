package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// ProviderController Provider 管理 API 控制器
type ProviderController struct {
	svc *service.ProviderService
}

// NewProviderController 创建 Provider Controller
func NewProviderController(svc *service.ProviderService) *ProviderController {
	return &ProviderController{svc: svc}
}

// CreateProviderRequest 创建 Provider 请求
type CreateProviderRequest struct {
	Name        string              `json:"name" binding:"required"`
	Type        string              `json:"type" binding:"required"`
	BaseURL     string              `json:"base_url" binding:"required"`
	Timeout     int                 `json:"timeout" binding:"required,min=1"`
	APIKey      string              `json:"api_key,omitempty"`
	ExtraConfig map[string]any      `json:"extra_config,omitempty"`
	Enabled     bool                `json:"enabled"`
}

// UpdateProviderRequest 更新 Provider 请求
type UpdateProviderRequest struct {
	Type        string              `json:"type,omitempty" binding:"required_with=BaseURL Timeout"`
	BaseURL     string              `json:"base_url,omitempty"`
	Timeout     int                 `json:"timeout,omitempty,min=1"`
	APIKey      string              `json:"api_key,omitempty"`
	ExtraConfig map[string]any      `json:"extra_config,omitempty"`
	Enabled     *bool               `json:"enabled,omitempty"`
}

// CreateProvider 创建 Provider
// POST /api/v1/providers
func (c *ProviderController) CreateProvider(ctx *gin.Context) {
	var req CreateProviderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider := &model.Provider{
		Name:        req.Name,
		Type:        req.Type,
		BaseURL:     req.BaseURL,
		Timeout:     req.Timeout,
		ExtraConfig: req.ExtraConfig,
		Enabled:     req.Enabled,
	}

	if req.APIKey != "" {
		provider.APIKey = &req.APIKey
	}

	if err := c.svc.CreateProvider(ctx, provider); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, provider)
}

// ListProviders 列出所有 Provider
// GET /api/v1/providers
func (c *ProviderController) ListProviders(ctx *gin.Context) {
	providers, err := c.svc.ListProviders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, providers)
}

// GetProvider 获取单个 Provider
// GET /api/v1/providers/:name
func (c *ProviderController) GetProvider(ctx *gin.Context) {
	name := ctx.Param("name")

	provider, err := c.svc.GetProvider(name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name": provider.Name(),
		"type": provider.Type(),
	})
}

// UpdateProvider 更新 Provider
// PUT /api/v1/providers/:name
func (c *ProviderController) UpdateProvider(ctx *gin.Context) {
	var req UpdateProviderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现更新逻辑
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DeleteProvider 删除 Provider
// DELETE /api/v1/providers/:name
func (c *ProviderController) DeleteProvider(ctx *gin.Context) {
	// TODO: 实现删除逻辑
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// RegisterRoutes 注册路由
func (c *ProviderController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/providers", c.CreateProvider)
	r.GET("/providers", c.ListProviders)
	r.GET("/providers/:name", c.GetProvider)
	r.PUT("/providers/:name", c.UpdateProvider)
	r.DELETE("/providers/:name", c.DeleteProvider)
}
