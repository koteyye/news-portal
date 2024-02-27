package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
)

func (s *Storage) SignUp(ctx context.Context, login string, hashPassword string) (uuid.UUID, error) {
	var userID uuid.UUID

	query1 := "insert into users (login, hashed_password) values ($1, $2) returning id;"

	err := s.db.QueryRowContext(ctx, query1, login, hashPassword).Scan(&userID)
	if err != nil {
		return uuid.Nil, errorHandle(err)
	}

	return userID, nil
}

func (s *Storage) SignIn(ctx context.Context, login string, hashPassword string) (*models.Profile, error) {
	//TOBE
	return nil, nil
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