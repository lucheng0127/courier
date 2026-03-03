package service

import (
	"os"
	"testing"

	"github.com/lucheng0127/courier/internal/model"
)

func setupJWTTest(t *testing.T) func() {
	// 设置测试环境变量
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	os.Setenv("JWT_ISSUER", "test-issuer")
	os.Setenv("JWT_ACCESS_TOKEN_EXPIRES_IN", "15m")
	os.Setenv("JWT_REFRESH_TOKEN_EXPIRES_IN", "168h")

	return func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_ISSUER")
		os.Unsetenv("JWT_ACCESS_TOKEN_EXPIRES_IN")
		os.Unsetenv("JWT_REFRESH_TOKEN_EXPIRES_IN")
	}
}

func TestNewJWTService(t *testing.T) {
	cleanup := setupJWTTest(t)
	defer cleanup()

	svc, err := NewJWTService()
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	if svc == nil {
		t.Error("NewJWTService() returned nil service")
	}
}

func TestNewJWTServiceMissingSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	_, err := NewJWTService()
	if err == nil {
		t.Error("NewJWTService() expected error when JWT_SECRET is missing")
	}
}

func TestGenerateAccessToken(t *testing.T) {
	cleanup := setupJWTTest(t)
	defer cleanup()

	svc, err := NewJWTService()
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	user := &model.User{
		ID:    123,
		Email: "test@example.com",
		Role:  "user",
	}

	token, err := svc.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateAccessToken() returned empty token")
	}

	// Token 应该包含三个部分（header.payload.signature）
	parts := splitToken(token)
	if len(parts) != 3 {
		t.Errorf("GenerateAccessToken() token has %d parts, want 3", len(parts))
	}
}

func TestValidateAccessToken(t *testing.T) {
	cleanup := setupJWTTest(t)
	defer cleanup()

	svc, err := NewJWTService()
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	user := &model.User{
		ID:    123,
		Email: "test@example.com",
		Role:  "user",
	}

	token, err := svc.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	claims, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("ValidateAccessToken() UserID = %v, want %v", claims.UserID, user.ID)
	}

	if claims.UserEmail != user.Email {
		t.Errorf("ValidateAccessToken() UserEmail = %v, want %v", claims.UserEmail, user.Email)
	}

	if claims.UserRole != user.Role {
		t.Errorf("ValidateAccessToken() UserRole = %v, want %v", claims.UserRole, user.Role)
	}
}

func TestValidateAccessTokenInvalid(t *testing.T) {
	cleanup := setupJWTTest(t)
	defer cleanup()

	svc, err := NewJWTService()
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid token format",
			token: "invalid.token",
		},
		{
			name:  "malformed token",
			token: "not-a-jwt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ValidateAccessToken(tt.token)
			if err == nil {
				t.Error("ValidateAccessToken() expected error for invalid token")
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	cleanup := setupJWTTest(t)
	defer cleanup()

	svc, err := NewJWTService()
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	userID := int64(123)

	token, err := svc.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateRefreshToken() returned empty token")
	}

	parts := splitToken(token)
	if len(parts) != 3 {
		t.Errorf("GenerateRefreshToken() token has %d parts, want 3", len(parts))
	}
}

func TestValidateRefreshToken(t *testing.T) {
	cleanup := setupJWTTest(t)
	defer cleanup()

	svc, err := NewJWTService()
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	userID := int64(123)

	token, err := svc.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	claims, err := svc.ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("ValidateRefreshToken() UserID = %v, want %v", claims.UserID, userID)
	}

	// Refresh token 不应包含敏感信息
	if claims.UserEmail != "" {
		t.Error("ValidateRefreshToken() UserEmail should be empty in refresh token")
	}
}

func TestAccessTokenCannotBeUsedAsRefreshToken(t *testing.T) {
	cleanup := setupJWTTest(t)
	defer cleanup()

	svc, err := NewJWTService()
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	user := &model.User{
		ID:    123,
		Email: "test@example.com",
		Role:  "user",
	}

	accessToken, err := svc.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	// 尝试用 refresh token 验证方法验证 access token
	_, err = svc.ValidateRefreshToken(accessToken)
	if err == nil {
		t.Error("ValidateRefreshToken() should reject access token")
	}
}

func TestGetAccessTokenExpiration(t *testing.T) {
	cleanup := setupJWTTest(t)
	defer cleanup()

	svc, err := NewJWTService()
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	expiration := svc.GetAccessTokenExpiration()

	// 默认 15 分钟 = 900 秒
	if expiration != 900 {
		t.Errorf("GetAccessTokenExpiration() = %v, want 900", expiration)
	}
}

// splitToken 辅助函数
func splitToken(token string) []string {
	parts := make([]string, 0)
	current := ""
	for _, c := range token {
		if c == '.' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
