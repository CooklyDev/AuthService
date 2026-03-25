package presentation

import (
	"errors"
	"net/http"
	"testing"

	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/domain"
)

func TestMapAppErrorReturnsBusinessRuleResponse(t *testing.T) {
	// Arrange
	err := domain.NewBusinessRuleError("invalid credentials")

	// Act
	status, code, message := MapAppError(err)

	// Assert
	if status != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, status)
	}
	if code != "BUSINESS_RULE_VIOLATION" {
		t.Fatalf("expected code %q, got %q", "BUSINESS_RULE_VIOLATION", code)
	}
	if message != "invalid credentials" {
		t.Fatalf("expected message %q, got %q", "invalid credentials", message)
	}
}

func TestMapAppErrorReturnsAdapterInvariantResponse(t *testing.T) {
	// Arrange
	err := adapters.NewInvariantError("get session", errors.New("invalid stored value"))

	// Act
	status, code, message := MapAppError(err)

	// Assert
	if status != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, status)
	}
	if code != "ADAPTER_INVARIANT_VIOLATION" {
		t.Fatalf("expected code %q, got %q", "ADAPTER_INVARIANT_VIOLATION", code)
	}
	if message != "request violates service invariants" {
		t.Fatalf("expected message %q, got %q", "request violates service invariants", message)
	}
}

func TestMapAppErrorReturnsServerResponse(t *testing.T) {
	// Arrange
	err := adapters.NewServerError("get session", errors.New("redis unavailable"))

	// Act
	status, code, message := MapAppError(err)

	// Assert
	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
	if code != "INTERNAL_SERVER_ERROR" {
		t.Fatalf("expected code %q, got %q", "INTERNAL_SERVER_ERROR", code)
	}
	if message != "internal server error" {
		t.Fatalf("expected message %q, got %q", "internal server error", message)
	}
}
