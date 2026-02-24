package service

import (
	"errors"
	"github.com/lucheng0127/courier/pkg/config"
)

// ModelService 模型服务接口
type ModelService interface {
	GetModels() []ModelInfo
	GetModelByName(name string) (*config.ModelConfig, error)
}

// ModelInfo 模型信息
type ModelInfo struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

// modelService 模型服务实现
type modelService struct {
	models []config.ModelConfig
}

// NewModelService 创建模型服务
func NewModelService(models []config.ModelConfig) ModelService {
	return &modelService{models: models}
}

func (s *modelService) GetModels() []ModelInfo {
	result := make([]ModelInfo, 0, len(s.models))
	for _, m := range s.models {
		result = append(result, ModelInfo{
			Name:     m.Name,
			Provider: m.Provider,
		})
	}
	return result
}

func (s *modelService) GetModelByName(name string) (*config.ModelConfig, error) {
	for _, m := range s.models {
		if m.Name == name {
			return &m, nil
		}
	}
	return nil, errors.New("模型不存在")
}
