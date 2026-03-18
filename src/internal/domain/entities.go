package domain

import (
	"strings"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
}

type Provider int

const (
	ProviderLocal Provider = iota
)

type AuthIdentity struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Provider     Provider
	ProviderID   string // for oauth providers only
	Email        string // for password provider only
	PasswordHash string // for password provider only
}

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func ValidatePassword(password string) bool {
	if len(password) < 5 || len(password) > 100 {
		return false
	}

	hasDigit := false
	hasLetter := false

	for _, ch := range password {
		if ch >= '0' && ch <= '9' {
			hasDigit = true
		}
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			hasLetter = true
		}
	}

	return hasDigit && hasLetter
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func NewUser(id uuid.UUID, username string) (*User, error) {
	username = strings.TrimSpace(username)

	if id == uuid.Nil {
		return nil, NewBusinessRuleError("invalid argument: user id is required")
	}
	if username == "" {
		return nil, NewBusinessRuleError("invalid argument: username is required")
	}

	return &User{
		ID:       id,
		Username: username,
	}, nil
}

func NewAuthIdentity(id uuid.UUID, userID uuid.UUID, provider Provider, providerID string, email string, passwordHash string) (*AuthIdentity, error) {
	if id == uuid.Nil {
		return nil, NewBusinessRuleError("invalid argument: auth identity id is required")
	}
	if userID == uuid.Nil {
		return nil, NewBusinessRuleError("invalid argument: user id is required")
	}
	switch provider {
	case ProviderLocal:
		if providerID != "" {
			return nil, NewBusinessRuleError("invalid argument: provider id should be empty for local provider")
		}
		if email == "" && provider == ProviderLocal {
			return nil, NewBusinessRuleError("invalid argument: email is required for password provider")
		}
		if passwordHash == "" && provider == ProviderLocal {
			return nil, NewBusinessRuleError("invalid argument: password hash is required for password provider")
		}
	default:
		return nil, NewBusinessRuleError("invalid argument: unsupported auth provider")
	}

	return &AuthIdentity{
		ID:           id,
		UserID:       userID,
		Provider:     provider,
		ProviderID:   providerID,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func NewSession(id uuid.UUID, userID uuid.UUID) (*Session, error) {
	if id == uuid.Nil {
		return nil, NewBusinessRuleError("invalid argument: session id is required")
	}
	if userID == uuid.Nil {
		return nil, NewBusinessRuleError("invalid argument: user id is required")
	}

	return &Session{
		ID:     id,
		UserID: userID,
	}, nil
}
