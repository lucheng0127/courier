package model

import "time"

// APIKey API Key 模型
type APIKey struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	Key        string     `gorm:"uniqueIndex;not null;size:64" json:"key"`
	Status     string     `gorm:"not null;default:active" json:"status"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (APIKey) TableName() string {
	return "api_keys"
}
