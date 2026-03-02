package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucheng0127/courier/internal/adapter"
	"github.com/lucheng0127/courier/internal/logger"
	"github.com/lucheng0127/courier/internal/middleware"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// ChatController Chat API 控制器
type ChatController struct {
	router       *service.RouterService
	retrySvc     *service.RetryService
	usageService *service.UsageService
}

// NewChatController 创建 Chat 控制器
func NewChatController(router *service.RouterService, usageService *service.UsageService) *ChatController {
	return &ChatController{
		router:       router,
		retrySvc:     service.NewRetryService(),
		usageService: usageService,
	}
}

// ChatCompletions Chat Completions 端点
// POST /v1/chat/completions
func (c *ChatController) ChatCompletions(ctx *gin.Context) {
	startTime := time.Now()
	traceID := middleware.GetTraceID(ctx)

	// 解析请求
	var req model.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request format", map[string]any{
			"trace_id": traceID,
			"error":    err.Error(),
		})
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
		c.handleModelError(ctx, err, traceID)
		return
	}

	// 生成请求 ID
	requestID := "chatcmpl-" + uuid.New().String()

	// 获取 Fallback 模型列表
	fallbackModels := c.getFallbackModels(ctx, modelInfo)

	// 设置超时（默认 30 秒，可从 Provider 配置读取）
	timeout := time.Duration(modelInfo.Provider.Timeout()) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	// 创建带超时的上下文
	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeout)
	defer cancel()

	// 使用重试服务处理请求
	result, err := c.retrySvc.RetryWithFallback(timeoutCtx, fallbackModels, func(ctx context.Context, modelName string) (any, error) {
		return c.callProvider(ctx, &req, modelInfo.ProviderName, modelName, requestID)
	})

	// 记录日志
	c.logRequestWithRetry(ctx, requestID, &req, modelInfo, result, err, time.Since(startTime).Milliseconds())

	if err != nil {
		c.handleProviderError(ctx, err, result)
		return
	}

	// 处理响应
	if req.Stream {
		c.handleStreamResponse(ctx, result.Response, requestID, startTime)
	} else {
		resp := result.Response.(*model.ChatResponse)
		ctx.JSON(http.StatusOK, resp)
	}
}

// callProvider 调用 Provider
func (c *ChatController) callProvider(ctx context.Context, req *model.ChatRequest, providerName, modelName, requestID string) (any, error) {
	// 获取 Provider
	provider, err := c.router.ResolveProvider(providerName)
	if err != nil {
		return nil, err
	}

	// 转换请求格式
	adapterReq := c.toAdapterRequest(req, modelName)

	// 添加 TraceID 到请求 Header
	if traceID := ctx.Value("trace_id"); traceID != nil {
		// 这里需要 Adapter 支持传递 Header
		// 暂时通过 context 传递
		ctx = context.WithValue(ctx, "trace_id", traceID)
	}

	// 处理流式/非流式响应
	if req.Stream {
		chunks, err := provider.ChatStream(ctx, adapterReq)
		if err != nil {
			return nil, err
		}
		return chunks, nil
	} else {
		resp, err := provider.Chat(ctx, adapterReq)
		if err != nil {
			return nil, err
		}
		// 转换响应格式
		return c.toChatResponse(resp, requestID, providerName+"/"+modelName, time.Now()), nil
	}
}

