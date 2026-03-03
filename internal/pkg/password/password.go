package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// bcryptCostFactor bcrypt cost factor
	// 推荐值：12-14，越高越安全但越慢
	bcryptCostFactor = 12
)

// HashPassword 对密码进行哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCostFactor)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword 验证密码是否匹配哈希
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
