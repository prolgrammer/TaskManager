package main

import (
	"TaskManager/cmd/app"
	_ "TaskManager/docs"
)

// @title           Task service
// @version         0.0.1
// @description     Service for processing tasks

// @host      localhost:8080
// @BasePath  /
func main() {
	app.Run()
}
