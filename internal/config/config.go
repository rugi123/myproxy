package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Version float64 `mapstructure:"version"`
	Port    int     `mapstructure:"port"`
}

type BaseConfig struct {
	App      AppConfig `mapstructure:"app"`
	LogLevel int       `mapstructure:"log_level"`
}

type ServerConfig struct {
	BaseConfig BaseConfig `mapstructure:",squash"`
	TunnelPort int        `mapstructure:"tunnel_port"`
}

type ClientConfig struct {
	BaseConfig BaseConfig `mapstructure:",squash"`
	Server     struct {
		IP   string `mapstructure:"ip"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
}

func loadConfig(path string, generateFn func(string) error, target interface{}) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := generateFn(path); err != nil {
				return fmt.Errorf("generate conf error: %w", err)
			}
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("read genconf error: %w", err)
			}
		}
	}

	if err := viper.Unmarshal(target); err != nil {
		return fmt.Errorf("unmarshal conf error: %w", err)
	}
	return nil
}

func LoadServer(path string) (*ServerConfig, error) {
	var conf ServerConfig
	if err := loadConfig(path, genServerConf, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func LoadClient(path string) (*ClientConfig, error) {
	var conf ClientConfig
	if err := loadConfig(path, genClientConf, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func genServerConf(path string) error {
	conf := `app:
  version: 0.1
  port: 8080
tunnel_port: 8080
log_level: 4`
	file, err := os.Create(path + "config.yaml")
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(conf))
	return err
}
func genClientConf(path string) error {
	conf := `app:
  version: 0.1
  port: 8080
server: 
  ip: 0
  port: 8080
log_level: 4`
	file, err := os.Create(path + "config.yaml")
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(conf))
	return err
}
