package usecases

import (
	"errors"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

func (service AuthService) Login(email string, password string) (*domain.Session, error) {
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

	userRepo := service.UoW.UserRepository()
	user, err := userRepo.GetByEmail(email)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"login failed: get user by email: email=%s error=%s",
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}
	if user == nil {
		service.Logger.Warn(
			fmt.Sprintf(
				"login failed: invalid credentials: email=%s",
				maskedEmail,
			),
		)

		return nil, errors.New("invalid credentials")
	}

	passwordMatches, err := service.Hasher.Compare(password, user.HashedPassword)
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

		return nil, errors.New("invalid credentials")
	}

	session, err := domain.NewSession(uuid.New(), user.ID)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"login failed: create session: user_id=%s email=%s error=%s",
				user.ID,
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
				"login failed: add session: user_id=%s email=%s error=%s",
				user.ID,
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
				user.ID,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	service.Logger.Info(
		fmt.Sprintf(
			"login completed: user_id=%s email=%s",
			user.ID,
			maskedEmail,
		),
	)

	return session, nil
}
