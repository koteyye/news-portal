package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v10"
)

// ErrConfigValidate обрабатываемая ошибка конфигурации
var ErrConfigValidate = errors.New("обязательные атрибуты конфигурации не заданы")

type configPath struct {
	path string `env:"CONFIG_PATH"`
}

// Config конфигурация сервиса
type Config struct {
	// RESTAddress адрес REST-сервера
	RESTAddress string `json:"rest_address"`

	// LogLevel уровень логирования
	LogLevel slog.Level `json:"log_level"`

	// DBDSN dsn хранилища
	DBDSN string `json:"db_dsn"`

	// S3Address адрес файлового хранилища
	S3Address string `json:"s3_address"`

	// CorsAllowed хосты, которым допускается вызывать REST-API сервиса
	CorsAllowed []string `json:"cors_allowed"`
}

func (c *Config) validate() error {
	if c.RESTAddress == "" {
		return fmt.Errorf("%w: не задан rest_address", ErrConfigValidate)
	}
	if c.DBDSN == "" {
		return fmt.Errorf("%w: не задан db_dsn", ErrConfigValidate)
	}
	if c.S3Address == "" {
		return fmt.Errorf("%w: не задан s3_address", ErrConfigValidate)
	}
	return nil
}

func (c *Config) configFromFile(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл конфигурации: %w", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return fmt.Errorf("не удалось десериализовать конфиг из файла: %w", err)
	}
	return nil
}

// GetConfig получить конфигурацию приложения
func GetConfig() (*Config, error) {
	var flagConfigPath string
	flag.StringVar(&flagConfigPath, "c", "", "file config path")
	flag.Parse()

	var envConfigPath configPath
	if err := env.Parse(&envConfigPath); err != nil {
		return nil, err
	}
	var config *Config
	if envConfigPath.path != "" {
		if err := config.configFromFile(envConfigPath.path); err != nil {
			return nil, err
		}

	} else {
		if err := config.configFromFile(flagConfigPath); err != nil {
			return nil, err
		}
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}
