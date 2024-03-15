package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adminresthandler "github.com/koteyye/news-portal/internal/user/admin_resthandler.go"
	"github.com/koteyye/news-portal/internal/user/config"
	"github.com/koteyye/news-portal/internal/user/grpchandler"
	"github.com/koteyye/news-portal/internal/user/resthandler"
	"github.com/koteyye/news-portal/internal/user/service"
	"github.com/koteyye/news-portal/pkg/s3"
	"github.com/koteyye/news-portal/pkg/signer"
	"github.com/koteyye/news-portal/pkg/storage/postgres"
	"github.com/koteyye/news-portal/server"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "github.com/koteyye/news-portal/proto"
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
	minio, err := s3.InitS3Handler(cfg.S3Address, cfg.S3KeyID, cfg.SecretKey, false)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var subnet *net.IPNet
	if cfg.TrustSubnet != "" {
		subnet, err = cfg.CIDR()
		if err != nil {
			logger.Error(err.Error())
		}
	}

	service := service.NewService(storage, minio, logger)
	signer := signer.New([]byte(cfg.SecretKey))
	restHandler := resthandler.NewRESTHandler(service, logger, cfg.CorsAllowed, signer)
	restAdminHandler := adminresthandler.NewAdminRESTHandler(service, logger, subnet)
	grpcHandler := grpchandler.InitGRPCHandlers(service, subnet)

	g.Go(func() error {
		runRESTServer(gCtx, cfg, restHandler, logger)
		return nil
	})

	g.Go(func() error {
		runRESTServer(gCtx, cfg, restAdminHandler, logger)
		return nil
	})

	g.Go(func() error {
		runGRPCServer(gCtx, cfg, grpcHandler, logger)
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

func runRESTServer(ctx context.Context, cfg *config.Config, restHandler any, log *slog.Logger) error {
	restServer := new(server.Server)
	switch handler := restHandler.(type) {
	case *resthandler.RESTHandler:
		go func() {
			log.Info(fmt.Sprintf("start rest server on %s", cfg.RESTAddress))
			if err := restServer.Run(cfg.RESTAddress, handler.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error(err.Error())
				return
			}
		}()
	case *adminresthandler.AdminRESTHandler:
		go func() {
			log.Info(fmt.Sprintf("start admin rest server on %s", cfg.AdminRESTAddress))
			if err := restServer.Run(cfg.AdminRESTAddress, handler.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error(err.Error())
				return
			}
		}()
	}

	<-ctx.Done()

	log.Info("shutting down rest server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := restServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("can't shutdown rest server: %w", err)
	}

	return nil
}

func runGRPCServer(ctx context.Context, cfg *config.Config, handler *grpchandler.GRPCHandler, log *slog.Logger) error {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			handler.SubnetInterceptor,
		),
	}
	s := grpc.NewServer(opts...)
	go func() {
		listen, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			log.Error(err.Error())
			return
		}
		pb.RegisterUserServer(s, handler)
		log.Info(fmt.Sprintf("start grpc server on %v", cfg.GRPCPort))
		if err = s.Serve(listen); err != nil {
			log.Error(err.Error())
			return
		}
	}()

	<-ctx.Done()
	log.Info("shutting down grpc server")
	s.GracefulStop()
	return nil
}
