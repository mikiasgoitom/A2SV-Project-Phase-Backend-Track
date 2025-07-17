package router

import (
	"task-management-api/data"

	"github.com/gin-gonic/gin"
)

func RunRouter() {
	router := gin.Default()

	router.GET("/tasks", data.GetAllTasks)
	router.GET("/tasks/:id", data.GetTaskByID)
	router.PUT("/tasks/:id", data.UpdateTask)
	router.DELETE("/tasks/:id", data.DeleteTask)
	router.POST("/tasks", data.CreateTask)

	router.Run("localhost:8080")
}
