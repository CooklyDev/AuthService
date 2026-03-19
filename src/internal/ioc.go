package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/application"
	applicationusecases "github.com/CooklyDev/AuthService/internal/application/usecases"
	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	logger       domain.Logger
	hasher       application.PasswordHasher
	postgresPool *pgxpool.Pool
	redisClient  *redis.Client
	sessionTTL   time.Duration
}

func NewContainer(
	logger domain.Logger,
	hasher application.PasswordHasher,
	ctx context.Context,
) (*Container, error) {
	var sessionTTL time.Duration

	sessionTTLStr, exists := LookupEnvOptional("SESSION_TTL")
	if !exists {
		sessionTTL = time.Hour * 24
	} else {
		parsedTTL, err := time.ParseDuration(sessionTTLStr)
		if err != nil {
			logger.Error(
				fmt.Sprintf(
					"init container failed: operation=parse_session_ttl error=%s",
					err.Error(),
				),
			)
			return nil, err
		}
		sessionTTL = parsedTTL
	}

	pool, err := newPostgresPool(ctx, logger)
	if err != nil {
		return nil, err
	}
	redisClient, err := newRedisClient(ctx, logger)
	if err != nil {
		return nil, err
	}

	return &Container{
		logger:       logger,
		hasher:       hasher,
		postgresPool: pool,
		redisClient:  redisClient,
		sessionTTL:   sessionTTL,
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

	uow := adapters.NewUnitOfWorkApp(c.postgresPool, c.redisClient, c.logger, c.sessionTTL)
	provider := adapters.NewIdProvider(sessionID, uow.SessionRepository())

	return c.CreateAuthService(uow, provider), nil
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

func newRedisClient(ctx context.Context, logger domain.Logger) (*redis.Client, error) {
	host := LookupEnvRequired("REDIS_HOST")
	portStr := LookupEnvRequired("REDIS_PORT")
	password, _ := LookupEnvOptional("REDIS_PASSWORD")

	var port uint16
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"init redis client failed: operation=parse_port error=%s",
				err.Error(),
			),
		)

		return nil, err
	}

	return adapters.NewRedisClient(ctx, logger, host, port, password)
}

func (c *Container) CreateAuthService(uow *adapters.UnitOfWorkApp, idProvider application.IdentityProvider) *applicationusecases.AuthService {
	return &applicationusecases.AuthService{
		Logger:      c.logger,
		Hasher:      c.hasher,
		UoW:         uow,
		IdProvdider: idProvider,
	}
}
