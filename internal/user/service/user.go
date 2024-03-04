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

	addRoles := setRoles(data.Profile.Roles)
	if len(addRoles) != 0 {
		err := s.storage.SetUserRoles(ctx, userID, data.Profile.Roles)
		if err != nil {
			return uuid.Nil, fmt.Errorf("can't set role: %w", err)
		}
	}
	return userID, nil
}

func (s *Service) EditUser(ctx context.Context, data *models.Profile) error {
	err := s.storage.EditUserByID(ctx, data)
	if err != nil {
		s.logger.Error(err.Error())
		return fmt.Errorf("can't edit user: %w", err)
	}
	addRoles := setRoles(data.Roles)
	if len(addRoles) != 0 {
		err := s.storage.SetUserRoles(ctx, uuid.FromStringOrNil(data.ID), data.Roles)
		if err != nil {
			return fmt.Errorf("can't set role: %w", err)
		}
	}
	return nil
}

func (s *Service) DeleteUsersByIDs(ctx context.Context, userIDs []string) error {
	userUUIDs := make([]uuid.UUID, 0, len(userIDs))
	for _, userID := range userIDs {
		userUUID, err := uuid.FromString(userID)
		if err != nil {
			s.logger.Error(err.Error())
			return errors.New("can't parse userID")
		}
		userUUIDs = append(userUUIDs, userUUID)
	}
	err := s.storage.DeleteUserByIDs(ctx, userUUIDs)
	if err != nil {
		return fmt.Errorf("can't delete: %w", err)
	}
	return nil
}

func setRoles(roles []string) []string {
	if len(roles) == 0 {
		return roles
	}

	var addRoles []string
	for _, role := range roles {
		if role != models.DefaultRole {
			addRoles = append(addRoles, role)
		}
	}
	return addRoles
}