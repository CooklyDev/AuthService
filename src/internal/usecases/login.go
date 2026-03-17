package usecases

import (
	"errors"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

func (service AuthService) Login(email string, password string) (*domain.Session, error) {
	maskedEmail := domain.MaskEmail(email)

	service.logger.Info(
		fmt.Sprintf(
			"login started: email=%s",
			maskedEmail,
		),
	)

	user, err := service.userRepo.GetByEmail(email)
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"login failed: get user by email: email=%s error=%s",
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}
	if user == nil {
		service.logger.Warn(
			fmt.Sprintf(
				"login failed: invalid credentials: email=%s",
				maskedEmail,
			),
		)

		return nil, errors.New("invalid credentials")
	}

	passwordMatches, err := service.hasher.Compare(password, user.HashedPassword)
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"login failed: password compare: email=%s error=%s",
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	if !passwordMatches {
		service.logger.Warn(
			fmt.Sprintf(
				"login failed: invalid credentials: email=%s",
				maskedEmail,
			),
		)

		return nil, errors.New("invalid credentials")
	}

	session, err := domain.NewSession(uuid.New(), user.ID)
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"login failed: create session: user_id=%s email=%s error=%s",
				user.ID,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	err = service.sessionRepo.Add(session)
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"login failed: add session: user_id=%s email=%s error=%s",
				user.ID,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	err = service.uow.Commit()
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"login failed: commit transaction: user_id=%s email=%s error=%s",
				user.ID,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	service.logger.Info(
		fmt.Sprintf(
			"login completed: user_id=%s email=%s",
			user.ID,
			maskedEmail,
		),
	)

	return session, nil
}
