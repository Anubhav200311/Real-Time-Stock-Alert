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

	// Init Kafka Producer (for publishing stock data)
	services.InitKafkaProducer()

	// Start background stock fetcher (produces to Kafka)
	services.StartFetcher()

	// Setup router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r)

	log.Println("🚀 API service starting on port 8080")
	// Start server
	r.Run(":8080")
}
