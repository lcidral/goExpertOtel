package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for Service A
type Config struct {
	Port           string
	ServiceBURL    string
	RequestTimeout time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		ServiceBURL:    getEnv("SERVICE_B_URL", "http://localhost:8081"),
		RequestTimeout: getEnvDuration("REQUEST_TIMEOUT", 30*time.Second),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvDuration gets a duration environment variable with a default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}