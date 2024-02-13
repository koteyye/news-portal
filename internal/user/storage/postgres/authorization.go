package postgres

import (
	"context"

	"github.com/gofrs/uuid"
)

func (s *Storage) SignUp(ctx context.Context, login string, hashPassword string) (uuid.UUID, error) {
	//TOBE
	return uuid.Nil, nil
}

func (s *Storage) SignIn(ctx context.Context, login string, hashPassword string) (uuid.UUID, error) {
	//TOBE
	return uuid.Nil, nil
}

func (s *Storage) ChangePassword(ctx context.Context, login string, oldPassword string, newPassword string) error {
	//TOBE
	return nil
}