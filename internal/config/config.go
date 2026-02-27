package config

import (
	"fmt"
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
	Host string `mapstructure:"host" yaml:"host" env:"SERVER_HOST" envDefault:"0.0.0.0"`
	Port int    `mapstructure:"port" yaml:"port" env:"SERVER_PORT" envDefault:"8080"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"DATABASE_HOST" yaml:"host" env:"DATABASE_HOST" envDefault:"localhost"`
	Port            string        `mapstructure:"DATABASE_PORT" yaml:"port" env:"DATABASE_PORT" envDefault:"5432"`
	User            string        `mapstructure:"DATABASE_USER" yaml:"user" env:"DATABASE_USER" envDefault:"postgres"`
	Password        string        `mapstructure:"DATABASE_PASSWORD" yaml:"password" env:"DATABASE_PASSWORD" envDefault:""`
	DatabaseName    string        `mapstructure:"DATABASE_NAME" yaml:"database_name" env:"DATABASE_NAME" envDefault:""`
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns" env:"DATABASE_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns" env:"DATABASE_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime" env:"DATABASE_CONN_MAX_LIFETIME" envDefault:"5m"`
}

func (dc DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dc.Host, dc.Port, dc.User, dc.Password, dc.DatabaseName)
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
	cfg.EnableEnv("")

	var appConfig Config
	if err := cfg.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	appConfig.Database.Host = os.Getenv("DATABASE_HOST")
	appConfig.Database.Port = os.Getenv("DATABASE_PORT")
	appConfig.Database.User = os.Getenv("DATABASE_USER")
	appConfig.Database.Password = os.Getenv("DATABASE_PASSWORD")
	appConfig.Database.DatabaseName = os.Getenv("DATABASE_NAME")

	return &appConfig, nil
}
