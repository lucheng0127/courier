package model

import "time"

// RequestLog 请求日志模型
type RequestLog struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	ModelName       string     `gorm:"not null;index" json:"model_name"`
	RequestMessages string     `gorm:"type:text" json:"request_messages"`
	ResponseContent string     `gorm:"type:text" json:"response_content"`
	PromptTokens    int        `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalTokens     int        `json:"total_tokens"`
	LatencyMs       int        `json:"latency_ms"`
	Status          string     `gorm:"not null" json:"status"`
	ErrorMessage    string     `gorm:"type:text" json:"error_message"`
	CreatedAt       time.Time  `json:"created_at"`
}

// TableName 指定表名
func (RequestLog) TableName() string {
	return "request_logs"
}
