package adapters

import (
	"context"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(
	ctx context.Context,
	logger domain.Logger,
	host string,
	port uint16,
	user string,
	password string,
	database string,
	sslMode string,
) (*pgxpool.Pool, error) {
	logger.Info(
		fmt.Sprintf(
			"postgres pool initialization started: dependency=postgres host=%s port=%d db=%s",
			host,
			port,
			database,
		),
	)

	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		database,
		sslMode,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"postgres pool initialization failed: dependency=postgres operation=parse_config error=%s",
				err.Error(),
			),
		)

		return nil, NewAdapterError("parse postgres config", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"postgres pool initialization failed: dependency=postgres operation=create_pool error=%s",
				err.Error(),
			),
		)

		return nil, NewAdapterError("create postgres pool", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.Error(
			fmt.Sprintf(
				"postgres pool initialization failed: dependency=postgres operation=ping error=%s",
				err.Error(),
			),
		)

		return nil, NewAdapterError("ping postgres pool", err)
	}

	logger.Info(
		fmt.Sprintf(
			"postgres pool initialization completed: dependency=postgres host=%s port=%d db=%s",
			host,
			port,
			database,
		),
	)

	return pool, nil
}
