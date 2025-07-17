package models

import "time"

type Task struct {
	ID          string    `json:id`
	Title       string    `json:title`
	Description string    `json:description`
	Status      string    `json:status`
	DueDate     time.Time `json:due_date`
}

var Tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}
