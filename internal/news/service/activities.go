package service

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	pb "github.com/koteyye/news-portal/proto"
	"google.golang.org/grpc/metadata"
)

func (s *Service) GetLikesByNewsID(ctx context.Context, newsID uuid.UUID) ([]*models.Like, error) {
	likes, err := s.storage.GetLikesByNewsID(ctx, newsID)
	if err != nil {
		return nil, fmt.Errorf("can't get likes: %w", err)
	}

	var likersIDs []string
	for _, like := range likes {
		likersIDs = append(likersIDs, like.Liker.ID)
	}
	md := metadata.New(map[string]string{"X-Real-IP": s.serverAddress})
	ctx = metadata.NewOutgoingContext(ctx, md)
	w, err := s.userClient.GetUserByIDs(ctx, &pb.UserByIDsRequest{Userids: likersIDs})
	if err != nil {
		return nil, fmt.Errorf("can't get user info: %w", err)
	}
	for _, like := range likes {
		for _, user := range w.Users {
			if like.Liker.ID == user.UserID {
				like.Liker = &models.Profile{
					ID:        user.UserID,
					UserName:  user.Username,
					FirstName: user.Firstname,
					LastName:  user.Lastname,
					SurName:   user.Surname,
				}
			}
		}
	}
	return likes, nil
}

func (s *Service) CreateLike(ctx context.Context, newsID uuid.UUID, likerID uuid.UUID) error {
	err := s.storage.CreateLike(ctx, newsID, likerID)
	if err != nil {
		return fmt.Errorf("can't create like: %w", err)
	}
	return nil
}

func (s *Service) DeleteLike(ctx context.Context, newsID uuid.UUID, likerID uuid.UUID) error {
	err := s.storage.DeleteLike(ctx, newsID, likerID)
	if err != nil {
		return fmt.Errorf("can't delete like: %w", err)
	}
	return nil
}

func (s *Service) CreateComment(ctx context.Context, newsID uuid.UUID, comment *models.Comment) (uuid.UUID, error) {
	commentID, err := s.storage.CreateComment(ctx, newsID, comment)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create comment: %w", err)
	}
	return commentID, nil
}

func (s *Service) EditComment(ctx context.Context, comment *models.Comment) error {
	err := s.storage.EditComment(ctx, comment)
	if err != nil {
		return fmt.Errorf("can't edit comment: %w", err)
	}
	return nil
}

func (s *Service) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	err := s.storage.DeleteComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("can't delete comment: %w", err)
	}
	return nil
}

func (s *Service) GetCommentsByNewsID(ctx context.Context, newsID uuid.UUID) ([]*models.Comment, error) {
	comments, err := s.storage.GetCommentsByNewsID(ctx, newsID)
	if err != nil {
		return nil, fmt.Errorf("can't get comments: %w", err)
	}
	var authors []string
	for _, comment := range comments {
		authors = append(authors, comment.Author.ID)
	}
	md := metadata.New(map[string]string{"X-Real-IP": s.serverAddress})
	ctx = metadata.NewOutgoingContext(ctx, md)
	w, err := s.userClient.GetUserByIDs(ctx, &pb.UserByIDsRequest{Userids: authors})
	if err != nil {
		return nil, fmt.Errorf("can't get user info: %w", err)
	}
	for _, comment := range comments {
		for _, user := range w.Users {
			if comment.Author.ID == user.UserID {
				comment.Author = &models.Profile{
					ID:        user.UserID,
					UserName:  user.Username,
					FirstName: user.Firstname,
					LastName:  user.Lastname,
					SurName:   user.Surname,
				}
			}
		}
	}

	return comments, nil
}
