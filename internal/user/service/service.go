package service

import (
	"log/slog"

	"github.com/koteyye/news-portal/internal/user/storage"
)

// Service структура сервисного слоя
type Service struct {
	storage storage.Storage
	logger *slog.Logger
}

func NewService(storage storage.Storage, logger *slog.Logger) *Service {
	return &Service{storage: storage, logger: logger}
}