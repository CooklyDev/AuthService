package usecases

import (
	"errors"
	"testing"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

func TestResolveSessionSuccess(t *testing.T) {
	// Arrange
	sessionID := uuid.New()
	userID := uuid.New()
	uow := newUoWStub()
	uow.sessionRepo.session = &domain.Session{
		ID:     sessionID,
		UserID: userID,
	}
	service := AuthService{
		UoW: uow,
	}

	// Act
	result, err := service.ResolveSession(sessionID)

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected resolve dto, got nil")
	}
	if result.SessionID != sessionID {
		t.Fatalf("expected session id %s, got %s", sessionID, result.SessionID)
	}
	if result.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, result.UserID)
	}
}

func TestResolveSessionReturnsError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("failed to load session")
	uow := newUoWStub()
	uow.sessionRepo.err = expectedErr
	service := AuthService{
		UoW: uow,
	}

	// Act
	result, err := service.ResolveSession(uuid.New())

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}
