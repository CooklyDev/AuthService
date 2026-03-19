package usecases

import (
	"fmt"
	"github.com/google/uuid"

	"github.com/CooklyDev/AuthService/internal/domain"
)

func (service AuthService) Logout(sessionID uuid.UUID) error {
	if err := service.UoW.Begin(); err != nil {
		return err
	}

	defer func() {
		rollbackErr := service.UoW.Rollback()
		if rollbackErr != nil {
			service.Logger.Debug(
				fmt.Sprintf(
					"logout rollback skipped or failed: error=%s",
					rollbackErr.Error(),
				),
			)
		}
	}()

	service.Logger.Info("logout started")

	userIdentity := service.IdProvdider.GetUserId()
	if userIdentity == nil {
		return domain.NewBusinessRuleError("user not authenticated")
	}

	sessionRepo := service.UoW.SessionRepository()
	userSessions, err := sessionRepo.GetUserSessions(*userIdentity)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"logout failed: get user sessions: error=%s",
				err.Error(),
			),
		)

		return err
	}

	sessionFound := false
	for _, session := range userSessions {
		if session.ID == sessionID {
			sessionFound = true
			break
		}
	}

	if !sessionFound {
		return nil
	}

	err = service.UoW.Commit()
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"logout failed: commit transaction: error=%s",
				err.Error(),
			),
		)

		return err
	}

	err = sessionRepo.Delete(sessionID)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"logout failed: delete session: error=%s",
				err.Error(),
			),
		)

		return err
	}

	service.Logger.Info("logout completed")

	return nil
}
