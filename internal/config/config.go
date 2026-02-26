package config

import (
	"os"
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	Logger   LoggerConfig   `mapstructure:"logger" yaml:"logger"`
}

type ServerConfig struct {
	Host string `mapstructure:"host" yaml:"host" env:"SERVER_HOST" envDefault:"localhost"`
	Port int    `mapstructure:"port" yaml:"port" env:"SERVER_PORT" envDefault:"8080"`
}

type DatabaseConfig struct {
	DSN             string        `mapstructure:"dsn" yaml:"dsn" env:"DATABASE_DSN" envDefault:"postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns" env:"DATABASE_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns" env:"DATABASE_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime" env:"DATABASE_CONN_MAX_LIFETIME" envDefault:"5m"`
}

type RedisConfig struct {
	Address      string        `mapstructure:"address" yaml:"address" env:"REDIS_ADDRESS" envDefault:"localhost:6379"`
	Password     string        `mapstructure:"password" yaml:"password" env:"REDIS_PASSWORD" envDefault:""`
	DB           int           `mapstructure:"db" yaml:"db" env:"REDIS_DB" envDefault:"0"`
	MaxRecordTTL time.Duration `mapstructure:"max_record_ttl" yaml:"max_record_ttl" env:"REDIS_MAX_RECORD_TTL" envDefault:"5m"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level" yaml:"level" env:"LOG_LEVEL" envDefault:"info"`
}

func LoadConfig() (*Config, error) {
	cfg := config.New()

	if err := cfg.LoadEnvFiles("./.env"); err != nil {
		return nil, err
	}

	if err := cfg.LoadConfigFiles("./config.yaml"); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var appConfig Config
	if err := cfg.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	return &appConfig, nil
}
