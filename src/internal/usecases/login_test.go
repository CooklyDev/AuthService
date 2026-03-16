package usecases

import (
	"strings"
	"testing"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

func TestLoginSuccess(t *testing.T) {
	// Arrange
	userID := uuid.New()
	repo := &userRepoStub{
		user: &domain.User{
			ID:             userID,
			Email:          "alice@example.com",
			HashedPassword: "hashed-password",
		},
	}
	sessionRepo := &sessionRepoStub{}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	service := AuthService{
		userRepo:    repo,
		sessionRepo: sessionRepo,
		logger:      logger,
		hasher:      hasher,
	}

	// Act
	session, err := service.Login("alice@example.com", "password")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
	if session.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, session.UserID)
	}
}

func TestLoginReturnsErrorWhenUserDoesNotExist(t *testing.T) {
	// Arrange
	repo := &userRepoStub{}
	sessionRepo := &sessionRepoStub{}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	service := AuthService{
		userRepo:    repo,
		sessionRepo: sessionRepo,
		logger:      logger,
		hasher:      hasher,
	}

	// Act
	_, err := service.Login("alice@example.com", "password")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestLoginReturnsErrorWhenPasswordIsInvalid(t *testing.T) {
	// Arrange
	repo := &userRepoStub{
		user: &domain.User{
			ID:             uuid.New(),
			Email:          "alice@example.com",
			HashedPassword: "hashed-secret",
		},
	}
	sessionRepo := &sessionRepoStub{}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	service := AuthService{
		userRepo:    repo,
		sessionRepo: sessionRepo,
		logger:      logger,
		hasher:      hasher,
	}

	// Act
	_, err := service.Login("alice@example.com", "password")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}
