package service

import (
	"context"
	"fmt"

	"github.com/koteyye/news-portal/pkg/models"
	"github.com/koteyye/news-portal/pkg/password"
)

// SignUp регистрация пользователя и наполнение профиля
func (s *Service) SignUp(ctx context.Context, input *models.UserData) (*models.Profile, error) {
	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	userID, err := s.storage.SignUp(ctx, input.Login, hashedPassword)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, fmt.Errorf("can't create user: %w", err)
	}
	
	var profile *models.Profile
	if input.Profile != nil {
		err := s.storage.CreateProfileByUserID(ctx, userID, input.Profile)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, err
		}
	}
	return profile, nil
}