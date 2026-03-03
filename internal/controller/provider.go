package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// ProviderService Provider 服务接口（用于依赖注入和测试）
type ProviderService interface {
	CreateProvider(ctx context.Context, provider *model.Provider) error
	GetProvider(name string) (adapter.Provider, error)
	ListProviders(ctx context.Context) ([]*service.ProviderInfo, error)
	UpdateProvider(ctx context.Context, name string, updates map[string]any) (*model.Provider, error)
	DeleteProvider(ctx context.Context, name string) error
	ReloadProvider(ctx context.Context, name string) error
	ReloadAllProviders(ctx context.Context) error
	EnableProvider(ctx context.Context, name string) error
	DisableProvider(ctx context.Context, name string) error
	InitProviders(ctx context.Context) error
}

// ProviderController Provider 管理 API 控制器
type ProviderController struct {
	svc ProviderService
}

// NewProviderController 创建 Provider Controller
func NewProviderController(svc ProviderService) *ProviderController {
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
	Type            string         `json:"type,omitempty"`
	BaseURL         string         `json:"base_url,omitempty"`
	Timeout         int            `json:"timeout,omitempty"`
	APIKey          string         `json:"api_key,omitempty"`
	ExtraConfig     map[string]any `json:"extra_config,omitempty"`
	Enabled         *bool          `json:"enabled,omitempty"`
	FallbackModels  []interface{}  `json:"fallback_models,omitempty"`
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

	if err := c.svc.CreateProvider(context.Background(), provider); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, provider)
}

// ListProviders 列出所有 Provider
// GET /api/v1/providers
func (c *ProviderController) ListProviders(ctx *gin.Context) {
	providers, err := c.svc.ListProviders(context.Background())
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
	name := ctx.Param("name")

	var req UpdateProviderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建更新映射
	updates := make(map[string]any)
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.BaseURL != "" {
		updates["base_url"] = req.BaseURL
	}
	if req.Timeout > 0 {
		updates["timeout"] = req.Timeout
	}
	if req.APIKey != "" {
		updates["api_key"] = req.APIKey
	}
	if req.ExtraConfig != nil {
		updates["extra_config"] = req.ExtraConfig
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.FallbackModels != nil {
		updates["fallback_models"] = req.FallbackModels
	}

	// 调用 Service 层更新
	provider, err := c.svc.UpdateProvider(context.Background(), name, updates)
	if err != nil {
		// 检查是否是"not found"错误
		if containsString(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, provider)
}

// DeleteProvider 删除 Provider
// DELETE /api/v1/providers/:name
func (c *ProviderController) DeleteProvider(ctx *gin.Context) {
	name := ctx.Param("name")

	if err := c.svc.DeleteProvider(context.Background(), name); err != nil {
		// 检查是否是"not found"错误
		if containsString(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// containsString 检查字符串是否包含子串
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && indexOfString(s, substr) >= 0)
}

func indexOfString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// RegisterRoutes 注册路由
func (c *ProviderController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/providers", c.CreateProvider)
	r.GET("/providers", c.ListProviders)
	r.GET("/providers/:name", c.GetProvider)
	r.PUT("/providers/:name", c.UpdateProvider)
	r.DELETE("/providers/:name", c.DeleteProvider)
}
