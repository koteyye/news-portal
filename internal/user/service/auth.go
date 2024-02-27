package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/koteyye/news-portal/pkg/password"
)

// SignUp регистрация пользователя и наполнение профиля
func (s *Service) SignUp(ctx context.Context, input *models.UserData) (*models.Profile, error) {
	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		s.logger.Info(err.Error())
		return nil, err
	}
	userID, err := s.storage.CreateLogin(ctx, input.Login, hashedPassword)
	if err != nil {
		s.logger.Info(err.Error())
		return nil, fmt.Errorf("can't create user: %w", err)
	}
	
	if input.Profile == nil {
		input.Profile = &models.Profile{UserName: input.Login}
	}
	err = s.storage.CreateProfileByUserID(ctx, userID, input.Profile)
		if err != nil {
			s.logger.Info(err.Error())
			return nil, err
		}
	input.Profile.ID = userID.String()
	return input.Profile, nil
}

func (s *Service) SignIn(ctx context.Context, input *models.UserData) (*models.Profile, error) {
	userID, hashedPassword, err := s.storage.GetHashedPasswordByLogin(ctx, input.Login)
	if err != nil {
		s.logger.Info(err.Error())
		return nil, err
	}
	ok := password.Compare(hashedPassword, input.Password)
	if !ok {
		return nil, errors.New("invalid login or password")
	}

	profile, err := s.storage.GetUserListByIDs(ctx, []uuid.UUID{userID})
	if err != nil {
		s.logger.Info(err.Error())
		return nil, err
	}
	return profile[0], nil
}