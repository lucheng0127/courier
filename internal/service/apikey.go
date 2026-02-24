package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/repository"
)

// APIKeyService API Key 服务接口
type APIKeyService interface {
	GenerateAPIKey(userID uint) (*model.APIKey, error)
	ListAPIKeys(userID uint) ([]model.APIKey, error)
	DeleteAPIKey(userID, keyID uint) error
	DisableAPIKey(userID, keyID uint) (*model.APIKey, error)
}

// apiKeyService API Key 服务实现
type apiKeyService struct {
	userRepo     repository.UserRepository
	apiKeyRepo   repository.APIKeyRepository
}

// NewAPIKeyService 创建 API Key 服务
func NewAPIKeyService(userRepo repository.UserRepository, apiKeyRepo repository.APIKeyRepository) APIKeyService {
	return &apiKeyService{
		userRepo:   userRepo,
		apiKeyRepo: apiKeyRepo,
	}
}

// generateKeyString 生成 API Key 字符串
func generateKeyString() (string, error) {
	bytes := make([]byte, 24)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return "ck_" + hex.EncodeToString(bytes), nil
}

func (s *apiKeyService) GenerateAPIKey(userID uint) (*model.APIKey, error) {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 生成 API Key
	keyStr, err := generateKeyString()
	if err != nil {
		return nil, err
	}

	apiKey := &model.APIKey{
		UserID: userID,
		Key:    keyStr,
		Status: "active",
	}

	err = s.apiKeyRepo.Create(apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (s *apiKeyService) ListAPIKeys(userID uint) ([]model.APIKey, error) {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return s.apiKeyRepo.FindByUserID(userID)
}

func (s *apiKeyService) DeleteAPIKey(userID, keyID uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查 API Key 是否存在且属于该用户
	apiKey, err := s.apiKeyRepo.FindByID(keyID)
	if err != nil {
		return errors.New("API Key 不存在")
	}
	if apiKey.UserID != userID {
		return errors.New("API Key 不存在")
	}

	return s.apiKeyRepo.Delete(keyID)
}

func (s *apiKeyService) DisableAPIKey(userID, keyID uint) (*model.APIKey, error) {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查 API Key 是否存在且属于该用户
	apiKey, err := s.apiKeyRepo.FindByID(keyID)
	if err != nil {
		return nil, errors.New("API Key 不存在")
	}
	if apiKey.UserID != userID {
		return nil, errors.New("API Key 不存在")
	}

	// 更新状态
	err = s.apiKeyRepo.UpdateStatus(keyID, "disabled")
	if err != nil {
		return nil, err
	}

	// 重新获取更新后的数据
	apiKey, err = s.apiKeyRepo.FindByID(keyID)
	if err != nil {
		return nil, fmt.Errorf("获取更新后的 API Key 失败: %w", err)
	}

	return apiKey, nil
}
