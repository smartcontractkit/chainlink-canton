package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                   int
	Host                   string
	LedgerAPIHost          string
	LedgerAPIPort          int
	JWTSecret              string
	JWTAudience            string
	EnvironmentsConfigPath string
}

// load configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:                   getEnvInt("PORT", 8090),
		Host:                   getEnv("HOST", "0.0.0.0"),
		LedgerAPIHost:          getEnv("LEDGER_API_HOST", "localhost"),
		LedgerAPIPort:          getEnvInt("LEDGER_API_PORT", 10001),
		JWTSecret:              getEnv("JWT_SECRET", ""),
		JWTAudience:            getEnv("JWT_AUDIENCE", "https://canton.network.global"),
		EnvironmentsConfigPath: getEnv("ENVIRONMENTS_CONFIG", "./environments.yaml"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}

	return defaultValue
}
