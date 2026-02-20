package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Version float64 `mapstructure:"version"`
	Server  struct {
		Version    float64 `mapstructure:"version"`
		ServerIP   string  `mapstructure:"ip"`
		ServerPort int     `mapstructure:"port"`
	} `mapstructure:"server"`
	LogLevel string `mapstructure:"version"`
}

func Load(path string) (error, *Config) {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read error: %v", err), nil
	}

	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("unmarshal error: %v", err), nil
	}

	return nil, &config
}
