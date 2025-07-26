package internal

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService представляет сервис для работы с паролями
type PasswordService struct{}

// NewPasswordService создает новый сервис для работы с паролями
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword хеширует пароль
func (p *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword проверяет пароль
func (p *PasswordService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
} 