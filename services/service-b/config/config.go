package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for Service B
type Config struct {
	Port            string
	WeatherAPIKey   string
	WeatherAPIURL   string
	OpenCEPURL      string
	RequestTimeout  time.Duration
	CacheTTL        time.Duration
	CacheCleanup    time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8081"),
		WeatherAPIKey:   getEnv("WEATHER_API_KEY", ""),
		WeatherAPIURL:   getEnv("WEATHER_API_URL", "http://api.weatherapi.com/v1"),
		OpenCEPURL:      getEnv("OPENCEP_API_URL", "https://opencep.com"),
		RequestTimeout:  getEnvDuration("REQUEST_TIMEOUT", 10*time.Second),
		CacheTTL:        getEnvDuration("CACHE_TTL", 1*time.Hour),
		CacheCleanup:    getEnvDuration("CACHE_CLEANUP", 10*time.Minute),
	}
}

// ValidateConfig validates required configuration
func (c *Config) ValidateConfig() error {
	if c.WeatherAPIKey == "" {
		return &ConfigError{Field: "WEATHER_API_KEY", Message: "é obrigatória"}
	}
	return nil
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

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// ConfigError represents configuration errors
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "configuração inválida: " + e.Field + " " + e.Message
}