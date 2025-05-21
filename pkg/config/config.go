package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config contiene la configuración de la aplicación
type Config struct {
	Port     string
	Host     string
	Database DatabaseConfig
}

// DatabaseConfig contiene la configuración de la base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Load carga la configuración desde el archivo .env
func Load() (*Config, error) {
	// Obtener el directorio de trabajo actual
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo directorio de trabajo: %v", err)
	}

	// Construir la ruta al archivo .env
	envPath := filepath.Join(workDir, ".env")

	// Leer el archivo .env manualmente
	envFile, err := os.ReadFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo .env: %v", err)
	}

	// Convertir el contenido a string y eliminar el BOM si existe
	envContent := strings.TrimPrefix(string(envFile), "\ufeff")

	// Crear un mapa temporal para las variables de entorno
	envMap := make(map[string]string)

	// Procesar cada línea del archivo
	for _, line := range strings.Split(envContent, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envMap[key] = value
		}
	}

	// Configurar las variables de entorno
	for key, value := range envMap {
		os.Setenv(key, value)
	}

	config := &Config{
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "localhost"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "finanzas"),
		},
	}

	// Validar que la contraseña no esté vacía
	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD no puede estar vacía. Por favor, configura una contraseña en el archivo .env")
	}

	return config, nil
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
