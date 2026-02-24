package repository

import (
	"time"

	"github.com/lucheng0127/courier/internal/model"
	"gorm.io/gorm"
)

// APIKeyRepository API Key 仓储接口
type APIKeyRepository interface {
	Create(apiKey *model.APIKey) error
	FindByID(id uint) (*model.APIKey, error)
	FindByUserID(userID uint) ([]model.APIKey, error)
	FindByKey(key string) (*model.APIKey, error)
	Delete(id uint) error
	UpdateStatus(id uint, status string) error
	UpdateLastUsedAt(id uint, lastUsedAt time.Time) error
}

// apiKeyRepository API Key 仓储实现
type apiKeyRepository struct {
	db *gorm.DB
}

// NewAPIKeyRepository 创建 API Key 仓储
func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(apiKey *model.APIKey) error {
	return r.db.Create(apiKey).Error
}

func (r *apiKeyRepository) FindByID(id uint) (*model.APIKey, error) {
	var apiKey model.APIKey
	err := r.db.First(&apiKey, id).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *apiKeyRepository) FindByUserID(userID uint) ([]model.APIKey, error) {
	var apiKeys []model.APIKey
	err := r.db.Where("user_id = ?", userID).Find(&apiKeys).Error
	return apiKeys, err
}

func (r *apiKeyRepository) FindByKey(key string) (*model.APIKey, error) {
	var apiKey model.APIKey
	err := r.db.Where("key = ?", key).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *apiKeyRepository) Delete(id uint) error {
	return r.db.Delete(&model.APIKey{}, id).Error
}

func (r *apiKeyRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.APIKey{}).Where("id = ?", id).Update("status", status).Error
}

func (r *apiKeyRepository) UpdateLastUsedAt(id uint, lastUsedAt time.Time) error {
	return r.db.Model(&model.APIKey{}).Where("id = ?", id).Update("last_used_at", lastUsedAt).Error
}
