package model

import "time"

// ChatLog Chat 请求日志
type ChatLog struct {
	RequestID        string        `json:"request_id"`
	TraceID          string        `json:"trace_id"`          // TraceID
	APIKey           string        `json:"api_key"`           // 脱敏后
	Model            string        `json:"model"`             // 完整的 provider/model_name
	ProviderName     string        `json:"provider_name"`     // Provider 名称
	ProviderType     string        `json:"provider_type"`     // Provider 类型
	ModelName        string        `json:"model_name"`        // 模型名称
	FallbackCount    int           `json:"fallback_count"`    // Fallback 次数
	FinalModelName   string        `json:"final_model_name"`  // 最终使用的模型
	AttemptDetails   []AttemptDetail `json:"attempt_details,omitempty"` // 尝试详情
	PromptTokens     int           `json:"prompt_tokens"`
	CompletionTokens int           `json:"completion_tokens"`
	TotalTokens      int           `json:"total_tokens"`
	LatencyMs        int64         `json:"latency_ms"`        // 请求耗时（毫秒）
	Status           string        `json:"status"`            // success, error
	Error            string        `json:"error,omitempty"`   // 错误信息
	Timestamp        time.Time     `json:"timestamp"`
}

// AttemptDetail 单次尝试详情（用于日志）
type AttemptDetail struct {
	ModelName string `json:"model_name"`
	ErrorType string `json:"error_type,omitempty"`
	DurationMs int64 `json:"duration_ms"`
}
