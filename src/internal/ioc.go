package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/application"
	applicationusecases "github.com/CooklyDev/AuthService/internal/application/usecases"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	logger         domain.Logger
	hasher         application.PasswordHasher
	postgresPool   *pgxpool.Pool
	redisClient    *redis.Client
	PostgresConfig *PostgresConfig
	RedisConfig    *RedisConfig
	AppConfig      *AppConfig
}

func NewContainer(
	logger domain.Logger,
	hasher application.PasswordHasher,
	ctx context.Context,
	postgresConfig *PostgresConfig,
	redisConfig *RedisConfig,
	appConfig *AppConfig,
) (*Container, error) {
	pool, err := newPostgresPool(ctx, logger, *postgresConfig)
	if err != nil {
		return nil, err
	}
	redisClient, err := newRedisClient(ctx, logger, *redisConfig)
	if err != nil {
		return nil, err
	}

	return &Container{
		logger:         logger,
		hasher:         hasher,
		postgresPool:   pool,
		redisClient:    redisClient,
		PostgresConfig: postgresConfig,
		RedisConfig:    redisConfig,
		AppConfig:      appConfig,
	}, nil
}

func (c *Container) Close() {
	if c.redisClient != nil {
		c.logger.Info("redis client closing: dependency=redis")
		if err := c.redisClient.Close(); err != nil {
			c.logger.Warn(
				fmt.Sprintf(
					"redis client close failed: dependency=redis error=%s",
					err.Error(),
				),
			)
		}
		c.redisClient = nil
		c.logger.Info("redis client closed: dependency=redis")
	}

	if c.postgresPool == nil {
		return
	}

	c.logger.Info("postgres pool closing: dependency=postgres")
	c.postgresPool.Close()
	c.postgresPool = nil
	c.logger.Info("postgres pool closed: dependency=postgres")
}

func (c *Container) GetAuthService(sessionID uuid.UUID) (*applicationusecases.AuthService, error) {
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

	uow := adapters.NewUnitOfWorkApp(c.postgresPool, c.redisClient, c.logger, c.AppConfig.SessionTTL, c.AppConfig.SessionPrefix)
	provider := adapters.NewIdProvider(sessionID, uow.SessionRepository())

	return c.CreateAuthService(uow, provider), nil
}

func newPostgresPool(ctx context.Context, logger domain.Logger, config PostgresConfig) (*pgxpool.Pool, error) {
	pool, err := adapters.NewPostgresPool(
		ctx,
		logger,
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func newRedisClient(ctx context.Context, logger domain.Logger, config RedisConfig) (*redis.Client, error) {
	return adapters.NewRedisClient(ctx, logger, config.Host, config.Port, config.Password)
}

func (c *Container) CreateAuthService(uow *adapters.UnitOfWorkApp, idProvider application.IdentityProvider) *applicationusecases.AuthService {
	return &applicationusecases.AuthService{
		Logger:      c.logger,
		Hasher:      c.hasher,
		UoW:         uow,
		IdProvdider: idProvider,
	}
}