// getFallbackModels 获取 Fallback 模型列表
func (c *ChatController) getFallbackModels(ctx *gin.Context, modelInfo *service.ModelInfo) []string {
	// 从 Provider 获取 Fallback 配置
	config := modelInfo.Provider.Config()
	if config == nil {
		// 没有 Fallback 配置，只使用当前模型
		return []string{modelInfo.ModelName}
	}

	// 获取 fallback_models 配置
	fallbackModelsRaw, ok := config["fallback_models"]
	if !ok {
		return []string{modelInfo.ModelName}
	}

	// 解析 Fallback 模型列表
	fallbackModels, ok := fallbackModelsRaw.([]string)
	if !ok || len(fallbackModels) == 0 {
		return []string{modelInfo.ModelName}
	}

	// 验证请求的模型是否在 Fallback 列表中
	found := false
	for _, m := range fallbackModels {
		if m == modelInfo.ModelName {
			found = true
			break
		}
	}

	if found {
		// 使用配置的 Fallback 列表
		return fallbackModels
	}

	// 请求的模型不在列表中，将当前模型作为第一个，然后是 Fallback 列表
	result := make([]string, 0, len(fallbackModels)+1)
	result = append(result, modelInfo.ModelName)
	result = append(result, fallbackModels...)
	return result
}

// handleModelError 处理模型解析错误
func (c *ChatController) handleModelError(ctx *gin.Context, err error, traceID string) {
	switch e := err.(type) {
	case *service.ModelFormatError:
		logger.Warn("Model format error", map[string]any{
			"trace_id": traceID,
			"error":    e.Error(),
		})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": e.Error(),
				"type":    "invalid_request_error",
			},
		})
	case *service.ProviderNotFoundError:
		logger.Warn("Provider not found", map[string]any{
			"trace_id":      traceID,
			"provider_name": e.ProviderName,
		})
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"message": e.Error(),
				"type":    "invalid_request_error",
			},
		})
	default:
		logger.Error("Failed to resolve model", map[string]any{
			"trace_id": traceID,
			"error":    err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to resolve model",
				"type":    "api_error",
			},
		})
	}
}

// handleProviderError 处理 Provider 错误
func (c *ChatController) handleProviderError(ctx *gin.Context, err error, result *service.RetryResult) {
	traceID := middleware.GetTraceID(ctx)

	if result != nil && len(result.AttemptDetails) > 0 {
		// Fallback 耗尽
		details := make([]gin.H, 0, len(result.AttemptDetails))
		for _, detail := range result.AttemptDetails {
			details = append(details, gin.H{
				"model":      detail.ModelName,
				"error_type": detail.ErrorType,
				"duration_ms": detail.Duration.Milliseconds(),
			})
		}

		logger.Error("All models failed", map[string]any{
			"trace_id":         traceID,
			"attempt_count":    len(result.AttemptDetails),
			"attempt_details":  result.AttemptDetails,
		})

		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"message": fmt.Sprintf("All models failed after %d attempts. Last error: %v", len(result.AttemptDetails), err),
				"type":    "service_unavailable",
				"details": details,
			},
		})
		return
	}

	// 其他错误
	logger.Error("Provider error", map[string]any{
		"trace_id": traceID,
		"error":    err.Error(),
	})

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"message": err.Error(),
			"type":    "api_error",
		},
	})
}

