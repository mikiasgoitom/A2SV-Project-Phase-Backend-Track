package route

import (
	controller "clean-architecture/delivary/controllers"
	"clean-architecture/delivary/middleware"
	"clean-architecture/usecase/contract"

	"github.com/gin-gonic/gin"
)

type Router struct {
	UserHandler *controller.UserHandler
	TaskHandler *controller.TaskHandler
	AuthHandler *controller.AuthHandler
	jwtService  contract.ITokenService
}

func NewRouter(userHandler *controller.UserHandler, taskHandler *controller.TaskHandler, authHandler *controller.AuthHandler, jwtService contract.ITokenService) *Router {
	return &Router{
		UserHandler: userHandler,
		TaskHandler: taskHandler,
		AuthHandler: authHandler,
		jwtService:  jwtService,
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	// User routes
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", r.UserHandler.HandleRegister)
		userGroup.PUT("/:id", r.UserHandler.UpdateUser)
		userGroup.DELETE("/:id", r.UserHandler.DeleteUser)
	}

	// Task routes
	taskGroup := router.Group("/tasks")
	taskGroup.Use(middleware.AuthMiddleware(r.jwtService))
	{
		taskGroup.GET("/overdue", r.TaskHandler.OverDueTasks)
		taskGroup.POST("/", r.TaskHandler.CreateNewTask)
		taskGroup.PUT("/:id", r.TaskHandler.UpdateTask)
		taskGroup.DELETE("/:id", r.TaskHandler.DeleteTask)
	}

	// auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", r.AuthHandler.HandleLogin)
		authGroup.POST("/logout", middleware.AuthMiddleware(r.jwtService), r.AuthHandler.HandleLogout)
	}

}
