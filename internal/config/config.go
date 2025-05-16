package config

import "time"

type AppConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	Rest     Rest
}

type Rest struct {
	ListenPort   string        `envconfig:"LISTEN_PORT" required:"true"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" required:"true"`
	ServerName   string        `envconfig:"SERVER_NAME" required:"true"`
	Token        string        `envconfig:"TOKEN" required:"true"`
}
