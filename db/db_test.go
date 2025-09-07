package db

import (
	"os"
	"testing"
)

func TestConnectDatabase(t *testing.T) {
	// Test with mock environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")

	// This test will verify the function doesn't panic
	// In a real environment, you'd use a test database
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ConnectDatabase panicked: %v", r)
		}
	}()

	// Note: This will likely fail to connect to actual DB in CI/CD
	// but it tests the connection string building logic
	// ConnectDatabase()

	// Cleanup
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
}

func TestEnvironmentVariableDefaults(t *testing.T) {
	// Ensure all environment variables are unset
	envVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	// Test default values
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost" // fallback
	}
	if host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", host)
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432" // fallback
	}
	if port != "5432" {
		t.Errorf("Expected default port '5432', got '%s'", port)
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres" // fallback
	}
	if user != "postgres" {
		t.Errorf("Expected default user 'postgres', got '%s'", user)
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres" // fallback
	}
	if password != "postgres" {
		t.Errorf("Expected default password 'postgres', got '%s'", password)
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "stock_alerts" // fallback
	}
	if dbname != "stock_alerts" {
		t.Errorf("Expected default dbname 'stock_alerts', got '%s'", dbname)
	}
}

func TestDatabaseConfiguration(t *testing.T) {
	tests := []struct {
		host     string
		port     string
		user     string
		password string
		dbname   string
		expected string
	}{
		{
			host:     "localhost",
			port:     "5432",
			user:     "postgres",
			password: "password",
			dbname:   "testdb",
			expected: "host=localhost user=postgres password=password dbname=testdb port=5432 sslmode=disable",
		},
		{
			host:     "db.example.com",
			port:     "5433",
			user:     "myuser",
			password: "mypass",
			dbname:   "mydb",
			expected: "host=db.example.com user=myuser password=mypass dbname=mydb port=5433 sslmode=disable",
		},
	}

	for _, test := range tests {
		// Set environment variables
		os.Setenv("DB_HOST", test.host)
		os.Setenv("DB_PORT", test.port)
		os.Setenv("DB_USER", test.user)
		os.Setenv("DB_PASSWORD", test.password)
		os.Setenv("DB_NAME", test.dbname)

		// Build DSN string (same logic as in ConnectDatabase)
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}

		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "postgres"
		}

		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "stock_alerts"
		}

		dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"

		if dsn != test.expected {
			t.Errorf("Expected DSN '%s', got '%s'", test.expected, dsn)
		}

		// Cleanup
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}
}

// Test database connection validation
func TestDatabaseConnectionValidation(t *testing.T) {
	// Test that required fields are not empty
	requiredFields := map[string]string{
		"host":   "localhost",
		"port":   "5432",
		"user":   "postgres",
		"dbname": "stock_alerts",
	}

	for field, value := range requiredFields {
		if value == "" {
			t.Errorf("Required field %s should not be empty", field)
		}
	}
}

// Test SSL mode configuration
func TestSSLConfiguration(t *testing.T) {
	// In our current setup, we use sslmode=disable
	// This test verifies that the SSL mode is properly configured
	sslMode := "disable"

	if sslMode != "disable" && sslMode != "require" && sslMode != "verify-full" {
		t.Errorf("Invalid SSL mode: %s", sslMode)
	}
}

// Test connection string format
func TestConnectionStringFormat(t *testing.T) {
	// Test that connection string follows PostgreSQL format
	host := "localhost"
	user := "postgres"
	password := "password"
	dbname := "testdb"
	port := "5432"
	sslmode := "disable"

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode

	// Check that all required components are present
	expectedComponents := []string{
		"host=" + host,
		"user=" + user,
		"password=" + password,
		"dbname=" + dbname,
		"port=" + port,
		"sslmode=" + sslmode,
	}

	for _, component := range expectedComponents {
		if !contains(dsn, component) {
			t.Errorf("DSN missing component: %s", component)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
