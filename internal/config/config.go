package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// AppConfig конфигурация приложения
type AppConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	Rest     Rest
	Database Database
}

// Rest конфигурация API
type Rest struct {
	ListenPort   string        `envconfig:"LISTEN_PORT" required:"true"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`
	ServerName   string        `envconfig:"SERVER_NAME" required:"true"`
	Token        string        `envconfig:"TOKEN" required:"true"`
}

// Database конфигурация базы данных
type Database struct {
	Host                string        `envconfig:"DB_HOST" required:"true"`
	Port                int           `envconfig:"DB_PORT" required:"true"`
	User                string        `envconfig:"DB_USER" required:"true"`
	Password            string        `envconfig:"DB_PWD" required:"true"`
	DBName              string        `envconfig:"DB_NAME" required:"true"`
	SSLMode             string        `envconfig:"DB_SSL_MODE" required:"true"`
	PoolMaxConns        int           `envconfig:"DB_POOL_MAX_CONNS" required:"true"`
	PoolMaxConnLifetime time.Duration `envconfig:"DB_POOL_MAX_CONN_LIFETIME" required:"true"`
	PoolMaxConnIdleTime time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" required:"true"`
}

// NewConfig загружает конфигурацию из файла .env
func NewConfig() *AppConfig {
	if err := godotenv.Load(".env"); err != nil {
		panic("error loading .env file")
	}

	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		panic("error loading environment variables")
	}

	return &cfg
}
