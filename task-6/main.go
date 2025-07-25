package main

import (
	"os"
	"task-6/data"
	"task-6/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("filed to load .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085" // Default port if not set
	}

	// Initialize the database connection
	if err := data.ConnectToMongoDB(); err != nil {
		panic("failed to connect to MongoDB: " + err.Error())
	}
	defer data.CloseMongoDBConnection()
	// Start the server

	// create index for username
	data.CreateUserIndex()
	server := router.RunRouter()
	if err := server.Run("localhost:" + port); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
