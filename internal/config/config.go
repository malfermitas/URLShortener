package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
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
	cfg.EnableEnv("")

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

	err := LoadFromEnv(&appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}

func LoadFromEnv(ptr interface{}) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("LoadFromEnv: expected pointer to struct, got %T", ptr)
	}

	// loadStruct рекурсивно обрабатывает структуру
	var loadStruct func(v reflect.Value) error
	loadStruct = func(v reflect.Value) error {
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			// Пропускаем неэкспортные поля
			if !value.CanSet() {
				continue
			}

			// Если это вложенная структура — спускаемся
			if field.Type.Kind() == reflect.Struct && value.Kind() == reflect.Struct {
				if err := loadStruct(value); err != nil {
					return err
				}
				continue
			}

			envVar := field.Tag.Get("env")
			if envVar == "" {
				continue
			}

			envVal := os.Getenv(envVar)

			if envVal == "" {
				envVal = field.Tag.Get("envDefault")
			}

			if err := setFieldValue(value, envVal); err != nil {
				return fmt.Errorf("env %s: %w", envVar, err)
			}
		}

		return nil
	}

	v = v.Elem()
	return loadStruct(v)
}

// setFieldValue приводит строку к типу поля и присваивает значение
func setFieldValue(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(u)
		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)
		return nil

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
		return nil

	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(d))
			return nil
		}
		return fmt.Errorf("unsupported struct type: %s", v.Type())

	default:
		return fmt.Errorf("unsupported kind: %s", v.Kind())
	}
}
