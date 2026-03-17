package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/CooklyDev/AuthService/internal/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	logger       domain.Logger
	hasher       usecases.PasswordHasher
	postgresPool *pgxpool.Pool
}

func NewContainer(
	logger domain.Logger,
	hasher usecases.PasswordHasher,
	ctx context.Context,
) (*Container, error) {
	pool, err := newPostgresPool(ctx, logger)
	if err != nil {
		return nil, err
	}

	return &Container{
		logger:       logger,
		hasher:       hasher,
		postgresPool: pool,
	}, nil
}

func (c *Container) Close() {
	if c.postgresPool == nil {
		return
	}

	c.logger.Info("postgres pool closing: dependency=postgres")
	c.postgresPool.Close()
	c.postgresPool = nil
	c.logger.Info("postgres pool closed: dependency=postgres")
}

func (c *Container) GetAuthService() (*usecases.AuthService, error) {
	if c.postgresPool == nil {
		err := errors.New("postgres pool is not initialized")
		c.logger.Error(
			fmt.Sprintf(
				"create auth service failed: dependency=postgres error=%s",
				err.Error(),
			),
		)

		return nil, err
	}

	uow := adapters.NewUnitOfWorkPostgres(c.postgresPool, c.logger)

	return c.CreateAuthService(uow), nil
}

func newPostgresPool(ctx context.Context, logger domain.Logger) (*pgxpool.Pool, error) {
	host := LookupEnvRequired("DB_HOST")
	portStr := LookupEnvRequired("DB_PORT")
	user := LookupEnvRequired("DB_USER")
	password := LookupEnvRequired("DB_PASSWORD")
	database := LookupEnvRequired("DB_NAME")

	var port uint16
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"init postgres pool failed: operation=parse_port error=%s",
				err.Error(),
			),
		)

		return nil, err
	}

	pool, err := adapters.NewPostgresPool(
		ctx,
		logger,
		host,
		port,
		user,
		password,
		database,
	)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (c *Container) CreateAuthService(uow *adapters.UnitOfWorkPostgres) *usecases.AuthService {
	return &usecases.AuthService{
		Logger: c.logger,
		Hasher: c.hasher,
		UoW:    uow,
	}
}
