package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/internal/user/models"
)


func (s *Storage) GetUserListByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error) {
	//TOBE
	return nil, nil
}

func (s *Storage) EditUserByID(ctx context.Context, userID uuid.UUID) error {
	//TOBE
	return nil
}

func (s *Storage) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	//TOBE
	return nil
}