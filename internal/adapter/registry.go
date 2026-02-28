package adapter

import (
	"fmt"
	"sync"

	"github.com/lucheng0127/courier/internal/model"
)

// AdapterType Provider 类型标识
type AdapterType string

const (
	AdapterTypeOpenAI    AdapterType = "openai"
	AdapterTypeAnthropic AdapterType = "anthropic"
	AdapterTypeVLLM      AdapterType = "vllm"
	AdapterTypeOllama    AdapterType = "ollama"
)

// AdapterFactory Adapter 工厂函数
type AdapterFactory func(config *model.Provider) (Provider, error)

var (
	registry      = map[AdapterType]AdapterFactory{}
	providerStore = &sync.Map{} // name -> Provider
	registryMutex sync.RWMutex
)

// RegisterAdapterType 注册 Adapter 类型
func RegisterAdapterType(adapterType AdapterType, factory AdapterFactory) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	registry[adapterType] = factory
}

// NewAdapter 创建 Provider 实例
func NewAdapter(config *model.Provider) (Provider, error) {
	registryMutex.RLock()
	factory, ok := registry[AdapterType(config.Type)]
	registryMutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown provider type: %s", config.Type)
	}

	return factory(config)
}

// RegisterProvider 注册 Provider 实例
func RegisterProvider(provider Provider) {
	providerStore.Store(provider.Name(), provider)
}

// UnregisterProvider 注销 Provider 实例
func UnregisterProvider(name string) {
	providerStore.Delete(name)
}

// GetProvider 获取 Provider 实例
func GetProvider(name string) (Provider, bool) {
	value, ok := providerStore.Load(name)
	if !ok {
		return nil, false
	}
	return value.(Provider), true
}

// ListProviders 列出所有已注册的 Provider
func ListProviders() []Provider {
	var providers []Provider
	providerStore.Range(func(key, value any) bool {
		providers = append(providers, value.(Provider))
		return true
	})
	return providers
}

// ReplaceProvider 替换 Provider 实例（原子操作）
func ReplaceProvider(provider Provider) {
	providerStore.Store(provider.Name(), provider)
}
