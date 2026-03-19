package adapters

import (
	"github.com/CooklyDev/AuthService/internal/application"
	"github.com/google/uuid"
)

type IdProvider struct {
	session     uuid.UUID
	sessionRepo application.SessionRepo
}

var _ application.IdentityProvider = (*IdProvider)(nil)

func NewIdProvider(session uuid.UUID, sessionRepo application.SessionRepo) *IdProvider {
	return &IdProvider{
		session:     session,
		sessionRepo: sessionRepo,
	}
}

func (p *IdProvider) GetUserId() *uuid.UUID {
	if p.session == uuid.Nil {
		return nil
	}

	session, err := p.sessionRepo.GetSession(p.session)
	if err != nil || session == nil {
		return nil
	}

	return &session.UserID
}
