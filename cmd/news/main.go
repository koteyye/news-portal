package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/koteyye/news-portal/pkg/s3"
	"github.com/koteyye/news-portal/pkg/signer"
	pb "github.com/koteyye/news-portal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koteyye/news-portal/internal/news/config"
	"github.com/koteyye/news-portal/internal/news/resthandler"
	"github.com/koteyye/news-portal/internal/news/service"
	"github.com/koteyye/news-portal/pkg/storage/postgres"
	"github.com/koteyye/news-portal/server"
	"golang.org/x/sync/errgroup"

	_ "github.com/lib/pq"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	g, gCtx := errgroup.WithContext(ctx)

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("can't get config")
	}
	logger := newLogger(cfg)
	slog.SetDefault(logger)

	storage, err := postgres.NewStorage(cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = storage.Up(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	minio, err := s3.InitS3Repo(cfg.S3Address, cfg.S3KeyID, cfg.SecretKey, false)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	connect, err := grpc.Dial(cfg.UserServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ip := GetLocalIP()
	userClient := pb.NewUserClient(connect)
	newService := service.NewService(storage, minio, logger, userClient, ip)
	newSigner := signer.New([]byte(cfg.SecretKey))
	restHandler := resthandler.NewRESTHandler(newService, logger, cfg.CorsAllowed, newSigner)

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
		log.Info(fmt.Sprintf("start rest server on %s", cfg.RESTAddress))
		if err := restServer.Run(cfg.RESTAddress, handler.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(err.Error())
		}
	}()

	<-ctx.Done()

	log.Info("shutting down rest server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := restServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("can't shutdown rest server: %w", err)
	}

	return nil
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
