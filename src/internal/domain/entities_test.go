package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "returns false when password is shorter than minimum length",
			password: "a1b",
			expected: false,
		},
		{
			name:     "returns false when password is longer than maximum length",
			password: strings.Repeat("a", 100) + "1",
			expected: false,
		},
		{
			name:     "returns false when password has no digits",
			password: "abcde",
			expected: false,
		},
		{
			name:     "returns false when password has no letters",
			password: "12345",
			expected: false,
		},
		{
			name:     "returns true when password has letters and digits within valid length",
			password: "abc12",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			password := test.password

			// Act
			actual := ValidatePassword(password)

			// Assert
			if actual != test.expected {
				t.Fatalf("expected %v, got %v for password %q", test.expected, actual, password)
			}
		})
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "trims leading and trailing whitespace and lowercases email",
			email:    "  Alice@Example.COM  ",
			expected: "alice@example.com",
		},
		{
			name:     "preserves internal symbols while lowercasing",
			email:    "First.Last+tag@Example.COM",
			expected: "first.last+tag@example.com",
		},
		{
			name:     "returns empty string for whitespace only email",
			email:    "   ",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			email := test.email

			// Act
			actual := NormalizeEmail(email)

			// Assert
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestNewUserSuccess(t *testing.T) {
	// Arrange
	id := uuid.New()

	// Act
	user, err := NewUser(id, "  alice  ")

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
}

func TestNewUserReturnsErrorWhenIDIsEmpty(t *testing.T) {
	// Arrange
	username := "alice"

	// Act
	user, err := NewUser(uuid.Nil, username)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func TestNewUserReturnsErrorWhenUsernameIsEmpty(t *testing.T) {
	// Arrange
	id := uuid.New()

	// Act
	user, err := NewUser(id, "   ")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func TestNewAuthIdentitySuccessForLocalProvider(t *testing.T) {
	// Arrange
	id := uuid.New()
	userID := uuid.New()

	// Act
	identity, err := NewAuthIdentity(id, userID, ProviderLocal, "", "alice@example.com", "hashed-password")

	// Assert
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if identity == nil {
		t.Fatal("expected auth identity, got nil")
	}
	if identity.ID != id {
		t.Fatalf("expected id %s, got %s", id, identity.ID)
	}
	if identity.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, identity.UserID)
	}
	if identity.Provider != ProviderLocal {
		t.Fatalf("expected provider %d, got %d", ProviderLocal, identity.Provider)
	}
	if identity.Email != "alice@example.com" {
		t.Fatalf("expected email to be preserved, got %q", identity.Email)
	}
	if identity.PasswordHash != "hashed-password" {
		t.Fatalf("expected password hash to be preserved, got %q", identity.PasswordHash)
	}
}

func TestNewAuthIdentityReturnsErrorWhenUserIDIsEmpty(t *testing.T) {
	// Arrange
	id := uuid.New()

	// Act
	identity, err := NewAuthIdentity(id, uuid.Nil, ProviderLocal, "", "alice@example.com", "hashed-password")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if identity != nil {
		t.Fatalf("expected nil auth identity, got %+v", identity)
	}
}

func TestNewAuthIdentityReturnsErrorWhenEmailIsEmptyForLocalProvider(t *testing.T) {
	// Arrange
	id := uuid.New()
	userID := uuid.New()

	// Act
	identity, err := NewAuthIdentity(id, userID, ProviderLocal, "", "", "hashed-password")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if identity != nil {
		t.Fatalf("expected nil auth identity, got %+v", identity)
	}
}

func TestNewAuthIdentityReturnsErrorWhenPasswordHashIsEmptyForLocalProvider(t *testing.T) {
	// Arrange
	id := uuid.New()
	userID := uuid.New()

	// Act
	identity, err := NewAuthIdentity(id, userID, ProviderLocal, "", "alice@example.com", "")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if identity != nil {
		t.Fatalf("expected nil auth identity, got %+v", identity)
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
