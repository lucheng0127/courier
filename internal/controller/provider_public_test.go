package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
	"github.com/stretchr/testify/assert"
)

// MockProviderServiceForPublic 模拟 ProviderService 用于公开查询接口测试
type MockProviderServiceForPublic struct {
	providers []*model.Provider
}

func (m *MockProviderServiceForPublic) CreateProvider(ctx context.Context, provider *model.Provider) error {
	return nil
}

func (m *MockProviderServiceForPublic) GetProvider(name string) (adapter.Provider, error) {
	return &MockProvider{name: name, typ: "openai"}, nil
}

func (m *MockProviderServiceForPublic) GetProviderByName(ctx context.Context, name string) (*model.Provider, error) {
	for _, p := range m.providers {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, assert.AnError
}

func (m *MockProviderServiceForPublic) ListProviders(ctx context.Context, enabledFilter *bool) ([]*service.ProviderInfo, error) {
	var result []*service.ProviderInfo
	for _, p := range m.providers {
		if enabledFilter != nil && p.Enabled != *enabledFilter {
			continue
		}
		result = append(result, &service.ProviderInfo{
			Provider:  p,
			IsRunning: p.Enabled,
		})
	}
	return result, nil
}

func (m *MockProviderServiceForPublic) UpdateProvider(ctx context.Context, name string, updates map[string]any) (*model.Provider, error) {
	return nil, nil
}

func (m *MockProviderServiceForPublic) DeleteProvider(ctx context.Context, name string) error {
	return nil
}

func (m *MockProviderServiceForPublic) ReloadProvider(ctx context.Context, name string) error {
	return nil
}

func (m *MockProviderServiceForPublic) ReloadAllProviders(ctx context.Context) error {
	return nil
}

func (m *MockProviderServiceForPublic) EnableProvider(ctx context.Context, name string) error {
	return nil
}

func (m *MockProviderServiceForPublic) DisableProvider(ctx context.Context, name string) error {
	return nil
}

func (m *MockProviderServiceForPublic) InitProviders(ctx context.Context) error {
	return nil
}

func setupTestRouter(providerSvc ProviderService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	providerCtrl := NewProviderController(providerSvc)
	router.GET("/providers", providerCtrl.ListProviders)
	router.GET("/providers/:name/models", providerCtrl.ListProviderModels)

	return router
}

// setUserContext 设置用户上下文
func setUserContext(c *gin.Context, role string) {
	c.Set("user_id", int64(1))
	c.Set("user_email", "test@example.com")
	c.Set("user_role", role)
	c.Next()
}

// TestListProviders_Admin 测试管理员获取完整 Provider 列表
func TestListProviders_Admin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	fallbackModels := make(model.JSON)
	fallbackModels["gpt-4"] = true
	fallbackModels["gpt-3.5-turbo"] = true

	providers := []*model.Provider{
		{
			Name:           "openai-main",
			Type:           "openai",
			BaseURL:        "https://api.openai.com/v1",
			Timeout:        60,
			Enabled:        true,
			FallbackModels: fallbackModels,
		},
	}

	mockSvc := &MockProviderServiceForPublic{providers: providers}
	providerCtrl := NewProviderController(mockSvc)

	router.GET("/providers", func(c *gin.Context) {
		setUserContext(c, "admin")
		c.Next()
	}, providerCtrl.ListProviders)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/providers", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// 管理员应该能看到 is_running 字段
	assert.Contains(t, w.Body.String(), "is_running")
}

// TestListProviders_User 测试普通用户获取简化 Provider 列表
func TestListProviders_User(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	fallbackModels := make(model.JSON)
	fallbackModels["gpt-4"] = true
	fallbackModels["gpt-3.5-turbo"] = true

	providers := []*model.Provider{
		{
			Name:           "openai-main",
			Type:           "openai",
			BaseURL:        "https://api.openai.com/v1",
			Timeout:        60,
			Enabled:        true,
			FallbackModels: fallbackModels,
		},
	}

	mockSvc := &MockProviderServiceForPublic{providers: providers}
	providerCtrl := NewProviderController(mockSvc)

	// 设置普通用户角色
	router.GET("/providers", func(c *gin.Context) {
		setUserContext(c, "user")
	}, providerCtrl.ListProviders)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/providers", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// 不包含敏感信息
	assert.NotContains(t, w.Body.String(), "is_running")
	assert.NotContains(t, w.Body.String(), "timeout")
	assert.NotContains(t, w.Body.String(), "api_key")
	// 包含允许的字段
	assert.Contains(t, w.Body.String(), "name")
	assert.Contains(t, w.Body.String(), "type")
	assert.Contains(t, w.Body.String(), "base_url")
	assert.Contains(t, w.Body.String(), "enabled")
	assert.Contains(t, w.Body.String(), "fallback_models")
}

// TestListProviders_EnabledFilter 测试启用状态过滤
func TestListProviders_EnabledFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	fallbackModels := make(model.JSON)

	providers := []*model.Provider{
		{
			Name:           "openai-enabled",
			Type:           "openai",
			BaseURL:        "https://api.openai.com/v1",
			Timeout:        60,
			Enabled:        true,
			FallbackModels: fallbackModels,
		},
		{
			Name:           "openai-disabled",
			Type:           "openai",
			BaseURL:        "https://api.openai.com/v1",
			Timeout:        60,
			Enabled:        false,
			FallbackModels: fallbackModels,
		},
	}

	mockSvc := &MockProviderServiceForPublic{providers: providers}
	providerCtrl := NewProviderController(mockSvc)

	router.GET("/providers", func(c *gin.Context) {
		setUserContext(c, "user")
	}, providerCtrl.ListProviders)

	// 测试 enabled=true
	t.Run("enabled=true", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/providers?enabled=true", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "openai-enabled")
		assert.NotContains(t, w.Body.String(), "openai-disabled")
	})

	// 测试 enabled=false
	t.Run("enabled=false", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/providers?enabled=false", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotContains(t, w.Body.String(), "openai-enabled")
		assert.Contains(t, w.Body.String(), "openai-disabled")
	})

	// 测试无效的 enabled 参数
	t.Run("enabled=invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/providers?enabled=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid enabled parameter")
	})
}

// TestListProviderModels 测试获取 Provider 模型列表
func TestListProviderModels(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	fallbackModels := make(model.JSON)
	fallbackModels["gpt-4"] = true
	fallbackModels["gpt-3.5-turbo"] = true

	providers := []*model.Provider{
		{
			Name:           "openai-main",
			Type:           "openai",
			BaseURL:        "https://api.openai.com/v1",
			Timeout:        60,
			Enabled:        true,
			FallbackModels: fallbackModels,
		},
	}

	mockSvc := &MockProviderServiceForPublic{providers: providers}
	providerCtrl := NewProviderController(mockSvc)

	router.GET("/providers/:name/models", providerCtrl.ListProviderModels)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/providers/openai-main/models", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "openai-main")
	assert.Contains(t, w.Body.String(), "openai")
	assert.Contains(t, w.Body.String(), "gpt-4")
	assert.Contains(t, w.Body.String(), "gpt-3.5-turbo")
	assert.Contains(t, w.Body.String(), "models")
}

// TestListProviderModels_NotFound 测试获取不存在的 Provider 模型列表
func TestListProviderModels_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockSvc := &MockProviderServiceForPublic{providers: []*model.Provider{}}
	providerCtrl := NewProviderController(mockSvc)

	router.GET("/providers/:name/models", providerCtrl.ListProviderModels)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/providers/not-exist/models", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
