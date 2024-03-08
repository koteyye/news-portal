package service

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"

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
	newsAttr.Author = userID

	newsID, err := s.storage.CreateNews(ctx, newsAttr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create user to storage: %w", err)
	}

	return newsID, nil
}
