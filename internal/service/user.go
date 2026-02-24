package service

import (
	"errors"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(name, email string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	ListUsers() ([]model.User, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(name, email string) (*model.User, error) {
	// 检查用户名是否已存在
	existing, err := s.userRepo.FindByName(name)
	if err == nil && existing != nil {
		return nil, errors.New("用户名已存在")
	}

	user := &model.User{
		Name:  name,
		Email: email,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) ListUsers() ([]model.User, error) {
	return s.userRepo.FindAll()
}
