package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// APIKeyHandler API Key 处理器
type APIKeyHandler struct {
	apiKeyService service.APIKeyService
}

// NewAPIKeyHandler 创建 API Key 处理器
func NewAPIKeyHandler(apiKeyService service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{apiKeyService: apiKeyService}
}

// GenerateAPIKey 生成 API Key
func (h *APIKeyHandler) GenerateAPIKey(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数无效"})
		return
	}

	apiKey, err := h.apiKeyService.GenerateAPIKey(uint(userID))
	if err != nil {
		if err.Error() == "用户不存在" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "生成 API Key 失败"})
		return
	}

	c.JSON(http.StatusCreated, toAPIKeyResponse(apiKey))
}

// ListAPIKeys 查询 API Key 列表
func (h *APIKeyHandler) ListAPIKeys(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数无效"})
		return
	}

	apiKeys, err := h.apiKeyService.ListAPIKeys(uint(userID))
	if err != nil {
		if err.Error() == "用户不存在" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "查询 API Key 列表失败"})
		return
	}

	response := make([]APIKeyResponse, 0, len(apiKeys))
	for _, apiKey := range apiKeys {
		response = append(response, toAPIKeyResponse(&apiKey))
	}

	c.JSON(http.StatusOK, response)
}

// DeleteAPIKey 删除 API Key
func (h *APIKeyHandler) DeleteAPIKey(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数无效"})
		return
	}

	keyIDStr := c.Param("keyid")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数无效"})
		return
	}

	err = h.apiKeyService.DeleteAPIKey(uint(userID), uint(keyID))
	if err != nil {
		if err.Error() == "用户不存在" || err.Error() == "API Key 不存在" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "删除 API Key 失败"})
		return
	}

	c.Status(http.StatusNoContent)
}

// DisableAPIKey 禁用 API Key
func (h *APIKeyHandler) DisableAPIKey(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数无效"})
		return
	}

	keyIDStr := c.Param("keyid")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "参数无效"})
		return
	}

	apiKey, err := h.apiKeyService.DisableAPIKey(uint(userID), uint(keyID))
	if err != nil {
		if err.Error() == "用户不存在" || err.Error() == "API Key 不存在" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "禁用 API Key 失败"})
		return
	}

	c.JSON(http.StatusOK, toAPIKeyResponse(apiKey))
}

func toAPIKeyResponse(apiKey *model.APIKey) APIKeyResponse {
	resp := APIKeyResponse{
		ID:         apiKey.ID,
		UserID:     apiKey.UserID,
		Key:        apiKey.Key,
		Status:     apiKey.Status,
		CreatedAt:  apiKey.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  apiKey.UpdatedAt.Format(time.RFC3339),
	}

	if apiKey.LastUsedAt != nil {
		resp.LastUsedAt = apiKey.LastUsedAt.Format(time.RFC3339)
	}
	if apiKey.ExpiresAt != nil {
		resp.ExpiresAt = apiKey.ExpiresAt.Format(time.RFC3339)
	}

	return resp
}
