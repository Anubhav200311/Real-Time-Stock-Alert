package main

import (
	"log"
	"stock-alerts/db"
	"stock-alerts/routes"
	"stock-alerts/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Connect DB
	db.ConnectDatabase()

	// Init Kafka
	services.InitKafkaProducer()

	// Start consumers
	services.StartConsumer()            // Alert consumer
	services.StartPersistenceConsumer() // Persistence consumer ✅

	// Start background stock fetcher (produces to Kafka)
	services.StartFetcher()

	// Setup router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r)

	// Start server
	r.Run(":8080")
}
