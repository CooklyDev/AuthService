package usecases

import "github.com/google/uuid"

type ResolveDTO struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
}

func (service AuthService) ResolveSession(sessionID uuid.UUID) (*ResolveDTO, error) {
	sessionRepo := service.UoW.SessionRepository()
	session, err := sessionRepo.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	return &ResolveDTO{
		SessionID: session.ID,
		UserID:    session.UserID,
	}, nil
}
