package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/middleware"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
)

// UsageController 使用统计控制器
type UsageController struct {
	usageSvc *service.UsageService
}

// NewUsageController 创建 Usage Controller
func NewUsageController(usageSvc *service.UsageService) *UsageController {
	return &UsageController{
		usageSvc: usageSvc,
	}
}

// RegisterRoutes 注册路由
func (c *UsageController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/usage", c.GetUsageStats)
	// 权限说明：
	// - 管理员：可查询任意用户或所有用户的统计（通过 user_id 参数过滤）
	// - 普通用户：只能查询自己的统计（自动过滤 user_id 参数）
}

// GetUsageStats 查询使用统计
// GET /api/v1/usage?user_id=<id>&start_date=<date>&end_date=<date>&group_by=<field>
// 权限：管理员可查询所有用户，普通用户只能查询自己
func (c *UsageController) GetUsageStats(ctx *gin.Context) {
	// 首先从上下文中获取用户信息
	userID, hasAuth := middleware.GetUserID(ctx)
	userRole, _ := middleware.GetUserRole(ctx)

	var req model.UsageStatsRequest
	// 手动绑定查询参数，不使用验证（因为 UserID 可能由上下文提供）
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			req.UserID = id
		}
	}

	// 权限检查：普通用户强制使用自己的 user_id，忽略传入的参数
	if hasAuth && userRole != "admin" {
		req.UserID = userID
	}

	// 如果 user_id 为 0（未设置），使用当前用户的 ID
	if req.UserID == 0 && hasAuth {
		req.UserID = userID
	}

	// 绑定 group_by 参数
	req.GroupBy = ctx.Query("group_by")

	// 设置默认的 group_by 为 day
	if req.GroupBy == "" {
		req.GroupBy = "day"
	}

	// 验证 group_by 值
	if req.GroupBy != "day" && req.GroupBy != "model" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid group_by parameter. Must be 'day' or 'model'",
			"type":    "invalid_request_error",
		})
		return
	}

	// 解析时间参数
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid start_date format. Use RFC3339 format",
				"type":    "invalid_request_error",
			})
			return
		}
		req.StartDate = &startDate
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid end_date format. Use RFC3339 format",
				"type":    "invalid_request_error",
			})
			return
		}
		req.EndDate = &endDate
	}

	stats, err := c.usageSvc.GetUsageStats(ctx, &req)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
				"type":    "invalid_request_error",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get usage stats",
			"type":    "api_error",
		})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}
