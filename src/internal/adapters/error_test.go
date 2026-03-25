package adapters

import (
	"errors"
	"testing"
)

func TestNewInvariantError(t *testing.T) {
	// Arrange
	cause := errors.New("invalid stored value")

	// Act
	err := NewInvariantError("get session", cause)

	// Assert
	if !errors.Is(err, ErrInvariant) {
		t.Fatalf("expected invariant error, got %v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause, got %v", err)
	}
}

func TestNewServerError(t *testing.T) {
	// Arrange
	cause := errors.New("redis unavailable")

	// Act
	err := NewServerError("get session", cause)

	// Assert
	if !errors.Is(err, ErrServer) {
		t.Fatalf("expected server error, got %v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause, got %v", err)
	}
}
