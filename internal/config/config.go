package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	// Redis configuration
	Redis struct {
		Host     string
		Port     string
		Password string
	}

	// Server configuration
	Server struct {
		Port         string
		AllowOrigins string
	}
	MaximumShortUrlCount int // Maximum number of short URLs
	Expiration           int // Default expiration time for cache entries in seconds
}

// viperInstance is a singleton instance of viper
var viperInstance *viper.Viper

// InitConfig initializes the viper configuration
func InitConfig() {
	viperInstance = viper.New()

	// Set default values
	setDefaults()

	// Configure viper to read from environment variables
	viperInstance.AutomaticEnv()

	// Configure viper to read from .env file
	viperInstance.SetConfigName(".env")
	viperInstance.SetConfigType("env")
	viperInstance.AddConfigPath(".")

	// Read the config file
	if err := viperInstance.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Warning: Error reading config file: %v\n", err)
		}
	}

	// Make environment variables case-insensitive
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// setDefaults sets default values for configuration
func setDefaults() {
	// Server defaults
	viperInstance.SetDefault("server.port", "8080")

	// Redis defaults
	viperInstance.SetDefault("redis.host", "localhost")
	viperInstance.SetDefault("redis.port", "6379")

}

// Load loads the configuration from viper
func Load() (*Config, error) {
	// Initialize viper if it hasn't been initialized yet
	if viperInstance == nil {
		InitConfig()
	}
	config := &Config{}

	// Redis configuration
	config.Redis.Host = viperInstance.GetString("REDIS_HOST")
	config.Redis.Port = viperInstance.GetString("REDIS_PORT")
	config.Redis.Password = viperInstance.GetString("REDIS_PASSWORD")

	// Server configuration
	config.Server.Port = viperInstance.GetString("PORT")
	config.Server.AllowOrigins = viperInstance.GetString("ALLOW_ORIGINS")
	config.MaximumShortUrlCount = viperInstance.GetInt("MAXIMUM_SHORT_URL_COUNT")
	config.Expiration = viperInstance.GetInt("EXPIRATION")
	return config, nil
}

// GetViper returns the viper instance
func GetViper() *viper.Viper {
	if viperInstance == nil {
		InitConfig()
	}
	return viperInstance
}
