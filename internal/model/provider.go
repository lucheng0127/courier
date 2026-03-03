package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSON 类型用于存储 extra_config
type JSON map[string]any

// Value 实现 driver.Valuer 接口
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSON) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Provider 配置模型
type Provider struct {
	ID            int64     `json:"id" db:"id" gorm:"primaryKey"`
	Name          string    `json:"name" db:"name" gorm:"uniqueIndex;not null"`               // 必填：唯一标识
	Type          string    `json:"type" db:"type" gorm:"index;not null"`               // 必填：openai, anthropic, vllm 等
	BaseURL       string    `json:"base_url" db:"base_url" gorm:"not null"`       // 必填：API 地址
	Timeout       int       `json:"timeout" db:"timeout" gorm:"default:300"`         // 必填：超时时间（秒），默认 300
	APIKey        *string   `json:"api_key,omitempty" db:"api_key" gorm:"type:varchar(2048)"` // 可选：SaaS 需要
	ExtraConfig   JSON      `json:"extra_config,omitempty" db:"extra_config" gorm:"type:jsonb"` // 可选：扩展配置
	Enabled       bool      `json:"enabled" db:"enabled" gorm:"default:true"`         // 启用状态，默认 true
	FallbackModels JSON     `json:"fallback_models,omitempty" db:"fallback_models" gorm:"type:jsonb"` // 可选：Fallback 模型列表 ["model-1", "model-2"]
	CreatedAt     time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime;default:NOW()"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime;default:NOW()"`
}

// TableName 指定表名
func (Provider) TableName() string {
	return "providers"
}
