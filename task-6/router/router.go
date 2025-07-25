package router

import (
	"task-6/controllers"
	"task-6/middleware"

	"github.com/gin-gonic/gin"
)

func RunRouter() *gin.Engine {
	router := gin.Default()
	// task
	router.GET("/tasks", controllers.GetAllTasks)
	router.GET("/tasks/:id", controllers.GetTaskByID)
	router.PUT("/tasks/:id", controllers.UpdateTask)
	router.DELETE("/tasks/:id", controllers.DeleteTask)
	router.POST("/tasks", controllers.CreateTask)
	// user
	router.GET("/users", middleware.AuthMiddleware("ADMIN"), controllers.GetAllUsers)
	router.POST("/login", controllers.LoginUser)
	router.POST("/register", controllers.CreateUser)
	router.GET("/user/:id", controllers.GetUserByID)
	router.PUT("/user/:id", middleware.AuthMiddleware("USER", "ADMIN"), controllers.UpdateUser)
	router.DELETE("/user/:id", middleware.AuthMiddleware("USER", "ADMIN"), controllers.DeleteUser)
	return router
}
