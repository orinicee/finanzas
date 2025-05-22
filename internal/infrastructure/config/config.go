package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config representa la configuración de la aplicación
type Config struct {
	Database DatabaseConfig
}

// DatabaseConfig representa la configuración de la base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load carga la configuración desde las variables de entorno
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error al cargar el archivo .env: %v", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "finanzas"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}, nil
}

// getEnv obtiene una variable de entorno o un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
