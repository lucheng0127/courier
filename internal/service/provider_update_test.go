package service

import (
	"context"
	"testing"

	"github.com/lucheng0127/courier/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProviderRepository 模拟 Provider Repository
type MockProviderRepository struct {
	mock.Mock
}

func (m *MockProviderRepository) Create(ctx context.Context, provider *model.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockProviderRepository) GetByID(ctx context.Context, id int64) (*model.Provider, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Provider), args.Error(1)
}

func (m *MockProviderRepository) GetByName(ctx context.Context, name string) (*model.Provider, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Provider), args.Error(1)
}

func (m *MockProviderRepository) List(ctx context.Context) ([]*model.Provider, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Provider), args.Error(1)
}

func (m *MockProviderRepository) Update(ctx context.Context, provider *model.Provider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockProviderRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProviderRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

// TestUpdateProvider_Success 测试成功更新 Provider
func TestUpdateProvider_Success(t *testing.T) {
	mockRepo := new(MockProviderRepository)
	svc := NewProviderService(mockRepo)
	ctx := context.Background()

	// 模拟 Provider
	existingProvider := &model.Provider{
		ID:      1,
		Name:    "test-provider",
		Type:    "openai",
		BaseURL: "https://api.openai.com/v1",
		Timeout: 60,
		Enabled: true,
	}

	// 设置期望
	mockRepo.On("GetByName", ctx, "test-provider").Return(existingProvider, nil)
	mockRepo.On("Update", ctx, mock.Anything).Return(nil)

	// 执行更新
	updates := map[string]any{
		"timeout": 120,
		"enabled": false,
	}

	result, err := svc.UpdateProvider(ctx, "test-provider", updates)

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 120, result.Timeout)
	assert.False(t, result.Enabled)

	mockRepo.AssertExpectations(t)
}

// TestUpdateProvider_NotFound 测试更新不存在的 Provider
func TestUpdateProvider_NotFound(t *testing.T) {
	mockRepo := new(MockProviderRepository)
	svc := NewProviderService(mockRepo)
	ctx := context.Background()

	// 设置期望 - Provider 不存在
	mockRepo.On("GetByName", ctx, "non-existent").Return(nil, assert.AnError)

	// 执行更新
	updates := map[string]any{
		"timeout": 120,
	}

	_, err := svc.UpdateProvider(ctx, "non-existent", updates)

	// 验证
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}

// TestUpdateProvider_DisableToEnable 测试从禁用变为启用
func TestUpdateProvider_DisableToEnable(t *testing.T) {
	mockRepo := new(MockProviderRepository)
	svc := NewProviderService(mockRepo)
	ctx := context.Background()

	// 模拟已禁用的 Provider
	existingProvider := &model.Provider{
		ID:      1,
		Name:    "test-provider",
		Type:    "openai",
		BaseURL: "https://api.openai.com/v1",
		Timeout: 60,
		Enabled: false,
	}

	mockRepo.On("GetByName", ctx, "test-provider").Return(existingProvider, nil)
	mockRepo.On("Update", ctx, mock.Anything).Return(nil)

	// 执行更新 - 启用 Provider
	updates := map[string]any{
		"enabled": true,
	}

	_, err := svc.UpdateProvider(ctx, "test-provider", updates)

	// 验证 - 由于我们没有真实的 adapter registry，这里主要验证不报错
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestUpdateProvider_EnableToDisable 测试从启用变为禁用
func TestUpdateProvider_EnableToDisable(t *testing.T) {
	mockRepo := new(MockProviderRepository)
	svc := NewProviderService(mockRepo)
	ctx := context.Background()

	// 模拟已启用的 Provider
	existingProvider := &model.Provider{
		ID:      1,
		Name:    "test-provider",
		Type:    "openai",
		BaseURL: "https://api.openai.com/v1",
		Timeout: 60,
		Enabled: true,
	}

	mockRepo.On("GetByName", ctx, "test-provider").Return(existingProvider, nil)
	mockRepo.On("Update", ctx, mock.Anything).Return(nil)

	// 执行更新 - 禁用 Provider
	updates := map[string]any{
		"enabled": false,
	}

	_, err := svc.UpdateProvider(ctx, "test-provider", updates)

	// 验证
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestDeleteProvider_Success 测试成功删除 Provider
func TestDeleteProvider_Success(t *testing.T) {
	mockRepo := new(MockProviderRepository)
	svc := NewProviderService(mockRepo)
	ctx := context.Background()

	// 模拟 Provider
	existingProvider := &model.Provider{
		ID:      1,
		Name:    "test-provider",
		Type:    "openai",
		BaseURL: "https://api.openai.com/v1",
		Timeout: 60,
		Enabled: false, // 未启用
	}

	mockRepo.On("GetByName", ctx, "test-provider").Return(existingProvider, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(nil)

	// 执行删除
	err := svc.DeleteProvider(ctx, "test-provider")

	// 验证
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestDeleteProvider_EnabledProvider 测试删除已启用的 Provider
func TestDeleteProvider_EnabledProvider(t *testing.T) {
	mockRepo := new(MockProviderRepository)
	svc := NewProviderService(mockRepo)
	ctx := context.Background()

	// 模拟已启用的 Provider
	existingProvider := &model.Provider{
		ID:      1,
		Name:    "test-provider",
		Type:    "openai",
		BaseURL: "https://api.openai.com/v1",
		Timeout: 60,
		Enabled: true, // 已启用
	}

	mockRepo.On("GetByName", ctx, "test-provider").Return(existingProvider, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(nil)

	// 执行删除
	err := svc.DeleteProvider(ctx, "test-provider")

	// 验证
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestDeleteProvider_NotFound 测试删除不存在的 Provider
func TestDeleteProvider_NotFound(t *testing.T) {
	mockRepo := new(MockProviderRepository)
	svc := NewProviderService(mockRepo)
	ctx := context.Background()

	// 设置期望 - Provider 不存在
	mockRepo.On("GetByName", ctx, "non-existent").Return(nil, assert.AnError)

	// 执行删除
	err := svc.DeleteProvider(ctx, "non-existent")

	// 验证
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}
