package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUserSuccess(t *testing.T) {
	// Arrange
	id := uuid.New()

	// Act
	user, err := NewUser(id, "  alice  ", "  alice@example.com  ", "hashed-password")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != id {
		t.Fatalf("expected id %s, got %s", id, user.ID)
	}
	if user.Username != "alice" {
		t.Fatalf("expected trimmed username, got %q", user.Username)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("expected trimmed email, got %q", user.Email)
	}
	if user.HashedPassword != "hashed-password" {
		t.Fatalf("expected hashed password to be preserved, got %q", user.HashedPassword)
	}
}

func TestNewUserReturnsErrorWhenUsernameIsEmpty(t *testing.T) {
	// Arrange
	id := uuid.New()

	// Act
	user, err := NewUser(id, "   ", "alice@example.com", "hashed-password")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func TestNewUserReturnsErrorWhenEmailIsEmpty(t *testing.T) {
	// Arrange
	id := uuid.New()

	// Act
	user, err := NewUser(id, "alice", "   ", "hashed-password")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func TestNewSessionSuccess(t *testing.T) {
	// Arrange
	id := uuid.New()
	userID := uuid.New()

	// Act
	session, err := NewSession(id, userID)

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
	if session.ID != id {
		t.Fatalf("expected id %s, got %s", id, session.ID)
	}
	if session.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, session.UserID)
	}
}

func TestNewSessionReturnsErrorWhenIDIsEmpty(t *testing.T) {
	// Arrange
	userID := uuid.New()

	// Act
	session, err := NewSession(uuid.Nil, userID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if session != nil {
		t.Fatalf("expected nil session, got %+v", session)
	}
}

func TestNewSessionReturnsErrorWhenUserIDIsEmpty(t *testing.T) {
	// Arrange
	id := uuid.New()

	// Act
	session, err := NewSession(id, uuid.Nil)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if session != nil {
		t.Fatalf("expected nil session, got %+v", session)
	}
}
