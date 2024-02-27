package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/lib/pq"
)


func (s *Storage) GetUserListByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error) {
	var profiles []*models.Profile
	query := "select user_id, username, first_name, last_name, sur_name from profile where user_id = ANY($1)"

	rows, err := s.db.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, errorHandle(err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var profile models.Profile
		err = rows.Scan(&profile.ID, &profile.UserName, &profile.FirstName, &profile.LastName, &profile.SurName)
		if err != nil {
			return nil, errorHandle(err)
		}
		profiles = append(profiles, &profile)
	}
	return profiles, nil
}

func (s *Storage) CreateProfileByUserID(ctx context.Context, userID uuid.UUID, profile *models.Profile) error {
	query := "insert into profile (user_id, username, first_name, last_name, sur_name) values ($1, $2, $3, $4, $5)"

	_, err := s.db.ExecContext(ctx, query, userID, profile.UserName, profile.FirstName, profile.LastName, profile.SurName)
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