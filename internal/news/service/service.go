package service

import (
	"github.com/koteyye/news-portal/pkg/s3"
	"github.com/koteyye/news-portal/pkg/storage"
	pb "github.com/koteyye/news-portal/proto"
	"log/slog"
)

type Service struct {
	storage    storage.Storage
	logger     *slog.Logger
	s3         *s3.S3repo
	userClient pb.UserClient
}

func NewService(storage storage.Storage, s3 *s3.S3repo, logger *slog.Logger, userClient pb.UserClient) *Service {
	return &Service{storage: storage, s3: s3, logger: logger, userClient: userClient}
}
