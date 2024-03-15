package service

import (
	"log/slog"

	"github.com/koteyye/news-portal/pkg/s3"
	"github.com/koteyye/news-portal/pkg/storage"
)

// Service структура сервисного слоя
type Service struct {
	storage storage.Storage
	s3      *s3.Handler
	logger  *slog.Logger
}

func NewService(storage storage.Storage, s3 *s3.Handler, logger *slog.Logger) *Service {
	return &Service{storage: storage, s3: s3, logger: logger}
}
