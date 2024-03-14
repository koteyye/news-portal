package news_cleaner

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/s3"
	"github.com/koteyye/news-portal/pkg/storage"
	"log/slog"
	"time"
)

// NewsCleaner уборщик хранилища s3
type NewsCleaner struct {
	ticker  *time.Ticker
	storage storage.Storage
	logger  *slog.Logger
	s3      *s3.Handler
}

const tick = time.Hour * 1

// InitCleaner возвращает новый экземпляр NewsCleaner
func InitCleaner(storage storage.Storage, logger *slog.Logger, s3 *s3.Handler) *NewsCleaner {
	return &NewsCleaner{ticker: time.NewTicker(time.Hour * 1), storage: storage, logger: logger, s3: s3}
}

// StartWorker запускает обработчик очистки s3
func (n *NewsCleaner) StartWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			files, err := n.storage.GetDeletingFiles(ctx)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					n.logger.Error(fmt.Errorf("can't get deleting files: %w", err).Error())
					return
				}
				n.logger.Info("no deleting files")
				return
			}
			var deletedFilesIDs []uuid.UUID
			for _, file := range files {
				err = n.s3.RemoveFile(ctx, file.BucketName, file.FileName)
				if err != nil {
					n.logger.Error(fmt.Errorf("can't remove file in s3: %w", err).Error())
					return
				}
				deletedFileUUID, err := uuid.FromString(file.ID)
				if err != nil {
					n.logger.Error(fmt.Errorf("can't parse fileID: %w", err).Error())
					return
				}
				deletedFilesIDs = append(deletedFilesIDs, deletedFileUUID)
			}
			err = n.storage.SetHardDeletedFilesByIDs(ctx, deletedFilesIDs)
			if err != nil {
				n.logger.Error(fmt.Errorf("can't set hard deleted files: %w", err).Error())
				return
			}
			return
		case <-n.ticker.C:
			files, err := n.storage.GetDeletingFiles(ctx)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					n.logger.Error(fmt.Errorf("can't get deleting files: %w", err).Error())
					return
				}
				n.logger.Info("no deleting files")
			}
			if len(files) > 0 {
				var deletedFilesIDs []uuid.UUID
				for _, file := range files {
					err = n.s3.RemoveFile(ctx, file.BucketName, file.FileName)
					if err != nil {
						n.logger.Error(fmt.Errorf("can't remove file in s3: %w", err).Error())
						break
					}
					deletedFileUUID, err := uuid.FromString(file.ID)
					if err != nil {
						n.logger.Error(fmt.Errorf("can't parse fileID: %w", err).Error())
						break
					}
					deletedFilesIDs = append(deletedFilesIDs, deletedFileUUID)
				}
				err = n.storage.SetHardDeletedFilesByIDs(ctx, deletedFilesIDs)
				if err != nil {
					n.logger.Error(fmt.Errorf("can't set hard deleted files: %w", err).Error())
					break
				}
			}
			n.ticker.Reset(tick)
		}
	}
}
