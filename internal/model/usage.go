package model

import "time"

// UsageRecord 使用记录模型
type UsageRecord struct {
	ID               int64     `json:"id" db:"id" gorm:"primaryKey"`
	UserID           int64     `json:"user_id" db:"user_id" gorm:"index"`
	APIKeyID         *int64    `json:"api_key_id,omitempty" db:"api_key_id" gorm:"index"`
	RequestID        string    `json:"request_id" db:"request_id" gorm:"index"`
	TraceID          string    `json:"trace_id" db:"trace_id" gorm:"index"`
	Model            string    `json:"model" db:"model" gorm:"index"`
	ProviderName     string    `json:"provider_name" db:"provider_name" gorm:"index"`
	PromptTokens     int       `json:"prompt_tokens" db:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens" db:"completion_tokens"`
	TotalTokens      int       `json:"total_tokens" db:"total_tokens"`
	LatencyMs        int64     `json:"latency_ms" db:"latency_ms"`
	Status           string    `json:"status" db:"status" gorm:"index"` // success, error
	ErrorType        string    `json:"error_type,omitempty" db:"error_type"`
	Timestamp        time.Time `json:"timestamp" db:"timestamp" gorm:"autoCreateTime;default:NOW()"`
}

// TableName 指定表名
func (UsageRecord) TableName() string {
	return "usage_records"
}

// UsageStatsRequest 使用统计查询请求
type UsageStatsRequest struct {
	UserID    int64     `form:"user_id"`
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
	GroupBy   string    `form:"group_by"` // day, model
}

// UsageStatsResponse 使用统计响应
type UsageStatsResponse struct {
	UserID        int64              `json:"user_id"`
	Period        TimePeriod         `json:"period"`
	Summary       UsageSummary       `json:"summary"`
	DailyBreakdown []DailyUsageStats `json:"daily_breakdown,omitempty"`
	ModelBreakdown []ModelUsageStats `json:"model_breakdown,omitempty"`
}

// TimePeriod 时间范围
type TimePeriod struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// UsageSummary 使用汇总
type UsageSummary struct {
	TotalRequests         int     `json:"total_requests"`
	TotalTokens           int64   `json:"total_tokens"`
	TotalPromptTokens     int64   `json:"total_prompt_tokens"`
	TotalCompletionTokens int64   `json:"total_completion_tokens"`
	AverageLatencyMs      float64 `json:"average_latency_ms"`
}

// DailyUsageStats 按天统计
type DailyUsageStats struct {
	Date               string  `json:"date"`
	Requests           int     `json:"requests"`
	Tokens             int64   `json:"tokens"`
	PromptTokens       int64   `json:"prompt_tokens"`
	CompletionTokens   int64   `json:"completion_tokens"`
	AverageLatencyMs   float64 `json:"average_latency_ms"`
}

// ModelUsageStats 按模型统计
type ModelUsageStats struct {
	Model             string  `json:"model"`
	Requests          int     `json:"requests"`
	Tokens            int64   `json:"tokens"`
	PromptTokens      int64   `json:"prompt_tokens"`
	CompletionTokens  int64   `json:"completion_tokens"`
	AverageLatencyMs  float64 `json:"average_latency_ms"`
}

// DailyStatsRow 数据库查询结果（按天统计）
type DailyStatsRow struct {
	Date               string `db:"date"`
	TotalRequests      int    `db:"total_requests"`
	TotalTokens        int64  `db:"total_tokens"`
	TotalPromptTokens  int64  `db:"total_prompt_tokens"`
	TotalCompletionTokens int64 `db:"total_completion_tokens"`
	AverageLatencyMs   float64 `db:"average_latency_ms"`
}

// ModelStatsRow 数据库查询结果（按模型统计）
type ModelStatsRow struct {
	Model             string  `db:"model"`
	TotalRequests     int     `db:"total_requests"`
	TotalTokens       int64   `db:"total_tokens"`
	TotalPromptTokens int64   `db:"total_prompt_tokens"`
	TotalCompletionTokens int64 `db:"total_completion_tokens"`
	AverageLatencyMs  float64 `db:"average_latency_ms"`
}

// SummaryRow 数据库查询结果（汇总统计）
type SummaryRow struct {
	TotalRequests         int     `db:"total_requests"`
	TotalTokens           int64   `db:"total_tokens"`
	TotalPromptTokens     int64   `db:"total_prompt_tokens"`
	TotalCompletionTokens int64   `db:"total_completion_tokens"`
	AverageLatencyMs      float64 `db:"average_latency_ms"`
}
