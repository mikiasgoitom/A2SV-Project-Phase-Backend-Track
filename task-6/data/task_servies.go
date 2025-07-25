package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"task-6/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// this function is commented out to prevent seeding the database every time the application starts
/*
func SeedDB() {
	var tasks = []models.Task{
		{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
		{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
		{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
	}
	for i := range tasks {
		_, err := taskCollection.InsertOne(context.TODO(), tasks[i])
		if err != nil {
			log.Println("error seeding task:", err)
		}
	}
} */

// this function gets all the tasks from the database
func GetAllTasks() []models.Task {
	var tasks []models.Task
	cursor, err := taskCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		log.Println("Error fetching tasks:", err)
		return tasks
	}

	for cursor.Next(context.TODO()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			log.Println("Error decoding task:", err)
			continue
		}
		tasks = append(tasks, task)
	}
	fmt.Println("Fetched tasks:", tasks)
	return tasks
}

// this function gets a task by its ID from the database
func GetTaskByID(id string) *models.Task {
	filter := bson.D{{Key: "id", Value: id}}
	var task models.Task
	err := taskCollection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		log.Println("Error finding task by ID:", err)
		return nil
	}
	return &task
}

func UpdateTask(id string, updatedTask models.Task) *models.Task {
	updateFields := bson.D{}

	if updatedTask.Title != "" {
		updateFields = append(updateFields, bson.E{Key: "title", Value: updatedTask.Title})
	}
	if updatedTask.Description != "" {
		updateFields = append(updateFields, bson.E{Key: "description", Value: updatedTask.Description})
	}
	if !updatedTask.DueDate.IsZero() {
		updateFields = append(updateFields, bson.E{Key: "due_date", Value: updatedTask.DueDate})
	}
	if updatedTask.Status != "" {
		updateFields = append(updateFields, bson.E{Key: "status", Value: updatedTask.Status})
	}

	if len(updateFields) == 0 {
		log.Println("No fields to update")
		return nil
	}

	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: updateFields},
	}

	_, err := taskCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Println("Error updating task:", err)
		return nil
	}

	var task models.Task
	err = taskCollection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		log.Println("Error finding updated task:", err)
		return nil
	}
	return &task
}

func CreateTask(newTask models.Task) *models.Task {
	newTask.ID = uuid.New().String()

	_, err := taskCollection.InsertOne(context.TODO(), newTask)

	if err != nil {
		log.Fatal("Error creating a task:", err)
	}
	return &newTask
}

func DeleteTask(id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	result, err := taskCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println("Error deleting task:", err)
	}

	if result.DeletedCount == 0 {
		log.Println("No task found with the given ID")
		return errors.New("no task found with the given ID")
	}

	log.Println("Task deleted successfully")
	return nil
}
