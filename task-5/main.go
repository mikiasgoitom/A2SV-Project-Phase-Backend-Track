package main

import (
	"task-management-api/data"
	"task-management-api/router"
)

func main() {
	if err := data.ConnectToMongoDB(); err != nil {
		panic("failed to connect to MongoDB: " + err.Error())
	}
	defer data.CloseMongoDBConnection()
	// Start the server
	server := router.RunRouter()
	if err := server.Run("localhost:8080"); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
