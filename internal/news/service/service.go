package service

import (
	"github.com/koteyye/news-portal/pkg/s3"
	"github.com/koteyye/news-portal/pkg/storage"
	pb "github.com/koteyye/news-portal/proto"
	"log/slog"
)

type Service struct {
	storage       storage.Storage
	logger        *slog.Logger
	s3            *s3.Handler
	userClient    pb.UserClient
	serverAddress string
}

func NewService(storage storage.Storage, s3 *s3.Handler, logger *slog.Logger, userClient pb.UserClient, serverAddress string) *Service {
	return &Service{storage: storage, s3: s3, logger: logger, userClient: userClient, serverAddress: serverAddress}
}
