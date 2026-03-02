package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/repository"
)

// AuthService 认证服务
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建 Auth Service
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
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
