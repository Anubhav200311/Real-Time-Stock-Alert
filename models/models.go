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
