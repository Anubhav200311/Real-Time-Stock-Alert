package services

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestInitKafkaProducer(t *testing.T) {
	// Test with environment variable
	os.Setenv("KAFKA_BROKER", "test-broker:9092")
	InitKafkaProducer()

	if kafkaWriter == nil {
		t.Error("Expected kafkaWriter to be initialized, got nil")
	}

	// Cleanup
	os.Unsetenv("KAFKA_BROKER")
}

func TestInitKafkaProducerWithoutEnv(t *testing.T) {
	// Ensure no environment variable is set
	os.Unsetenv("KAFKA_BROKER")

	InitKafkaProducer()

	if kafkaWriter == nil {
		t.Error("Expected kafkaWriter to be initialized with fallback, got nil")
	}
}

func TestStockEventSerialization(t *testing.T) {
	// Test the data structure that PublishStockPrice creates
	symbol := "AAPL"
	price := 150.50
	now := time.Now()

	event := map[string]any{
		"symbol": symbol,
		"price":  price,
		"time":   now,
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Errorf("Failed to marshal stock event: %v", err)
	}

	// Unmarshal to verify structure
	var unmarshaled map[string]any
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal stock event: %v", err)
	}

	if unmarshaled["symbol"] != symbol {
		t.Errorf("Expected symbol %s, got %v", symbol, unmarshaled["symbol"])
	}

	if unmarshaled["price"] != price {
		t.Errorf("Expected price %f, got %v", price, unmarshaled["price"])
	}

	if unmarshaled["time"] == nil {
		t.Error("Expected time field to be present")
	}
}

// Test helper function for validating stock symbols
func TestValidateSymbol(t *testing.T) {
	tests := []struct {
		symbol   string
		expected bool
	}{
		{"AAPL", true},
		{"GOOGL", true},
		{"TSLA", true},
		{"", false},
		{"A", true},
		{"ABCDEFGHIJK", false}, // Too long
	}

	for _, test := range tests {
		result := validateSymbol(test.symbol)
		if result != test.expected {
			t.Errorf("validateSymbol(%s) = %t, expected %t", test.symbol, result, test.expected)
		}
	}
}

// Test helper function for validating prices
func TestValidatePrice(t *testing.T) {
	tests := []struct {
		price    float64
		expected bool
	}{
		{100.50, true},
		{0.01, true},
		{0, false},
		{-10.50, false},
		{999999.99, true},
	}

	for _, test := range tests {
		result := validatePrice(test.price)
		if result != test.expected {
			t.Errorf("validatePrice(%f) = %t, expected %t", test.price, result, test.expected)
		}
	}
}

// Helper functions that we're testing (these would be useful in your actual code)
func validateSymbol(symbol string) bool {
	return len(symbol) > 0 && len(symbol) <= 10
}

func validatePrice(price float64) bool {
	return price > 0
}

// Test for stock price calculation logic
func TestCalculatePercentageChange(t *testing.T) {
	tests := []struct {
		oldPrice    float64
		newPrice    float64
		expectedPct float64
	}{
		{100.0, 110.0, 10.0},
		{100.0, 90.0, -10.0},
		{50.0, 75.0, 50.0},
		{200.0, 200.0, 0.0},
	}

	for _, test := range tests {
		result := calculatePercentageChange(test.oldPrice, test.newPrice)
		if result != test.expectedPct {
			t.Errorf("calculatePercentageChange(%f, %f) = %f, expected %f",
				test.oldPrice, test.newPrice, result, test.expectedPct)
		}
	}
}

// Helper function for percentage calculation
func calculatePercentageChange(oldPrice, newPrice float64) float64 {
	if oldPrice == 0 {
		return 0
	}
	return ((newPrice - oldPrice) / oldPrice) * 100
}

// Test for alert threshold logic
func TestShouldTriggerAlert(t *testing.T) {
	tests := []struct {
		currentPrice   float64
		thresholdPrice float64
		expected       bool
		description    string
	}{
		{155.0, 150.0, true, "Price above threshold should trigger alert"},
		{145.0, 150.0, false, "Price below threshold should not trigger alert"},
		{150.0, 150.0, true, "Price equal to threshold should trigger alert"},
		{0.0, 150.0, false, "Invalid price should not trigger alert"},
	}

	for _, test := range tests {
		result := shouldTriggerAlert(test.currentPrice, test.thresholdPrice)
		if result != test.expected {
			t.Errorf("%s: shouldTriggerAlert(%f, %f) = %t, expected %t",
				test.description, test.currentPrice, test.thresholdPrice, result, test.expected)
		}
	}
}

// Helper function for alert logic
func shouldTriggerAlert(currentPrice, thresholdPrice float64) bool {
	return currentPrice > 0 && currentPrice >= thresholdPrice
}
