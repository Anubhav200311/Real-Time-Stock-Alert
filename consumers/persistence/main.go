package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"stock-alerts/db"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

// StockPriceRecord represents the stock price data to be persisted
type StockPriceRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Symbol    string `gorm:"size:10;index"`
	Price     float64
	Timestamp time.Time
}

type StockEvent struct {
	Symbol string    `json:"symbol"`
	Price  float64   `json:"price"`
	Time   time.Time `json:"time"`
}

// TableName overrides the default table name
func (StockPriceRecord) TableName() string {
	return "stock_price_records"
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Connect DB
	db.ConnectDatabase()

	// Ensure the table exists
	if err := db.DB.AutoMigrate(&StockPriceRecord{}); err != nil {
		log.Printf("⚠️  AutoMigrate stock_price_records table failed: %v\n", err)
	}

	log.Println("💾 Persistence Consumer starting...")

	// Get Kafka broker from environment variable
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "127.0.0.1:9093" // fallback for local development
	}

	// Start consuming messages for persistence
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "stock_prices",
		GroupID:  "persistence-consumer-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("❌ Kafka read error:", err)
			continue
		}

		var event StockEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Println("❌ JSON parse error:", err)
			continue
		}

		// Create stock price record
		rec := StockPriceRecord{
			Symbol:    event.Symbol,
			Price:     event.Price,
			Timestamp: event.Time,
		}

		if res := db.DB.Create(&rec); res.Error != nil {
			log.Printf("❌ Failed to store price for %s: %v\n", rec.Symbol, res.Error)
			continue
		}
		log.Printf("✅ Stored stock price: %s -> %.2f at %s\n", rec.Symbol, rec.Price, rec.Timestamp.Format(time.RFC3339))
	}
}
