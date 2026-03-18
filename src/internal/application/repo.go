package application

import (
	"github.com/CooklyDev/AuthService/internal/domain"
)

type UserRepo interface {
	Add(user *domain.User) error
}

type AuthIdentityRepo interface {
	Add(identity *domain.AuthIdentity) error
	GetByEmail(email string) (*domain.AuthIdentity, error)
}

type SessionRepo interface {
	Add(session *domain.Session) error
	Delete(sessionID string) error
}
