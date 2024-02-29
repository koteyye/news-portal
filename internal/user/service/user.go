package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/koteyye/news-portal/pkg/password"
)

func (s *Service) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error) {
	profiles, err := s.storage.GetUserListByIDs(ctx, userIDs)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New("can't get user list by ID")
	}
	return profiles, nil
}

func (s *Service) CreateUser(ctx context.Context, data *models.UserData) (uuid.UUID, error) {
	hashedPassword, err := password.Hash(data.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't hash password: %w", err)
	}
	userID, err := s.storage.CreateLogin(ctx, data.Login, hashedPassword)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create login: %w", err)
	}
	err = s.storage.CreateProfileByUserID(ctx, userID, data.Profile)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create profile: %w", err)
	}
	var addRoles []string
	for _, role := range data.Profile.Roles {
		if role != models.DefaultRole {
			addRoles = append(addRoles, role)
		}
	}
	if len(addRoles) != 0 {
		err := s.storage.SetUserRoles(ctx, userID, addRoles)
		if err != nil {
			return uuid.Nil, fmt.Errorf("can't set role: %w", err)
		}
	}
	return userID, nil
}