package models

import (
	"testing"
	"time"
)

func TestUserModel(t *testing.T) {
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected user name to be 'John Doe', got %s", user.Name)
	}

	if user.Email != "john@example.com" {
		t.Errorf("Expected user email to be 'john@example.com', got %s", user.Email)
	}
}

func TestPortfolioModel(t *testing.T) {
	portfolio := Portfolio{
		ID:     1,
		UserID: 1,
	}

	if portfolio.UserID != 1 {
		t.Errorf("Expected portfolio UserID to be 1, got %d", portfolio.UserID)
	}
}

func TestStockModel(t *testing.T) {
	stock := Stock{
		ID:             1,
		PortfolioID:    1,
		StockSymbol:    "AAPL",
		ThresholdPrice: 150.00,
	}

	if stock.StockSymbol != "AAPL" {
		t.Errorf("Expected stock symbol to be 'AAPL', got %s", stock.StockSymbol)
	}

	if stock.ThresholdPrice != 150.00 {
		t.Errorf("Expected threshold price to be 150.00, got %f", stock.ThresholdPrice)
	}
}

func TestAlertModel(t *testing.T) {
	now := time.Now()
	alert := Alert{
		ID:          1,
		UserID:      1,
		StockSymbol: "AAPL",
		Price:       155.50,
		Timestamp:   now,
	}

	if alert.StockSymbol != "AAPL" {
		t.Errorf("Expected alert stock symbol to be 'AAPL', got %s", alert.StockSymbol)
	}

	if alert.Price != 155.50 {
		t.Errorf("Expected alert price to be 155.50, got %f", alert.Price)
	}

	if alert.Timestamp != now {
		t.Errorf("Expected alert timestamp to match, got %v", alert.Timestamp)
	}
}

func TestStockPriceModel(t *testing.T) {
	now := time.Now()
	stockPrice := StockPrice{
		ID:        1,
		Symbol:    "TSLA",
		Price:     250.75,
		Timestamp: now,
	}

	if stockPrice.Symbol != "TSLA" {
		t.Errorf("Expected stock price symbol to be 'TSLA', got %s", stockPrice.Symbol)
	}

	if stockPrice.Price != 250.75 {
		t.Errorf("Expected stock price to be 250.75, got %f", stockPrice.Price)
	}
}

func TestStockAnalyticsModel(t *testing.T) {
	now := time.Now()
	analytics := StockAnalytics{
		ID:          1,
		Symbol:      "GOOGL",
		Avg5:        2500.00,
		Avg20:       2450.00,
		Signal:      "BULLISH",
		GeneratedAt: now,
	}

	if analytics.Symbol != "GOOGL" {
		t.Errorf("Expected analytics symbol to be 'GOOGL', got %s", analytics.Symbol)
	}

	if analytics.Signal != "BULLISH" {
		t.Errorf("Expected analytics signal to be 'BULLISH', got %s", analytics.Signal)
	}

	if analytics.Avg5 <= analytics.Avg20 {
		t.Errorf("Expected Avg5 (%f) to be greater than Avg20 (%f) for BULLISH signal", analytics.Avg5, analytics.Avg20)
	}
}

// Test helper function to validate stock symbols
func TestValidateStockSymbol(t *testing.T) {
	tests := []struct {
		symbol string
		valid  bool
	}{
		{"AAPL", true},
		{"TSLA", true},
		{"GOOGL", true},
		{"", false},
		{"TOOLONGNAME", false}, // More than 10 chars would be invalid in DB
	}

	for _, test := range tests {
		isValid := len(test.symbol) > 0 && len(test.symbol) <= 10
		if isValid != test.valid {
			t.Errorf("Expected symbol '%s' to be valid: %t, got %t", test.symbol, test.valid, isValid)
		}
	}
}
