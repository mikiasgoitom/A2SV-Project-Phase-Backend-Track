package data

import (
	"net/http"
	"task-management-api/models"

	"github.com/gin-gonic/gin"
)

func GetAllTasks(c *gin.Context) {
	if models.Tasks != nil {
		c.IndentedJSON(http.StatusOK, models.Tasks)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not tasks found"})
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	for _, task := range models.Tasks {
		if task.ID == id {
			c.IndentedJSON(http.StatusOK, task)
			return
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "task not found"})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
	}

	for _, task := range models.Tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				task.Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				task.Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				task.Status = updatedTask.Status
			}
			c.IndentedJSON(http.StatusOK, task)
			return
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "task not found"})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	for i, task := range models.Tasks {
		if task.ID == id {
			models.Tasks = append(models.Tasks[:i], models.Tasks[i+1:]...)
			return
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "task not found"})
}

/*
-GET /tasks: Get a list of all tasks.
-GET /tasks/:id: Get the details of a specific task.
-PUT /tasks/:id: Update a specific task. This endpoint should accept a JSON body with the new details of the task.
- DELETE /tasks/:id: Delete a specific task.
- POST /tasks: Create a new task. This endpoint should accept a JSON body with the task's title, description, due date, and status.
*/

func CreateTask(c *gin.Context) {
	var newTask models.Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.Tasks = append(models.Tasks, newTask)
	c.IndentedJSON(http.StatusAccepted, gin.H{"message": "task created successfully"})
}
