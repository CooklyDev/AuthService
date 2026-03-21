package usecases

import (
	"fmt"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

type ResolveDTO struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
}

func (service AuthService) ResolveSession(sessionID uuid.UUID) (*ResolveDTO, error) {
	service.Logger.Info(
		fmt.Sprintf(
			"resolve session started: session_id=%s",
			sessionID,
		),
	)

	sessionRepo := service.UoW.SessionRepository()
	session, err := sessionRepo.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, domain.NewBusinessRuleError("session not found")
	}

	service.Logger.Info(
		fmt.Sprintf(
			"resolve session: session_id=%s user_id=%s",
			session.ID,
			session.UserID,
		),
	)

	return &ResolveDTO{
		SessionID: session.ID,
		UserID:    session.UserID,
	}, nil
}
