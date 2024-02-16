package service

import (
	"log/slog"

	"github.com/koteyye/news-portal/internal/news/storage"
)

type Service struct {
	storage storage.Storage
	logger *slog.Logger
}

func NewService(storage storage.Storage, logger *slog.Logger) *Service {
	return &Service{storage: storage, logger: logger}
}