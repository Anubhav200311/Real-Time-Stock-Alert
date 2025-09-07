package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"stock-alerts/db"
	"stock-alerts/models"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

type StockEvent struct {
	Symbol string    `json:"symbol"`
	Price  float64   `json:"price"`
	Time   time.Time `json:"time"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Connect DB
	db.ConnectDatabase()

	log.Println("🔔 Alert Consumer starting...")

	// Get Kafka broker from environment variable
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "127.0.0.1:9093" // fallback for local development
	}

	// Start consuming messages for alerts
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "stock_prices",
		GroupID:  "stock-alerts-consumer",
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

		processAlertEvent(event)
	}
}

func processAlertEvent(e StockEvent) {
	var stocks []models.Stock
	db.DB.Where("stock_symbol = ?", e.Symbol).Find(&stocks)

	for _, stock := range stocks {
		if e.Price >= stock.ThresholdPrice {
			alert := models.Alert{
				UserID:      getUserIDFromPortfolio(stock.PortfolioID),
				StockSymbol: e.Symbol,
				Price:       e.Price,
				Timestamp:   e.Time,
			}
			db.DB.Create(&alert)
			log.Printf("🚨 Alert created for %s at %.2f\n", e.Symbol, e.Price)
		}
	}
}

func getUserIDFromPortfolio(portfolioID uint) uint {
	var portfolio models.Portfolio
	db.DB.First(&portfolio, portfolioID)
	return portfolio.UserID
}
