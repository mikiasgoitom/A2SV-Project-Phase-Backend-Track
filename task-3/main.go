package main

import (
	"libraryManagement/controllers"
	"libraryManagement/services"
)

func main() {
	var service services.Library
	lib := service.NewLibrary()
	controllers.Welcome()
	controllers.Menu(lib)
}
