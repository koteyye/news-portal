package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/lib/pq"
)

func (s *Storage) GetUserListByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error) {
	profiles := make([]*models.Profile, 0, len(userIDs))

	query1 := "select user_id, username, first_name, last_name, sur_name from profile where user_id = ANY($1) and deleted_at is null;"
	query2 := `select user_id, role_name from user_roles ur 
	left join roles r on r.id = ur.role_id 
	where ur.user_id = ANY($1);`

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
		if len(profiles) == 0 {
			return errors.New("value not found")
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
	query1 := "insert into profile (user_id, username, first_name, last_name, sur_name) values ($1, $2, $3, $4, $5);"
	query2 := "insert into user_roles (user_id, role_id) values ($1, (select id from roles where role_name = $2));"

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		_, err := s.db.ExecContext(ctx, query1, userID, profile.UserName, profile.FirstName, profile.LastName, profile.SurName)
		if err != nil {
			return fmt.Errorf("can't insert profile: %w", err)
		}
		_, err = s.db.ExecContext(ctx, query2, userID, models.DefaultRole)
		if err != nil {
			return fmt.Errorf("can't insert user_roles: %w", err)
		}
		return nil
	})
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) EditUserByID(ctx context.Context, profile *models.Profile) error {
	query := `update profile
	set username = $1, first_name = $2, last_name = $3, sur_name = $4
	where user_id = $5;`

	_, err := s.db.ExecContext(ctx, query, profile.UserName, profile.FirstName, profile.LastName, profile.SurName, profile.ID)
	if err != nil {
		return errorHandle(err)
	}

	return nil
}

func (s *Storage) DeleteUserByIDs(ctx context.Context, userIDs []uuid.UUID) error {
	query1 := "update users set deleted_at = now() where id = ANY($1);"
	query2 := "update profile set deleted_at = now() where user_id = ANY($1)"

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		_, err := s.db.ExecContext(ctx, query1, pq.Array(userIDs))
		if err != nil {
			return fmt.Errorf("can't delete user: %w", err)
		}

		_, err = s.db.ExecContext(ctx, query2, pq.Array(userIDs))
		if err != nil {
			return fmt.Errorf("can't delete profile: %w", err)
		}
		return nil
	})
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) SetUserRoles(ctx context.Context, userID uuid.UUID, roles []string) error {
	roleIDs := make([]uuid.UUID, 0, len(roles))

	query1 := "select id from roles where role_name = ANY($1)"

	query := `insert into user_roles (user_id, role_id) values `

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		rows, err := s.db.QueryContext(ctx, query1, pq.Array(roles))
		if err != nil {
			return fmt.Errorf("can't get role id from role_name: %w", err)
		}
		for rows.Next() {
			var roleID uuid.UUID
			err = rows.Scan(&roleID)
			if err != nil {
				return fmt.Errorf("can't scan role ID: %w", err)
			}
			roleIDs = append(roleIDs, roleID)
		}

		var values []interface{}
		for i, roleID := range roleIDs {
			values = append(values, userID, roleID)

			numFields := 2
			n := i * numFields

			query += `(`
			for j := 0; j < numFields; j++ {
				query += `$` + strconv.Itoa(n+j+1) + `,`
			}
			query = query[:len(query)-1] + `),`
		}
		query = query[:len(query)-1]

		_, err = s.db.ExecContext(ctx, query, values...)
		if err != nil {
			return fmt.Errorf("can't insert user_roles: %w", err)
		}
		return nil
	})
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) EditRoles(ctx context.Context, userID uuid.UUID, roles []string) error {
	roleIDs := make([]uuid.UUID, 0, len(roles))

	queryDel := `delete from user_roles where user_id = $1`

	query1 := "select id from roles where role_name = ANY($1)"

	query := `insert into user_roles (user_id, role_id) values`

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		_, err := s.db.ExecContext(ctx, queryDel, userID)
		if err != nil {
			return fmt.Errorf("can't delete current user_roles: %w", err)
		}

		rows, err := s.db.QueryContext(ctx, query1, pq.Array(roles))
		if err != nil {
			return fmt.Errorf("can't get role id from role_name: %w", err)
		}
		for rows.Next() {
			var roleID uuid.UUID
			err = rows.Scan(&roleID)
			if err != nil {
				return fmt.Errorf("can't scan role ID: %w", err)
			}
			roleIDs = append(roleIDs, roleID)
		}

		var values []interface{}
		for i, roleID := range roleIDs {
			values = append(values, userID, roleID)

			numFields := 2
			n := i * numFields

			query += `(`
			for j := 0; j < numFields; j++ {
				query += `$` + strconv.Itoa(n+j+1) + `,`
			}
			query = query[:len(query)-1] + `),`
		}
		query = query[:len(query)-1]

		_, err = s.db.ExecContext(ctx, query, values...)
		if err != nil {
			return fmt.Errorf("can't insert user_roles: %w", err)
		}
		return nil
	})
	if err != nil {
		return errorHandle(err)
	}
	return nil
}
