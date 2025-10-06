package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	APIKey      string
	Environment string
}

func Load() (*Config, error) {
	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        getEnvWithDefault("PORT", "8080"),
		APIKey:      os.Getenv("API_KEY"),
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	requiredEnvs := map[string]string{
		"DATABASE_URL": c.DatabaseURL,
		"JWT_SECRET":   c.JWTSecret,
		"API_KEY":      c.APIKey,
	}

	var missingEnvs []string
	for envName, envValue := range requiredEnvs {
		if envValue == "" {
			missingEnvs = append(missingEnvs, envName)
		}
	}

	if len(missingEnvs) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missingEnvs)
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadWithFallback() (*Config, error) {
	config, err := Load()
	if err != nil {
		log.Printf("Warning: %v", err)
		log.Println("Please ensure all required environment variables are set")
		return nil, err
	}
	return config, nil
}
