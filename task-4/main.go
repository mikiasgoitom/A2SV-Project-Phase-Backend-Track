package main

import "task-management-api/router"

func main() {
	server := router.RunRouter()
	if err := server.Run("localhost:8080"); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
