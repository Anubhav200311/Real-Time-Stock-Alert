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

// StockAnalytics represents aggregated stock analytics
type StockAnalytics struct {
	ID           uint      `gorm:"primaryKey"`
	Symbol       string    `gorm:"size:10;uniqueIndex:idx_symbol_date"`
	Date         time.Time `gorm:"uniqueIndex:idx_symbol_date"`
	MinPrice     float64
	MaxPrice     float64
	AvgPrice     float64
	TotalVolume  int64
	PriceChanges int // number of price updates in a day
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type StockEvent struct {
	Symbol string    `json:"symbol"`
	Price  float64   `json:"price"`
	Time   time.Time `json:"time"`
}

// TableName overrides the default table name
func (StockAnalytics) TableName() string {
	return "stock_analytics"
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
	if err := db.DB.AutoMigrate(&StockAnalytics{}); err != nil {
		log.Printf("⚠️  AutoMigrate stock_analytics table failed: %v\n", err)
	}

	log.Println("📊 Analytics Consumer starting...")

	// Get Kafka broker from environment variable
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "127.0.0.1:9093" // fallback for local development
	}

	// Start consuming messages for analytics
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "stock_prices",
		GroupID:  "analytics-consumer-group",
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

		updateAnalytics(event)
	}
}

func updateAnalytics(event StockEvent) {
	// Get date for grouping (without time component)
	date := time.Date(event.Time.Year(), event.Time.Month(), event.Time.Day(), 0, 0, 0, 0, event.Time.Location())

	var analytics StockAnalytics

	// Find or create analytics record for this symbol and date
	result := db.DB.Where("symbol = ? AND date = ?", event.Symbol, date).First(&analytics)

	if result.Error != nil {
		// Create new analytics record
		analytics = StockAnalytics{
			Symbol:       event.Symbol,
			Date:         date,
			MinPrice:     event.Price,
			MaxPrice:     event.Price,
			AvgPrice:     event.Price,
			TotalVolume:  1,
			PriceChanges: 1,
		}

		if err := db.DB.Create(&analytics).Error; err != nil {
			log.Printf("❌ Failed to create analytics for %s: %v\n", event.Symbol, err)
			return
		}
		log.Printf("📊 Created analytics record for %s on %s\n", event.Symbol, date.Format("2006-01-02"))
	} else {
		// Update existing analytics record
		newAvg := (analytics.AvgPrice*float64(analytics.PriceChanges) + event.Price) / float64(analytics.PriceChanges+1)

		updates := map[string]interface{}{
			"avg_price":     newAvg,
			"price_changes": analytics.PriceChanges + 1,
		}

		if event.Price < analytics.MinPrice {
			updates["min_price"] = event.Price
		}
		if event.Price > analytics.MaxPrice {
			updates["max_price"] = event.Price
		}

		if err := db.DB.Model(&analytics).Updates(updates).Error; err != nil {
			log.Printf("❌ Failed to update analytics for %s: %v\n", event.Symbol, err)
			return
		}
		log.Printf("📊 Updated analytics for %s: price=%.2f, avg=%.2f, changes=%d\n",
			event.Symbol, event.Price, newAvg, analytics.PriceChanges+1)
	}
}
