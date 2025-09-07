package services

import (
	"context"
	"encoding/json"
	"log"

	"stock-alerts/db"
	"stock-alerts/models"

	"github.com/segmentio/kafka-go"
)

func StartPersistenceConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"127.0.0.1:9093"}, // force IPv4 to avoid [::1]
		Topic:    "stock_prices",
		GroupID:  "persistence-consumer-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	go func() {
		log.Println("▶️ Persistence consumer started...")
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Println("❌ Kafka read error:", err)
				continue
			}

			var event StockEvent // already defined in consumer.go
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Println("❌ JSON parse error:", err)
				continue
			}

			rec := models.StockPrice{
				Symbol:    event.Symbol,
				Price:     event.Price,
				Timestamp: event.Time,
			}

			if res := db.DB.Create(&rec); res.Error != nil {
				log.Printf("❌ Failed to store price for %s: %v\n", rec.Symbol, res.Error)
				continue
			}
			log.Printf("✅ Stored stock price: %s -> %.2f\n", rec.Symbol, rec.Price)
		}
	}()
}
