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

var API_KEY = os.Getenv("ALPHA_VANTAGE_API_KEY")

type GlobalQuote struct {
	Symbol string `json:"01. symbol"`
	Price  string `json:"05. price"`
}

type ApiResponse struct {
	Quote GlobalQuote `json:"Global Quote"`
}

// FetchPrice calls Alpha Vantage API
func FetchPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, API_KEY)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result ApiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(result.Quote.Price, 64)
	if err != nil {
		return 0, err
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
			log.Printf("🚨 Alert created for %s at price %.2f\n", stock.StockSymbol, price)
		}
	}
}

func getUserIDFromPortfolio(portfolioID uint) uint {
	var portfolio models.Portfolio
	db.DB.First(&portfolio, portfolioID)
	return portfolio.UserID
}
