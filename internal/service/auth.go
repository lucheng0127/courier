package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/repository"
	passwordpkg "github.com/lucheng0127/courier/internal/pkg/password"
)

// AuthService 认证服务
type AuthService struct {
	userRepo repository.UserRepository
	jwtSvc   JWTService
}

// NewAuthService 创建 Auth Service
func NewAuthService(userRepo repository.UserRepository, jwtSvc JWTService) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

// ValidateAPIKey 验证 API Key 并返回关联的 API Key 记录
func (s *AuthService) ValidateAPIKey(ctx context.Context, apiKey string) (*model.APIKey, error) {
	// 计算哈希
	keyHash := repository.HashAPIKey(apiKey)

	// 从数据库查询
	keyRecord, err := s.userRepo.GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, fmt.Errorf("invalid api key")
	}

	// 检查状态
	if keyRecord.Status != "active" {
		return nil, fmt.Errorf("api key is %s", keyRecord.Status)
	}

	// 检查过期时间
	if keyRecord.ExpiresAt != nil && keyRecord.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("api key has expired")
	}

	return keyRecord, nil
}

// GetUserByID 获取用户信息
func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// GetUserByEmail 获取用户信息
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

// CreateUser 创建用户
func (s *AuthService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	user := &model.User{
		Name:   req.Name,
		Email:  req.Email,
		Status: "active",
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// 验证密码强度（至少 8 个字符）
	if len(req.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	// 哈希密码
	passwordHash, err := passwordpkg.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         "user", // 注册用户默认为普通用户
		Status:       "active",
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// CreateAPIKey 创建 API Key
func (s *AuthService) CreateAPIKey(ctx context.Context, userID int64, req *model.CreateAPIKeyRequest) (*model.CreateAPIKeyResponse, error) {
	// 验证用户存在且状态为 active
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if user.Status != "active" {
		return nil, fmt.Errorf("user is not active")
	}

	// 生成 API Key
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate api key: %w", err)
	}

	keyPrefix := apiKey[:10] // 前10位用于识别
	keyHash := repository.HashAPIKey(apiKey)

	keyRecord := &model.APIKey{
		UserID:    userID,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		Name:      req.Name,
		Status:    "active",
		ExpiresAt: req.ExpiresAt,
	}

	if err := s.userRepo.CreateAPIKey(ctx, keyRecord); err != nil {
		return nil, fmt.Errorf("failed to create api key: %w", err)
	}

	return &model.CreateAPIKeyResponse{
		ID:        keyRecord.ID,
		Key:       apiKey,
		KeyPrefix: keyPrefix,
		Name:      req.Name,
		Status:    "active",
		ExpiresAt: req.ExpiresAt,
		CreatedAt: keyRecord.CreatedAt,
	}, nil
}

// ListAPIKeys 列出用户的所有 API Key
func (s *AuthService) ListAPIKeys(ctx context.Context, userID int64) ([]*model.APIKey, error) {
	// 验证用户存在
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.userRepo.ListAPIKeysByUserID(ctx, userID)
}

// RevokeAPIKey 撤销 API Key
func (s *AuthService) RevokeAPIKey(ctx context.Context, userID, keyID int64) error {
	// 获取 API Key
	key, err := s.userRepo.GetAPIKeyByID(ctx, keyID)
	if err != nil {
		return fmt.Errorf("api key not found")
	}

	// 验证属于该用户
	if key.UserID != userID {
		return fmt.Errorf("api key does not belong to user")
	}

	return s.userRepo.UpdateAPIKeyStatus(ctx, keyID, "revoked")
}

// EnableAPIKey 启用 API Key
func (s *AuthService) EnableAPIKey(ctx context.Context, userID, keyID int64) (*model.APIKey, error) {
	// 获取 API Key
	key, err := s.userRepo.GetAPIKeyByID(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("api key not found")
	}

	// 验证属于该用户
	if key.UserID != userID {
		return nil, fmt.Errorf("api key does not belong to user")
	}

	// 更新状态为 active
	if err := s.userRepo.UpdateAPIKeyStatus(ctx, keyID, "active"); err != nil {
		return nil, fmt.Errorf("failed to enable api key: %w", err)
	}

	// 重新获取更新后的 API Key
	return s.userRepo.GetAPIKeyByID(ctx, keyID)
}

// DisableAPIKey 禁用 API Key
func (s *AuthService) DisableAPIKey(ctx context.Context, userID, keyID int64) (*model.APIKey, error) {
	// 获取 API Key
	key, err := s.userRepo.GetAPIKeyByID(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("api key not found")
	}

	// 验证属于该用户
	if key.UserID != userID {
		return nil, fmt.Errorf("api key does not belong to user")
	}

	// 更新状态为 disabled
	if err := s.userRepo.UpdateAPIKeyStatus(ctx, keyID, "disabled"); err != nil {
		return nil, fmt.Errorf("failed to disable api key: %w", err)
	}

	// 重新获取更新后的 API Key
	return s.userRepo.GetAPIKeyByID(ctx, keyID)
}

// DeleteAPIKey 删除 API Key（硬删除）
func (s *AuthService) DeleteAPIKey(ctx context.Context, userID, keyID int64) error {
	// 获取 API Key
	key, err := s.userRepo.GetAPIKeyByID(ctx, keyID)
	if err != nil {
		return fmt.Errorf("api key not found")
	}

	// 验证属于该用户
	if key.UserID != userID {
		return fmt.Errorf("api key does not belong to user")
	}

	return s.userRepo.DeleteAPIKey(ctx, keyID)
}

// UpdateKeyLastUsed 更新 API Key 最后使用时间（异步）
func (s *AuthService) UpdateKeyLastUsed(ctx context.Context, keyID int64) {
	// 使用独立的 context，避免请求取消影响更新
	updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = s.userRepo.UpdateKeyLastUsed(updateCtx, keyID)
}

// generateAPIKey 生成 API Key
// 格式: sk-<32位随机字符>
func generateAPIKey() (string, error) {
	// 生成 16 字节随机数 (32 个十六进制字符)
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	randomPart := hex.EncodeToString(bytes)
	return "sk-" + randomPart, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// 获取用户（包含密码哈希）
	user, err := s.userRepo.GetUserByEmailWithPassword(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// 验证密码
	if user.PasswordHash == "" || !passwordpkg.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, fmt.Errorf("user account is %s", user.Status)
	}

	// 生成 Access Token
	accessToken, err := s.jwtSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成 Refresh Token
	refreshToken, err := s.jwtSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtSvc.GetAccessTokenExpiration(),
	}, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error) {
	// 验证 Refresh Token
	claims, err := s.jwtSvc.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// 获取用户信息
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, fmt.Errorf("user account is %s", user.Status)
	}

	// 生成新的 Access Token
	accessToken, err := s.jwtSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成新的 Refresh Token
	refreshToken, err := s.jwtSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &model.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtSvc.GetAccessTokenExpiration(),
	}, nil
}

// CreateAdminUser 创建管理员用户（用于初始化）
func (s *AuthService) CreateAdminUser(ctx context.Context, name, email, password string) (*model.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return existingUser, nil // 已存在则直接返回
	}

	// 哈希密码
	passwordHash, err := passwordpkg.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "admin",
		Status:       "active",
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return user, nil
}

// EnsureInitialAdmin 确保存在初始管理员用户
func (s *AuthService) EnsureInitialAdmin(ctx context.Context) error {
	// 检查是否已存在管理员用户
	users, err := s.userRepo.ListUsers(ctx, nil, 1, 0)
	if err == nil && len(users) > 0 {
		// 已有用户，跳过初始化
		return nil
	}

	// 从环境变量读取初始管理员信息
	email := os.Getenv("INITIAL_ADMIN_EMAIL")
	password := os.Getenv("INITIAL_ADMIN_PASSWORD")

	if email == "" || password == "" {
		return nil // 未配置则跳过
	}

	_, err = s.CreateAdminUser(ctx, "Admin", email, password)
	return err
}
