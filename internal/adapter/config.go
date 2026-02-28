package adapter

import (
	"time"

	"github.com/lucheng0127/courier/internal/model"
)

// ProviderConfig Adapter 配置
type ProviderConfig struct {
	Name        string            // Provider 实例名称
	Type        string            // Provider 类型
	BaseURL     string            // API 地址
	Timeout     time.Duration     // 超时时间
	APIKey      string            // API Key（可选）
	ExtraConfig map[string]any    // 扩展配置
}

// NewProviderConfig 从 model.Provider 创建 ProviderConfig
func NewProviderConfig(provider *model.Provider) *ProviderConfig {
	timeout := time.Duration(provider.Timeout) * time.Second
	if timeout == 0 {
		timeout = 300 * time.Second // 默认 300 秒
	}

	config := &ProviderConfig{
		Name:        provider.Name,
		Type:        provider.Type,
		BaseURL:     provider.BaseURL,
		Timeout:     timeout,
		ExtraConfig: provider.ExtraConfig,
	}

	if provider.APIKey != nil {
		config.APIKey = *provider.APIKey
	}

	return config
}
