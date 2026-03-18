package usecases

import (
	"fmt"
	"strings"

	"github.com/CooklyDev/AuthService/internal/application"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
)

type AuthService struct {
	Logger domain.Logger
	Hasher application.PasswordHasher
	UoW    application.UnitOfWork
}

func (service AuthService) LocalRegister(username string, email string, password string) (*domain.Session, error) {
	email = domain.NormalizeEmail(email)
	password = strings.TrimSpace(password)
	maskedEmail := domain.MaskEmail(email)

	if err := service.UoW.Begin(); err != nil {
		return nil, err
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

	authIdentityRepo := service.UoW.AuthIdentityRepository()
	pastAuth, err := authIdentityRepo.GetByEmail(email)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"register failed: get auth identity by email: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}
	if pastAuth != nil {
		service.Logger.Warn(
			fmt.Sprintf(
				"register failed: email already exists: username=%s email=%s",
				username,
				maskedEmail,
			),
		)

		return nil, domain.NewBusinessRuleError("invalid email: user with such email already exists")
	}

	if !(domain.ValidatePassword(password)) {
		service.Logger.Warn(
			fmt.Sprintf(
				"register failed: invalid password: username=%s email=%s",
				username,
				maskedEmail,
			),
		)

		return nil, domain.NewBusinessRuleError("invalid password")
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

		return nil, err
	}

	user, err := domain.NewUser(uuid.New(), username)
	if err != nil {
		service.Logger.Warn(
			fmt.Sprintf(
				"register failed: invalid user data: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}
	authIdentity, err := domain.NewAuthIdentity(uuid.New(), user.ID, domain.ProviderLocal, "", email, hashedPassword)
	if err != nil {
		service.Logger.Warn(
			fmt.Sprintf(
				"register failed: invalid auth identity data: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}

	userRepo := service.UoW.UserRepository()
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

		return nil, err
	}
	err = authIdentityRepo.Add(authIdentity)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"register failed: add auth identity: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}
	sessionRepo := service.UoW.SessionRepository()
	session, err := domain.NewSession(uuid.New(), user.ID)
	if err != nil {
		service.Logger.Warn(
			fmt.Sprintf(
				"register failed: invalid session data: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)

		return nil, err
	}
	err = sessionRepo.Add(session)
	if err != nil {
		service.Logger.Error(
			fmt.Sprintf(
				"register failed: create session: username=%s email=%s error=%s",
				username,
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
				"register failed: commit transaction: username=%s email=%s error=%s",
				username,
				maskedEmail,
				err.Error(),
			),
		)
		return nil, err
	}

	service.Logger.Info(
		fmt.Sprintf(
			"register completed: user_id=%s username=%s email=%s",
			user.ID,
			user.Username,
			maskedEmail,
		),
	)

	return session, nil
}
