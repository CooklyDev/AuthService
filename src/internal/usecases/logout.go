package usecases

import (
	"fmt"
	"strings"
)

func (service AuthService) Logout(sessionID string) error {
	service.logger.Info("logout started")

	sessionID = strings.TrimSpace(sessionID)

	err := service.sessionRepo.Delete(sessionID)
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"logout failed: delete session: error=%s",
				err.Error(),
			),
		)

		return err
	}

	err = service.uow.Commit()
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"logout failed: commit transaction: error=%s",
				err.Error(),
			),
		)

		return err
	}

	service.logger.Info("logout completed")

	return nil
}
