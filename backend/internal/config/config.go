package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server  ServerConfig  `json:"server"`
	MongoDB MongoDBConfig `json:"mongodb"`
	JWT     JWTConfig     `json:"jwt"`
	Broker  BrokerConfig  `json:"broker"`
	Logging LoggingConfig `json:"logging"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port            string        `json:"port"`
	ReadTimeout     time.Duration `json:"readTimeout"`
	WriteTimeout    time.Duration `json:"writeTimeout"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`
	AllowedOrigins  []string      `json:"allowedOrigins"`
	AllowedMethods  []string      `json:"allowedMethods"`
	AllowedHeaders  []string      `json:"allowedHeaders"`
}

// MongoDBConfig represents the MongoDB configuration
type MongoDBConfig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
}

// JWTConfig represents the JWT configuration
type JWTConfig struct {
	Secret           string        `json:"secret"`
	ExpirationTime   time.Duration `json:"expirationTime"`
	RefreshSecret    string        `json:"refreshSecret"`
	RefreshExpiryTime time.Duration `json:"refreshExpiryTime"`
}

// BrokerConfig represents the broker configuration
type BrokerConfig struct {
	DefaultBroker string                 `json:"defaultBroker"`
	Brokers       map[string]interface{} `json:"brokers"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level      string `json:"level"`
	File       string `json:"file"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
	Compress   bool   `json:"compress"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            "8080",
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 15 * time.Second,
			AllowedOrigins:  []string{"*"},
			AllowedMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:  []string{"Content-Type", "Authorization"},
		},
		MongoDB: MongoDBConfig{
			URI:      "mongodb://localhost:27017",
			Database: "trading_platform",
		},
		JWT: JWTConfig{
			Secret:           "your-secret-key",
			ExpirationTime:   24 * time.Hour,
			RefreshSecret:    "your-refresh-secret-key",
			RefreshExpiryTime: 7 * 24 * time.Hour,
		},
		Broker: BrokerConfig{
			DefaultBroker: "simulator",
			Brokers: map[string]interface{}{
				"simulator": map[string]interface{}{
					"enabled": true,
				},
			},
		},
		Logging: LoggingConfig{
			Level:      "info",
			File:       "logs/app.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		},
	}
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
