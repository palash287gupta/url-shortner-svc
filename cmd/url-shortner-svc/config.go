package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Port    string
	BaseURL string
}

func LoadConfig() *Config {
	viper.SetDefault("port", "8080")
	viper.SetDefault("base_url", "http://localhost:8080")

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	viper.BindEnv("port", "PORT")
	viper.BindEnv("base_url", "BASE_URL")

	cfg := &Config{
		Port:    viper.GetString("port"),
		BaseURL: viper.GetString("base_url"),
	}

	log.WithFields(log.Fields{
		"port":    cfg.Port,
		"baseURL": cfg.BaseURL,
	}).Info("Configuration loaded")

	return cfg
}
