package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/CooklyDev/AuthService/internal"
)

func main() {
	postgresConfig := internal.NewPostgresConfig()
	redisConfig := internal.NewRedisConfig()
	appConfig := internal.NewAppConfig()
	logger := internal.NewConsoleLogger()

	publisher, err := internal.NewPublisher(
		context.Background(),
		logger,
		postgresConfig,
		redisConfig,
		appConfig,
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := publisher.Run(ctx); err != nil {
		panic(err)
	}
}
