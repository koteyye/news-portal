package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koteyye/news-portal/internal/news/config"
	"github.com/koteyye/news-portal/internal/news/resthandler"
	"github.com/koteyye/news-portal/internal/news/service"
	"github.com/koteyye/news-portal/internal/news/storage/postgres"
	"github.com/koteyye/news-portal/server"
	"golang.org/x/sync/errgroup"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	g, gCtx := errgroup.WithContext(ctx)

	cfg, err := config.GetConfig()
	logger := newLogger(cfg)
	slog.SetDefault(logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	storage, err := postgres.NewStorage(cfg)
	if err != nil {
		logger.Error(err.Error())
	}
	service := service.NewService(storage, logger)
	restHandler := resthandler.NewRESTHandler(service, logger, cfg.CorsAllowed)

	g.Go(func() error {
		runRESTServer(gCtx, cfg, restHandler, logger)
		return nil
	})

	if err = g.Wait(); err != nil {
		logger.Error(err.Error())
	}
}

func newLogger(c *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: c.LogLevel}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

func runRESTServer(ctx context.Context, cfg *config.Config, handler *resthandler.RESTHandler, log *slog.Logger) error {
	restServer := new(server.Server)
	go func() {
		log.Info("выполняется запуск rest сервера")
		if err := restServer.Run(cfg.RESTAddress, handler.InitRoutes()); err != nil {
			log.Error(err.Error())
		}
	}()

	<-ctx.Done()

	log.Info("выполняется выключение rest сервера")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := restServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("не удалось отключить rest сервер: %w", err)
	}

	return nil
}
