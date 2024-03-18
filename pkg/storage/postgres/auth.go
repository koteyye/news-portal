package postgres

import (
	"context"

	"github.com/gofrs/uuid"
)

func (s *Storage) CreateLogin(ctx context.Context, login string, hashPassword string) (uuid.UUID, error) {
	var userID uuid.UUID

	query1 := "insert into users (login, hashed_password) values ($1, $2) returning id;"

	err := s.db.QueryRowContext(ctx, query1, login, hashPassword).Scan(&userID)
	if err != nil {
		return uuid.Nil, errorHandle(err)
	}

	return userID, nil
}

func (s *Storage) GetHashedPasswordByLogin(ctx context.Context, login string) (uuid.UUID, string, error) {
	var hashedPassword string
	var userID uuid.UUID

	query := "select id, hashed_password from users where login = $1"

	err := s.db.QueryRowContext(ctx, query, login).Scan(&userID, &hashedPassword)
	if err != nil {
		return uuid.Nil, "", errorHandle(err)
	}
	return userID, hashedPassword, nil
}
