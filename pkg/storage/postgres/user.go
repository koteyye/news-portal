package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/lib/pq"
)


func (s *Storage) GetUserListByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error) {
	profiles := make([]*models.Profile, 0, len(userIDs))

	query1 := "select user_id, username, first_name, last_name, sur_name from profile where user_id = ANY($1)"
	query2 := `select user_id, role_name from user_roles ur 
	left join roles r on r.id = ur.role_id 
	where ur.user_id = ANY($1)`

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		rows, err := s.db.QueryContext(ctx, query1, pq.Array(userIDs))
		if err != nil {
			return fmt.Errorf("can't get profile: %w", err)
		}
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			var profile models.Profile
			err = rows.Scan(&profile.ID, &profile.UserName, &profile.FirstName, &profile.LastName, &profile.SurName)
			if err != nil {
				return fmt.Errorf("can't scan profile: %w", err)
			}
			profiles = append(profiles, &profile)
		}

		rows, err = s.db.QueryContext(ctx, query2, pq.Array(userIDs))
		if err != nil {
			return fmt.Errorf("can't get roles: %w", err)
		}
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			var userID uuid.UUID
			var role string
			err := rows.Scan(&userID, &role)
			if err != nil {
				return fmt.Errorf("can't scan roles: %w", err)
			}
			for i := range profiles {
				if profiles[i].ID == userID.String() {
					profiles[i].Roles = append(profiles[i].Roles, role)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errorHandle(err)
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