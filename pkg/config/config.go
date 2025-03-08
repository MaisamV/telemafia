package config

import (
	"encoding/json"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig `json:"server"`
	DB     DBConfig     `json:"db"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port string `json:"port"`
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// Load loads configuration from environment variables or falls back to config file
func Load() (*Config, error) {
	config := &Config{}

	// Try to load from environment variables first
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.DB.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		config.DB.Port = dbPort
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.DB.User = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		config.DB.Password = dbPass
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.DB.DBName = dbName
	}

	// If environment variables are not set, try loading from config file
	if config.Server.Port == "" {
		file, err := os.Open("config.json")
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(config); err != nil {
			return nil, err
		}
	}

	// Set defaults if values are still empty
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.DB.Host == "" {
		config.DB.Host = "localhost"
	}
	if config.DB.Port == "" {
		config.DB.Port = "5432"
	}

	return config, nil
} 