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

var priceHistory = make(map[string][]float64) // symbol → recent prices

func StartAnalyticsConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"127.0.0.1:9093"},
		Topic:   "stock_prices",
		GroupID: "analytics-consumer-group",
	})

	go func() {
		log.Println("▶️ Analytics consumer started...")
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
	}()
}

func updateAnalytics(e StockEvent) {
	// Append to history
	history := priceHistory[e.Symbol]
	history = append(history, e.Price)
	if len(history) > 20 {
		history = history[len(history)-20:] // keep last 20
	}
	priceHistory[e.Symbol] = history

	// Calculate averages
	avg5 := avg(history, 5)
	avg20 := avg(history, 20)

	// Determine signal
	signal := "NEUTRAL"
	if avg5 > avg20 {
		signal = "BULLISH"
	} else if avg5 < avg20 {
		signal = "BEARISH"
	}

	// Save to DB
	analytics := models.StockAnalytics{
		Symbol:      e.Symbol,
		Avg5:        avg5,
		Avg20:       avg20,
		Signal:      signal,
		GeneratedAt: time.Now(),
	}
	if res := db.DB.Create(&analytics); res.Error != nil {
		log.Printf("❌ Failed to save analytics: %v\n", res.Error)
	} else {
		log.Printf("📊 Analytics for %s: avg5=%.2f, avg20=%.2f, signal=%s\n",
			e.Symbol, avg5, avg20, signal)
	}
}

func avg(arr []float64, n int) float64 {
	if len(arr) < n {
		return 0
	}
	sum := 0.0
	for _, v := range arr[len(arr)-n:] {
		sum += v
	}
	return sum / float64(n)
}
