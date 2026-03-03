package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/lucheng0127/courier/internal/logger"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/repository"
)

// UsageService 使用统计服务
type UsageService struct {
	usageRepo repository.UsageRepository
	userRepo  repository.UserRepository
	recordCh  chan *model.UsageRecord
	wg        sync.WaitGroup
	once      sync.Once
	stopCh    chan struct{}
}

// NewUsageService 创建 Usage Service
func NewUsageService(usageRepo repository.UsageRepository, userRepo repository.UserRepository) *UsageService {
	s := &UsageService{
		usageRepo: usageRepo,
		userRepo:  userRepo,
		recordCh:  make(chan *model.UsageRecord, 1000),
		stopCh:    make(chan struct{}),
	}
	s.startBackgroundWorkers()
	return s
}

// startBackgroundWorkers 启动后台处理协程
func (s *UsageService) startBackgroundWorkers() {
	s.once.Do(func() {
		// 启动 3 个处理协程
		for i := 0; i < 3; i++ {
			s.wg.Add(1)
			go s.processRecords()
		}
	})
}

// processRecords 处理使用记录
func (s *UsageService) processRecords() {
	defer s.wg.Done()

	batchSize := 50
	batch := make([]*model.UsageRecord, 0, batchSize)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// 批量写入
		for _, record := range batch {
			if err := s.usageRepo.CreateUsageRecord(ctx, record); err != nil {
				logger.L.Error("Failed to create usage record",
					zap.String("request_id", record.RequestID),
					zap.Int64("user_id", record.UserID),
					zap.Error(err))
			}
		}
		batch = batch[:0]
	}

	for {
		select {
		case record := <-s.recordCh:
			batch = append(batch, record)
			if len(batch) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-s.stopCh:
			flush()
			return
		}
	}
}

// RecordUsage 记录使用量（异步）
func (s *UsageService) RecordUsage(ctx context.Context, record *model.UsageRecord) error {
	select {
	case s.recordCh <- record:
		return nil
	default:
		// channel 满了，同步写入
		if err := s.usageRepo.CreateUsageRecord(ctx, record); err != nil {
			return fmt.Errorf("failed to record usage: %w", err)
		}
		return nil
	}
}

// GetUsageStats 获取使用统计
func (s *UsageService) GetUsageStats(ctx context.Context, req *model.UsageStatsRequest) (*model.UsageStatsResponse, error) {
	// 验证用户存在
	_, err := s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 设置默认时间范围（最近 30 天）
	startDate := req.StartDate
	if startDate == nil {
		t := time.Now().AddDate(0, 0, -30)
		startDate = &t
	}

	endDate := req.EndDate
	if endDate == nil {
		t := time.Now()
		endDate = &t
	}

	// 获取汇总统计
	summary, err := s.usageRepo.GetUsageSummary(ctx, req.UserID, *startDate, *endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage summary: %w", err)
	}

	response := &model.UsageStatsResponse{
		UserID: req.UserID,
		Period: model.TimePeriod{
			Start: *startDate,
			End:   *endDate,
		},
		Summary: model.UsageSummary{
			TotalRequests:         summary.TotalRequests,
			TotalTokens:           summary.TotalTokens,
			TotalPromptTokens:     summary.TotalPromptTokens,
			TotalCompletionTokens: summary.TotalCompletionTokens,
			AverageLatencyMs:      summary.AverageLatencyMs,
		},
	}

	// 根据 group_by 参数获取详细数据
	switch req.GroupBy {
	case "model":
		modelStats, err := s.usageRepo.AggregateUsageByModel(ctx, req.UserID, *startDate, *endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get model stats: %w", err)
		}
		response.ModelBreakdown = make([]model.ModelUsageStats, len(modelStats))
		for i, row := range modelStats {
			response.ModelBreakdown[i] = model.ModelUsageStats{
				Model:            row.Model,
				Requests:         row.TotalRequests,
				Tokens:           row.TotalTokens,
				PromptTokens:     row.TotalPromptTokens,
				CompletionTokens: row.TotalCompletionTokens,
				AverageLatencyMs: row.AverageLatencyMs,
			}
		}
	default: // 按天聚合
		dailyStats, err := s.usageRepo.AggregateUsageByDay(ctx, req.UserID, *startDate, *endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get daily stats: %w", err)
		}
		response.DailyBreakdown = make([]model.DailyUsageStats, len(dailyStats))
		for i, row := range dailyStats {
			response.DailyBreakdown[i] = model.DailyUsageStats{
				Date:             row.Date,
				Requests:         row.TotalRequests,
				Tokens:           row.TotalTokens,
				PromptTokens:     row.TotalPromptTokens,
				CompletionTokens: row.TotalCompletionTokens,
				AverageLatencyMs: row.AverageLatencyMs,
			}
		}
	}

	return response, nil
}

// Close 关闭服务
func (s *UsageService) Close() error {
	close(s.stopCh)
	s.wg.Wait()
	return nil
}
