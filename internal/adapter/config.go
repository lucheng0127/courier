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
// 支持两种格式：
// 1. 数组格式（数据库直接存储）：{"0": "model-1", "1": "model-2"}
// 2. Set 格式（Service 层转换）：{"model-1": true, "model-2": true}
func parseJSONToStringArray(j model.JSON) []string {
	result := make([]string, 0, len(j))

	// 优先尝试 Set 格式（key 是模型名，value 是 true）
	isSetFormat := true
	for _, v := range j {
		if boolVal, ok := v.(bool); !ok || !boolVal {
			isSetFormat = false
			break
		}
	}

	if isSetFormat && len(j) > 0 {
		// Set 格式：key 是模型名
		for k := range j {
			result = append(result, k)
		}
		return result
	}

	// 数组格式：value 是模型名（按索引顺序）
	// 需要按数字索引排序以保证顺序
	maxIdx := -1
	for k := range j {
		// 检查 key 是否是数字索引
		idx := -1
		// 尝试解析为整数
		allDigits := true
		for _, c := range k {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && len(k) > 0 {
			idx = atoi(k)
		}
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	if maxIdx >= 0 {
		// 按索引顺序收集
		result = make([]string, 0, maxIdx+1)
		for i := 0; i <= maxIdx; i++ {
			if v, ok := j[itoa(i)]; ok {
				if str, ok := v.(string); ok {
					result = append(result, str)
				}
			}
		}
		return result
	}

	// 兜底：直接取所有字符串值
	for _, v := range j {
		if str, ok := v.(string); ok {
			result = append(result, str)
		}
	}

	return result
}

// 简单的字符串转整数
func atoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

// 整数转字符串
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for i > 0 {
		buf = append(buf, byte('0'+i%10))
		i /= 10
	}
	// 反转
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
