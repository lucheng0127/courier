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
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`               // 必填：唯一标识
	Type        string    `json:"type" db:"type"`               // 必填：openai, anthropic, vllm 等
	BaseURL     string    `json:"base_url" db:"base_url"`       // 必填：API 地址
	Timeout     int       `json:"timeout" db:"timeout"`         // 必填：超时时间（秒），默认 300
	APIKey      *string   `json:"api_key,omitempty" db:"api_key"` // 可选：SaaS 需要
	ExtraConfig JSON      `json:"extra_config,omitempty" db:"extra_config"` // 可选：扩展配置
	Enabled     bool      `json:"enabled" db:"enabled"`         // 启用状态，默认 true
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
