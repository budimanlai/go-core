package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	ServerHost string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret          string
	JWTExpirationHours int
	JWTIssuer          string

	Environment string
	LogLevel    string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8083"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "gocore"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		JWTIssuer:          getEnv("JWT_ISSUER", "go-core"),

		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
