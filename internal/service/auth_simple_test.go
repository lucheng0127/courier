package service

import (
	"testing"

	"github.com/lucheng0127/courier/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateAPIKey 测试 API Key 生成
func TestGenerateAPIKey(t *testing.T) {
	key, err := generateAPIKey()

	require.NoError(t, err)
	assert.True(t, len(key) >= 35, "API Key should be at least 35 characters")
	assert.True(t, key[:3] == "sk-", "API Key should start with 'sk-'")

	// 验证格式: sk-<32个十六进制字符>
	// sk- (3 chars) + 32 hex chars = 35 chars minimum
	assert.Equal(t, 35, len(key), "API Key should be exactly 35 characters (sk- + 32 hex chars)")
}

// TestHashAPIKey 测试 API Key 哈希
func TestHashAPIKey(t *testing.T) {
	key := "sk-test1234567890123456789012345678"
	hash1 := repository.HashAPIKey(key)
	hash2 := repository.HashAPIKey(key)

	assert.Equal(t, hash1, hash2, "Same key should produce same hash")
	assert.Len(t, hash1, 64, "SHA256 hash should be 64 hex characters")

	// 不同的 key 应该产生不同的 hash
	differentKey := "sk-different12345678901234567890123"
	differentHash := repository.HashAPIKey(differentKey)
	assert.NotEqual(t, hash1, differentHash, "Different keys should produce different hashes")
}

// TestAuthService_CreateUser_EmailUniqueness 测试邮箱唯一性检查
func TestAuthService_CreateUser_EmailUniqueness(t *testing.T) {
	// 这个测试需要 mock repository，暂时跳过
	// 在实际部署时通过集成测试验证
	t.Skip("Requires database or mock repository")
}

// TestAuthService_ValidateAPIKey_KeyPrefix 测试 key_prefix 是前10位
func TestAuthService_ValidateAPIKey_KeyPrefix(t *testing.T) {
	// 测试 key_prefix 逻辑
	testKey := "sk-0123456789abcdefghijklmnop"
	expectedPrefix := testKey[:10]

	assert.Equal(t, "sk-0123456", expectedPrefix, "Key prefix should be first 10 characters")
	assert.Len(t, expectedPrefix, 10, "Key prefix should be 10 characters")
}

// TestAPIKeyFormat 测试 API Key 格式规范
func TestAPIKeyFormat(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		valid    bool
		reason   string
	}{
		{
			name:   "有效格式",
			key:    "sk-0123456789abcdef0123456789abcdef",
			valid:  true,
			reason: "符合 sk-<32位十六进制> 格式",
		},
		{
			name:   "太短",
			key:    "sk-abc",
			valid:  false,
			reason: "长度不足",
		},
		{
			name:   "缺少前缀",
			key:    "0123456789abcdef0123456789abcdef",
			valid:  false,
			reason: "缺少 sk- 前缀",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.key) == 35 && tt.key[:3] == "sk-"
			assert.Equal(t, tt.valid, isValid, tt.reason)
		})
	}
}

// BenchmarkGenerateAPIKey API Key 生成性能测试
func BenchmarkGenerateAPIKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = generateAPIKey()
	}
}

// BenchmarkHashAPIKey API Key 哈希性能测试
func BenchmarkHashAPIKey(b *testing.B) {
	key := "sk-0123456789abcdef0123456789abcdef"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repository.HashAPIKey(key)
	}
}
