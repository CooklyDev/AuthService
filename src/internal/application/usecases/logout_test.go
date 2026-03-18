package usecases

import (
	"testing"
)

func TestLogoutSuccess(t *testing.T) {
	// Arrange
	uow := newUoWStub()
	logger := &loggerStub{}
	service := AuthService{
		Logger: logger,
		Hasher: &hasherStub{},
		UoW:    uow,
	}

	// Act
	err := service.Logout("some-session-id")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestLogoutSuccessWithWhitespaceSessionID(t *testing.T) {
	// Arrange
	uow := newUoWStub()
	logger := &loggerStub{}
	service := AuthService{
		Logger: logger,
		Hasher: &hasherStub{},
		UoW:    uow,
	}

	// Act
	err := service.Logout("  some-session-id  ")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
