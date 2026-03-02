package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

// AttemptDetail 单次尝试详情
type AttemptDetail struct {
	ModelName string        `json:"model_name"`
	Error     error         `json:"-"`
	ErrorType string        `json:"error_type"`
	Duration  time.Duration `json:"duration_ms"`
}

// RetryResult 重试结果
type RetryResult struct {
	Success         bool             `json:"success"`
	FallbackCount   int              `json:"fallback_count"`
	FinalModelName  string           `json:"final_model_name"`
	AttemptDetails  []AttemptDetail  `json:"attempt_details"`
	TotalDuration   time.Duration    `json:"total_duration_ms"`
	Response        any              `json:"-"` // 成功时的响应
}

// RetryableFunc 可重试的函数类型
type RetryableFunc func(ctx context.Context, modelName string) (any, error)

// RetryService 重试服务
type RetryService struct{}

// NewRetryService 创建重试服务
func NewRetryService() *RetryService {
	return &RetryService{}
}

// IsRetryableError 判断错误是否可重试
func (s *RetryService) IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// 网络错误
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	// 系统调用错误（如连接拒绝）
	var sysErr *os.SyscallError
	if errors.As(err, &sysErr) {
		if sysErr.Err == syscall.ECONNREFUSED {
			return true
		}
	}

	// 超时错误
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// 检查错误消息内容
	errMsg := err.Error()
	// 5xx 服务器错误
	if containsStatusCode(errMsg, 500, 502, 503, 504) {
		return true
	}

	// 连接失败、超时、DNS 解析失败等
	if containsSubstring(errMsg, "timeout", "connection refused", "connection reset", "dns", "no such host") {
		return true
	}

	return false
}

// RetryWithFallback 带 Fallback 的重试逻辑
func (s *RetryService) RetryWithFallback(
	ctx context.Context,
	fallbackModels []string,
	retryableFunc RetryableFunc,
) (*RetryResult, error) {
	startTime := time.Now()
	result := &RetryResult{
		AttemptDetails: make([]AttemptDetail, 0, len(fallbackModels)),
	}

	// 如果没有 Fallback 模型列表，直接执行
	if len(fallbackModels) == 0 {
		return nil, errors.New("no fallback models provided")
	}

	// 依次尝试每个模型
	for i, modelName := range fallbackModels {
		attemptStart := time.Now()
		detail := AttemptDetail{
			ModelName: modelName,
		}

		// 执行函数
		resp, err := retryableFunc(ctx, modelName)
		detail.Duration = time.Since(attemptStart)

		if err == nil {
			// 成功
			detail.ErrorType = ""
			result.Success = true
			result.FallbackCount = i
			result.FinalModelName = modelName
			result.Response = resp
			result.AttemptDetails = append(result.AttemptDetails, detail)
			result.TotalDuration = time.Since(startTime)
			return result, nil
		}

		// 失败
		detail.Error = err
		detail.ErrorType = classifyError(err)
		result.AttemptDetails = append(result.AttemptDetails, detail)

		// 判断是否可重试
		if !s.IsRetryableError(err) {
			// 不可重试错误，直接返回
			result.TotalDuration = time.Since(startTime)
			return result, fmt.Errorf("non-retryable error with model %s: %w", modelName, err)
		}

		// 记录日志
		log.Printf("[WARN] Model %s failed (%s), trying next fallback model", modelName, err)
	}

	// 所有模型都失败
	result.TotalDuration = time.Since(startTime)
	return result, fmt.Errorf("all models failed after %d attempts", len(fallbackModels))
}

// containsStatusCode 检查错误消息是否包含指定的 HTTP 状态码
func containsStatusCode(msg string, codes ...int) bool {
	for _, code := range codes {
		if strings.Contains(msg, fmt.Sprintf("%d", code)) {
			return true
		}
	}
	return false
}

// containsSubstring 检查字符串是否包含子串（忽略大小写）
func containsSubstring(s string, substrs ...string) bool {
	sLower := strings.ToLower(s)
	for _, sub := range substrs {
		if strings.Contains(sLower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

// classifyError 分类错误类型
func classifyError(err error) string {
	if err == nil {
		return "unknown"
	}

	errMsg := strings.ToLower(err.Error())

	// 超时错误
	if errors.Is(err, context.DeadlineExceeded) || strings.Contains(errMsg, "timeout") {
		return "timeout"
	}

	// 连接错误
	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "connection reset") || strings.Contains(errMsg, "econnrefused") {
		return "connection_error"
	}

	// DNS 错误
	if strings.Contains(errMsg, "no such host") || strings.Contains(errMsg, "dns") {
		return "dns_error"
	}

	// 5xx 错误
	if containsStatusCode(errMsg, 500, 502, 503, 504) {
		return "server_error"
	}

	// 4xx 错误
	if containsStatusCode(errMsg, 400, 401, 403, 404) {
		return "client_error"
	}

	return "unknown"
}
