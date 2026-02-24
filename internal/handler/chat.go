package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/client"
	"github.com/lucheng0127/courier/internal/middleware"
	"github.com/lucheng0127/courier/internal/service"
)

// ChatHandler 聊天处理器
type ChatHandler struct {
	chatService service.ChatService
}

// NewChatHandler 创建聊天处理器
func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Messages []ChatMessage `json:"messages" binding:"required"`
	Stream   bool          `json:"stream"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// Chat 发起模型对话
func (h *ChatHandler) Chat(c *gin.Context) {
	// 获取模型名称
	modelName := c.Param("model")

	// 获取用户信息
	_, ok := middleware.GetUserInfo(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	userID, _ := middleware.GetUserID(c)

	// 解析请求
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 验证消息
	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "messages 不能为空"})
		return
	}

	// 转换为服务请求
	chatReq := &client.ChatRequest{
		Messages: make([]client.ChatMessage, len(req.Messages)),
		Stream:   req.Stream,
	}
	for i, msg := range req.Messages {
		chatReq.Messages[i] = client.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 判断是否为流式请求
	if req.Stream {
		h.handleStreamChat(c, userID, modelName, chatReq)
	} else {
		h.handleNormalChat(c, userID, modelName, chatReq)
	}
}

// handleNormalChat 处理非流式对话
func (h *ChatHandler) handleNormalChat(c *gin.Context, userID uint, modelName string, req *client.ChatRequest) {
	resp, err := h.chatService.Chat(c.Request.Context(), userID, modelName, req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "上游模型请求失败"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// handleStreamChat 处理流式对话
func (h *ChatHandler) handleStreamChat(c *gin.Context, userID uint, modelName string, req *client.ChatRequest) {
	// 设置 SSE 响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 调用流式服务
	chunkChan, errChan, _ := h.chatService.ChatStream(c.Request.Context(), userID, modelName, req)

	// 刷新器
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "不支持流式响应"})
		return
	}

	// 发送流式响应
	for {
		select {
		case chunk, ok := <-chunkChan:
			if !ok {
				// 发送结束标记
				c.SSEvent("", "[DONE]")
				flusher.Flush()
				return
			}

			// 发送数据块
			c.SSEvent("", chunk)
			flusher.Flush()

		case err, ok := <-errChan:
			if !ok {
				return
			}
			// 发送错误
			c.SSEvent("", gin.H{"error": err.Error()})
			flusher.Flush()
			return

		case <-c.Request.Context().Done():
			return
		}
	}
}
