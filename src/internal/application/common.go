package application

import (
	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password string, hashedPassword string) (bool, error)
}

type UnitOfWork interface {
	Begin() error
	Commit() error
	Rollback() error
	UserRepository() UserRepo
	AuthIdentityRepository() AuthIdentityRepo
	SessionRepository() SessionRepo
}

type IdentityProvider interface {
	GetUserId() *uuid.UUID
}
