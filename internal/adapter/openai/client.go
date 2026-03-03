package openai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lucheng0127/courier/internal/adapter"
)

// Client OpenAI API 客户端
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient 创建 OpenAI 客户端
func NewClient(baseURL, apiKey string, timeout int) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// ChatRequest OpenAI API 请求格式
type ChatRequest struct {
	Model       string              `json:"model"`
	Messages    []ChatMessage       `json:"messages"`
	Temperature *float64            `json:"temperature,omitempty"`
	MaxTokens   *int                `json:"max_tokens,omitempty"`
	TopP        *float64            `json:"top_p,omitempty"`
	Stream      bool                `json:"stream,omitempty"`
}

// ChatMessage OpenAI API 消息格式
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse OpenAI API 响应格式
type ChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []ChatChoice   `json:"choices"`
	Usage   ChatUsage      `json:"usage"`
}

// ChatChoice OpenAI API 选择项
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatUsage OpenAI API 使用量
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamChunk OpenAI SSE 流式响应块
type StreamChunk struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []StreamChoice    `json:"choices"`
}

// StreamChoice 流式选择项
type StreamChoice struct {
	Index        int            `json:"index"`
	Delta        StreamDelta   `json:"delta"`
	FinishReason *string       `json:"finish_reason"`
}

// StreamDelta 流式增量
type StreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// ErrorResponse OpenAI API 错误响应
type ErrorResponse struct {
	ErrorDetail ErrorDetail `json:"error"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

func (e *ErrorResponse) Error() string {
	return e.ErrorDetail.Message
}

// DoChatRequest 执行非流式聊天请求（导出供其他 Adapter 使用）
func (c *Client) DoChatRequest(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	req.Stream = false

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// TraceID 透传
	if traceID := getTraceID(ctx); traceID != "" {
		httpReq.Header.Set("X-Trace-ID", traceID)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 检查 HTTP 状态码
	if httpResp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.ErrorDetail.Message != "" {
			return nil, &errResp
		}
		return nil, fmt.Errorf("request failed with status %d: %s", httpResp.StatusCode, string(respBody))
	}

	var resp ChatResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// DoChatStreamRequest 执行流式聊天请求（导出供其他 Adapter 使用）
func (c *Client) DoChatStreamRequest(ctx context.Context, req *ChatRequest, respChan chan<- *adapter.ChatStreamChunk) error {
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// TraceID 透传
	if traceID := getTraceID(ctx); traceID != "" {
		httpReq.Header.Set("X-Trace-ID", traceID)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// 检查 HTTP 状态码
	if httpResp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(httpResp.Body)
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.ErrorDetail.Message != "" {
			return &errResp
		}
		return fmt.Errorf("request failed with status %d: %s", httpResp.StatusCode, string(respBody))
	}

	// 解析 SSE 流
	scanner := bufio.NewScanner(httpResp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// 解析 data: 行
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// 检查结束标记
		if data == "[DONE]" {
			break
		}

		// 解析 JSON
		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // 跳过无效数据
		}

		// 转换为内部格式并发送
		internalChunk := convertStreamChunk(&chunk)
		select {
		case respChan <- internalChunk:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read stream: %w", err)
	}

	return nil
}

// getTraceID 从 context 获取 TraceID
func getTraceID(ctx context.Context) string {
	// 从 context 中获取 TraceID
	// 假设 TraceID 存储在 context 中
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

// ConvertChatRequest 将内部请求格式转换为 OpenAI 格式（导出供其他 Adapter 使用）
func ConvertChatRequest(req *adapter.ChatRequest, defaultConfig map[string]any) *ChatRequest {
	openaiReq := &ChatRequest{
		Model:    req.Model,
		Messages: make([]ChatMessage, len(req.Messages)),
	}

	// 转换消息
	for i, msg := range req.Messages {
		openaiReq.Messages[i] = ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 优先使用请求参数
	if req.Temperature != nil {
		openaiReq.Temperature = req.Temperature
	} else if temp, ok := defaultConfig["temperature"].(float64); ok {
		openaiReq.Temperature = &temp
	}

	if req.MaxTokens != nil {
		openaiReq.MaxTokens = req.MaxTokens
	} else if maxTokens, ok := defaultConfig["max_tokens"].(float64); ok {
		val := int(maxTokens)
		openaiReq.MaxTokens = &val
	}

	if topP, ok := defaultConfig["top_p"].(float64); ok {
		openaiReq.TopP = &topP
	}

	return openaiReq
}

// ConvertChatResponse 将 OpenAI 响应格式转换为内部格式（导出供其他 Adapter 使用）
func ConvertChatResponse(resp *ChatResponse) *adapter.ChatResponse {
	choices := make([]adapter.Choice, len(resp.Choices))
	for i, c := range resp.Choices {
		choices[i] = adapter.Choice{
			Index: c.Index,
			Message: adapter.Message{
				Role:    c.Message.Role,
				Content: c.Message.Content,
			},
			FinishReason: c.FinishReason,
		}
	}

	return &adapter.ChatResponse{
		ID:      resp.ID,
		Model:   resp.Model,
		Choices: choices,
		Usage: adapter.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}
}

// convertStreamChunk 将 OpenAI 流式块转换为内部格式
func convertStreamChunk(chunk *StreamChunk) *adapter.ChatStreamChunk {
	choices := make([]adapter.StreamChoice, len(chunk.Choices))
	for i, c := range chunk.Choices {
		choices[i] = adapter.StreamChoice{
			Index: c.Index,
			Delta: adapter.MessageDelta{
				Role:    c.Delta.Role,
				Content: c.Delta.Content,
			},
			FinishReason: c.FinishReason,
		}
	}

	return &adapter.ChatStreamChunk{
		ID:      chunk.ID,
		Model:   chunk.Model,
		Choices: choices,
	}
}
