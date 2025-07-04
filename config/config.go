package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	App      App      `yaml:"app"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Database struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Name            string `yaml:"name"`
	SSLMode         string `yaml:"sslmode"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

type App struct {
	Secret        string `yaml:"secret"`
	TokenTtlHours int    `yaml:"tokenTtlHours"`
}

func MustLoad() (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read config file: %w", err))
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to unmarshal config: %w", err))
	}

	return config, nil
}
