package usecases

import (
	"errors"
	"strings"
	"testing"

	"github.com/CooklyDev/AuthService/internal/domain"
)

func TestLocalRegisterSuccess(t *testing.T) {
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
	session, err := service.LocalRegister("alice", "alice@example.com", "password1")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
}

func TestLocalRegisterReturnsErrorWhenEmailAlreadyExists(t *testing.T) {
	// Arrange
	uow := newUoWStub()
	uow.authIdentityRepo.identity = &domain.AuthIdentity{Email: "alice@example.com"}
	hasher := &hasherStub{}
	logger := &loggerStub{}
	service := AuthService{
		Logger: logger,
		Hasher: hasher,
		UoW:    uow,
	}

	// Act
	session, err := service.LocalRegister("alice", "alice@example.com", "password1")

	// Assert
	if !errors.Is(err, domain.ErrBusinessRule) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
	if session != nil {
		t.Fatalf("expected nil session, got %+v", session)
	}
}

func TestLocalRegisterNormalizesPasswordBeforeValidation(t *testing.T) {
	// Arrange
	password := " " + strings.Repeat("a", 99) + "1" + " "
	hasher := &hasherStub{}
	logger := &loggerStub{}
	uow := newUoWStub()
	service := AuthService{
		Logger: logger,
		Hasher: hasher,
		UoW:    uow,
	}

	// Act
	session, err := service.LocalRegister("alice", "alice@example.com", password)

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
}
