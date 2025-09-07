package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"stock-alerts/db"
	"stock-alerts/models"

	"github.com/segmentio/kafka-go"
)

type StockEvent struct {
	Symbol string    `json:"symbol"`
	Price  float64   `json:"price"`
	Time   time.Time `json:"time"`
}

func StartConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"127.0.0.1:9093"},
		Topic:    "stock_prices",
		GroupID:  "stock-alerts-consumer",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	go func() {
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

			processEvent(event)
		}
	}()
}

func processEvent(e StockEvent) {
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
			log.Printf("🚨 Alert for %s at %.2f\n", e.Symbol, e.Price)
		}
	}
}
