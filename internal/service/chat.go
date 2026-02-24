package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lucheng0127/courier/internal/client"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/repository"
)

// ChatService 聊天服务接口
type ChatService interface {
	Chat(ctx context.Context, userID uint, modelName string, req *client.ChatRequest) (*ChatResponse, error)
	ChatStream(ctx context.Context, userID uint, modelName string, req *client.ChatRequest) (<-chan *client.ChatChunk, <-chan error, func() (*ChatStreamResult, error))
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

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Delta 增量消息
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

// ChatStreamResult 流式对话结果
type ChatStreamResult struct {
	FullContent string
	Usage       *Usage
}

// chatService 聊天服务实现
type chatService struct {
	modelService    ModelService
	modelClient     client.ModelClient
	requestLogRepo  repository.RequestLogRepository
	logWriter       *asyncLogWriter
}

// NewChatService 创建聊天服务
func NewChatService(
	modelService ModelService,
	modelClient client.ModelClient,
	requestLogRepo repository.RequestLogRepository,
) ChatService {
	return &chatService{
		modelService:   modelService,
		modelClient:    modelClient,
		requestLogRepo: requestLogRepo,
		logWriter:      newAsyncLogWriter(requestLogRepo),
	}
}

func (s *chatService) Chat(ctx context.Context, userID uint, modelName string, req *client.ChatRequest) (*ChatResponse, error) {
	// 获取模型配置
	modelConfig, err := s.modelService.GetModelByName(modelName)
	if err != nil {
		return nil, errors.New("模型不存在")
	}

	// 记录开始时间
	startTime := time.Now()

	// 设置模型名称
	req.Model = modelName

	// 序列化请求消息
	messagesJSON, _ := json.Marshal(req.Messages)

	// 调用模型
	resp, err := s.modelClient.Chat(ctx, modelConfig, req)
	if err != nil {
		// 记录错误日志
		s.logWriter.Write(&model.RequestLog{
			UserID:          userID,
			ModelName:       modelName,
			RequestMessages: string(messagesJSON),
			Status:          "error",
			ErrorMessage:    err.Error(),
			LatencyMs:       int(time.Since(startTime).Milliseconds()),
		})
		return nil, err
	}

	// 构建响应内容
	var responseContent string
	if len(resp.Choices) > 0 {
		responseContent = resp.Choices[0].Message.Content
	}

	// 记录日志
	logEntry := &model.RequestLog{
		UserID:           userID,
		ModelName:        modelName,
		RequestMessages:  string(messagesJSON),
		ResponseContent:  responseContent,
		PromptTokens:     0,
		CompletionTokens: 0,
		TotalTokens:      0,
		LatencyMs:        int(time.Since(startTime).Milliseconds()),
		Status:           "success",
	}

	if resp.Usage != nil {
		logEntry.PromptTokens = resp.Usage.PromptTokens
		logEntry.CompletionTokens = resp.Usage.CompletionTokens
		logEntry.TotalTokens = resp.Usage.TotalTokens
	}

	s.logWriter.Write(logEntry)

	// 转换为服务响应
	return &ChatResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: convertChoices(resp.Choices),
		Usage:   convertUsage(resp.Usage),
	}, nil
}

