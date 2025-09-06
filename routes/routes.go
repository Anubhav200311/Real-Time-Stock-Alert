package routes

import (
	"fmt"
	"net/http"
	"stock-alerts/db"
	"stock-alerts/models"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes adds endpoints
func RegisterRoutes(r *gin.Engine) {
	// User routes
	r.POST("/users", createUser)
	r.GET("/users", listUsers)

	// Portfolio routes
	r.POST("/users/:id/portfolio", createPortfolio)
	r.GET("/users/:id/portfolio", getPortfolio)

	// Stock routes
	r.POST("/portfolio/:id/stocks", addStock)
	r.GET("/portfolio/:id/stocks", listStocks)

	// Alerts
	r.GET("/users/:id/alerts", getAlerts)
}

// ----------------- User Handlers -----------------
func createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&user)
	c.JSON(http.StatusOK, user)
}

func listUsers(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

// ----------------- Portfolio Handlers -----------------
func createPortfolio(c *gin.Context) {
	userID := c.Param("id")
	var portfolio models.Portfolio
	portfolio.UserID = uint(parseID(userID))
	db.DB.Create(&portfolio)
	c.JSON(http.StatusOK, portfolio)
}

func getPortfolio(c *gin.Context) {
	userID := c.Param("id")
	var portfolio models.Portfolio
	err := db.DB.Preload("Stocks").Where("user_id = ?", userID).First(&portfolio).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "portfolio not found"})
		return
	}
	c.JSON(http.StatusOK, portfolio)
}

// ----------------- Stock Handlers -----------------
func addStock(c *gin.Context) {
	portfolioID := c.Param("id")
	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stock.PortfolioID = uint(parseID(portfolioID))
	db.DB.Create(&stock)
	c.JSON(http.StatusOK, stock)
}

func listStocks(c *gin.Context) {
	portfolioID := c.Param("id")
	var stocks []models.Stock
	db.DB.Where("portfolio_id = ?", portfolioID).Find(&stocks)
	c.JSON(http.StatusOK, stocks)
}

// ----------------- Alerts Handler -----------------
func getAlerts(c *gin.Context) {
	userID := c.Param("id")
	var alerts []models.Alert
	db.DB.Where("user_id = ?", userID).Find(&alerts)
	c.JSON(http.StatusOK, alerts)
}

// ----------------- Helper -----------------
func parseID(id string) uint {
	var val uint
	fmt.Sscanf(id, "%d", &val)
	return val
}
