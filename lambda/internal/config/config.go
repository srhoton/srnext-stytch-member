package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	StytchProjectID string
	StytchSecret    string
	StytchEnv       string
	LogLevel        string
	AWSRegion       string
	RequestTimeout  int
	MaxRetries      int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		StytchProjectID: os.Getenv("STYTCH_PROJECT_ID"),
		StytchSecret:    os.Getenv("STYTCH_SECRET"),
		StytchEnv:       getEnvWithDefault("STYTCH_ENV", "test"),
		LogLevel:        getEnvWithDefault("LOG_LEVEL", "info"),
		AWSRegion:       getEnvWithDefault("AWS_REGION", "us-east-1"),
		RequestTimeout:  getEnvAsIntWithDefault("REQUEST_TIMEOUT", 30),
		MaxRetries:      getEnvAsIntWithDefault("MAX_RETRIES", 3),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.StytchProjectID == "" {
		return fmt.Errorf("STYTCH_PROJECT_ID is required")
	}

	if c.StytchSecret == "" {
		return fmt.Errorf("STYTCH_SECRET is required")
	}

	if c.StytchEnv != "test" && c.StytchEnv != "live" {
		return fmt.Errorf("STYTCH_ENV must be either 'test' or 'live'")
	}

	if c.RequestTimeout <= 0 {
		return fmt.Errorf("REQUEST_TIMEOUT must be positive")
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("MAX_RETRIES cannot be negative")
	}

	return nil
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntWithDefault gets an environment variable as int with a default value
func getEnvAsIntWithDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var intValue int
	if _, err := fmt.Sscanf(value, "%d", &intValue); err != nil {
		return defaultValue
	}

	return intValue
}
