package db

import (
	"fmt"
	"log"
	"os"

	"stock-alerts/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Get database configuration from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost" // fallback for local development
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "stock_alerts"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate tables (include StockPrice now)
	database.AutoMigrate(
		&models.User{},
		&models.Portfolio{},
		&models.Stock{},
		&models.Alert{},
		&models.StockPrice{}, // ✅ added,
		&models.StockAnalytics{},
	)

	DB = database
	fmt.Println("✅ Database connected & migrated")
}
