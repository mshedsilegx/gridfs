package config

import (
	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	MongoURI             string `mapstructure:"MONGO_URI"`
	MongoUser            string `mapstructure:"MONGO_USER"`
	MongoPass            string `mapstructure:"MONGO_PASS"`
	MongoDB              string `mapstructure:"MONGO_DB"`
	MongoGridFSPrefix     string `mapstructure:"MONGO_GRIDFS_PREFIX"`
	NumWorkers           int    `mapstructure:"NUM_WORKERS"`
	LargeFileThresholdMB int    `mapstructure:"LARGE_FILE_THRESHOLD_MB"`
}

// LoadConfig reads configuration from a file and returns a Config struct.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("properties")

	// Set default values
	viper.SetDefault("NUM_WORKERS", 10)
	viper.SetDefault("LARGE_FILE_THRESHOLD_MB", 20)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
