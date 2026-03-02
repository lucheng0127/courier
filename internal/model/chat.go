package model

// ChatRequest Chat 请求（OpenAI 兼容格式）
type ChatRequest struct {
	Model       string             `json:"model" binding:"required"`        // provider/model_name 格式
	Messages    []ChatMessage      `json:"messages" binding:"required,min=1"`
	Stream      bool               `json:"stream"`                          // 是否流式响应
	Temperature *float64           `json:"temperature,omitempty"`
	MaxTokens   *int               `json:"max_tokens,omitempty"`
	TopP        *float64           `json:"top_p,omitempty"`
	N           *int               `json:"n,omitempty"`
	Stop        *string            `json:"stop,omitempty"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role" binding:"required"`    // system, user, assistant
	Content string `json:"content" binding:"required"`
}

// ChatResponse Chat 非流式响应（OpenAI 兼容格式）
type ChatResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`        // chat.completion
	Created int64              `json:"created"`
	Model   string             `json:"model"`         // 客户端请求的 model 参数
	Choices []ChatChoice       `json:"choices"`
	Usage   ChatUsage          `json:"usage"`
}

// ChatChoice 响应选项
type ChatChoice struct {
	Index        int            `json:"index"`
	Message      ChatMessage    `json:"message"`
	FinishReason string         `json:"finish_reason"` // stop, length, content_filter
}

// ChatUsage Token 使用统计
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatStreamResponse Chat 流式响应（SSE 格式）
type ChatStreamResponse struct {
	ID      string                `json:"id"`
	Object  string                `json:"object"`        // chat.completion.chunk
	Created int64                 `json:"created"`
	Model   string                `json:"model"`
	Choices []ChatStreamChoice    `json:"choices"`
}

// ChatStreamChoice 流式响应选项
type ChatStreamChoice struct {
	Index        int             `json:"index"`
	Delta        ChatMessageDelta `json:"delta"`
	FinishReason *string         `json:"finish_reason"`
}

// ChatMessageDelta 流式消息增量
type ChatMessageDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}
