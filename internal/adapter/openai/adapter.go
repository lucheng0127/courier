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
	// 创建客户端
	client := NewClient(a.config.BaseURL, a.config.APIKey, a.config.TimeoutSeconds)

	// 转换请求格式
	openaiReq := ConvertChatRequest(req, a.config.ExtraConfig)

	// 设置超时
	if a.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, a.config.Timeout)
		defer cancel()
	}

	// 发送请求
	resp, err := client.DoChatRequest(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	return ConvertChatResponse(resp), nil
}

// ChatStream 流式对话调用
func (a *Adapter) ChatStream(ctx context.Context, req *adapter.ChatRequest) (<-chan *adapter.ChatStreamChunk, error) {
	// 创建客户端
	client := NewClient(a.config.BaseURL, a.config.APIKey, a.config.TimeoutSeconds)

	// 转换请求格式
	openaiReq := ConvertChatRequest(req, a.config.ExtraConfig)

	// 创建响应 channel
	respChan := make(chan *adapter.ChatStreamChunk, 10)

	// 启动 goroutine 处理流式请求
	go func() {
		defer close(respChan)

		// 设置超时
		streamCtx := ctx
		if a.config.Timeout > 0 {
			var cancel context.CancelFunc
			streamCtx, cancel = context.WithTimeout(ctx, a.config.Timeout)
			defer cancel()
		}

		// 发送流式请求
		err := client.DoChatStreamRequest(streamCtx, openaiReq, respChan)
		if err != nil {
			// 错误处理：可以选择发送错误块或忽略
			// 这里简单地记录错误后退出
			_ = err
		}
	}()

	return respChan, nil
}

// Type 返回 Provider 类型
func (a *Adapter) Type() string {
	return string(adapter.AdapterTypeOpenAI)
}

// Name 返回 Provider 实例名称
func (a *Adapter) Name() string {
	return a.config.Name
}

// Timeout 返回超时时间（秒）
func (a *Adapter) Timeout() int {
	return a.config.TimeoutSeconds
}

// Config 返回配置信息
func (a *Adapter) Config() map[string]any {
	return a.config.GetConfig()
}

func init() {
	adapter.RegisterAdapterType(adapter.AdapterTypeOpenAI, NewAdapter)
}
