package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100"`
	Email     string    `gorm:"unique"`
	Portfolio Portfolio `gorm:"foreignKey:UserID"`
}

type Portfolio struct {
	ID     uint    `gorm:"primaryKey"`
	UserID uint    `gorm:"unique"` // 1 user = 1 portfolio
	Stocks []Stock `gorm:"foreignKey:PortfolioID"`
}

type Stock struct {
	ID             uint `gorm:"primaryKey"`
	PortfolioID    uint
	StockSymbol    string `gorm:"size:10"`
	ThresholdPrice float64
}

type Alert struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	StockSymbol string `gorm:"size:10"`
	Price       float64
	Timestamp   time.Time
}

// ✅ New model for persistence consumer
type StockPrice struct {
	ID        uint   `gorm:"primaryKey"`
	Symbol    string `gorm:"size:10;index"`
	Price     float64
	Timestamp time.Time
}

type StockAnalytics struct {
	ID          uint   `gorm:"primaryKey"`
	Symbol      string `gorm:"size:10;index"`
	Avg5        float64
	Avg20       float64
	Signal      string // "BULLISH", "BEARISH", "NEUTRAL"
	GeneratedAt time.Time
}
