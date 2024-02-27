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


// id uuid not null default gen_random_uuid() primary key unique,
//     user_id uuid not null,
//     username varchar(512) not null,
//     first_name varchar(512),
//     last_name varchar(512),
//     sur_name varchar(512),
//     created_at timestamp default now(),
//     updated_at timestamp default now(),
//     deleted_at timestamp
//     avatar_id uuid,
//     foreign key (user_id) references users(id),
//     foreign key (avatar_id) references avatar(id)