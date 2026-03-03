package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lucheng0127/courier/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// ErrUserNotFound 用户未找到错误
	ErrUserNotFound = errors.New("user not found")
	// ErrAPIKeyNotFound API Key 未找到错误
	ErrAPIKeyNotFound = errors.New("api key not found")
)

// MockUserRepository 用于测试的 mock repository
type MockUserRepository struct {
	users         map[int64]*model.User
	emailToID     map[string]int64
	nextID        int64
	passwordHashes map[string]string // email -> password hash
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:         make(map[int64]*model.User),
		emailToID:     make(map[string]int64),
		nextID:        1,
		passwordHashes: make(map[string]string),
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.emailToID[user.Email] = user.ID
	if user.PasswordHash != "" {
		m.passwordHashes[user.Email] = user.PasswordHash
	}
	return nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	id, ok := m.emailToID[email]
	if !ok {
		return nil, ErrUserNotFound
	}
	return m.users[id], nil
}

func (m *MockUserRepository) GetUserByEmailWithPassword(ctx context.Context, email string) (*model.User, error) {
	id, ok := m.emailToID[email]
	if !ok {
		return nil, ErrUserNotFound
	}
	return m.users[id], nil
}

func (m *MockUserRepository) ListUsers(ctx context.Context, status *string, limit, offset int) ([]*model.User, error) {
	users := make([]*model.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) UpdateUserStatus(ctx context.Context, id int64, status string) error {
	if user, ok := m.users[id]; ok {
		user.Status = status
	}
	return nil
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	if user, ok := m.users[id]; ok {
		user.PasswordHash = passwordHash
	}
	return nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int64) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) CreateAPIKey(ctx context.Context, key *model.APIKey) error {
	return nil
}

func (m *MockUserRepository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*model.APIKey, error) {
	return nil, ErrAPIKeyNotFound
}

func (m *MockUserRepository) GetAPIKeyByID(ctx context.Context, id int64) (*model.APIKey, error) {
	return nil, ErrAPIKeyNotFound
}

func (m *MockUserRepository) ListAPIKeysByUserID(ctx context.Context, userID int64) ([]*model.APIKey, error) {
	return nil, nil
}

func (m *MockUserRepository) UpdateAPIKeyStatus(ctx context.Context, id int64, status string) error {
	return nil
}

func (m *MockUserRepository) UpdateKeyLastUsed(ctx context.Context, id int64) error {
	return nil
}

// TestAuthService_Register_Success 测试成功注册新用户
func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockJWT := &MockJWTService{}
	authSvc := NewAuthService(mockRepo, mockJWT)

	req := &model.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	user, err := authSvc.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "user", user.Role) // 默认角色为 user
	assert.Equal(t, "active", user.Status)
	assert.NotEmpty(t, user.PasswordHash, "密码应该被哈希存储")
	assert.NotEqual(t, "password123", user.PasswordHash, "密码不应该明文存储")
}

// TestAuthService_Register_EmailExists 测试邮箱已存在时注册失败
func TestAuthService_Register_EmailExists(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockJWT := &MockJWTService{}
	authSvc := NewAuthService(mockRepo, mockJWT)

	// 先注册一个用户
	req1 := &model.RegisterRequest{
		Name:     "First User",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, _ = authSvc.Register(context.Background(), req1)

	// 尝试用相同邮箱注册
	req2 := &model.RegisterRequest{
		Name:     "Second User",
		Email:    "test@example.com",
		Password: "different456",
	}
	_, err := authSvc.Register(context.Background(), req2)

	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())
}

// TestAuthService_Register_PasswordTooShort 测试密码过短时注册失败
func TestAuthService_Register_PasswordTooShort(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockJWT := &MockJWTService{}
	authSvc := NewAuthService(mockRepo, mockJWT)

	req := &model.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "short",
	}

	_, err := authSvc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "password must be at least 8 characters", err.Error())
}

// TestAuthService_Register_PasswordExactly8Chars 测试密码正好8个字符时注册成功
func TestAuthService_Register_PasswordExactly8Chars(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockJWT := &MockJWTService{}
	authSvc := NewAuthService(mockRepo, mockJWT)

	req := &model.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "12345678", // 正好 8 个字符
	}

	user, err := authSvc.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.PasswordHash)
}

// MockJWTService 用于测试的 mock JWT service
type MockJWTService struct{}

func (m *MockJWTService) GenerateAccessToken(user *model.User) (string, error) {
	return "mock-access-token", nil
}

func (m *MockJWTService) GenerateRefreshToken(userID int64) (string, error) {
	return "mock-refresh-token", nil
}

func (m *MockJWTService) ValidateAccessToken(token string) (*model.JWTClaims, error) {
	return &model.JWTClaims{UserID: 1, UserRole: "user"}, nil
}

func (m *MockJWTService) ValidateRefreshToken(token string) (*model.JWTClaims, error) {
	return &model.JWTClaims{UserID: 1}, nil
}

func (m *MockJWTService) GetAccessTokenExpiration() int {
	return 3600
}
