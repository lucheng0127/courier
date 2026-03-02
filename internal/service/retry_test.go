package service

import (
	"context"
	"errors"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

// TestIsRetryableError 测试错误分类
func TestIsRetryableError(t *testing.T) {
	svc := NewRetryService()

	tests := []struct {
		name     string
		err      error
		retryable bool
	}{
		{
			name:     "网络错误",
			err:      &net.OpError{Op: "dial", Err: errors.New("connection refused")},
			retryable: true,
		},
		{
			name:     "超时错误",
			err:      context.DeadlineExceeded,
			retryable: true,
		},
		{
			name:     "系统调用错误 - ECONNREFUSED",
			err:      &net.OpError{Err: &os.SyscallError{Err: syscall.ECONNREFUSED}},
			retryable: true,
		},
		{
			name:     "5xx 错误消息",
			err:      errors.New("HTTP 500: Internal Server Error"),
			retryable: true,
		},
		{
			name:     "timeout 错误消息",
			err:      errors.New("request timeout after 30s"),
			retryable: true,
		},
		{
			name:     "4xx 客户端错误",
			err:      errors.New("HTTP 400: Bad Request"),
			retryable: false,
		},
		{
			name:     "认证失败",
			err:      errors.New("authentication failed: invalid API key"),
			retryable: false,
		},
		{
			name:     "nil 错误",
			err:      nil,
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.IsRetryableError(tt.err)
			if result != tt.retryable {
				t.Errorf("IsRetryableError() = %v, want %v", result, tt.retryable)
			}
		})
	}
}

// TestRetryWithFallback_Success 测试 Fallback 成功场景
func TestRetryWithFallback_Success(t *testing.T) {
	svc := NewRetryService()
	ctx := context.Background()

	// 模拟：第一个模型失败，第二个成功
	attemptCount := 0
	mockFunc := func(ctx context.Context, modelName string) (any, error) {
		attemptCount++
		if modelName == "model-1" {
			return nil, errors.New("timeout")
		}
		return "success", nil
	}

	result, err := svc.RetryWithFallback(ctx, []string{"model-1", "model-2"}, mockFunc)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Success {
		t.Error("expected success")
	}

	if result.FallbackCount != 1 {
		t.Errorf("expected FallbackCount = 1, got %d", result.FallbackCount)
	}

	if result.FinalModelName != "model-2" {
		t.Errorf("expected FinalModelName = model-2, got %s", result.FinalModelName)
	}

	if attemptCount != 2 {
		t.Errorf("expected 2 attempts, got %d", attemptCount)
	}
}

// TestRetryWithFallback_AllFailed 测试所有模型都失败
func TestRetryWithFallback_AllFailed(t *testing.T) {
	svc := NewRetryService()
	ctx := context.Background()

	mockFunc := func(ctx context.Context, modelName string) (any, error) {
		return nil, errors.New("timeout")
	}

	result, err := svc.RetryWithFallback(ctx, []string{"model-1", "model-2", "model-3"}, mockFunc)

	if err == nil {
		t.Fatal("expected error when all models fail")
	}

	if result.Success {
		t.Error("expected failure")
	}

	if len(result.AttemptDetails) != 3 {
		t.Errorf("expected 3 attempt details, got %d", len(result.AttemptDetails))
	}

	// 验证所有尝试都记录了错误
	for i, detail := range result.AttemptDetails {
		if detail.Error == nil {
			t.Errorf("attempt %d: expected error", i)
		}
	}
}

// TestRetryWithFallback_NonRetryableError 测试不可重试错误
func TestRetryWithFallback_NonRetryableError(t *testing.T) {
	svc := NewRetryService()
	ctx := context.Background()

	attemptCount := 0
	mockFunc := func(ctx context.Context, modelName string) (any, error) {
		attemptCount++
		// 第一个模型返回不可重试错误
		return nil, errors.New("HTTP 400: Bad Request")
	}

	result, err := svc.RetryWithFallback(ctx, []string{"model-1", "model-2", "model-3"}, mockFunc)

	if err == nil {
		t.Fatal("expected error for non-retryable error")
	}

	// 不可重试错误应该只尝试一次
	if attemptCount != 1 {
		t.Errorf("expected 1 attempt for non-retryable error, got %d", attemptCount)
	}

	if len(result.AttemptDetails) != 1 {
		t.Errorf("expected 1 attempt detail, got %d", len(result.AttemptDetails))
	}
}

// TestRetryWithFallback_Timeout 测试超时取消
func TestRetryWithFallback_Timeout(t *testing.T) {
	svc := NewRetryService()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	mockFunc := func(ctx context.Context, modelName string) (any, error) {
		// 检查 context 是否已取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return nil, nil
		}
	}

	result, err := svc.RetryWithFallback(ctx, []string{"model-1", "model-2"}, mockFunc)

	if err == nil {
		t.Fatal("expected timeout error")
	}

	if result.TotalDuration < 100*time.Millisecond {
		t.Error("expected total duration >= timeout")
	}
}

// TestRetryWithFallback_FirstSuccess 测试第一次就成功
func TestRetryWithFallback_FirstSuccess(t *testing.T) {
	svc := NewRetryService()
	ctx := context.Background()

	attemptCount := 0
	mockFunc := func(ctx context.Context, modelName string) (any, error) {
		attemptCount++
		return "success", nil
	}

	result, err := svc.RetryWithFallback(ctx, []string{"model-1", "model-2", "model-3"}, mockFunc)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.FallbackCount != 0 {
		t.Errorf("expected FallbackCount = 0, got %d", result.FallbackCount)
	}

	if result.FinalModelName != "model-1" {
		t.Errorf("expected FinalModelName = model-1, got %s", result.FinalModelName)
	}

	if attemptCount != 1 {
		t.Errorf("expected 1 attempt, got %d", attemptCount)
	}
}

// TestClassifyError 测试错误分类
func TestClassifyError(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{context.DeadlineExceeded, "timeout"},
		{errors.New("request timeout"), "timeout"},
		{errors.New("connection refused"), "connection_error"},
		{errors.New("no such host"), "dns_error"},
		{errors.New("HTTP 500"), "server_error"},
		{errors.New("HTTP 400"), "client_error"},
		{errors.New("unknown error"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			result := classifyError(tt.err)
			if result != tt.expected {
				t.Errorf("classifyError() = %s, want %s", result, tt.expected)
			}
		})
	}
}

// BenchmarkRetryWithFallback 基准测试
func BenchmarkRetryWithFallback(b *testing.B) {
	svc := NewRetryService()
	ctx := context.Background()

	mockFunc := func(ctx context.Context, modelName string) (any, error) {
		return "success", nil
	}

	models := []string{"model-1", "model-2", "model-3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.RetryWithFallback(ctx, models, mockFunc)
	}
}
