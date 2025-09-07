package services

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

// InitKafkaProducer sets up the Kafka writer
func InitKafkaProducer() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "127.0.0.1:9093" // fallback for local development
	}

	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "stock_prices",
		Balancer: &kafka.LeastBytes{},
	}
}

// PublishStockPrice sends stock data to Kafka
func PublishStockPrice(symbol string, price float64) {
	event := map[string]any{
		"symbol": symbol,
		"price":  price,
		"time":   time.Now(),
	}
	data, _ := json.Marshal(event)

	err := kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{Value: data},
	)
	if err != nil {
		log.Println("❌ Kafka write failed:", err)
	} else {
		log.Printf("✅ Published to Kafka: %s %.2f\n", symbol, price)
	}
}
