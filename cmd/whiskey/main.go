package main

import (
	"context"
	"flag"
	"github.com/pkg/errors"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/delivery/api"
	"github.com/whiskey-back/internal/repository"
	"github.com/whiskey-back/internal/server"
	"github.com/whiskey-back/internal/service"
	"github.com/whiskey-back/pkg/database/postgres"
	"github.com/whiskey-back/pkg/jwt"
	"github.com/whiskey-back/pkg/logger"
	"github.com/whiskey-back/pkg/sms"
	"github.com/whiskey-back/pkg/yoomoney"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configPath := new(string)

	flag.StringVar(configPath, "config-path", "config/config-local.yaml", "specify path to yaml")
	flag.Parse()

	configFile, err := os.Open(*configPath)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with os.Open config"))
	}

	cfg := config.Config{}
	if err := yaml.NewDecoder(configFile).Decode(&cfg); err != nil {
		logger.LogFatal(errors.Wrap(err, "err with Decode config"))
	}

	if err = logger.NewLogger(cfg.Telegram); err != nil {
		logger.LogFatal(err)
	}

	postgresClient, err := postgres.NewPostgres(cfg.PostgresDsn)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewPostgres"))
	}

	db, err := postgresClient.Database()
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with Gorm"))
	}

	//для запуска необходимо прописать в ./config,
	//если это тестовый запуск то необходимо закомментировать иницаилизацию NewMinio
	newMinio, err := minio.NewMinio(cfg.FileStorage)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with newMinio"))
	}

	yooMoney, err := yoomoney.NewYoomoney(db, cfg.Yoomoney)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewYoomoney"))
	}

	apiSMS, err := sms.NewSMS(cfg.SMS)
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewSMS"))
	}

	tokenManager, err := jwt.NewManager()
	if err != nil {
		logger.LogFatal(errors.Wrap(err, "err with token Manager"))
	}

	repos := repository.NewRepositories(db)

	services := service.NewServices(
		&cfg,
		repos,
		tokenManager,
		apiSMS,
		nil,
		yooMoney,
	)

	endpoints := api.NewHandlers(&cfg, services, tokenManager)

	srv := server.NewServer(&cfg, endpoints)

	go func() {

		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.LogFatal(errors.Wrap(err, "err with Run"))
		}

	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err = srv.Shutdown(ctx); err != nil {
		logger.LogFatal(errors.Wrap(err, "failed to stop server"))
	}

	if err = postgresClient.Close(); err != nil {
		logger.LogFatal(errors.Wrap(err, "failed to stop db"))
	}
}
