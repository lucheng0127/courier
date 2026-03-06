package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/courier/internal/model"
	"github.com/lucheng0127/courier/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// ErrUserNotFound 用户未找到错误
	ErrUserNotFound = errors.New("user not found")
	// ErrAPIKeyNotFound API Key 未找到错误
	ErrAPIKeyNotFound = errors.New("api key not found")
)

// setupTestAuthRouter 设置测试路由
func setupTestAuthRouter(authSvc *service.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	authCtrl := NewAuthController(authSvc)

	api := router.Group("/api/v1")
	authCtrl.RegisterRoutes(api)

	return router
}

// TestAuthController_Register_Success 测试成功注册
func TestAuthController_Register_Success(t *testing.T) {
	mockRepo := NewMockUserRepositoryForController()
	mockJWT := &MockJWTServiceForController{}
	authSvc := service.NewAuthService(mockRepo, mockJWT)
	router := setupTestAuthRouter(authSvc)

	reqBody := model.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.RegisterResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "Test User", resp.Name)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "user", resp.Role)
	assert.Equal(t, "active", resp.Status)
	assert.Greater(t, resp.ID, int64(0))
}

// TestAuthController_Register_MissingName 测试缺少姓名字段
func TestAuthController_Register_MissingName(t *testing.T) {
	mockRepo := NewMockUserRepositoryForController()
	mockJWT := &MockJWTServiceForController{}
	authSvc := service.NewAuthService(mockRepo, mockJWT)
	router := setupTestAuthRouter(authSvc)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid_request_error", resp["type"])
}

// TestAuthController_Register_InvalidEmail 测试无效邮箱格式
func TestAuthController_Register_InvalidEmail(t *testing.T) {
	mockRepo := NewMockUserRepositoryForController()
	mockJWT := &MockJWTServiceForController{}
	authSvc := service.NewAuthService(mockRepo, mockJWT)
	router := setupTestAuthRouter(authSvc)

	reqBody := map[string]string{
		"name":     "Test User",
		"email":    "invalid-email",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAuthController_Register_EmailExists 测试邮箱已存在
func TestAuthController_Register_EmailExists(t *testing.T) {
	mockRepo := NewMockUserRepositoryForController()
	mockJWT := &MockJWTServiceForController{}
	authSvc := service.NewAuthService(mockRepo, mockJWT)
	router := setupTestAuthRouter(authSvc)

	// 先注册一个用户
	reqBody1 := model.RegisterRequest{
		Name:     "First User",
		Email:    "test@example.com",
		Password: "password123",
	}
	body1, _ := json.Marshal(reqBody1)
	req1, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// 尝试用相同邮箱再次注册
	reqBody2 := model.RegisterRequest{
		Name:     "Second User",
		Email:    "test@example.com",
		Password: "different456",
	}
	body2, _ := json.Marshal(reqBody2)
	req2, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, "invalid_request_error", resp["type"])
}

// TestAuthController_Register_PasswordTooShort 测试密码过短
func TestAuthController_Register_PasswordTooShort(t *testing.T) {
	mockRepo := NewMockUserRepositoryForController()
	mockJWT := &MockJWTServiceForController{}
	authSvc := service.NewAuthService(mockRepo, mockJWT)
	router := setupTestAuthRouter(authSvc)

	reqBody := model.RegisterRequest{
		Name:     "Test User",
		Email:    "short-pwd@example.com",
		Password: "short",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.168.1.2:1234" // 使用不同的 IP 避免速率限制

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid_request_error", resp["type"])
}

// TestAuthController_Register_MissingPassword 测试缺少密码字段
func TestAuthController_Register_MissingPassword(t *testing.T) {
	mockRepo := NewMockUserRepositoryForController()
	mockJWT := &MockJWTServiceForController{}
	authSvc := service.NewAuthService(mockRepo, mockJWT)
	router := setupTestAuthRouter(authSvc)

	reqBody := map[string]string{
		"name":  "Test User",
		"email": "no-pwd@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.168.1.3:1234" // 使用不同的 IP 避免速率限制

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ========== Mock 类型（用于控制器测试） ==========

// MockUserRepositoryForController 用于控制器测试
type MockUserRepositoryForController struct {
	users     map[int64]*model.User
	emailToID map[string]int64
	nextID    int64
}

func NewMockUserRepositoryForController() *MockUserRepositoryForController {
	return &MockUserRepositoryForController{
		users:     make(map[int64]*model.User),
		emailToID: make(map[string]int64),
		nextID:    1,
	}
}

func (m *MockUserRepositoryForController) CreateUser(ctx context.Context, user *model.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.emailToID[user.Email] = user.ID
	return nil
}

func (m *MockUserRepositoryForController) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepositoryForController) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	id, ok := m.emailToID[email]
	if !ok {
		return nil, ErrUserNotFound
	}
	return m.users[id], nil
}

func (m *MockUserRepositoryForController) GetUserByEmailWithPassword(ctx context.Context, email string) (*model.User, error) {
	return m.GetUserByEmail(ctx, email)
}

func (m *MockUserRepositoryForController) ListUsers(ctx context.Context, status *string, limit, offset int) ([]*model.User, error) {
	return nil, nil
}

func (m *MockUserRepositoryForController) UpdateUser(ctx context.Context, user *model.User) error {
	return nil
}

func (m *MockUserRepositoryForController) UpdateUserStatus(ctx context.Context, id int64, status string) error {
	return nil
}

func (m *MockUserRepositoryForController) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	return nil
}

func (m *MockUserRepositoryForController) DeleteUser(ctx context.Context, id int64) error {
	return nil
}

func (m *MockUserRepositoryForController) CreateAPIKey(ctx context.Context, key *model.APIKey) error {
	return nil
}

func (m *MockUserRepositoryForController) GetAPIKeyByHash(ctx context.Context, keyHash string) (*model.APIKey, error) {
	return nil, ErrAPIKeyNotFound
}

func (m *MockUserRepositoryForController) GetAPIKeyByID(ctx context.Context, id int64) (*model.APIKey, error) {
	return nil, ErrAPIKeyNotFound
}

func (m *MockUserRepositoryForController) ListAPIKeysByUserID(ctx context.Context, userID int64) ([]*model.APIKey, error) {
	return nil, nil
}

func (m *MockUserRepositoryForController) UpdateAPIKeyStatus(ctx context.Context, id int64, status string) error {
	return nil
}

func (m *MockUserRepositoryForController) UpdateKeyLastUsed(ctx context.Context, id int64) error {
	return nil
}

func (m *MockUserRepositoryForController) DeleteAPIKey(ctx context.Context, id int64) error {
	return nil
}

// MockJWTServiceForController 用于控制器测试
type MockJWTServiceForController struct{}

func (m *MockJWTServiceForController) GenerateAccessToken(user *model.User) (string, error) {
	return "mock-token", nil
}

func (m *MockJWTServiceForController) GenerateRefreshToken(userID int64) (string, error) {
	return "mock-refresh", nil
}

func (m *MockJWTServiceForController) ValidateAccessToken(token string) (*model.JWTClaims, error) {
	return &model.JWTClaims{UserID: 1}, nil
}

func (m *MockJWTServiceForController) ValidateRefreshToken(token string) (*model.JWTClaims, error) {
	return &model.JWTClaims{UserID: 1}, nil
}

func (m *MockJWTServiceForController) GetAccessTokenExpiration() int {
	return 3600
}
