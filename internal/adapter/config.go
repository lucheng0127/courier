package adapter

import (
	"time"

	"github.com/lucheng0127/courier/internal/model"
)

// ProviderConfig Adapter 配置
type ProviderConfig struct {
	Name           string            // Provider 实例名称
	Type           string            // Provider 类型
	BaseURL        string            // API 地址
	Timeout        time.Duration     // 超时时间
	TimeoutSeconds int               // 超时时间（秒）
	APIKey         string            // API Key（可选）
	ExtraConfig    map[string]any    // 扩展配置
	FallbackModels []string          // Fallback 模型列表
}

// NewProviderConfig 从 model.Provider 创建 ProviderConfig
func NewProviderConfig(provider *model.Provider) *ProviderConfig {
	timeoutSeconds := provider.Timeout
	if timeoutSeconds == 0 {
		timeoutSeconds = 300 // 默认 300 秒
	}

	config := &ProviderConfig{
		Name:           provider.Name,
		Type:           provider.Type,
		BaseURL:        provider.BaseURL,
		Timeout:        time.Duration(timeoutSeconds) * time.Second,
		TimeoutSeconds: timeoutSeconds,
		ExtraConfig:    provider.ExtraConfig,
	}

	if provider.APIKey != nil {
		config.APIKey = *provider.APIKey
	}

	// 解析 FallbackModels
	// fallback_models 是 JSONB 类型，存储为数组 ["model-1", "model-2"]
	// sqlx 会将其解析为 JSON (map[string]any) 类型
	// 我们需要将整个 map 转换为字符串数组
	if len(provider.FallbackModels) > 0 {
		// 尝试直接解析为 []string
		// 由于 JSON 数组在 Go 中可能被解析为不同类型，我们需要处理多种情况
		config.FallbackModels = parseJSONToStringArray(provider.FallbackModels)
	}

	// 也可以从 ExtraConfig 获取 fallback_models（兼容旧配置）
	if len(config.FallbackModels) == 0 && provider.ExtraConfig != nil {
		if fbModels, ok := provider.ExtraConfig["fallback_models"].([]interface{}); ok {
			config.FallbackModels = make([]string, 0, len(fbModels))
			for _, m := range fbModels {
				if str, ok := m.(string); ok {
					config.FallbackModels = append(config.FallbackModels, str)
				}
			}
		}
	}

	return config
}

// GetConfig 返回完整配置（包括 FallbackModels）
func (c *ProviderConfig) GetConfig() map[string]any {
	cfg := make(map[string]any)
	for k, v := range c.ExtraConfig {
		cfg[k] = v
	}
	// 添加 fallback_models
	if len(c.FallbackModels) > 0 {
		cfg["fallback_models"] = c.FallbackModels
	}
	return cfg
}

// parseJSONToStringArray 从 JSON map 解析字符串数组
func parseJSONToStringArray(j model.JSON) []string {
	result := make([]string, 0)

	// 尝试各种可能的类型断言
	for _, v := range j {
		if str, ok := v.(string); ok {
			result = append(result, str)
		} else if arr, ok := v.([]interface{}); ok {
			// 嵌套数组
			for _, item := range arr {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
		}
	}

	return result
}
