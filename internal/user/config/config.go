package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/caarlos0/env/v10"
)

// ErrConfigValidate обрабатываемая ошибка конфигурации
var ErrConfigValidate = errors.New("required config attributes is empty")

type configPath struct {
	path string `env:"CONFIG_PATH"`
}

// Config конфигурация сервиса
type Config struct {
	// RESTAddress адрес REST-сервера
	RESTAddress string `json:"rest_address"`

	// GRPCPort порт GRPC-сервера
	GRPCPort string `json:"grpc_port"`

	// LogLevel уровень логирования
	LogLevel slog.Level `json:"log_level"`

	// DBDSN dsn хранилища
	DBDSN string `json:"db_dsn"`

	// S3Address адрес файлового хранилища
	S3Address string `json:"s3_address"`

	// S3KeyID идентификатор ключа доступа в S3
	S3KeyID string `json:"s3_keyID"`

	// S3SecretKey ключ доступа в S3
	S3SecretKey string `json:"s3_secretKey"`

	// Subnet адрес доверенной подсети
	TrustSubnet string `json:"trust_subnet"`

	// CorsAllowed хосты, которым допускается вызывать REST-API сервиса
	CorsAllowed []string `json:"cors_allowed"`

	// SecretKey ключ для подписания
	SecretKey string `json:"secret_key"`

	// AdminRESTAddress адрес REST-сервера админки
	AdminRESTAddress string `json:"admin_address"`

	// AdminTrustSubnet адрес довернной подсети админки
	AdminTrustSubnet string `json:"admin_trust_subnet"`
}

func (c *Config) validate() error {
	if c.RESTAddress == "" {
		return fmt.Errorf("%w: rest_address", ErrConfigValidate)
	}
	if c.GRPCPort == "" {
		return fmt.Errorf("%w: grpc_port", ErrConfigValidate)
	}
	if c.DBDSN == "" {
		return fmt.Errorf("%w: db_dsn", ErrConfigValidate)
	}
	if c.S3Address == "" {
		return fmt.Errorf("%w: s3_address", ErrConfigValidate)
	}
	if c.S3KeyID == "" {
		return fmt.Errorf("%w: s3_keyID", ErrConfigValidate)
	}
	if c.TrustSubnet == "" {
		return fmt.Errorf("%w: trust_subnet", ErrConfigValidate)
	}
	if c.SecretKey == "" {
		return fmt.Errorf("%w: secret_key", ErrConfigValidate)
	}
	if c.AdminRESTAddress == "" {
		return fmt.Errorf("%w: admin_address", ErrConfigValidate)
	}
	if c.AdminTrustSubnet == "" {
		return fmt.Errorf("%w: admin_trust_subnet", ErrConfigValidate)
	}
	return nil
}

func configFromFile(filepath string) (*Config, error) {
	var c Config
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл конфигурации: %w", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return nil, fmt.Errorf("не удалось десериализовать конфиг из файла: %w", err)
	}
	return &c, nil
}

// CIDR получение *IPNet из конфигурации
func (c *Config) CIDR() (*net.IPNet, error) {
	_, subnet, err := net.ParseCIDR(c.TrustSubnet)
	return subnet, err
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
		c, err := configFromFile(envConfigPath.path)
		if err != nil {
			return nil, err
		}
		config = c

	} else if flagConfigPath != "" {
		c, err := configFromFile(flagConfigPath)
		if err != nil {
			return nil, err
		}
		config = c
	} else {
		return nil, errors.New("empty configpath")
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}
