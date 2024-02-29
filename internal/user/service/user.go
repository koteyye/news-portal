package service

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
)

func (s *Service) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error) {
	profiles, err := s.storage.GetUserListByIDs(ctx, userIDs)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, errors.New("can't get user list by ID")
	}
	return profiles, nil
}