func (s *chatService) ChatStream(ctx context.Context, userID uint, modelName string, req *client.ChatRequest) (<-chan *client.ChatChunk, <-chan error, func() (*ChatStreamResult, error)) {
	// 获取模型配置
	modelConfig, err := s.modelService.GetModelByName(modelName)
	if err != nil {
		errChan := make(chan error, 1)
		close(errChan)
		errChan <- errors.New("模型不存在")
		return nil, errChan, nil
	}

	// 记录开始时间
	startTime := time.Now()

	// 设置模型名称
	req.Model = modelName

	// 序列化请求消息
	messagesJSON, _ := json.Marshal(req.Messages)

	// 调用流式模型
	chunkChan, modelErrChan := s.modelClient.ChatStream(ctx, modelConfig, req)

	// 创建输出 channel
	outChunkChan := make(chan *client.ChatChunk, 16)
	outErrChan := make(chan error, 1)

	// 用于收集完整内容
	var fullContent strings.Builder
	var usage *client.Usage

	// 启动 goroutine 处理
	go func() {
		defer close(outChunkChan)
		defer close(outErrChan)

		for {
			select {
			case chunk, ok := <-chunkChan:
				if !ok {
					// 流结束，记录日志
					logEntry := &model.RequestLog{
						UserID:           userID,
						ModelName:        modelName,
						RequestMessages:  string(messagesJSON),
						ResponseContent:  fullContent.String(),
						PromptTokens:     0,
						CompletionTokens: 0,
						TotalTokens:      0,
						LatencyMs:        int(time.Since(startTime).Milliseconds()),
						Status:           "success",
					}

					if usage != nil {
						logEntry.PromptTokens = usage.PromptTokens
						logEntry.CompletionTokens = usage.CompletionTokens
						logEntry.TotalTokens = usage.TotalTokens
					}

					s.logWriter.Write(logEntry)
					return
				}

				// 收集内容
				if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
					fullContent.WriteString(chunk.Choices[0].Delta.Content)
				}

				// 保存 usage（通常在最后一个块中）
				if chunk.Usage != nil {
					usage = chunk.Usage
				}

				// 转发到输出
				select {
				case outChunkChan <- chunk:
				case <-ctx.Done():
					return
				}

			case err, ok := <-modelErrChan:
				if !ok {
					return
				}
				// 记录错误日志
				s.logWriter.Write(&model.RequestLog{
					UserID:          userID,
					ModelName:       modelName,
					RequestMessages: string(messagesJSON),
					ResponseContent: fullContent.String(),
					Status:          "error",
					ErrorMessage:    err.Error(),
					LatencyMs:       int(time.Since(startTime).Milliseconds()),
				})
				outErrChan <- err
				return

			case <-ctx.Done():
				// 客户端取消，记录中断日志
				s.logWriter.Write(&model.RequestLog{
					UserID:          userID,
					ModelName:       modelName,
					RequestMessages: string(messagesJSON),
					ResponseContent: fullContent.String(),
					Status:          "interrupted",
					ErrorMessage:    "client disconnected",
					LatencyMs:       int(time.Since(startTime).Milliseconds()),
				})
				return
			}
		}
	}()

	// 返回结果获取函数
	resultFunc := func() (*ChatStreamResult, error) {
		return &ChatStreamResult{
			FullContent: fullContent.String(),
			Usage:       convertUsage(usage),
		}, nil
	}

	return outChunkChan, outErrChan, resultFunc
}

// asyncLogWriter 异步日志写入器
type asyncLogWriter struct {
	repo  repository.RequestLogRepository
	queue chan *model.RequestLog
}

// newAsyncLogWriter 创建异步日志写入器
func newAsyncLogWriter(repo repository.RequestLogRepository) *asyncLogWriter {
	w := &asyncLogWriter{
		repo:  repo,
		queue: make(chan *model.RequestLog, 100),
	}
	go w.process()
	return w
}

// Write 写入日志（异步）
func (w *asyncLogWriter) Write(log *model.RequestLog) {
	select {
	case w.queue <- log:
	default:
		// 队列满，丢弃（可以记录到系统日志）
	}
}

// process 处理日志队列
func (w *asyncLogWriter) process() {
	for log := range w.queue {
		if err := w.repo.Create(log); err != nil {
			// 记录到系统日志
			fmt.Printf("写入请求日志失败: %v\n", err)
		}
	}
}

func convertChoices(choices []client.Choice) []Choice {
	result := make([]Choice, len(choices))
	for i, c := range choices {
		result[i] = Choice{
			Index:        c.Index,
			Message:      ChatMessage(c.Message),
			Delta:        convertDelta(c.Delta),
			FinishReason: c.FinishReason,
		}
	}
	return result
}

func convertDelta(delta *client.Delta) *Delta {
	if delta == nil {
		return nil
	}
	return &Delta{
		Role:    delta.Role,
		Content: delta.Content,
	}
}

func convertUsage(usage *client.Usage) *Usage {
	if usage == nil {
		return nil
	}
	return &Usage{
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
	}
}
