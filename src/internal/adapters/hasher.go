package adapters

import (
	"github.com/CooklyDev/AuthService/internal/application"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

var _ application.PasswordHasher = (*BcryptHasher)(nil)

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (hasher *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (hasher *BcryptHasher) Compare(password string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, nil
}