// handleNonStreamResponse 处理非流式响应（已合并到 callProvider）
// handleStreamResponse 处理流式响应
func (c *ChatController) handleStreamResponse(ctx *gin.Context, resp any, requestID string, startTime time.Time) {
	chunks := resp.(<-chan *adapter.ChatStreamChunk)

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
	for chunk := range chunks {
		// 检查客户端是否断开
		select {
		case <-ctx.Request.Context().Done():
			return
		default:
		}

		// 转换为 SSE 格式
		sseChunk := c.toStreamChunk(chunk, requestID)

		// 发送数据
		data, _ := json.Marshal(sseChunk)
		fmt.Fprintf(ctx.Writer, "data: %s\n", data)
		fmt.Fprint(ctx.Writer, "\n")
		flusher.Flush()
	}

	// 发送结束标记
	fmt.Fprint(ctx.Writer, "data: [DONE]\n\n")
	flusher.Flush()
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
func (c *ChatController) toChatResponse(adapterResp *adapter.ChatResponse, requestID, modelName string, created time.Time) *model.ChatResponse {
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

	return &model.ChatResponse{
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
func (c *ChatController) toStreamChunk(chunk *adapter.ChatStreamChunk, requestID string) model.ChatStreamResponse {
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
		Model:   chunk.Model,
		Choices: choices,
	}
}

// logRequestWithRetry 记录带重试信息的请求日志
func (c *ChatController) logRequestWithRetry(ctx *gin.Context, requestID string, req *model.ChatRequest, modelInfo *service.ModelInfo, result *service.RetryResult, err error, latencyMs int64) {
	traceID := middleware.GetTraceID(ctx)
	apiKeyMasked, _ := ctx.Get("api_key_masked")

	// 获取用户信息（从中间件注入）
	userID, hasUserID := ctx.Get("user_id")
	apiKeyID, hasAPIKeyID := ctx.Get("api_key_id")

	status := "success"
	errorMsg := ""

	if err != nil {
		status = "error"
		errorMsg = err.Error()
	}

	// 构建尝试详情
	attemptDetails := make([]model.AttemptDetail, 0)
	if result != nil {
		for _, detail := range result.AttemptDetails {
			attemptDetails = append(attemptDetails, model.AttemptDetail{
				ModelName: detail.ModelName,
				ErrorType: detail.ErrorType,
				DurationMs: detail.Duration.Milliseconds(),
			})
		}
	}

	log := model.ChatLog{
		RequestID:        requestID,
		TraceID:          traceID,
		APIKey:           fmt.Sprintf("%v", apiKeyMasked),
		Model:            req.Model,
		ProviderName:     modelInfo.ProviderName,
		ProviderType:     modelInfo.Provider.Type(),
		ModelName:        modelInfo.ModelName,
		FallbackCount:    0,
		FinalModelName:   modelInfo.ModelName,
		AttemptDetails:   attemptDetails,
		PromptTokens:     0,
		CompletionTokens: 0,
		TotalTokens:      0,
		LatencyMs:        latencyMs,
		Status:           status,
		Error:            errorMsg,
		Timestamp:        time.Now(),
	}

	if result != nil && result.Success {
		log.FallbackCount = result.FallbackCount
		log.FinalModelName = result.FinalModelName

		// 如果是响应中有 Usage 信息
		if resp, ok := result.Response.(*model.ChatResponse); ok {
			log.PromptTokens = resp.Usage.PromptTokens
			log.CompletionTokens = resp.Usage.CompletionTokens
			log.TotalTokens = resp.Usage.TotalTokens
		}
	}

	// 使用结构化日志
	if status == "success" {
		logger.Info("Chat request completed", map[string]any{
			"trace_id":        log.TraceID,
			"request_id":      log.RequestID,
			"model":           log.Model,
			"provider":        log.ProviderName,
			"fallback_count":  log.FallbackCount,
			"final_model":     log.FinalModelName,
			"latency_ms":      log.LatencyMs,
			"status":          log.Status,
		})
	} else {
		logger.Error("Chat request failed", map[string]any{
			"trace_id":        log.TraceID,
			"request_id":      log.RequestID,
			"model":           log.Model,
			"provider":        log.ProviderName,
			"attempt_count":   len(log.AttemptDetails),
			"error":           log.Error,
		})
	}

	// 如果有用户信息，记录使用量到数据库
	if hasUserID && hasAPIKeyID && c.usageService != nil {
		record := &model.UsageRecord{
			UserID:           userID.(int64),
			APIKeyID:         apiKeyID.(int64),
			RequestID:        requestID,
			TraceID:          traceID,
			Model:            req.Model,
			ProviderName:     modelInfo.ProviderName,
			PromptTokens:     log.PromptTokens,
			CompletionTokens: log.CompletionTokens,
			TotalTokens:      log.TotalTokens,
			LatencyMs:        latencyMs,
			Status:           status,
			ErrorType:        errorMsg,
		}

		// 异步记录使用量（使用独立 context）
		if err := c.usageService.RecordUsage(context.Background(), record); err != nil {
			logger.Error("Failed to record usage", map[string]any{
				"request_id": requestID,
				"user_id":    userID,
				"error":      err.Error(),
			})
		}
	}
}
