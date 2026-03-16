package usecases

import (
	"errors"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    UserRepo
	sessionRepo SessionRepo
	logger      domain.Logger
	hasher      PasswordHasher
}

func (service AuthService) Register(username string, email string, password string) error {
	maskedEmail := domain.MaskEmail(email)

	service.logger.Info(
		fmt.Sprintf(
			"register started: username=%s email=%s",
			username,
			maskedEmail,
		),
	)

	if !(domain.ValidatePassword(password)) {
		service.logger.Warn(
			fmt.Sprintf(
				"register failed: invalid password: username=%s email=%s",
				username,
				maskedEmail,
			),
		)

		return errors.New("invalid password")
	}

	hashedPassword, err := service.hasher.Hash(password)
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"register failed: password hashing: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return err
	}

	oldUser, err := service.userRepo.GetByEmail(email)
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"register failed: get user by email: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return err
	}
	if oldUser != nil {
		service.logger.Warn(
			fmt.Sprintf(
				"register failed: email already exists: username=%s email=%s",
				username,
				maskedEmail,
			),
		)

		return errors.New("invalid email: user with such email already exists")
	}

	user, err := domain.NewUser(uuid.New(), username, email, hashedPassword)
	if err != nil {
		service.logger.Warn(
			fmt.Sprintf(
				"register failed: invalid user data: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return err
	}

	err = service.userRepo.Add(user)
	if err != nil {
		service.logger.Error(
			fmt.Sprintf(
				"register failed: add user: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return err
	}

	service.logger.Info(
		fmt.Sprintf(
			"register completed: user_id=%s username=%s email=%s",
			user.ID,
			user.Username,
			maskedEmail,
		),
	)

	return nil
}
