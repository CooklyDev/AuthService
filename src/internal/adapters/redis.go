package adapters

import (
	"context"
	"fmt"

	"github.com/CooklyDev/AuthService/internal/domain"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(
	ctx context.Context,
	logger domain.Logger,
	host string,
	port uint16,
	password string,
) (*redis.Client, error) {
	logger.Info(
		fmt.Sprintf(
			"redis client initialization started: dependency=redis host=%s port=%d",
			host,
			port,
		),
	)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error(
			fmt.Sprintf(
				"redis client initialization failed: dependency=redis operation=ping error=%s",
				err.Error(),
			),
		)

		return nil, NewAdapterError("ping redis", err)
	}

	return client, nil
}
