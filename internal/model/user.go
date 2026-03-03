package model

import "time"

// User 用户模型
type User struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // 密码哈希，不输出到 JSON
	Role         string    `json:"role" db:"role"`       // user, admin
	Status       string    `json:"status" db:"status"`   // active, disabled
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// APIKey API Key 模型
type APIKey struct {
	ID         int64      `json:"id" db:"id"`
	UserID     int64      `json:"user_id" db:"user_id"`
	KeyHash    string     `json:"-" db:"key_hash"`    // SHA256 哈希存储，不输出到 JSON
	KeyPrefix  string     `json:"key_prefix" db:"key_prefix"` // 前8位用于识别
	Name       string     `json:"name" db:"name"`      // 用户定义的名称
	Status     string     `json:"status" db:"status"`  // active, disabled, revoked
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// CreateAPIKeyRequest 创建 API Key 请求
type CreateAPIKeyRequest struct {
	Name      string    `json:"name" binding:"required"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// CreateAPIKeyResponse 创建 API Key 响应（包含完整 Key，仅在创建时返回）
type CreateAPIKeyResponse struct {
	ID        int64     `json:"id"`
	Key       string    `json:"key"`       // 完整的 API Key，仅此一次
	KeyPrefix string    `json:"key_prefix"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// APIKeyListItem API Key 列表项（不包含完整 Key）
type APIKeyListItem struct {
	ID         int64      `json:"id"`
	KeyPrefix  string     `json:"key_prefix"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // Bearer
	ExpiresIn    int    `json:"expires_in"` // 秒
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse 刷新 Token 响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
