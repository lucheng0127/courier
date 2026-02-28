package service

import (
	"context"
	"fmt"
	"log"

	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/repository"
)

// ProviderService Provider 管理服务
type ProviderService struct {
	repo repository.ProviderRepository
}

// NewProviderService 创建 Provider Service
func NewProviderService(repo repository.ProviderRepository) *ProviderService {
	return &ProviderService{repo: repo}
}

// CreateProvider 创建 Provider
func (s *ProviderService) CreateProvider(ctx context.Context, provider *model.Provider) error {
	// 检查 name 唯一性
	exists, err := s.repo.ExistsByName(ctx, provider.Name)
	if err != nil {
		return fmt.Errorf("failed to check provider name: %w", err)
	}
	if exists {
		return fmt.Errorf("provider name already exists: %s", provider.Name)
	}

	// 创建 Provider
	if err := s.repo.Create(ctx, provider); err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// 如果启用，初始化并注册
	if provider.Enabled {
		return s.initAndRegisterProvider(provider)
	}

	return nil
}

// GetProvider 获取 Provider 实例
func (s *ProviderService) GetProvider(name string) (adapter.Provider, error) {
	provider, ok := adapter.GetProvider(name)
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", name)
	}
	return provider, nil
}

// ListProviders 列出所有 Provider 及状态
func (s *ProviderService) ListProviders(ctx context.Context) ([]*ProviderInfo, error) {
	providers, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	var infos []*ProviderInfo
	for _, p := range providers {
		instance, ok := adapter.GetProvider(p.Name)
		info := &ProviderInfo{
			Provider:  p,
			IsRunning: ok && instance != nil,
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// ReloadProvider 重新加载指定 Provider
func (s *ProviderService) ReloadProvider(ctx context.Context, name string) error {
	// 从数据库加载最新配置
	provider, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	// 如果未启用，注销并退出
	if !provider.Enabled {
		adapter.UnregisterProvider(name)
		log.Printf("Provider %s disabled and unregistered", name)
		return nil
	}

	// 创建新 Adapter 实例
	newProvider, err := adapter.NewAdapter(provider)
	if err != nil {
		return fmt.Errorf("failed to create adapter: %w", err)
	}

	// 原子替换 Registry 中的实例
	adapter.ReplaceProvider(newProvider)
	log.Printf("Provider %s reloaded successfully", name)

	return nil
}

// ReloadAllProviders 重新加载所有 Provider
func (s *ProviderService) ReloadAllProviders(ctx context.Context) error {
	providers, err := s.repo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list providers: %w", err)
	}

	var lastErr error
	for _, provider := range providers {
		if err := s.ReloadProvider(ctx, provider.Name); err != nil {
			log.Printf("Failed to reload provider %s: %v", provider.Name, err)
			lastErr = err
		}
	}

	return lastErr
}

// EnableProvider 启用 Provider
func (s *ProviderService) EnableProvider(ctx context.Context, name string) error {
	provider, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	if provider.Enabled {
		return fmt.Errorf("provider already enabled: %s", name)
	}

	provider.Enabled = true
	if err := s.repo.Update(ctx, provider); err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}

	return s.ReloadProvider(ctx, name)
}

// DisableProvider 禁用 Provider
func (s *ProviderService) DisableProvider(ctx context.Context, name string) error {
	provider, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	if !provider.Enabled {
		return fmt.Errorf("provider already disabled: %s", name)
	}

	provider.Enabled = false
	if err := s.repo.Update(ctx, provider); err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}

	adapter.UnregisterProvider(name)
	log.Printf("Provider %s disabled and unregistered", name)

	return nil
}

// InitProviders 初始化所有 Provider
func (s *ProviderService) InitProviders(ctx context.Context) error {
	providers, err := s.repo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list providers: %w", err)
	}

	var initErrs []error
	for _, provider := range providers {
		if err := s.initAndRegisterProvider(provider); err != nil {
			log.Printf("Failed to initialize provider %s: %v", provider.Name, err)
			initErrs = append(initErrs, err)
		}
	}

	if len(initErrs) > 0 {
		return fmt.Errorf("failed to initialize %d providers", len(initErrs))
	}

	return nil
}

// initAndRegisterProvider 初始化并注册 Provider
func (s *ProviderService) initAndRegisterProvider(provider *model.Provider) error {
	// 跳过未启用的 Provider
	if !provider.Enabled {
		log.Printf("Provider %s is disabled, skipping initialization", provider.Name)
		return nil
	}

	// 创建 Adapter 实例
	instance, err := adapter.NewAdapter(provider)
	if err != nil {
		return fmt.Errorf("failed to create adapter: %w", err)
	}

	// 注册到 Registry
	adapter.RegisterProvider(instance)
	log.Printf("Provider %s (type: %s) initialized successfully", provider.Name, provider.Type)

	return nil
}

// ProviderInfo Provider 信息
type ProviderInfo struct {
	Provider  *model.Provider `json:"provider"`
	IsRunning bool            `json:"is_running"`
}
