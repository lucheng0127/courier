package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lucheng0127/courier/pkg/config"
)

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Stream      bool          `json:"stream,omitempty"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []Choice      `json:"choices"`
	Usage   *Usage        `json:"usage,omitempty"`
}

// Choice 选择
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message,omitempty"`
	Delta        *Delta      `json:"delta,omitempty"`
	FinishReason string      `json:"finish_reason,omitempty"`
}

// Delta 增量消息（流式）
type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// Usage Token 使用量
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatChunk 流式响应块
type ChatChunk struct {
	ID      string  `json:"id"`
	Object  string  `json:"object"`
	Created int64   `json:"created"`
	Model   string  `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage  `json:"usage,omitempty"`
}

// ModelClient 模型客户端接口
type ModelClient interface {
	Chat(ctx context.Context, modelConfig *config.ModelConfig, req *ChatRequest) (*ChatResponse, error)
	ChatStream(ctx context.Context, modelConfig *config.ModelConfig, req *ChatRequest) (<-chan *ChatChunk, <-chan error)
}

// modelClient 模型客户端实现
type modelClient struct {
	httpClient *http.Client
}

// NewModelClient 创建模型客户端
func NewModelClient(timeout time.Duration) ModelClient {
	return &modelClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *modelClient) Chat(ctx context.Context, modelConfig *config.ModelConfig, req *ChatRequest) (*ChatResponse, error) {
	// 构建请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建 HTTP 请求
	url := modelConfig.BaseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+modelConfig.APIKey)

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("上游返回错误: %s, status: %d", string(body), resp.StatusCode)
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &chatResp, nil
}

func (c *modelClient) ChatStream(ctx context.Context, modelConfig *config.ModelConfig, req *ChatRequest) (<-chan *ChatChunk, <-chan error) {
	chunkChan := make(chan *ChatChunk, 16)
	errChan := make(chan error, 1)

	go func() {
		defer close(chunkChan)
		defer close(errChan)

		// 设置流式请求
		req.Stream = true
		reqBody, err := json.Marshal(req)
		if err != nil {
			errChan <- fmt.Errorf("序列化请求失败: %w", err)
			return
		}

		// 创建 HTTP 请求
		url := modelConfig.BaseURL + "/chat/completions"
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
		if err != nil {
			errChan <- fmt.Errorf("创建请求失败: %w", err)
			return
		}

		// 设置请求头
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+modelConfig.APIKey)
		httpReq.Header.Set("Accept", "text/event-stream")

		// 发送请求
		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			errChan <- fmt.Errorf("发送请求失败: %w", err)
			return
		}
		defer resp.Body.Close()

		// 检查状态码
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errChan <- fmt.Errorf("上游返回错误: %s, status: %d", string(body), resp.StatusCode)
			return
		}

		// 读取流式响应
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// 跳过空行和注释
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// 解析 SSE 格式
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			// 检查结束标记
			if data == "[DONE]" {
				break
			}

			// 解析 JSON
			var chunk ChatChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				errChan <- fmt.Errorf("解析响应块失败: %w", err)
				return
			}

			// 发送到 channel
			select {
			case chunkChan <- &chunk:
			case <-ctx.Done():
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("读取流失败: %w", err)
			return
		}
	}()

	return chunkChan, errChan
}
