package vllm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/model"
)

// TestAdapter_Chat_Success 测试 vLLM Chat 方法成功场景
func TestAdapter_Chat_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("expected path /v1/chat/completions, got %s", r.URL.Path)
		}

		// vLLM 可能有也可能没有 API Key
		_ = r.Header.Get("Authorization")

		// 返回模拟响应
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"id": "vllm-test-id",
			"object": "chat.completion",
			"created": 1234567890,
			"model": "qwen-turbo",
			"choices": [{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "Hello from vLLM!"
				},
				"finish_reason": "stop"
			}],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 5,
				"total_tokens": 15
			}
		}`))
	}))
	defer server.Close()

	// 创建 Provider 配置（使用包含 /v1 路径的 base URL）
	provider := &model.Provider{
		Name:     "test-vllm",
		Type:     "vllm",
		BaseURL:  server.URL + "/v1",
		Timeout:  30,
		Enabled:  true,
	}

	vllmAdapter, err := NewAdapter(provider)
	if err != nil {
		t.Fatalf("failed to create adapter: %v", err)
	}

	// 测试 Chat 方法
	req := &adapter.ChatRequest{
		Model: "qwen-turbo",
		Messages: []adapter.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := vllmAdapter.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "vllm-test-id" {
		t.Errorf("expected ID vllm-test-id, got %s", resp.ID)
	}

	if resp.Model != "qwen-turbo" {
		t.Errorf("expected model qwen-turbo, got %s", resp.Model)
	}

	if len(resp.Choices) != 1 {
		t.Errorf("expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "Hello from vLLM!" {
		t.Errorf("expected content 'Hello from vLLM!', got %s", resp.Choices[0].Message.Content)
	}
}

// TestAdapter_Chat_NoAPIKey 测试无 API Key 的场景（本地 vLLM）
func TestAdapter_Chat_NoAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证没有 Authorization 头
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Errorf("expected no Authorization header, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"id": "local-vllm-id",
			"object": "chat.completion",
			"created": 1234567890,
			"model": "local-model",
			"choices": [{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "Response from local vLLM"
				},
				"finish_reason": "stop"
			}],
			"usage": {
				"prompt_tokens": 8,
				"completion_tokens": 4,
				"total_tokens": 12
			}
		}`))
	}))
	defer server.Close()

	// 创建没有 API Key 的 Provider 配置
	provider := &model.Provider{
		Name:    "local-vllm",
		Type:    "vllm",
		BaseURL: server.URL,
		Timeout: 30,
		Enabled: true,
		// APIKey 为 nil
	}

	vllmAdapter, err := NewAdapter(provider)
	if err != nil {
		t.Fatalf("failed to create adapter: %v", err)
	}

	req := &adapter.ChatRequest{
		Model: "local-model",
		Messages: []adapter.Message{
			{Role: "user", Content: "Test"},
		},
	}

	resp, err := vllmAdapter.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Choices[0].Message.Content != "Response from local vLLM" {
		t.Errorf("unexpected content: %s", resp.Choices[0].Message.Content)
	}
}

// TestAdapter_ChatStream_Success 测试流式请求
func TestAdapter_ChatStream_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		chunks := []string{
			`data: {"id":"stream-1","object":"chat.completion.chunk","created":1234567890,"model":"qwen-turbo","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}`,
			`data: {"id":"stream-2","object":"chat.completion.chunk","created":1234567890,"model":"qwen-turbo","choices":[{"index":0,"delta":{"content":"Stream"},"finish_reason":null}]}`,
			`data: {"id":"stream-3","object":"chat.completion.chunk","created":1234567890,"model":"qwen-turbo","choices":[{"index":0,"delta":{"content":" response"},"finish_reason":null}]}`,
			`data: {"id":"stream-4","object":"chat.completion.chunk","created":1234567890,"model":"qwen-turbo","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
			`data: [DONE]`,
		}

		for _, chunk := range chunks {
			w.Write([]byte(chunk + "\n\n"))
		}
	}))
	defer server.Close()

	provider := &model.Provider{
		Name:     "stream-vllm",
		Type:     "vllm",
		BaseURL:  server.URL,
		Timeout:  30,
		Enabled:  true,
	}

	vllmAdapter, err := NewAdapter(provider)
	if err != nil {
		t.Fatalf("failed to create adapter: %v", err)
	}

	req := &adapter.ChatRequest{
		Model: "qwen-turbo",
		Messages: []adapter.Message{
			{Role: "user", Content: "Test streaming"},
		},
	}

	respChan, err := vllmAdapter.ChatStream(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 收集所有 chunk
	var chunks []*adapter.ChatStreamChunk
	for chunk := range respChan {
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 4 {
		t.Errorf("expected 4 chunks, got %d", len(chunks))
	}

	// 验证最后一个 chunk 有 finish_reason
	if chunks[3].Choices[0].FinishReason == nil || *chunks[3].Choices[0].FinishReason != "stop" {
		t.Error("expected last chunk to have finish_reason=stop")
	}
}

// TestAdapter_Config 测试 Adapter 配置方法
func TestAdapter_Config(t *testing.T) {
	provider := &model.Provider{
		Name:     "config-test",
		Type:     "vllm",
		BaseURL:  "http://localhost:8000",
		Timeout:  60,
		Enabled:  true,
	}

	vllmAdapter, err := NewAdapter(provider)
	if err != nil {
		t.Fatalf("failed to create adapter: %v", err)
	}

	if vllmAdapter.Type() != "vllm" {
		t.Errorf("expected type vllm, got %s", vllmAdapter.Type())
	}

	if vllmAdapter.Name() != "config-test" {
		t.Errorf("expected name config-test, got %s", vllmAdapter.Name())
	}

	if vllmAdapter.Timeout() != 60 {
		t.Errorf("expected timeout 60, got %d", vllmAdapter.Timeout())
	}

	_ = vllmAdapter.Config()
	// GetConfig() 只返回 ExtraConfig 和 FallbackModels
	// name 字段通过 Name() 方法获取
	if vllmAdapter.Name() != "config-test" {
		t.Errorf("expected name config-test, got %s", vllmAdapter.Name())
	}
}

// TestAdapter_Chat_WithDefaultParams 测试使用默认参数
func TestAdapter_Chat_WithDefaultParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 读取请求体，验证默认参数
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		// 验证请求包含默认参数
		bodyStr := string(body)
		if !strings.Contains(bodyStr, `"temperature":0.8`) {
			t.Error("expected default temperature 0.8 in request")
		}
		if !strings.Contains(bodyStr, `"max_tokens":1500`) {
			t.Error("expected default max_tokens 1500 in request")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"id": "test-id",
			"object": "chat.completion",
			"created": 1234567890,
			"model": "test-model",
			"choices": [{
				"index": 0,
				"message": {"role": "assistant", "content": "OK"},
				"finish_reason": "stop"
			}],
			"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
		}`))
	}))
	defer server.Close()

	provider := &model.Provider{
		Name:     "default-params-test",
		Type:     "vllm",
		BaseURL:  server.URL,
		Timeout:  30,
		Enabled:  true,
		ExtraConfig: model.JSON{
			"temperature": 0.8,
			"max_tokens":  1500.0,
			"top_p":       0.95,
		},
	}

	vllmAdapter, err := NewAdapter(provider)
	if err != nil {
		t.Fatalf("failed to create adapter: %v", err)
	}

	req := &adapter.ChatRequest{
		Model: "test-model",
		Messages: []adapter.Message{
			{Role: "user", Content: "Test"},
		},
		// 不提供 Temperature 和 MaxTokens，使用默认值
	}

	_, err = vllmAdapter.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
