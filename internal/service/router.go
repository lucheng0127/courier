package service

import (
	"fmt"
	"strings"

	"github.com/lucheng0127/courier/internal/adapter"
)

// RouterService 模型路由服务
type RouterService struct{}

// NewRouterService 创建路由服务
func NewRouterService() *RouterService {
	return &RouterService{}
}

// ModelInfo 模型信息
type ModelInfo struct {
	ProviderName string // Provider 名称
	ModelName    string // 模型名称
	Provider     adapter.Provider
}

// ParseModel 解析模型参数 `provider/model_name`
func (s *RouterService) ParseModel(model string) (*ModelInfo, error) {
	parts := strings.Split(model, "/")
	if len(parts) != 2 {
		return nil, &ModelFormatError{Model: model}
	}

	providerName := parts[0]
	modelName := parts[1]

	if providerName == "" || modelName == "" {
		return nil, &ModelFormatError{Model: model}
	}

	return &ModelInfo{
		ProviderName: providerName,
		ModelName:    modelName,
	}, nil
}

// ResolveProvider 根据 provider 名称获取 Provider 实例
func (s *RouterService) ResolveProvider(providerName string) (adapter.Provider, error) {
	provider, ok := adapter.GetProvider(providerName)
	if !ok {
		return nil, &ProviderNotFoundError{ProviderName: providerName}
	}

	return provider, nil
}

// ResolveModel 解析模型并获取 Provider
func (s *RouterService) ResolveModel(model string) (*ModelInfo, error) {
	info, err := s.ParseModel(model)
	if err != nil {
		return nil, err
	}

	provider, err := s.ResolveProvider(info.ProviderName)
	if err != nil {
		return nil, err
	}

	info.Provider = provider
	return info, nil
}

// GetAvailableModels 获取所有可用模型列表
func (s *RouterService) GetAvailableModels() []string {
	var models []string

	providers := adapter.ListProviders()
	for _, provider := range providers {
		// 格式: provider-name/*
		models = append(models, fmt.Sprintf("%s/*", provider.Name()))
	}

	return models
}

// === 错误类型 ===

// ModelFormatError 模型格式错误
type ModelFormatError struct {
	Model string
}

func (e *ModelFormatError) Error() string {
	return fmt.Sprintf("invalid model format: %s (expected format: provider/model_name)", e.Model)
}

// ProviderNotFoundError Provider 不存在错误
type ProviderNotFoundError struct {
	ProviderName string
}

func (e *ProviderNotFoundError) Error() string {
	return fmt.Sprintf("provider not found: %s", e.ProviderName)
}

// ProviderDisabledError Provider 未启用错误
type ProviderDisabledError struct {
	ProviderName string
}

func (e *ProviderDisabledError) Error() string {
	return fmt.Sprintf("provider is disabled: %s", e.ProviderName)
}
