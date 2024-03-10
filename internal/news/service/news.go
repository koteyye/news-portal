package service

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	pb "github.com/koteyye/news-portal/proto"
	"google.golang.org/grpc/metadata"

	"github.com/koteyye/news-portal/pkg/models"
	"mime/multipart"
	"path/filepath"
)

const (
	newsBucketName    = "news"
	previewBucketName = "newspreviewimg"
)

func (s *Service) CreateNews(ctx context.Context,
	newsAttr *models.NewsAttributes,
	newsFile multipart.File,
	newsFileHeader *multipart.FileHeader,
	previewFile multipart.File,
	previewFileHeader *multipart.FileHeader,
	userID string,
) (uuid.UUID, error) {
	fileID, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't generate uuid: %w", err)
	}
	fileName := filepath.Ext(newsFileHeader.Filename)
	newFileName := fileID.String() + fileName

	_, mimeType, err := s.s3.UploadFile(ctx, newsFile, newsBucketName, newFileName, newsFileHeader.Size)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't upload file to s3: %w", err)
	}

	newsAttr.Content = &models.File{
		ID:         fileID.String(),
		MimeType:   mimeType,
		BucketName: newsBucketName,
		FileName:   newFileName,
	}

	if previewFile != nil {
		previewFileID, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, fmt.Errorf("can't generate uuid: %w", err)
		}
		previewFileName := filepath.Ext(previewFileHeader.Filename)
		previewNewFileName := previewFileID.String() + previewFileName

		_, previewMimeType, err := s.s3.UploadFile(ctx, previewFile, previewBucketName, previewNewFileName, previewFileHeader.Size)
		if err != nil {
			return uuid.Nil, fmt.Errorf("can't upload file to s3: %w", err)
		}

		newsAttr.Preview = &models.File{
			ID:         previewFileID.String(),
			MimeType:   previewMimeType,
			BucketName: previewBucketName,
			FileName:   previewFileName,
		}
	}
	author := models.Profile{ID: userID}
	newsAttr.Author = &author

	newsID, err := s.storage.CreateNews(ctx, newsAttr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create user to storage: %w", err)
	}

	return newsID, nil
}

func (s *Service) EditNews(ctx context.Context,
	newsID uuid.UUID,
	newsAttr *models.NewsAttributes,
	newsFile multipart.File,
	newsFileHeader *multipart.FileHeader,
	previewFile multipart.File,
	previewFileHeader *multipart.FileHeader,
	userID uuid.UUID) error {
	fileID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("can't generate uuid: %w", err)
	}
	fileName := filepath.Ext(newsFileHeader.Filename)
	newFileName := fileID.String() + fileName

	_, mimeType, err := s.s3.UploadFile(ctx, newsFile, newsBucketName, newFileName, newsFileHeader.Size)
	if err != nil {
		return fmt.Errorf("can't upload file to s3: %w", err)
	}

	newsAttr.Content = &models.File{
		ID:         fileID.String(),
		MimeType:   mimeType,
		BucketName: newsBucketName,
		FileName:   newFileName,
	}

	if previewFile != nil {
		previewFileID, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("can't generate uuid: %w", err)
		}
		previewFileName := filepath.Ext(previewFileHeader.Filename)
		previewNewFileName := previewFileID.String() + previewFileName

		_, previewMimeType, err := s.s3.UploadFile(ctx, previewFile, previewBucketName, previewNewFileName, previewFileHeader.Size)
		if err != nil {
			return fmt.Errorf("can't upload file to s3: %w", err)
		}

		newsAttr.Preview = &models.File{
			ID:         previewFileID.String(),
			MimeType:   previewMimeType,
			BucketName: previewBucketName,
			FileName:   previewFileName,
		}
	}

	err = s.storage.EditNewsByID(ctx, newsID, userID, newsAttr)
	if err != nil {
		return fmt.Errorf("can't edit news: %w", err)
	}
	return nil
}

func (s *Service) GetNewsByIDs(ctx context.Context, newsIDs []uuid.UUID) ([]*models.NewsAttributes, error) {
	news, err := s.storage.GetNewsByIDs(ctx, newsIDs)
	if err != nil {
		return nil, fmt.Errorf("can't get news: %w", err)
	}

	var userIDs []string
	for _, newsItem := range news {
		userIDs = append(userIDs, newsItem.Author.ID)
		if newsItem.UserCreated.ID != newsItem.Author.ID {
			userIDs = append(userIDs, newsItem.UserCreated.ID)
		}
		if newsItem.UserUpdated.ID != newsItem.Author.ID {
			userIDs = append(userIDs, newsItem.UserUpdated.ID)
		}
	}
	md := metadata.New(map[string]string{"X-Real-IP": s.serverAddress})
	ctx = metadata.NewOutgoingContext(ctx, md)
	w, err := s.userClient.GetUserByIDs(ctx, &pb.UserByIDsRequest{Userids: userIDs})
	if err != nil {
		return nil, fmt.Errorf("can't get user info: %w", err)
	}
	news = userAppend(news, w.Users)

	return news, nil
}

func (s *Service) GetNewsList(ctx context.Context, limit int, offset int) ([]*models.NewsAttributes, error) {
	news, err := s.storage.GetNewsList(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("can't get news list: %w", err)
	}

	var userIDs []string
	for _, newsItem := range news {
		userIDs = append(userIDs, newsItem.Author.ID)
		if newsItem.UserCreated.ID != newsItem.Author.ID {
			userIDs = append(userIDs, newsItem.UserCreated.ID)
		}
		if newsItem.UserUpdated.ID != newsItem.Author.ID {
			userIDs = append(userIDs, newsItem.UserUpdated.ID)
		}
	}
	md := metadata.New(map[string]string{"X-Real-IP": s.serverAddress})
	ctx = metadata.NewOutgoingContext(ctx, md)
	w, err := s.userClient.GetUserByIDs(ctx, &pb.UserByIDsRequest{Userids: userIDs})
	if err != nil {
		return nil, fmt.Errorf("can't get user info: %w", err)
	}
	news = userAppend(news, w.Users)

	return news, nil
}

func (s *Service) DeleteNewsByID(ctx context.Context, newsID uuid.UUID) error {
	err := s.storage.DeleteNewsByID(ctx, newsID)
	if err != nil {
		return fmt.Errorf("can't delete news: %w", err)
	}
	return nil
}

func userAppend(news []*models.NewsAttributes, users []*pb.Users) []*models.NewsAttributes {
	for _, newsItem := range news {
		for _, user := range users {
			user := models.Profile{
				ID:        user.UserID,
				UserName:  user.Username,
				FirstName: user.Firstname,
				LastName:  user.Lastname,
				SurName:   user.Surname,
				AvatarID:  uuid.FromStringOrNil(user.Avatar),
			}
			if newsItem.Author.ID == user.ID {
				newsItem.Author = &user
			}
			if newsItem.UserCreated.ID == user.ID {
				newsItem.UserCreated = &user
			}
			if newsItem.UserUpdated.ID == user.ID {
				newsItem.UserUpdated = &user
			}
		}
	}
	return news
}
