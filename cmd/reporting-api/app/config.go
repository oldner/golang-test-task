// Package /////////////////////////////////////////////////////////////////////
package app

// Imports /////////////////////////////////////////////////////////////////////
import (
	"github.com/spf13/viper"
)

// Types ///////////////////////////////////////////////////////////////////////

// Config holds the configuration for the Reporting API service.
type Config struct {
	Port     string `mapstructure:"REPORTING_API_PORT"`
	RedisURL string `mapstructure:"REDIS_URL"`
}

// Functions ///////////////////////////////////////////////////////////////////

// LoadConfig reads the configuration from the .env file and environment variables.
// It returns a pointer to a Config struct and any error encountered.
func LoadConfig() (*Config, error) {
	// Set the config file name
	viper.SetConfigFile(".env")

	// Enable Viper to read environment variables
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

////////////////////////////////////////////////////////////////////////////////
