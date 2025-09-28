package main

import (
	"log"
	"unit-management-be/internal/db"
	"unit-management-be/internal/routes"

	"github.com/joho/godotenv"
)

// @title Unit Management API
// @version 1.0
// @description API to manage unit including Creation, Edit, View and Delete
// @termsOfService http://swagger.io/terms/

// @contact.name Sutanto Adi Nugroho
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:5000
// @BasePath /api
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, falling back to system environment variables")
	}

	// try to connect database
	db.ConnectDatabase()

	// run migrations
	db.RunMigrations()

	routes.Run()
}
