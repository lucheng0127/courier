package password

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}

	if hash == password {
		t.Error("HashPassword() returned unhashed password")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "wrong password",
			password: "wrongPassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "password against empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyPassword(tt.password, tt.hash); got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "testPassword123"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("HashPassword() error = %v, %v", err1, err2)
	}

	// 每次哈希应该产生不同的结果（因为 bcrypt 使用随机 salt）
	if hash1 == hash2 {
		t.Error("HashPassword() produced identical hashes (salt should be random)")
	}

	// 但两个哈希都应该能验证原密码
	if !VerifyPassword(password, hash1) {
		t.Error("VerifyPassword() failed for hash1")
	}
	if !VerifyPassword(password, hash2) {
		t.Error("VerifyPassword() failed for hash2")
	}
}
