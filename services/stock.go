package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"stock-alerts/db"
	"stock-alerts/models"
	"strconv"
	"time"
)

type GlobalQuote struct {
	Symbol string `json:"01. symbol"`
	Price  string `json:"05. price"`
}

type ApiResponse struct {
	Quote GlobalQuote `json:"Global Quote"`
}

// FetchPrice calls Alpha Vantage API
func FetchPrice(symbol string) (float64, error) {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("ALPHA_VANTAGE_API_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Debug: print the API response
	log.Printf("API Response for %s: %s", symbol, string(body))

	var result ApiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Check if price is empty
	if result.Quote.Price == "" {
		return 0, fmt.Errorf("empty price returned for symbol %s - check if symbol is valid or API limit reached", symbol)
	}

	price, err := strconv.ParseFloat(result.Quote.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price '%s': %v", result.Quote.Price, err)
	}
	return price, nil
}

// StartFetcher runs a background loop
func StartFetcher() {
	ticker := time.NewTicker(60 * time.Second) // every 1 min
	go func() {
		for {
			<-ticker.C
			checkStocks()
		}
	}()
}

func checkStocks() {
	var stocks []models.Stock
	db.DB.Find(&stocks)

	for _, stock := range stocks {
		price, err := FetchPrice(stock.StockSymbol)
		if err != nil {
			log.Println("Error fetching price:", err)
			continue
		}

		if price >= stock.ThresholdPrice {
			alert := models.Alert{
				UserID:      getUserIDFromPortfolio(stock.PortfolioID),
				StockSymbol: stock.StockSymbol,
				Price:       price,
				Timestamp:   time.Now(),
			}
			db.DB.Create(&alert)
			PublishStockPrice(stock.StockSymbol, price)
			log.Printf("🚨 Alert created for %s at price %.2f\n", stock.StockSymbol, price)
		}
	}
}

func getUserIDFromPortfolio(portfolioID uint) uint {
	var portfolio models.Portfolio
	db.DB.First(&portfolio, portfolioID)
	return portfolio.UserID
}
