package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// UsageController 使用统计控制器
type UsageController struct {
	usageService *service.UsageService
}

// NewUsageController 创建 Usage Controller
func NewUsageController(usageService *service.UsageService) *UsageController {
	return &UsageController{
		usageService: usageService,
	}
}

// GetUsageStats 查询使用统计
// GET /v1/usage?user_id=<id>&start_date=<date>&end_date=<date>&group_by=<field>
func (c *UsageController) GetUsageStats(ctx *gin.Context) {
	var req model.UsageStatsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 设置默认的 group_by 为 day
	if req.GroupBy == "" {
		req.GroupBy = "day"
	}

	// 验证 group_by 值
	if req.GroupBy != "day" && req.GroupBy != "model" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid group_by parameter. Must be 'day' or 'model'",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	// 解析时间参数
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
	 startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": "Invalid start_date format. Use RFC3339 format",
					"type":    "invalid_request_error",
				},
			})
			return
		}
		req.StartDate = &startDate
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": "Invalid end_date format. Use RFC3339 format",
					"type":    "invalid_request_error",
				},
			})
			return
		}
		req.EndDate = &endDate
	}

	stats, err := c.usageService.GetUsageStats(ctx, &req)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"message": "User not found",
					"type":    "invalid_request_error",
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to get usage stats",
				"type":    "api_error",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}
