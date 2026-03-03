package model

import "time"

// JWTClaims JWT 声明
type JWTClaims struct {
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// TokenType Token 类型
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// TokenInfo Token 信息
type TokenInfo struct {
	Token     string    `json:"token"`
	Type      TokenType `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
}
