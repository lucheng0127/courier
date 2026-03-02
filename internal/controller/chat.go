package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// ChatController Chat API 控制器
type ChatController struct {
	router *service.RouterService
}

// NewChatController 创建 Chat 控制器
func NewChatController(router *service.RouterService) *ChatController {
	return &ChatController{router: router}
}

// ChatCompletions Chat Completions 端点
// POST /v1/chat/completions
func (c *ChatController) ChatCompletions(ctx *gin.Context) {
	startTime := time.Now()

	// 解析请求
	var req model.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 解析模型参数
	modelInfo, err := c.router.ResolveModel(req.Model)
	if err != nil {
		switch e := err.(type) {
		case *service.ModelFormatError:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": e.Error(),
					"type":    "invalid_request_error",
				},
			})
		case *service.ProviderNotFoundError:
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": e.Error(),
					"type":    "invalid_request_error",
				},
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message": "Failed to resolve model",
					"type":    "api_error",
				},
			})
		}
		return
	}

	// 生成请求 ID
	requestID := "chatcmpl-" + uuid.New().String()

	// 处理流式/非流式响应
	if req.Stream {
		c.handleStreamResponse(ctx, &req, modelInfo, requestID, startTime)
	} else {
		c.handleNonStreamResponse(ctx, &req, modelInfo, requestID, startTime)
	}
}

// handleNonStreamResponse 处理非流式响应
func (c *ChatController) handleNonStreamResponse(ctx *gin.Context, req *model.ChatRequest, modelInfo *service.ModelInfo, requestID string, startTime time.Time) {
	// 转换请求格式
	adapterReq := c.toAdapterRequest(req, modelInfo.ModelName)

	// 调用 Provider
	adapterResp, err := modelInfo.Provider.Chat(ctx.Request.Context(), adapterReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "api_error",
			},
		})
		return
	}

	// 转换响应格式
	resp := c.toChatResponse(adapterResp, requestID, req.Model, time.Now())

	// 记录日志
	c.logRequest(ctx, requestID, req, modelInfo, &resp.Usage, time.Since(startTime).Milliseconds(), "success", "")

	ctx.JSON(http.StatusOK, resp)
}

// handleStreamResponse 处理流式响应
func (c *ChatController) handleStreamResponse(ctx *gin.Context, req *model.ChatRequest, modelInfo *service.ModelInfo, requestID string, startTime time.Time) {
	// 转换请求格式
	adapterReq := c.toAdapterRequest(req, modelInfo.ModelName)

	// 调用 Provider
	chunks, err := modelInfo.Provider.ChatStream(ctx.Request.Context(), adapterReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "api_error",
			},
		})
		return
	}

	// 设置 SSE Header
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	// 获取响应写入器
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Streaming not supported",
				"type":    "api_error",
			},
		})
		return
	}

	// 流式发送数据
	var totalTokens int
	for chunk := range chunks {
		// 检查客户端是否断开
		select {
		case <-ctx.Request.Context().Done():
			return
		default:
		}

		// 转换为 SSE 格式
		sseChunk := c.toStreamChunk(chunk, requestID, req.Model)

		// 发送数据
		data, _ := json.Marshal(sseChunk)
		fmt.Fprintf(ctx.Writer, "data: %s\n", data)
		fmt.Fprint(ctx.Writer, "\n")
		flusher.Flush()

		totalTokens += len(chunk.Choices[0].Delta.Content)
	}

	// 发送结束标记
	fmt.Fprint(ctx.Writer, "data: [DONE]\n\n")
	flusher.Flush()

	// 记录日志
	usage := &model.ChatUsage{
		PromptTokens:     0, // 流式响应无法精确统计
		CompletionTokens: totalTokens,
		TotalTokens:      totalTokens,
	}
	c.logRequest(ctx, requestID, req, modelInfo, usage, time.Since(startTime).Milliseconds(), "success", "")
}

// toAdapterRequest 转换为 Adapter 请求格式
func (c *ChatController) toAdapterRequest(req *model.ChatRequest, modelName string) *adapter.ChatRequest {
	messages := make([]adapter.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = adapter.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	return &adapter.ChatRequest{
		Messages:    messages,
		Model:       modelName,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}
}

// toChatResponse 转换为 Chat 响应格式
func (c *ChatController) toChatResponse(adapterResp *adapter.ChatResponse, requestID, modelName string, created time.Time) model.ChatResponse {
	choices := make([]model.ChatChoice, len(adapterResp.Choices))
	for i, choice := range adapterResp.Choices {
		choices[i] = model.ChatChoice{
			Index: choice.Index,
			Message: model.ChatMessage{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: choice.FinishReason,
		}
	}

	return model.ChatResponse{
		ID:      requestID,
		Object:  "chat.completion",
		Created: created.Unix(),
		Model:   modelName,
		Choices: choices,
		Usage: model.ChatUsage{
			PromptTokens:     adapterResp.Usage.PromptTokens,
			CompletionTokens: adapterResp.Usage.CompletionTokens,
			TotalTokens:      adapterResp.Usage.TotalTokens,
		},
	}
}

// toStreamChunk 转换为流式响应块
func (c *ChatController) toStreamChunk(chunk *adapter.ChatStreamChunk, requestID, modelName string) model.ChatStreamResponse {
	choices := make([]model.ChatStreamChoice, len(chunk.Choices))
	for i, choice := range chunk.Choices {
		choices[i] = model.ChatStreamChoice{
			Index: choice.Index,
			Delta: model.ChatMessageDelta{
				Role:    choice.Delta.Role,
				Content: choice.Delta.Content,
			},
			FinishReason: choice.FinishReason,
		}
	}

	return model.ChatStreamResponse{
		ID:      requestID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   modelName,
		Choices: choices,
	}
}

// logRequest 记录请求日志
func (c *ChatController) logRequest(ctx *gin.Context, requestID string, req *model.ChatRequest, modelInfo *service.ModelInfo, usage *model.ChatUsage, latencyMs int64, status, errorMsg string) {
	apiKeyMasked, _ := ctx.Get("api_key_masked")

	log := model.ChatLog{
		RequestID:        requestID,
		APIKey:           fmt.Sprintf("%v", apiKeyMasked),
		Model:            req.Model,
		ProviderName:     modelInfo.ProviderName,
		ProviderType:     modelInfo.Provider.Type(),
		ModelName:        modelInfo.ModelName,
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
		LatencyMs:        latencyMs,
		Status:           status,
		Error:            errorMsg,
		Timestamp:        time.Now(),
	}

	// JSON 格式输出日志
	logData, _ := json.Marshal(log)
	gin.DefaultWriter.Write(logData)
	gin.DefaultWriter.Write([]byte("\n"))
}
