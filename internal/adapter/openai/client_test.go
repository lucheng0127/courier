package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucheng0127/courier/internal/adapter"
)

// TestConvertChatRequest 测试请求格式转换
func TestConvertChatRequest(t *testing.T) {
	defaultConfig := map[string]any{
		"temperature": 0.7,
		"max_tokens":  2000.0,
		"top_p":       0.9,
	}

	temp := 0.5
	maxTokens := 1000

	req := &adapter.ChatRequest{
		Model:       "gpt-4",
		Messages:    []adapter.Message{{Role: "user", Content: "Hello"}},
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}

	result := ConvertChatRequest(req, defaultConfig)

	if result.Model != "gpt-4" {
		t.Errorf("expected model gpt-4, got %s", result.Model)
	}

	if len(result.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(result.Messages))
	}

	if *result.Temperature != 0.5 {
		t.Errorf("expected temperature 0.5, got %f", *result.Temperature)
	}

	if *result.MaxTokens != 1000 {
		t.Errorf("expected max_tokens 1000, got %d", *result.MaxTokens)
	}

	// top_p 应该从默认配置读取
	if *result.TopP != 0.9 {
		t.Errorf("expected top_p 0.9, got %f", *result.TopP)
	}
}

// TestConvertChatRequest_DefaultConfig 测试使用默认配置
func TestConvertChatRequest_DefaultConfig(t *testing.T) {
	defaultConfig := map[string]any{
		"temperature": 0.7,
		"max_tokens":  2000.0,
	}

	req := &adapter.ChatRequest{
		Model:    "gpt-4",
		Messages: []adapter.Message{{Role: "user", Content: "Hello"}},
		// 不提供 Temperature 和 MaxTokens
	}

	result := ConvertChatRequest(req, defaultConfig)

	if *result.Temperature != 0.7 {
		t.Errorf("expected temperature 0.7 from config, got %f", *result.Temperature)
	}

	if *result.MaxTokens != 2000 {
		t.Errorf("expected max_tokens 2000 from config, got %d", *result.MaxTokens)
	}
}

// TestConvertChatResponse 测试响应格式转换
func TestConvertChatResponse(t *testing.T) {
	resp := &ChatResponse{
		ID:     "chatcmpl-123",
		Model:  "gpt-4",
		Object: "chat.completion",
		Choices: []ChatChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: "Hello!",
				},
				FinishReason: "stop",
			},
		},
		Usage: ChatUsage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}

	result := ConvertChatResponse(resp)

	if result.ID != "chatcmpl-123" {
		t.Errorf("expected ID chatcmpl-123, got %s", result.ID)
	}

	if result.Model != "gpt-4" {
		t.Errorf("expected model gpt-4, got %s", result.Model)
	}

	if len(result.Choices) != 1 {
		t.Errorf("expected 1 choice, got %d", len(result.Choices))
	}

	if result.Choices[0].Message.Content != "Hello!" {
		t.Errorf("expected content 'Hello!', got %s", result.Choices[0].Message.Content)
	}

	if result.Usage.TotalTokens != 30 {
		t.Errorf("expected total_tokens 30, got %d", result.Usage.TotalTokens)
	}
}

// TestDoChatRequest_Success 测试成功的非流式请求
func TestDoChatRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("expected path /v1/chat/completions, got %s", r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization 'Bearer test-key', got %s", auth)
		}

		// 返回模拟响应
		resp := ChatResponse{
			ID:     "test-id",
			Model:  "gpt-4",
			Object: "chat.completion",
			Choices: []ChatChoice{
				{
					Index: 0,
					Message: ChatMessage{
						Role:    "assistant",
						Content: "Test response",
					},
					FinishReason: "stop",
				},
			},
			Usage: ChatUsage{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", 30)

	req := &ChatRequest{
		Model:    "gpt-4",
		Messages: []ChatMessage{{Role: "user", Content: "Hello"}},
	}

	resp, err := client.DoChatRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", resp.ID)
	}
}

// TestDoChatRequest_Error 测试错误响应
func TestDoChatRequest_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{
				"message": "Invalid API key",
				"type":    "invalid_request_error",
				"code":    "invalid_api_key",
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "invalid-key", 30)

	req := &ChatRequest{
		Model:    "gpt-4",
		Messages: []ChatMessage{{Role: "user", Content: "Hello"}},
	}

	_, err := client.DoChatRequest(context.Background(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "Invalid API key") {
		t.Errorf("expected error containing 'Invalid API key', got %v", err)
	}
}

// TestDoChatRequest_NoAPIKey 测试无 API Key 的请求
func TestDoChatRequest_NoAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证没有 Authorization 头
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Errorf("expected no Authorization header, got %s", auth)
		}

		resp := ChatResponse{
			ID:    "test-id",
			Model: "gpt-4",
			Choices: []ChatChoice{
				{
					Index: 0,
					Message: ChatMessage{
						Role:    "assistant",
						Content: "Test response",
					},
					FinishReason: "stop",
				},
			},
			Usage: ChatUsage{TotalTokens: 15},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 不提供 API Key
	client := NewClient(server.URL, "", 30)

	req := &ChatRequest{
		Model:    "gpt-4",
		Messages: []ChatMessage{{Role: "user", Content: "Hello"}},
	}

	resp, err := client.DoChatRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", resp.ID)
	}
}

// TestDoChatStreamRequest_Success 测试成功的流式请求
func TestDoChatStreamRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置 SSE 响应头
		w.Header().Set("Content-Type", "text/event-stream")

		// 发送 SSE 数据
		chunks := []string{
			`data: {"id":"chunk-1","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}`,
			`data: {"id":"chunk-2","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
			`data: {"id":"chunk-3","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":"stop"}]}`,
			`data: [DONE]`,
		}

		for _, chunk := range chunks {
			w.Write([]byte(chunk + "\n\n"))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", 30)

	req := &ChatRequest{
		Model:    "gpt-4",
		Messages: []ChatMessage{{Role: "user", Content: "Hello"}},
	}

	respChan := make(chan *adapter.ChatStreamChunk, 10)

	err := client.DoChatStreamRequest(context.Background(), req, respChan)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	close(respChan)

	// 收集所有 chunk
	var chunks []*adapter.ChatStreamChunk
	for chunk := range respChan {
		chunks = append(chunks, chunk)
	}

	// 验证收到 3 个数据块
	if len(chunks) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(chunks))
	}
}

// TestDoChatStreamRequest_ContextCancel 测试 context 取消
func TestDoChatStreamRequest_ContextCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		// 持续发送数据
		for i := 0; ; i++ {
			chunk := `data: {"id":"chunk-` + string(rune(i)) + `","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"data"},"finish_reason":null}]}`
			w.Write([]byte(chunk + "\n\n"))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", 30)

	req := &ChatRequest{
		Model:    "gpt-4",
		Messages: []ChatMessage{{Role: "user", Content: "Hello"}},
	}

	respChan := make(chan *adapter.ChatStreamChunk, 10)

	// 创建可取消的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 启动 goroutine，然后立即取消
	go func() {
		cancel()
	}()

	err := client.DoChatStreamRequest(ctx, req, respChan)

	if err == nil {
		t.Error("expected context canceled error, got nil")
	}

	close(respChan)
}
