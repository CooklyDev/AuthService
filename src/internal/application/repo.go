package application

import (
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

type UserRepo interface {
	Add(user *domain.User) error
}

type AuthIdentityRepo interface {
	Add(identity *domain.AuthIdentity) error
	GetByEmail(email string) (*domain.AuthIdentity, error)
}

type SessionRepo interface {
	GetUserSessions(userID uuid.UUID) ([]*domain.Session, error)
	GetSession(sessionID uuid.UUID) (*domain.Session, error)
}

type SessionOutboxRepo interface {
	AddSessionCreated(session *domain.Session) error
	AddSessionDeleted(session *domain.Session) error
}
