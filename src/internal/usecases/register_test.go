package usecases

import (
	"strings"
	"testing"

	"github.com/CooklyDev/AuthService/internal/domain"
)

func TestRegisterSuccess(t *testing.T) {
	// Arrange
	repo := &userRepoStub{}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	uow := &uowStub{}
	service := AuthService{
		userRepo: repo,
		logger:   logger,
		hasher:   hasher,
		uow:      uow,
	}

	// Act
	err := service.Register("alice", "alice@example.com", "password")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRegisterReturnsErrorWhenEmailAlreadyExists(t *testing.T) {
	// Arrange
	repo := &userRepoStub{
		user: &domain.User{Email: "alice@example.com"},
	}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	uow := &uowStub{}
	service := AuthService{
		userRepo: repo,
		logger:   logger,
		hasher:   hasher,
		uow:      uow,
	}

	// Act
	err := service.Register("alice", "alice@example.com", "password")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid email") {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}
