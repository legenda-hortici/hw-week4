package main

import (
	"log"
	"os"
	"os/signal"
	"restapi/internal/api"
	"restapi/internal/config"
	"restapi/internal/logger"
	"restapi/internal/repo"
	"restapi/internal/service"
	"syscall"

	"github.com/pkg/errors"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading .env file")
	}

	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Config: ", cfg)

	repo := repo.NewRepo(log)

	service := service.NewService(log, repo)

	app := api.NewRouters(&api.Routers{Service: service}, cfg.Rest.Token, log)

	go func() {
		log.Info("Server started on port: ", cfg.Rest.ListenPort)
		if err := app.Listen(":" + cfg.Rest.ListenPort); err != nil {
			log.Fatal(errors.Wrap(err, "error starting server"))
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	log.Info("Shutting down server...")
}
