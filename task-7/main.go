package main

import (
	controller "clean-architecture/delivary/controllers"
	route "clean-architecture/delivary/routers"
	"clean-architecture/infrastructure/database"
	jwtManager "clean-architecture/infrastructure/jwt"
	"clean-architecture/infrastructure/logger"
	passwordservice "clean-architecture/infrastructure/password_service"
	"clean-architecture/infrastructure/repository"
	"clean-architecture/usecase"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// jwt secret
	jwtSecret := "very_secret"

	// initialize service
	logger := logger.NewLogger()
	pwdService := passwordservice.NewPasswordService()

	// initialize the db
	ctx := context.Background()
	dbName := "task7"
	url := "mongodb://localhost:27017"
	client, db, err := database.InitMongoDB(ctx, url, dbName, logger)
	if err != nil {
		logger.Error(fmt.Errorf("failed to initialize database: %w", err).Error())
		return
	}
	defer database.CloseMongoDB(ctx, client)

	// initializing repositories
	userRepo := repository.NewUserRepository(db.Collection("users"), logger)
	taskRepo := repository.NewTaskRepository(db.Collection("tasks"), logger)
	jwtManager := jwtManager.NewJWTService(jwtSecret, userRepo)

	// initializing usecases
	userUsecase := usecase.NewUserUseCase(userRepo, pwdService)
	taskUsecase := usecase.NewTaskUseCase(taskRepo, logger)
	authUsecase := usecase.NewAuthUseCase(userRepo, pwdService, jwtManager, logger)

	// initializing controllers
	userHandler := controller.NewUserHandler(userUsecase)
	taskHandler := controller.NewTaskHandler(taskUsecase)
	authHandler := controller.NewAuthHandler(authUsecase)

	// initialize the router
	router := gin.Default()

	appRouter := route.NewRouter(userHandler, taskHandler, authHandler, jwtManager)
	appRouter.SetupRoutes(router)

	// start the server
	router.Run(":8080")
}
