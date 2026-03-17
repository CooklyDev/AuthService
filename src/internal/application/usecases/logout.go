package usecases

import (
	"fmt"
	"strings"
)

func (service AuthService) Logout(sessionID string) error {
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

	sessionID = strings.TrimSpace(sessionID)

	sessionRepo := service.UoW.SessionRepository()
	err := sessionRepo.Delete(sessionID)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"logout failed: delete session: error=%s",
				err.Error(),
			),
		)

		return err
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

	service.Logger.Info("logout completed")

	return nil
}
