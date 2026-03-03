package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
	"github.com/stretchr/testify/assert"
)

// MockProvider 模拟 adapter.Provider
type MockProvider struct {
	name string
	typ  string
}

func (m *MockProvider) Chat(ctx context.Context, req *adapter.ChatRequest) (*adapter.ChatResponse, error) {
	return nil, nil
}

func (m *MockProvider) ChatStream(ctx context.Context, req *adapter.ChatRequest) (<-chan *adapter.ChatStreamChunk, error) {
	return nil, nil
}

func (m *MockProvider) Type() string {
	return m.typ
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Timeout() int {
	return 60
}

func (m *MockProvider) Config() map[string]any {
	return nil
}

// MockProviderServiceForUpdate 模拟 ProviderService 用于 Update 测试
type MockProviderServiceForUpdate struct {
	updateFunc  func(ctx context.Context, name string, updates map[string]any) (*model.Provider, error)
	deleteFunc  func(ctx context.Context, name string) error
	getProvider func(name string) (adapter.Provider, error)
}

func (m *MockProviderServiceForUpdate) CreateProvider(ctx context.Context, provider *model.Provider) error {
	return nil
}

func (m *MockProviderServiceForUpdate) GetProvider(name string) (adapter.Provider, error) {
	if m.getProvider != nil {
		return m.getProvider(name)
	}
	return &MockProvider{name: name, typ: "openai"}, nil
}

func (m *MockProviderServiceForUpdate) ListProviders(ctx context.Context) ([]*service.ProviderInfo, error) {
	return nil, nil
}

func (m *MockProviderServiceForUpdate) ReloadProvider(ctx context.Context, name string) error {
	return nil
}

func (m *MockProviderServiceForUpdate) ReloadAllProviders(ctx context.Context) error {
	return nil
}

func (m *MockProviderServiceForUpdate) EnableProvider(ctx context.Context, name string) error {
	return nil
}

func (m *MockProviderServiceForUpdate) DisableProvider(ctx context.Context, name string) error {
	return nil
}

func (m *MockProviderServiceForUpdate) InitProviders(ctx context.Context) error {
	return nil
}

func (m *MockProviderServiceForUpdate) UpdateProvider(ctx context.Context, name string, updates map[string]any) (*model.Provider, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, name, updates)
	}
	return &model.Provider{Name: name}, nil
}

func (m *MockProviderServiceForUpdate) DeleteProvider(ctx context.Context, name string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, name)
	}
	return nil
}

// TestUpdateProvider_Success 测试成功更新 Provider
func TestUpdateProvider_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProviderServiceForUpdate{
		updateFunc: func(ctx context.Context, name string, updates map[string]any) (*model.Provider, error) {
			return &model.Provider{
				Name:    name,
				Type:    "openai",
				BaseURL: "https://api.openai.com/v1",
				Timeout: 120,
				Enabled: false,
			}, nil
		},
	}
	controller := NewProviderController(mockSvc)

	router := gin.New()
	router.PUT("/providers/:name", controller.UpdateProvider)

	body := `{"timeout": 120, "enabled": false}`
	req, _ := http.NewRequest("PUT", "/providers/test-provider", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Provider
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-provider", response.Name)
	assert.Equal(t, 120, response.Timeout)
	assert.False(t, response.Enabled)
}

// TestUpdateProvider_NotFound 测试更新不存在的 Provider
func TestUpdateProvider_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProviderServiceForUpdate{
		updateFunc: func(ctx context.Context, name string, updates map[string]any) (*model.Provider, error) {
		return nil, fmt.Errorf("provider not found: %s", name)
	},
	}
	controller := NewProviderController(mockSvc)

	router := gin.New()
	router.PUT("/providers/:name", controller.UpdateProvider)

	body := `{"timeout": 120}`
	req, _ := http.NewRequest("PUT", "/providers/non-existent", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestUpdateProvider_InvalidJSON 测试无效的 JSON
func TestUpdateProvider_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProviderServiceForUpdate{}
	controller := NewProviderController(mockSvc)

	router := gin.New()
	router.PUT("/providers/:name", controller.UpdateProvider)

	body := `{"timeout": "invalid"}`
	req, _ := http.NewRequest("PUT", "/providers/test-provider", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDeleteProvider_Success 测试成功删除 Provider
func TestDeleteProvider_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProviderServiceForUpdate{
		deleteFunc: func(ctx context.Context, name string) error {
			return nil
		},
	}
	controller := NewProviderController(mockSvc)

	router := gin.New()
	router.DELETE("/providers/:name", controller.DeleteProvider)

	req, _ := http.NewRequest("DELETE", "/providers/test-provider", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// TestDeleteProvider_NotFound 测试删除不存在的 Provider
func TestDeleteProvider_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProviderServiceForUpdate{
		deleteFunc: func(ctx context.Context, name string) error {
		return fmt.Errorf("provider not found: %s", name)
	},
	}
	controller := NewProviderController(mockSvc)

	router := gin.New()
	router.DELETE("/providers/:name", controller.DeleteProvider)

	req, _ := http.NewRequest("DELETE", "/providers/non-existent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
