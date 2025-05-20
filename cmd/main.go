package main

import (
	"context"
	"os"
	"os/signal"
	"restapi/internal/api"
	"restapi/internal/config"
	"restapi/internal/logger"
	"restapi/internal/repo/db"
	"restapi/internal/service"
	"syscall"

	"github.com/pkg/errors"
)

func main() {
	// Загрузка конфигурации
	cfg := config.NewConfig()

	// Инициализация логгера
	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error creating logger"))
	}

	log.Info("Config: ", cfg)

	// Инициализация репозитория
	repo, err := db.NewRepo(context.Background(), log, *cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error creating repository"))
	}

	// repo, err := memory.NewRepo(context.Background(), log, *cfg) // Для подключения репозитория in-memory

	// Инициализация сервиса
	service := service.NewService(log, repo)

	// Инициализация API
	app := api.NewRouters(&api.Routers{Service: service}, cfg.Rest.Token, log)

	// Запуск сервера в горутине
	go func() {
		log.Info("Server started on port: ", cfg.Rest.ListenPort)
		if err := app.Listen(":" + cfg.Rest.ListenPort); err != nil {
			log.Fatal(errors.Wrap(err, "error starting server"))
		}
	}()

	// Ожидание завершения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	app.Shutdown()
	log.Info("Shutting down server...")
}