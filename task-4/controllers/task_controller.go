package controllers

import (
	"net/http"
	"task-management-api/data"
	"task-management-api/models"

	"github.com/gin-gonic/gin"
)

func GetAllTasks(c *gin.Context) {
	tasks := data.GetAllTasks()
	if tasks != nil {
		c.IndentedJSON(http.StatusOK, tasks)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not tasks found"})
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task := data.GetTaskByID(id)
	if task != nil {
		c.IndentedJSON(http.StatusOK, task)
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	task := data.UpdateTask(id, updatedTask)
	if task != nil {
		c.IndentedJSON(http.StatusOK, task)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := data.DeleteTask(id); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

func CreateTask(c *gin.Context) {
	var newTask models.Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data.CreateTask(newTask)

	c.IndentedJSON(http.StatusAccepted, gin.H{"message": "task created successfully"})
}
