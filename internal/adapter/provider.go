package adapter

import (
	"context"
)

// Provider LLM Provider 接口
type Provider interface {
	// Chat 完成对话调用（非流式）
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// ChatStream 流式对话调用
	ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatStreamChunk, error)

	// Type 返回 Provider 类型标识
	Type() string

	// Name 返回 Provider 实例名称
	Name() string
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature *float64  `json:"temperature,omitempty"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	// 其他可选参数...
}

// Message 聊天消息
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// ChatResponse 聊天响应（非流式）
type ChatResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 响应选项
type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatStreamChunk 流式响应块
type ChatStreamChunk struct {
	ID      string         `json:"id"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice 流式选项
type StreamChoice struct {
	Index        int          `json:"index"`
	Delta        MessageDelta `json:"delta"`
	FinishReason *string      `json:"finish_reason"`
}

// MessageDelta 消息增量
type MessageDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// Usage 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
