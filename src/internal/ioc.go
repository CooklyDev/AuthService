package internal

import (
	"context"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/CooklyDev/AuthService/internal/usecases"
)

type Container struct {
	logger domain.Logger
	hasher usecases.PasswordHasher
}

func NewContainer(
	logger domain.Logger,
	hasher usecases.PasswordHasher,
) *Container {
	return &Container{
		logger: logger,
		hasher: hasher,
	}
}

func (c *Container) initPostgresUnitOfWork(ctx context.Context) (*adapters.UnitOfWorkPostgres, error) {
	host := LookupEnvRequired("DB_HOST")
	portStr := LookupEnvRequired("DB_PORT")
	user := LookupEnvRequired("DB_USER")
	password := LookupEnvRequired("DB_PASSWORD")
	database := LookupEnvRequired("DB_NAME")

	var port uint16
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		c.logger.Error(
			fmt.Sprintf(
				"init postgres pool failed: operation=parse_port error=%s",
				err.Error(),
			),
		)

		return nil, err
	}

	pool, err := adapters.NewPostgresPool(
		ctx,
		c.logger,
		host,
		port,
		user,
		password,
		database,
	)
	if err != nil {
		return nil, err
	}

	return adapters.NewUnitOfWorkPostgres(pool, c.logger), nil
}

func (c *Container) CreateAuthService(ctx context.Context) (*usecases.AuthService, error) {
	uow, err := c.initPostgresUnitOfWork(ctx)
	if err != nil {
		return nil, err
	}

	return &usecases.AuthService{
		Logger: c.logger,
		Hasher: c.hasher,
		UoW:    uow,
	}, nil
}
