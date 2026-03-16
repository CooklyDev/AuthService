package domain

import "testing"

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "masks email with local part longer than one character",
			email:    "alice@example.com",
			expected: "a***@example.com",
		},
		{
			name:     "returns stars when local part is one character",
			email:    "a@example.com",
			expected: "***",
		},
		{
			name:     "returns stars when email has no at sign",
			email:    "alice.example.com",
			expected: "***",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			email := test.email

			// Act
			maskedEmail := MaskEmail(email)

			// Assert
			if maskedEmail != test.expected {
				t.Fatalf("expected masked email %q, got %q", test.expected, maskedEmail)
			}
		})
	}
}
