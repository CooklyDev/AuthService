package usecases

import (
	"fmt"

	"github.com/CooklyDev/AuthService/internal/application"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

type AuthService struct {
	Logger domain.Logger
	Hasher application.PasswordHasher
	UoW    application.UnitOfWork
}

func (service AuthService) Register(username string, email string, password string) error {
	maskedEmail := domain.MaskEmail(email)

	if err := service.UoW.Begin(); err != nil {
		return err
	}

	defer func() {
		rollbackErr := service.UoW.Rollback()
		if rollbackErr != nil {
			service.Logger.Debug(
				fmt.Sprintf(
					"register rollback skipped or failed: username=%s email=%s error=%s",
					username,
					maskedEmail,
					rollbackErr.Error(),
				),
			)
		}
	}()

	service.Logger.Info(
		fmt.Sprintf(
			"register started: username=%s email=%s",
			username,
			maskedEmail,
		),
	)

	if !(domain.ValidatePassword(password)) {
		service.Logger.Warn(
			fmt.Sprintf(
				"register failed: invalid password: username=%s email=%s",
				username,
				maskedEmail,
			),
		)

		return domain.NewBusinessRuleError("invalid password")
	}

	hashedPassword, err := service.Hasher.Hash(password)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"register failed: password hashing: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return err
	}

	userRepo := service.UoW.UserRepository()
	oldUser, err := userRepo.GetByEmail(email)
	if err != nil {
		service.Logger.Error(
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
		service.Logger.Warn(
			fmt.Sprintf(
				"register failed: email already exists: username=%s email=%s",
				username,
				maskedEmail,
			),
		)

		return domain.NewBusinessRuleError("invalid email: user with such email already exists")
	}

	user, err := domain.NewUser(uuid.New(), username, email, hashedPassword)
	if err != nil {
		service.Logger.Warn(
			fmt.Sprintf(
				"register failed: invalid user data: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return err
	}

	err = userRepo.Add(user)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"register failed: add user: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return err
	}

	err = service.UoW.Commit()
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"register failed: commit transaction: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)
		return err
	}

	service.Logger.Info(
		fmt.Sprintf(
			"register completed: user_id=%s username=%s email=%s",
			user.ID,
			user.Username,
			maskedEmail,
		),
	)

	return nil
}
