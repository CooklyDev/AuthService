package domain

import (
	"strings"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Username       string
	Email          string
	HashedPassword string
}

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func ValidatePassword(password string) bool { // TODO: Implement validation
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

func NewUser(id uuid.UUID, username string, email string, password string) (*User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if username == "" {
		return nil, NewBusinessRuleError("invalid argument: username is required")
	}
	if email == "" {
		return nil, NewBusinessRuleError("invalid argument: email is required")
	}

	return &User{
		ID:             id,
		Username:       username,
		Email:          email,
		HashedPassword: password,
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
