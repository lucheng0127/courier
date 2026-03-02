package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lucheng0127/courier/internal/model"
)

// JWTService JWT 服务接口
type JWTService interface {
	// GenerateAccessToken 生成访问令牌
	GenerateAccessToken(user *model.User) (string, error)

	// GenerateRefreshToken 生成刷新令牌
	GenerateRefreshToken(userID int64) (string, error)

	// ValidateAccessToken 验证访问令牌
	ValidateAccessToken(token string) (*model.JWTClaims, error)

	// ValidateRefreshToken 验证刷新令牌
	ValidateRefreshToken(token string) (*model.JWTClaims, error)

	// GetAccessTokenExpiration 获取 Access Token 过期时间（秒）
	GetAccessTokenExpiration() int
}

// jwtService JWT 服务实现
type jwtService struct {
	secretKey              []byte
	issuer                 string
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

// NewJWTService 创建 JWT 服务
func NewJWTService() (JWTService, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "courier-gateway"
	}

	// 默认 Access Token 15 分钟
	accessExpiration := 15 * time.Minute
	if exp := os.Getenv("JWT_ACCESS_TOKEN_EXPIRES_IN"); exp != "" {
		if d, err := time.ParseDuration(exp); err == nil {
			accessExpiration = d
		}
	}

	// 默认 Refresh Token 7 天
	refreshExpiration := 7 * 24 * time.Hour
	if exp := os.Getenv("JWT_REFRESH_TOKEN_EXPIRES_IN"); exp != "" {
		if d, err := time.ParseDuration(exp); err == nil {
			refreshExpiration = d
		}
	}

	return &jwtService{
		secretKey:              []byte(secret),
		issuer:                 issuer,
		accessTokenExpiration:  accessExpiration,
		refreshTokenExpiration: refreshExpiration,
	}, nil
}

// GenerateAccessToken 生成访问令牌
func (s *jwtService) GenerateAccessToken(user *model.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessTokenExpiration)

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"user_email": user.Email,
		"user_role":  user.Role,
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"iss":        s.issuer,
		"type":       model.TokenTypeAccess,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// GenerateRefreshToken 生成刷新令牌
func (s *jwtService) GenerateRefreshToken(userID int64) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.refreshTokenExpiration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     expiresAt.Unix(),
		"iss":     s.issuer,
		"type":    model.TokenTypeRefresh,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateAccessToken 验证访问令牌
func (s *jwtService) ValidateAccessToken(tokenString string) (*model.JWTClaims, error) {
	return s.validateToken(tokenString, model.TokenTypeAccess)
}

// ValidateRefreshToken 验证刷新令牌
func (s *jwtService) ValidateRefreshToken(tokenString string) (*model.JWTClaims, error) {
	return s.validateToken(tokenString, model.TokenTypeRefresh)
}

// validateToken 验证 Token
func (s *jwtService) validateToken(tokenString string, tokenType model.TokenType) (*model.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// 验证发行者
	if iss, ok := claims["iss"].(string); !ok || iss != s.issuer {
		return nil, fmt.Errorf("invalid token issuer")
	}

	// 验证 Token 类型
	if claimsType, ok := claims["type"].(string); !ok || model.TokenType(claimsType) != tokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	// 提取 claims
	userID, _ := claims["user_id"].(float64)
	userEmail, _ := claims["user_email"].(string)
	userRole, _ := claims["user_role"].(string)
	iat, _ := claims["iat"].(float64)
	exp, _ := claims["exp"].(float64)

	return &model.JWTClaims{
		UserID:    int64(userID),
		UserEmail: userEmail,
		UserRole:  userRole,
		IssuedAt:  int64(iat),
		ExpiresAt: int64(exp),
	}, nil
}

// GetAccessTokenExpiration 获取 Access Token 过期时间（秒）
func (s *jwtService) GetAccessTokenExpiration() int {
	return int(s.accessTokenExpiration.Seconds())
}
