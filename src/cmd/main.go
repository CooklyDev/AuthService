package main

import (
	"context"

	"github.com/CooklyDev/AuthService/internal"
	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/presentation"

	"github.com/gin-gonic/gin"
)

func initContainer() *internal.Container {
	// Initialize the container and its dependencies
	logger := internal.NewNoopLogger()
	hasher := adapters.NewStubHasher()

	container, err := internal.NewContainer(
		logger,
		hasher,
		context.Background(),
	)
	if err != nil {
		panic(err)
	}
	return container
}

func main() {
	container := initContainer()
	defer container.Close()

	router := gin.Default()
	router.Use(presentation.GetContainer(container))
	router.GET("/", presentation.Hello)

	if err := router.Run(); err != nil {
		panic(err)
	}
}
