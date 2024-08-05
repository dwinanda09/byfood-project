package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	CORS     CORSConfig     `yaml:"cors"`
	Logging  LoggingConfig  `yaml:"logging"`
	API      APIConfig      `yaml:"api"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Driver         string `yaml:"driver"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Name           string `yaml:"name"`
	SSLMode        string `yaml:"ssl_mode"`
	MaxConnections int    `yaml:"max_connections"`
	MaxIdle        int    `yaml:"max_idle_connections"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

type LoggingConfig struct {
	Level       string `yaml:"level"`
	Environment string `yaml:"environment"`
}

type APIConfig struct {
	BasePath      string `yaml:"base_path"`
	EnableSwagger bool   `yaml:"enable_swagger"`
	SwaggerPath   string `yaml:"swagger_path"`
}

func (db *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode)
}

func Load() (*Config, error) {
	var config Config

	// Try to load from .env.yaml first, then fallback to .env.yaml.example
	configPaths := []string{".env.yaml", ".env.yaml.example"}
	var configData []byte
	var err error

	for _, path := range configPaths {
		configData, err = ioutil.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// Override with environment variables if present
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("HOST"); host != "" {
		config.Server.Host = host
	}
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		// Simple parsing for DATABASE_URL override
		// In production, you might want to use url.Parse
	}

	return &config, nil
}