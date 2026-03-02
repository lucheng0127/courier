package openai

import (
	"context"
	"fmt"

	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/model"
)

// Adapter OpenAI Adapter
type Adapter struct {
	config *adapter.ProviderConfig
}

// NewAdapter 创建 OpenAI Adapter
func NewAdapter(provider *model.Provider) (adapter.Provider, error) {
	config := adapter.NewProviderConfig(provider)

	// 验证配置
	if config.BaseURL == "" {
		return nil, fmt.Errorf("openai adapter requires base_url")
	}

	return &Adapter{config: config}, nil
}

// Chat 完成对话调用（非流式）
func (a *Adapter) Chat(ctx context.Context, req *adapter.ChatRequest) (*adapter.ChatResponse, error) {
	// TODO: 后续实现实际调用逻辑
	return nil, fmt.Errorf("not implemented yet")
}

// ChatStream 流式对话调用
func (a *Adapter) ChatStream(ctx context.Context, req *adapter.ChatRequest) (<-chan *adapter.ChatStreamChunk, error) {
	// TODO: 后续实现实际调用逻辑
	return nil, fmt.Errorf("not implemented yet")
}

// Type 返回 Provider 类型
func (a *Adapter) Type() string {
	return string(adapter.AdapterTypeOpenAI)
}

// Name 返回 Provider 实例名称
func (a *Adapter) Name() string {
	return a.config.Name
}

func init() {
	adapter.RegisterAdapterType(adapter.AdapterTypeOpenAI, NewAdapter)
}
