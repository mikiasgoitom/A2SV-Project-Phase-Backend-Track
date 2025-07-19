package data

import (
	"errors"
	"strconv"
	"task-management-api/models"
	"time"
)

var tasks = []models.Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

var tasksCount = 3

func GetAllTasks() []models.Task {
	return tasks
}

func GetTaskByID(id string) *models.Task {
	for i, task := range tasks {
		if task.ID == id {
			return &tasks[i]
		}
	}
	return nil
}

func UpdateTask(id string, updatedTask models.Task) *models.Task {
	for i := range tasks {
		if tasks[i].ID == id {
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				tasks[i].Status = updatedTask.Status
			}
			return &tasks[i]
		}
	}
	return nil
}

func CreateTask(newTask models.Task) models.Task {
	newTask.ID = strconv.Itoa(tasksCount + 1)
	tasksCount += 1
	tasks = append(tasks, newTask)

	return newTask
}

func DeleteTask(id string) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
