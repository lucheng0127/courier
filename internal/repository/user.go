package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lucheng0127/courier/internal/model"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *model.User) error

	// GetUserByID 按 ID 查询用户
	GetUserByID(ctx context.Context, id int64) (*model.User, error)

	// GetUserByEmail 按 email 查询用户
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// GetUserByEmailWithPassword 按 email 查询用户（包含密码哈希，用于登录验证）
	GetUserByEmailWithPassword(ctx context.Context, email string) (*model.User, error)

	// ListUsers 列出用户（支持分页和状态过滤）
	ListUsers(ctx context.Context, status *string, limit, offset int) ([]*model.User, error)

	// UpdateUserStatus 更新用户状态
	UpdateUserStatus(ctx context.Context, id int64, status string) error

	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, user *model.User) error

	// UpdatePassword 更新用户密码
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error

	// CreateAPIKey 创建 API Key
	CreateAPIKey(ctx context.Context, key *model.APIKey) error

	// GetAPIKeyByHash 按 key_hash 查询 API Key
	GetAPIKeyByHash(ctx context.Context, keyHash string) (*model.APIKey, error)

	// GetAPIKeyByID 按 ID 查询 API Key
	GetAPIKeyByID(ctx context.Context, id int64) (*model.APIKey, error)

	// ListAPIKeysByUserID 查询用户的所有 API Key
	ListAPIKeysByUserID(ctx context.Context, userID int64) ([]*model.APIKey, error)

	// UpdateAPIKeyStatus 更新 API Key 状态
	UpdateAPIKeyStatus(ctx context.Context, id int64, status string) error

	// UpdateKeyLastUsed 更新 API Key 最后使用时间
	UpdateKeyLastUsed(ctx context.Context, id int64) error
}

// userRepository 用户数据访问实现
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository 创建 User Repository
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser 创建用户
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Status,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID 按 ID 查询用户
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, role, status, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

// GetUserByEmail 按 email 查询用户
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, role, status, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetUserByEmailWithPassword 按 email 查询用户（包含密码哈希，用于登录验证）
func (r *userRepository) GetUserByEmailWithPassword(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, password_hash, role, status, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email with password: %w", err)
	}
	return &user, nil
}

// ListUsers 列出用户
func (r *userRepository) ListUsers(ctx context.Context, status *string, limit, offset int) ([]*model.User, error) {
	var users []*model.User
	query := `SELECT id, name, email, role, status, created_at, updated_at FROM users`
	args := []interface{}{}

	if status != nil {
		query += ` WHERE status = $1`
		args = append(args, *status)
	}

	query += ` ORDER BY created_at DESC`

	if limit > 0 {
		query += ` LIMIT $` + fmt.Sprint(len(args)+1)
		args = append(args, limit)
	}

	if offset > 0 {
		query += ` OFFSET $` + fmt.Sprint(len(args)+1)
		args = append(args, offset)
	}

	err := r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// UpdateUserStatus 更新用户状态
func (r *userRepository) UpdateUserStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE users
		SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at
	`
	var updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, status, id).Scan(&updatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

// UpdateUser 更新用户信息
func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, role = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.Role,
		user.ID,
	).Scan(&user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, passwordHash, id)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// CreateAPIKey 创建 API Key
func (r *userRepository) CreateAPIKey(ctx context.Context, key *model.APIKey) error {
	query := `
		INSERT INTO api_keys (user_id, key_hash, key_prefix, name, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query,
		key.UserID,
		key.KeyHash,
		key.KeyPrefix,
		key.Name,
		key.Status,
		key.ExpiresAt,
	).Scan(&key.ID, &key.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create api key: %w", err)
	}
	return nil
}

// GetAPIKeyByHash 按 key_hash 查询 API Key
func (r *userRepository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*model.APIKey, error) {
	var key model.APIKey
	query := `
		SELECT id, user_id, key_hash, key_prefix, name, status, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE key_hash = $1
	`
	err := r.db.GetContext(ctx, &key, query, keyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get api key by hash: %w", err)
	}
	return &key, nil
}

// GetAPIKeyByID 按 ID 查询 API Key
func (r *userRepository) GetAPIKeyByID(ctx context.Context, id int64) (*model.APIKey, error) {
	var key model.APIKey
	query := `
		SELECT id, user_id, key_hash, key_prefix, name, status, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &key, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get api key by id: %w", err)
	}
	return &key, nil
}

// ListAPIKeysByUserID 查询用户的所有 API Key
func (r *userRepository) ListAPIKeysByUserID(ctx context.Context, userID int64) ([]*model.APIKey, error) {
	var keys []*model.APIKey
	query := `
		SELECT id, user_id, key_hash, key_prefix, name, status, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.SelectContext(ctx, &keys, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list api keys: %w", err)
	}
	return keys, nil
}

// UpdateAPIKeyStatus 更新 API Key 状态
func (r *userRepository) UpdateAPIKeyStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE api_keys
		SET status = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update api key status: %w", err)
	}
	return nil
}

// UpdateKeyLastUsed 更新 API Key 最后使用时间
func (r *userRepository) UpdateKeyLastUsed(ctx context.Context, id int64) error {
	query := `
		UPDATE api_keys
		SET last_used_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update api key last used: %w", err)
	}
	return nil
}

// HashAPIKey 生成 API Key 的 SHA256 哈希
func HashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return fmt.Sprintf("%x", hash)
}
