package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Task struct {
	ID          string    `json:id`
	Title       string    `json:title`
	Description string    `json:description`
	DueDate     time.Time `json:due_date`
	Status      string    `json:status`
}

// Mock data for tasks
var tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func main() {
	router := gin.Default()
	router.GET("/ping", pinger)
	router.GET("/tasks", getAllTasks)
	router.GET("/tasks/:id", getTaskByID)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", removeTask)
	router.POST("/tasks", createTask)
	router.Run("localhost:8080")
}

func pinger(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "pong"})
}

func getAllTasks(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{"tasks": tasks})
}

func getTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, task := range tasks {
		if id == task.ID {
			ctx.IndentedJSON(http.StatusOK, gin.H{"task": task})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func updateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedTask Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if id == task.ID {
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "task updated"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func removeTask(ctx *gin.Context) {
	id := ctx.Param("id")

	for i, task := range tasks {
		if id == task.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "task removed"})
			return
		}
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func createTask(ctx *gin.Context) {
	var newTask Task

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "task created"})
}
