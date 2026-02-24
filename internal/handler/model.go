package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/service"
)

// ModelHandler 模型处理器
type ModelHandler struct {
	modelService service.ModelService
}

// NewModelHandler 创建模型处理器
func NewModelHandler(modelService service.ModelService) *ModelHandler {
	return &ModelHandler{modelService: modelService}
}

// ListModels 查询可用模型列表
func (h *ModelHandler) ListModels(c *gin.Context) {
	models := h.modelService.GetModels()
	c.JSON(http.StatusOK, models)
}
