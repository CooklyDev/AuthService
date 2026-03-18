package usecases

import (
	"github.com/CooklyDev/AuthService/internal/application"
	"github.com/CooklyDev/AuthService/internal/domain"
)

type userRepoStub struct {
}

func (stub *userRepoStub) Add(*domain.User) error {
	return nil
}

type authIdentityRepoStub struct {
	identity *domain.AuthIdentity
}

func (stub *authIdentityRepoStub) Add(*domain.AuthIdentity) error {
	return nil
}

func (stub *authIdentityRepoStub) GetByEmail(string) (*domain.AuthIdentity, error) {
	return stub.identity, nil
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

type uowStub struct {
	userRepo         *userRepoStub
	authIdentityRepo *authIdentityRepoStub
	sessionRepo      *sessionRepoStub
}

func newUoWStub() *uowStub {
	return &uowStub{
		userRepo:         &userRepoStub{},
		authIdentityRepo: &authIdentityRepoStub{},
		sessionRepo:      &sessionRepoStub{},
	}
}

func (stub *uowStub) Begin() error {
	return nil
}

func (stub *uowStub) Commit() error {
	return nil
}

func (stub *uowStub) Rollback() error {
	return nil
}

func (stub *uowStub) UserRepository() application.UserRepo {
	return stub.userRepo
}

func (stub *uowStub) AuthIdentityRepository() application.AuthIdentityRepo {
	return stub.authIdentityRepo
}

func (stub *uowStub) SessionRepository() application.SessionRepo {
	return stub.sessionRepo
}
