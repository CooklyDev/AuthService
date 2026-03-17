package adapters

import "github.com/CooklyDev/AuthService/internal/application"

type StubHasher struct{}

var _ application.PasswordHasher = (*StubHasher)(nil)

func NewStubHasher() *StubHasher {
	return &StubHasher{}
}

func (hasher *StubHasher) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (hasher *StubHasher) Compare(password string, hashedPassword string) (bool, error) {
	return "hashed-"+password == hashedPassword, nil
}
