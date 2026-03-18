package usecases

import (
	"errors"
	"testing"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

func TestLoginSuccess(t *testing.T) {
	// Arrange
	userID := uuid.New()
	uow := newUoWStub()
	uow.userRepo.user = &domain.User{
		ID:             userID,
		Email:          "alice@example.com",
		HashedPassword: "hashed-password",
	}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	service := AuthService{
		Logger: logger,
		Hasher: hasher,
		UoW:    uow,
	}

	// Act
	sessionId, err := service.Login("alice@example.com", "password")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if sessionId == nil {
		t.Fatal("expected session ID, got nil")
	}
}

func TestLoginReturnsErrorWhenUserDoesNotExist(t *testing.T) {
	// Arrange
	uow := newUoWStub()
	hasher := &hasherStub{}
	logger := &loggerStub{}
	service := AuthService{
		Logger: logger,
		Hasher: hasher,
		UoW:    uow,
	}

	// Act
	_, err := service.Login("alice@example.com", "password")

	// Assert
	if !errors.Is(err, domain.ErrBusinessRule) {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestLoginReturnsErrorWhenPasswordIsInvalid(t *testing.T) {
	// Arrange
	uow := newUoWStub()
	uow.userRepo.user = &domain.User{
		ID:             uuid.New(),
		Email:          "alice@example.com",
		HashedPassword: "hashed-secret",
	}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	service := AuthService{
		Logger: logger,
		Hasher: hasher,
		UoW:    uow,
	}

	// Act
	_, err := service.Login("alice@example.com", "password")

	// Assert
	if !errors.Is(err, domain.ErrBusinessRule) {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}
