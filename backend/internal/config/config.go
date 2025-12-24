package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	Port string
	Env  string // development, production

	// Database
	DatabaseURL string

	// Clerk Authentication
	ClerkSecretKey     string
	ClerkPublishableKey string
	ClerkJWTKey        string

	// CORS
	AllowedOrigins []string

	// Sync Settings
	SyncBatchSize       int
	SyncRetryAttempts   int
	SyncRetryDelay      time.Duration

	// Logging
	LogLevel  string
	LogFormat string // json, text
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		Env:                getEnv("ENVIRONMENT", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		ClerkSecretKey:     getEnv("CLERK_SECRET_KEY", ""),
		ClerkPublishableKey: getEnv("CLERK_PUBLISHABLE_KEY", ""),
		ClerkJWTKey:        getEnv("CLERK_JWT_KEY", ""),
		AllowedOrigins:     getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:5173", "http://127.0.0.1:5173"}),
		SyncBatchSize:      getEnvInt("SYNC_BATCH_SIZE", 50),
		SyncRetryAttempts:  getEnvInt("SYNC_RETRY_ATTEMPTS", 3),
		SyncRetryDelay:     getEnvDuration("SYNC_RETRY_DELAY", 5*time.Second),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFormat:          getEnv("LOG_FORMAT", "json"),
	}

	// Validate required settings
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate checks that all required configuration is present
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.ClerkSecretKey == "" {
		return fmt.Errorf("CLERK_SECRET_KEY is required")
	}
	if c.ClerkPublishableKey == "" {
		return fmt.Errorf("CLERK_PUBLISHABLE_KEY is required")
	}
	return nil
}

// IsDevelopment returns true if the app is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if the app is running in production mode
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// GetServerAddr returns the server address in host:port format
func (c *Config) GetServerAddr() string {
	return ":" + c.Port
}

// Helper functions

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if dur, err := time.ParseDuration(val); err == nil {
			return dur
		}
	}
	return defaultVal
}

func getEnvSlice(key string, defaultVal []string) []string {
	if val := os.Getenv(key); val != "" {
		// Simple split by comma - for production you might want a more robust parser
		return []string{val}
	}
	return defaultVal
}
