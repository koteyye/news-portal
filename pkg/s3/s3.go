package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/gabriel-vasile/mimetype"
)

// S3repo структура хранилища S3
type S3repo struct {
	client *minio.Client
}

// InitS3Repo возвращает новый экземпляр S3repo
func InitS3Repo(endpoint string, accessKeyID string, secretKey string, useSSL bool) (*S3repo, error) {
	client, err :=  minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretKey, ""),
		Secure: useSSL,
	})
	return &S3repo{client: client}, err
}

// UploadFile загружает файл в хранилище
func (s *S3repo) UploadFile(ctx context.Context, reader io.Reader, bucketName, filename string, fileSize int64) (minio.UploadInfo, string, error) {
	info, err := s.client.PutObject(ctx, bucketName, filename, reader, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return minio.UploadInfo{}, "", fmt.Errorf("не удалось загрузить файл в хранилище: %w", err)
	}
	mimeType, err := mimetype.DetectReader(reader)
	if err != nil {
		return  minio.UploadInfo{}, "", fmt.Errorf("не удалось определить тип загруженного файла: %w", err)
	}
	return info, mimeType.String(), nil
}

// RemoveFile удаляет файл из хранилища
func (s *S3repo) RemoveFile(ctx context.Context, bucketName, filename string) error {
	return s.client.RemoveObject(ctx, bucketName, filename, minio.RemoveObjectOptions{ForceDelete: true})
}

// GetFile получить файл из хранилища
func (s *S3repo) GetFile(ctx context.Context, bucketName, filename string) (*minio.Object, error) {
	return s.client.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
}