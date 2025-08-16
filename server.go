package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/h4rdc0m/aurora-api/api/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load environment variables: " + err.Error())
	}

	// Initialize the application

	err = bootstrap.RootApp.Execute()
	if err != nil {
		panic("Failed to start the application: " + err.Error())
	}
}
