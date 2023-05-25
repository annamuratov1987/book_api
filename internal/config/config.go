package config

import (
	"github.com/spf13/viper"
	"time"
)
import "github.com/kelseyhightower/envconfig"

type Config struct {
	DB Postgres

	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Auth struct {
		TokenTTL time.Duration `mapstructure:"token_ttl"`
	} `mapstructure:"auth"`
}

type Postgres struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

func New(folder, filename string) (*Config, error) {
	cnf := new(Config)

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cnf); err != nil {
		return nil, err
	}

	if err := envconfig.Process("db", &cnf.DB); err != nil {
		return nil, err
	}

	return cnf, nil
}
