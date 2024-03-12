package news_cleaner

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/koteyye/news-portal/internal/news/config"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/koteyye/news-portal/pkg/s3"
	mock_storage "github.com/koteyye/news-portal/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
	"time"
)

const (
	testBucketName = "testbucket"
	testFileName   = "test_file.jpeg"
)

var testCfg = config.Config{
	S3Address:   "127.0.0.1:9001",
	S3KeyID:     "I6Htx3aLeTs6lhz4",
	S3SecretKey: "6sb7DoFAlJ60Epi1d6FgktHRB8u9zgyJ",
}

func TestNewsCleaner_StartWorker(t *testing.T) {
	t.Run("clean", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		storage := mock_storage.NewMockStorage(c)
		minio, err := s3.InitS3Repo(testCfg.S3Address, testCfg.S3KeyID, testCfg.S3SecretKey, false)
		assert.NoError(t, err)
		err = minio.Ping(context.Background())
		if errors.Is(err, s3.ErrPing) {
			return
		}
		assert.NoError(t, err)

		file, err := os.Open("./" + testFileName)
		assert.NoError(t, err)
		fileUUID, err := uuid.NewV4()
		assert.NoError(t, err)

		info, mimetype, err := minio.UploadFile(context.Background(), file, testBucketName, testFileName, 10)
		assert.NoError(t, err)
		assert.NotNil(t, info)

		opts := &slog.HandlerOptions{Level: slog.LevelInfo}
		handler := slog.NewTextHandler(os.Stdout, opts)
		logger := slog.New(handler)

		storage.EXPECT().GetDeletingFiles(gomock.Any()).Return([]*models.File{
			{
				ID:         fileUUID.String(),
				MimeType:   mimetype,
				BucketName: testBucketName,
				FileName:   testFileName,
			},
		}, error(nil))

		storage.EXPECT().SetHardDeletedFilesByIDs(gomock.Any(), []uuid.UUID{fileUUID}).Return(error(nil))

		cleaner := &NewsCleaner{
			ticker:  time.NewTicker(time.Second * 1),
			storage: storage,
			logger:  logger,
			s3:      minio,
		}

		go func() {
			cleaner.StartWorker(context.Background())
		}()

		timer := time.NewTicker(time.Second * 5)
		<-timer.C
		_, err = minio.GetFile(context.Background(), testBucketName, testFileName)
		assert.NoError(t, err)
	})
}
