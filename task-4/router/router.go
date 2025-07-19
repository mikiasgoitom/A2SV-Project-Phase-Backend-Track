package router

import (
	"task-management-api/controllers"

	"github.com/gin-gonic/gin"
)

func RunRouter() {
	router := gin.Default()

	router.GET("/tasks", controllers.GetAllTasks)
	router.GET("/tasks/:id", controllers.GetTaskByID)
	router.PUT("/tasks/:id", controllers.UpdateTask)
	router.DELETE("/tasks/:id", controllers.DeleteTask)
	router.POST("/tasks", controllers.CreateTask)

	router.Run("localhost:8080")
}
