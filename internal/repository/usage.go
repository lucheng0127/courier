package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lucheng0127/courier/internal/model"
)

// UsageRepository 使用记录数据访问接口
type UsageRepository interface {
	// CreateUsageRecord 创建使用记录
	CreateUsageRecord(ctx context.Context, record *model.UsageRecord) error

	// QueryUsageByUserAndTimeRange 查询用户在指定时间范围的使用记录
	QueryUsageByUserAndTimeRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*model.UsageRecord, error)

	// AggregateUsageByDay 按天聚合使用统计
	AggregateUsageByDay(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*model.DailyStatsRow, error)

	// AggregateUsageByModel 按模型聚合使用统计
	AggregateUsageByModel(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*model.ModelStatsRow, error)

	// GetUsageSummary 获取使用汇总统计
	GetUsageSummary(ctx context.Context, userID int64, startDate, endDate time.Time) (*model.SummaryRow, error)
}

// usageRepository 使用记录数据访问实现
type usageRepository struct {
	db *sqlx.DB
}

// NewUsageRepository 创建 Usage Repository
func NewUsageRepository(db *sqlx.DB) UsageRepository {
	return &usageRepository{db: db}
}

// CreateUsageRecord 创建使用记录
func (r *usageRepository) CreateUsageRecord(ctx context.Context, record *model.UsageRecord) error {
	query := `
		INSERT INTO usage_records (
			user_id, api_key_id, request_id, trace_id, model, provider_name,
			prompt_tokens, completion_tokens, total_tokens, latency_ms, status, error_type
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, timestamp
	`
	err := r.db.QueryRowContext(ctx, query,
		record.UserID,
		record.APIKeyID,
		record.RequestID,
		record.TraceID,
		record.Model,
		record.ProviderName,
		record.PromptTokens,
		record.CompletionTokens,
		record.TotalTokens,
		record.LatencyMs,
		record.Status,
		record.ErrorType,
	).Scan(&record.ID, &record.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to create usage record: %w", err)
	}
	return nil
}

// QueryUsageByUserAndTimeRange 查询用户在指定时间范围的使用记录
func (r *usageRepository) QueryUsageByUserAndTimeRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*model.UsageRecord, error) {
	var records []*model.UsageRecord
	query := `
		SELECT id, user_id, api_key_id, request_id, trace_id, model, provider_name,
			prompt_tokens, completion_tokens, total_tokens, latency_ms, status, error_type, timestamp
		FROM usage_records
		WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3
		ORDER BY timestamp DESC
	`
	err := r.db.SelectContext(ctx, &records, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query usage records: %w", err)
	}
	return records, nil
}

// AggregateUsageByDay 按天聚合使用统计
func (r *usageRepository) AggregateUsageByDay(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*model.DailyStatsRow, error) {
	var rows []*model.DailyStatsRow
	query := `
		SELECT
			DATE(timestamp) as date,
			COUNT(*) as total_requests,
			COALESCE(SUM(total_tokens), 0) as total_tokens,
			COALESCE(SUM(prompt_tokens), 0) as total_prompt_tokens,
			COALESCE(SUM(completion_tokens), 0) as total_completion_tokens,
			COALESCE(AVG(latency_ms), 0) as average_latency_ms
		FROM usage_records
		WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3
		GROUP BY DATE(timestamp)
		ORDER BY date DESC
	`
	err := r.db.SelectContext(ctx, &rows, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate usage by day: %w", err)
	}
	return rows, nil
}

// AggregateUsageByModel 按模型聚合使用统计
func (r *usageRepository) AggregateUsageByModel(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*model.ModelStatsRow, error) {
	var rows []*model.ModelStatsRow
	query := `
		SELECT
			model,
			COUNT(*) as total_requests,
			COALESCE(SUM(total_tokens), 0) as total_tokens,
			COALESCE(SUM(prompt_tokens), 0) as total_prompt_tokens,
			COALESCE(SUM(completion_tokens), 0) as total_completion_tokens,
			COALESCE(AVG(latency_ms), 0) as average_latency_ms
		FROM usage_records
		WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3
		GROUP BY model
		ORDER BY total_tokens DESC
	`
	err := r.db.SelectContext(ctx, &rows, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate usage by model: %w", err)
	}
	return rows, nil
}

// GetUsageSummary 获取使用汇总统计
func (r *usageRepository) GetUsageSummary(ctx context.Context, userID int64, startDate, endDate time.Time) (*model.SummaryRow, error) {
	var row model.SummaryRow
	query := `
		SELECT
			COUNT(*) as total_requests,
			COALESCE(SUM(total_tokens), 0) as total_tokens,
			COALESCE(SUM(prompt_tokens), 0) as total_prompt_tokens,
			COALESCE(SUM(completion_tokens), 0) as total_completion_tokens,
			COALESCE(AVG(latency_ms), 0) as average_latency_ms
		FROM usage_records
		WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3
	`
	err := r.db.GetContext(ctx, &row, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage summary: %w", err)
	}
	return &row, nil
}
