package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = getEnv("SERVER_PORT", "8080")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "stockflow"),
		DBPassword: getEnv("DB_PASSWORD", "stockflow_dev_password"),
		DBName:     getEnv("DB_NAME", "stockflow"),
		ServerPort: port,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
