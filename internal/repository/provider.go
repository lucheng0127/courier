package repository

import (
	"context"
	"fmt"

	"github.com/lucheng0127/courier/internal/model"
	"github.com/jmoiron/sqlx"
)

// ProviderRepository Provider 数据访问接口
type ProviderRepository interface {
	// Create 创建 Provider
	Create(ctx context.Context, provider *model.Provider) error

	// GetByID 按 ID 查询
	GetByID(ctx context.Context, id int64) (*model.Provider, error)

	// GetByName 按 name 查询
	GetByName(ctx context.Context, name string) (*model.Provider, error)

	// List 列出所有 Provider
	List(ctx context.Context) ([]*model.Provider, error)

	// Update 更新 Provider
	Update(ctx context.Context, provider *model.Provider) error

	// Delete 删除 Provider
	Delete(ctx context.Context, id int64) error

	// ExistsByName 检查 name 是否存在
	ExistsByName(ctx context.Context, name string) (bool, error)
}

// providerRepository Provider 数据访问实现
type providerRepository struct {
	db *sqlx.DB
}

// NewProviderRepository 创建 Provider Repository
func NewProviderRepository(db *sqlx.DB) ProviderRepository {
	return &providerRepository{db: db}
}

// Create 创建 Provider
func (r *providerRepository) Create(ctx context.Context, provider *model.Provider) error {
	query := `
		INSERT INTO providers (name, type, base_url, timeout, api_key, extra_config, enabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		provider.Name,
		provider.Type,
		provider.BaseURL,
		provider.Timeout,
		provider.APIKey,
		provider.ExtraConfig,
		provider.Enabled,
	).Scan(&provider.ID, &provider.CreatedAt, &provider.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}
	return nil
}

// GetByID 按 ID 查询
func (r *providerRepository) GetByID(ctx context.Context, id int64) (*model.Provider, error) {
	var provider model.Provider
	query := `SELECT id, name, type, base_url, timeout, api_key, extra_config, enabled, created_at, updated_at FROM providers WHERE id = $1`
	err := r.db.GetContext(ctx, &provider, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider by id: %w", err)
	}
	return &provider, nil
}

// GetByName 按 name 查询
func (r *providerRepository) GetByName(ctx context.Context, name string) (*model.Provider, error) {
	var provider model.Provider
	query := `SELECT id, name, type, base_url, timeout, api_key, extra_config, enabled, created_at, updated_at FROM providers WHERE name = $1`
	err := r.db.GetContext(ctx, &provider, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider by name: %w", err)
	}
	return &provider, nil
}

// List 列出所有 Provider
func (r *providerRepository) List(ctx context.Context) ([]*model.Provider, error) {
	var providers []*model.Provider
	query := `SELECT id, name, type, base_url, timeout, api_key, extra_config, enabled, created_at, updated_at FROM providers ORDER BY name`
	err := r.db.SelectContext(ctx, &providers, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	return providers, nil
}

// Update 更新 Provider
func (r *providerRepository) Update(ctx context.Context, provider *model.Provider) error {
	query := `
		UPDATE providers
		SET type = $1, base_url = $2, timeout = $3, api_key = $4, extra_config = $5, enabled = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		provider.Type,
		provider.BaseURL,
		provider.Timeout,
		provider.APIKey,
		provider.ExtraConfig,
		provider.Enabled,
		provider.ID,
	).Scan(&provider.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}
	return nil
}

// Delete 删除 Provider
func (r *providerRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM providers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}
	return nil
}

// ExistsByName 检查 name 是否存在
func (r *providerRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM providers WHERE name = $1)`
	err := r.db.GetContext(ctx, &exists, query, name)
	if err != nil {
		return false, fmt.Errorf("failed to check provider name existence: %w", err)
	}
	return exists, nil
}
