package usecases

import (
	"errors"
	"testing"

	"github.com/CooklyDev/AuthService/internal/domain"
)

func TestRegisterSuccess(t *testing.T) {
	// Arrange
	hasher := &hasherStub{}
	logger := &loggerStub{}
	uow := newUoWStub()
	service := AuthService{
		Logger: logger,
		Hasher: hasher,
		UoW:    uow,
	}

	// Act
	err := service.Register("alice", "alice@example.com", "password1")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRegisterReturnsErrorWhenEmailAlreadyExists(t *testing.T) {
	// Arrange
	uow := newUoWStub()
	uow.userRepo.user = &domain.User{Email: "alice@example.com"}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	service := AuthService{
		Logger: logger,
		Hasher: hasher,
		UoW:    uow,
	}

	// Act
	err := service.Register("alice", "alice@example.com", "password1")

	// Assert
	if !errors.Is(err, domain.ErrBusinessRule) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}
