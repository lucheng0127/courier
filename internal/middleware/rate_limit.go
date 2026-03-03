package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time // IP -> 请求时间列表
	limit    int                    // 时间窗口内最大请求数
	window   time.Duration          // 时间窗口大小
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// 启动清理 goroutine，定期过期旧记录
	go rl.cleanup()

	return rl
}

// cleanup 定期清理过期的请求记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, times := range rl.requests {
			// 过滤掉时间窗口外的记录
			validTimes := make([]time.Time, 0, len(times))
			for _, t := range times {
				if now.Sub(t) < rl.window {
					validTimes = append(validTimes, t)
				}
			}
			if len(validTimes) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validTimes
			}
		}
		rl.mu.Unlock()
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// 获取该 IP 的请求记录
	times, exists := rl.requests[ip]
	if !exists {
		rl.requests[ip] = []time.Time{now}
		return true
	}

	// 过滤掉时间窗口外的记录
	validTimes := make([]time.Time, 0, len(times))
	for _, t := range times {
		if now.Sub(t) < rl.window {
			validTimes = append(validTimes, t)
		}
	}

	// 检查是否超过限制
	if len(validTimes) >= rl.limit {
		rl.requests[ip] = validTimes
		return false
	}

	// 添加当前请求
	rl.requests[ip] = append(validTimes, now)
	return true
}

// 注册速率限制器（全局实例）
var (
	registerLimiter     *RateLimiter
	registerLimiterOnce sync.Once
)

// getRegisterLimiter 获取注册速率限制器实例
// 同一 IP 每小时最多 5 次注册请求
func getRegisterLimiter() *RateLimiter {
	registerLimiterOnce.Do(func() {
		registerLimiter = NewRateLimiter(5, time.Hour)
	})
	return registerLimiter
}

// RegisterRateLimit 注册速率限制中间件
func RegisterRateLimit() gin.HandlerFunc {
	limiter := getRegisterLimiter()

	return func(ctx *gin.Context) {
		// 获取客户端 IP
		ip := ctx.ClientIP()

		// 检查是否允许请求
		if !limiter.Allow(ip) {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many registration attempts, please try again later",
				"type":    "rate_limit_error",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
