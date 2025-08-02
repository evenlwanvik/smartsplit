package config

import (
	"github.com/spf13/viper"
)

type Environment string

const (
	DevelopmentEnvironment Environment = "development"
	TestingEnvironment     Environment = "testing"
	StagingEnvironment     Environment = "staging"
	ProductionEnvironment  Environment = "production"
)

type Config struct {
	App *AppConfig `json:"app"`
}

type AppConfig struct {
	Env     Environment    `json:"env"`
	Port    int            `json:"port"`
	Limiter *LimiterConfig `json:"limiter"`
}

type LimiterConfig struct {
	RPS     float64 `json:"rps"`
	Burst   int     `json:"burst"`
	Enabled bool    `json:"enabled"`
}

func New() (*Config, error) {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(false)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
