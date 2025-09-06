package main

import (
	"stock-alerts/db"
	"stock-alerts/routes"
	"stock-alerts/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect DB
	db.ConnectDatabase()

	// Start background stock fetcher
	services.StartFetcher()

	// Setup router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r)

	// Start server
	r.Run(":8080")
}
