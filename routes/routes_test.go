package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"stock-alerts/models"
	"testing"

	"github.com/gin-gonic/gin"
)

// Setup test router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	RegisterRoutes(router)
	return router
}

// Test that routes are properly registered
func TestRouteRegistration(t *testing.T) {
	router := setupTestRouter()

	// Test that all expected routes are registered
	routes := router.Routes()

	expectedRoutes := []string{
		"POST /users",
		"GET /users",
		"POST /users/:id/portfolio",
		"GET /users/:id/portfolio",
		"POST /portfolio/:id/stocks",
		"GET /portfolio/:id/stocks",
		"GET /users/:id/alerts",
	}

	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeKey := route.Method + " " + route.Path
		routeMap[routeKey] = true
	}

	for _, expected := range expectedRoutes {
		if !routeMap[expected] {
			t.Errorf("Expected route %s to be registered", expected)
		}
	}
}

// Test basic request handling (without database)
func TestCreateUserWithInvalidData(t *testing.T) {
	router := setupTestRouter()

	// Send invalid JSON
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(`{"invalid": json}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Test JSON structure validation
func TestUserModelJSONSerialization(t *testing.T) {
	user := models.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Failed to marshal user: %v", err)
	}

	var unmarshaledUser models.User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Errorf("Failed to unmarshal user: %v", err)
	}

	if unmarshaledUser.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, unmarshaledUser.Name)
	}

	if unmarshaledUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, unmarshaledUser.Email)
	}
}

// Test portfolio model serialization
func TestPortfolioModelJSONSerialization(t *testing.T) {
	portfolio := models.Portfolio{
		ID:     1,
		UserID: 1,
	}

	jsonData, err := json.Marshal(portfolio)
	if err != nil {
		t.Errorf("Failed to marshal portfolio: %v", err)
	}

	var unmarshaledPortfolio models.Portfolio
	err = json.Unmarshal(jsonData, &unmarshaledPortfolio)
	if err != nil {
		t.Errorf("Failed to unmarshal portfolio: %v", err)
	}

	if unmarshaledPortfolio.UserID != portfolio.UserID {
		t.Errorf("Expected UserID %d, got %d", portfolio.UserID, unmarshaledPortfolio.UserID)
	}
}

// Test stock model serialization
func TestStockModelJSONSerialization(t *testing.T) {
	stock := models.Stock{
		ID:             1,
		PortfolioID:    1,
		StockSymbol:    "AAPL",
		ThresholdPrice: 150.00,
	}

	jsonData, err := json.Marshal(stock)
	if err != nil {
		t.Errorf("Failed to marshal stock: %v", err)
	}

	var unmarshaledStock models.Stock
	err = json.Unmarshal(jsonData, &unmarshaledStock)
	if err != nil {
		t.Errorf("Failed to unmarshal stock: %v", err)
	}

	if unmarshaledStock.StockSymbol != stock.StockSymbol {
		t.Errorf("Expected symbol %s, got %s", stock.StockSymbol, unmarshaledStock.StockSymbol)
	}

	if unmarshaledStock.ThresholdPrice != stock.ThresholdPrice {
		t.Errorf("Expected threshold %f, got %f", stock.ThresholdPrice, unmarshaledStock.ThresholdPrice)
	}
}

// Test HTTP request structure validation
func TestHTTPRequestValidation(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		method      string
		path        string
		contentType string
		body        string
		expectCode  int
		description string
	}{
		{"POST", "/users", "application/json", `{}`, 500, "Empty JSON should cause server error (no DB)"},
		{"POST", "/users", "application/json", `{"invalid"}`, 400, "Invalid JSON should return 400"},
		{"GET", "/users", "", "", 500, "GET request should cause server error (no DB)"},
		{"POST", "/nonexistent", "application/json", `{}`, 404, "Non-existent route should return 404"},
	}

	for _, test := range tests {
		var req *http.Request
		if test.body != "" {
			req, _ = http.NewRequest(test.method, test.path, bytes.NewBufferString(test.body))
		} else {
			req, _ = http.NewRequest(test.method, test.path, nil)
		}

		if test.contentType != "" {
			req.Header.Set("Content-Type", test.contentType)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != test.expectCode {
			t.Errorf("%s: Expected status %d, got %d", test.description, test.expectCode, w.Code)
		}
	}
}

// Test parameter extraction
func TestParameterExtraction(t *testing.T) {
	router := setupTestRouter()

	// Test with numeric parameter
	req, _ := http.NewRequest("GET", "/users/123/portfolio", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should handle parameter correctly (though DB will fail)
	if w.Code == http.StatusBadRequest {
		t.Error("Parameter parsing should not return 400")
	}
}

// Test middleware and recovery
func TestGinRecoveryMiddleware(t *testing.T) {
	router := setupTestRouter()

	// This should trigger the recovery middleware due to nil DB
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 500 due to panic, but not crash
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected recovery middleware to return 500, got %d", w.Code)
	}
}
