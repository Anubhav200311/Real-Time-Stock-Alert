package db

import (
	"fmt"
	"log"

	"stock-alerts/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=stock_alerts port=5432 sslmode=disable"
	// 👉 Replace with your actual DB creds

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
		&models.StockPrice{}, // ✅ added
	)

	DB = database
	fmt.Println("✅ Database connected & migrated")
}
