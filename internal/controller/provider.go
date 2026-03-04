package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/middleware"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// ProviderService Provider 服务接口（用于依赖注入和测试）
type ProviderService interface {
	CreateProvider(ctx context.Context, provider *model.Provider) error
	GetProvider(name string) (adapter.Provider, error)
	GetProviderByName(ctx context.Context, name string) (*model.Provider, error)
	ListProviders(ctx context.Context, enabledFilter *bool) ([]*service.ProviderInfo, error)
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
	Name           string         `json:"name" binding:"required"`
	Type           string         `json:"type" binding:"required"`
	BaseURL        string         `json:"base_url" binding:"required"`
	Timeout        int            `json:"timeout" binding:"required,min=1"`
	APIKey         string         `json:"api_key,omitempty"`
	ExtraConfig    map[string]any `json:"extra_config,omitempty"`
	Enabled        bool           `json:"enabled"`
	FallbackModels []interface{}  `json:"fallback_models,omitempty"`
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

// PublicProviderInfo 普通用户可见的 Provider 信息（不包含敏感信息）
type PublicProviderInfo struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	BaseURL        string   `json:"base_url"`
	Enabled        bool     `json:"enabled"`
	FallbackModels []string `json:"fallback_models"`
}

// ProviderModelsResponse Provider 模型列表响应
type ProviderModelsResponse struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Models []string `json:"models"`
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

	// 处理 FallbackModels
	if req.FallbackModels != nil {
		fallbackJSON := make(model.JSON)
		for _, v := range req.FallbackModels {
			if str, ok := v.(string); ok {
				fallbackJSON[str] = true
			}
		}
		provider.FallbackModels = fallbackJSON
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
	// 解析 enabled 查询参数
	var enabledFilter *bool
	if enabledStr := ctx.Query("enabled"); enabledStr != "" {
		if enabledStr == "true" {
			enabled := true
			enabledFilter = &enabled
		} else if enabledStr == "false" {
			enabled := false
			enabledFilter = &enabled
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid enabled parameter. Must be 'true' or 'false'",
				"type":    "invalid_request_error",
			})
			return
		}
	}

	providers, err := c.svc.ListProviders(context.Background(), enabledFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查用户角色
	role, _ := middleware.GetUserRole(ctx)

	// 管理员返回完整信息，普通用户返回简化信息
	if role == "admin" {
		ctx.JSON(http.StatusOK, gin.H{"providers": providers})
	} else {
		// 转换为普通用户可见的格式
		publicProviders := make([]PublicProviderInfo, 0, len(providers))
		for _, p := range providers {
			fallbackModels := make([]string, 0)
			if p.Provider.FallbackModels != nil {
				for model := range p.Provider.FallbackModels {
					fallbackModels = append(fallbackModels, model)
				}
			}
			publicProviders = append(publicProviders, PublicProviderInfo{
				Name:           p.Provider.Name,
				Type:           p.Provider.Type,
				BaseURL:        p.Provider.BaseURL,
				Enabled:        p.Provider.Enabled,
				FallbackModels: fallbackModels,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"providers": publicProviders})
	}
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

// ListProviderModels 获取 Provider 支持的模型列表
// GET /api/v1/providers/:name/models
func (c *ProviderController) ListProviderModels(ctx *gin.Context) {
	name := ctx.Param("name")

	provider, err := c.svc.GetProviderByName(context.Background(), name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	// 从 FallbackModels JSON 字段提取模型列表
	models := make([]string, 0)
	if provider.FallbackModels != nil {
		for model := range provider.FallbackModels {
			models = append(models, model)
		}
	}

	ctx.JSON(http.StatusOK, ProviderModelsResponse{
		Name:   provider.Name,
		Type:   provider.Type,
		Models: models,
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
