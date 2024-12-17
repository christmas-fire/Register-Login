package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type PostgreSQLConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Sslmode  string `yaml:"sslmode"`
}

func LoadConfig(path string) (*PostgreSQLConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault("user", "postgres")
	viper.SetDefault("password", "290605")
	viper.SetDefault("database", "register-login")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", "5432")
	viper.SetDefault("sslmode", "disable")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	var config PostgreSQLConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ошибка декодирования файла конфигурации: %w", err)
	}

	return &config, nil
}
