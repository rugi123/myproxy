package config

import (
	"fmt"
	"os"

	"github.com/rugi123/myproxy/client/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Version float64 `mapstructure:"version"`
	Server  struct {
		Version    float64 `mapstructure:"version"`
		ServerIP   string  `mapstructure:"ip"`
		ServerPort int     `mapstructure:"port"`
	} `mapstructure:"server"`
	LogLevel string `mapstructure:"log_level"`
}

func Load(path string) (error, *Config) {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		err = generateConf(path)
		if err != nil {
			return fmt.Errorf("generate conf: %v", err), nil
		}
		logger.Info("new config generated")
		Load(path)
	}

	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("unmarshal error: %v", err), nil
	}

	return nil, &config
}

func generateConf(path string) error {
	conf := []byte(`version: 0.1
server:
  ip: "192.168.0.12"
  port: 8080
log_level: "extra"`)

	file, err := os.Create(path + "config.yaml")
	if err != nil {
		return fmt.Errorf("create conf err: %v", err)
	}

	if _, err = file.Write(conf); err != nil {
		return fmt.Errorf("write conf err: %v", err)
	}

	return nil
}
