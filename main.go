package main

import (
	"car-service/database"
	"car-service/routes"
	"car-service/seeders"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// main function
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	//created a default gin router
	r := gin.Default()

	routes.Router(r)
	//connect the database
	database.ConnectDB()

	//seeders to create admin
	seeders.SeedersToAdmin()

	// server running
	r.Run(":8080")
	log.Fatal("Port running on http//localhost:8080")
}
