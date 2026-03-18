package main

import (
	"context"

	"github.com/CooklyDev/AuthService/internal"
	"github.com/CooklyDev/AuthService/internal/adapters"
	"github.com/CooklyDev/AuthService/internal/presentation"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/CooklyDev/AuthService/docs"
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

// @title Cookly Auth Service API
// @version 1.0
// @description REST API for user registration and login in Cookly Auth Service.
// @BasePath /api/v1
// @schemes http
func main() {
	container := initContainer()
	defer container.Close()

	router := gin.Default()
	v1 := router.Group("/api/v1")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.Use(presentation.GetContainer(container))

	v1.POST("/register", presentation.Register)
	v1.POST("/login", presentation.Login)
	v1.POST("/logout", presentation.Logout)

	if err := router.Run(); err != nil {
		panic(err)
	}
}
