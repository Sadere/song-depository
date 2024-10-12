package config

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	DefaultAddress = ":8080"
	DefaultLogLevel = "debug"
)

// App config struct
type Config struct {
	Address          string `mapstructure:"ADDRESS"`
	LogLevel         string `mapstructure:"LOG_LEVEL"`
	PostgresDSN      string `mapstructure:"DSN"`
	MusicInfoAddress string `mapstructure:"MUSIC_INFO_ADDRESS"`
}

// Reads config from app.env file and returns config struct
func NewConfig(path string) (*Config, error) {
	cfg := &Config{
		Address: DefaultAddress,
		LogLevel: DefaultLogLevel,
	}
	
	log.Printf("loading config from %s", path)

	viper.AddConfigPath(path)

	viper.SetConfigName("app")

	viper.SetConfigType("env")


	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, nil
}