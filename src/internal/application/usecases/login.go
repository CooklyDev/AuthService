package usecases

import (
	"fmt"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

func (service AuthService) LocalLogin(email string, password string) (*uuid.UUID, error) {
	email = domain.NormalizeEmail(email)
	maskedEmail := domain.MaskEmail(email)

	if err := service.UoW.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		rollbackErr := service.UoW.Rollback()
		if rollbackErr != nil {
			service.Logger.Debug(
				fmt.Sprintf(
					"login rollback skipped or failed: email=%s error=%s",
					maskedEmail,
					rollbackErr.Error(),
				),
			)
		}
	}()

	service.Logger.Info(
		fmt.Sprintf(
			"login started: email=%s",
			maskedEmail,
		),
	)

	authIdentityRepo := service.UoW.AuthIdentityRepository()
	authIdentity, err := authIdentityRepo.GetByEmail(email)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"login failed: get auth identity by email: email=%s error=%s",
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}
	if authIdentity == nil {
		service.Logger.Warn(
			fmt.Sprintf(
				"login failed: invalid credentials: email=%s",
				maskedEmail,
			),
		)

		return nil, domain.NewBusinessRuleError("invalid credentials")
	}

	passwordMatches, err := service.Hasher.Compare(password, authIdentity.PasswordHash)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"login failed: password compare: email=%s error=%s",
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	if !passwordMatches {
		service.Logger.Warn(
			fmt.Sprintf(
				"login failed: invalid credentials: email=%s",
				maskedEmail,
			),
		)

		return nil, domain.NewBusinessRuleError("invalid credentials")
	}

	session, err := domain.NewSession(uuid.New(), authIdentity.UserID)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"login failed: create session: user_id=%s email=%s error=%s",
				authIdentity.UserID,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	err = service.UoW.Commit()
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"login failed: commit transaction: user_id=%s email=%s error=%s",
				authIdentity.UserID,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	sessionRepo := service.UoW.SessionRepository()
	err = sessionRepo.Add(session)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"login failed: add session after commit: user_id=%s email=%s error=%s",
				authIdentity.UserID,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	service.Logger.Info(
		fmt.Sprintf(
			"login completed: user_id=%s email=%s",
			authIdentity.UserID,
			maskedEmail,
		),
	)

	return &session.ID, nil
}
