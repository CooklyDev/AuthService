package application

import (
	"github.com/CooklyDev/AuthService/internal/domain"
)

type UserRepo interface {
	Add(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
}

type SessionRepo interface {
	Add(session *domain.Session) error
	Delete(sessionID string) error
}
