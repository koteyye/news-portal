package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
)


func (s *Storage) GetUserListByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error) {
	//TOBE
	return nil, nil
}

func (s *Storage) CreateProfileByUserID(ctx context.Context, userID uuid.UUID, profile *models.Profile) error {
	query := "insert into profile (user_id, first_name, last_name, sur_name) values ($1, $2, $3, $4)"

	_, err := s.db.ExecContext(ctx, query, userID, profile.FirstName, profile.LastName, profile.SurName)
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) EditUserByID(ctx context.Context, userID uuid.UUID, profile *models.Profile) error {
	//TOBE

	return nil
}

func (s *Storage) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	//TOBE
	return nil
}