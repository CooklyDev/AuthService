package usecases

import "github.com/CooklyDev/AuthService/internal/domain"

type userRepoStub struct {
	user *domain.User
}

func (stub *userRepoStub) Add(*domain.User) error {
	return nil
}

func (stub *userRepoStub) GetByEmail(string) (*domain.User, error) {
	return stub.user, nil
}

type sessionRepoStub struct{}

func (stub *sessionRepoStub) Add(session *domain.Session) error {
	return nil
}

func (stub *sessionRepoStub) Delete(string) error {
	return nil
}

type hasherStub struct{}

func (stub *hasherStub) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (stub *hasherStub) Compare(password string, hashedPassword string) (bool, error) {
	return "hashed-"+password == hashedPassword, nil
}

type loggerStub struct{}

func (stub *loggerStub) Debug(string) {}

func (stub *loggerStub) Info(string) {}

func (stub *loggerStub) Warn(string) {}

func (stub *loggerStub) Error(string) {}

type uowStub struct{}

func (stub *uowStub) Commit() error {
	return nil
}

func (stub *uowStub) Rollback() error {
	return nil
}
