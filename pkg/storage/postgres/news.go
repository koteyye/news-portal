package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/lib/pq"
	"strconv"
)

func (s *Storage) CreateNews(ctx context.Context, newsAttr *models.NewsAttributes) (uuid.UUID, error) {
	var newsID uuid.UUID

	query1 := "insert into files (id, file_name, bucket_name, mime_type) values (:id, :file_name, :bucket_name, :mime_type)"
	query2 := `insert into news (title, author, description, content_id, preview_id, state, user_created, user_updated)
	values ($1, $2, $3, $4, $5, $6, $7, $8) returning id;`

	countFiles := 1
	if newsAttr.Preview != nil {
		countFiles += 1
	}
	files := make([]*models.File, 0, countFiles)
	files = append(files, newsAttr.Content)
	if newsAttr.Preview != nil {
		files = append(files, newsAttr.Preview)
	}

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		_, err := s.db.NamedExecContext(ctx, query1, files)
		if err != nil {
			return fmt.Errorf("can't insert files: %w", err)
		}

		err = s.db.QueryRowContext(
			ctx,
			query2,
			newsAttr.Title,
			newsAttr.AuthorInfo.ID,
			newsAttr.Description,
			newsAttr.Content.ID,
			newsAttr.Preview.ID,
			newsAttr.State,
			newsAttr.AuthorInfo.ID,
			newsAttr.AuthorInfo.ID).Scan(&newsID)
		if err != nil {
			return fmt.Errorf("can't insert news: %w", err)
		}

		return nil
	})

	if err != nil {
		return uuid.Nil, errorHandle(err)
	}

	return newsID, nil
}

func (s *Storage) EditNewsByID(ctx context.Context, newsID uuid.UUID, userUpdated uuid.UUID, newsAttr *models.NewsAttributes) error {
	query1 := "insert into files (id, file_name, bucket_name, mime_type) values"
	query2 := "select content_id, preview_id from news where id = $1"
	query3 := "update news set title = $1, description = $2, content_id = $3, preview_id = $4, user_updated = $5, updated_at = now() where id = $6"
	query4 := "update files set deleted_at = now() where id = any($1)"

	countFiles := 1
	if newsAttr.Preview != nil {
		countFiles += 1
	}
	files := make([]*models.File, 0, countFiles)
	files = append(files, newsAttr.Content)
	if newsAttr.Preview != nil {
		files = append(files, newsAttr.Preview)
	}

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		var values []interface{}
		for i, file := range files {
			values = append(values, file.ID, file.FileName, file.BucketName, file.MimeType)
			numFields := 4
			n := i * numFields

			query1 += `(`
			for j := 0; j < numFields; j++ {
				query1 += `$` + strconv.Itoa(n+j+1) + `,`
			}
			query1 = query1[:len(query1)-1] + `),`
		}
		query1 = query1[:len(query1)-1]

		_, err := s.db.ExecContext(ctx, query1, values...)
		if err != nil {
			return fmt.Errorf("can't insert files: %w", err)
		}

		var currentContentID uuid.UUID
		var currentPreviewID uuid.UUID
		err = s.db.QueryRowContext(ctx, query2, newsID).Scan(&currentContentID, &currentPreviewID)
		if err != nil {
			return fmt.Errorf("can't scan current contentID or current previewID: %w", err)
		}

		_, err = s.db.ExecContext(
			ctx,
			query3,
			newsAttr.Title,
			newsAttr.Description,
			newsAttr.Content.ID,
			newsAttr.Preview.ID,
			userUpdated,
			newsID)
		if err != nil {
			return fmt.Errorf("can't insert news: %w", err)
		}

		_, err = s.db.ExecContext(ctx, query4, pq.Array([]uuid.UUID{currentPreviewID, currentPreviewID}))
		if err != nil {
			return fmt.Errorf("can't delete previewes files: %w", err)
		}

		return nil
	})

	if err != nil {
		return errorHandle(err)
	}

	return nil
}

func (s *Storage) GetNewsByIDs(ctx context.Context, newsIDs []uuid.UUID) ([]*models.NewsAttributes, error) {
	var news []*models.NewsAttributes
	query1 := "select id, title, author, description, content_id, preview_id, state, created_at, updated_at, user_created, user_updated from news where id = any($1) and deleted_at is null"
	query2 := "select id, file_name, bucket_name, mime_type from files where id = any($1)"
	err := s.transaction(ctx, func(tx *sql.Tx) error {
		err := s.db.SelectContext(ctx, &news, query1, pq.Array(newsIDs))
		if err != nil {
			return fmt.Errorf("can't get news: %w", err)
		}
		filesIDs := make([]string, 0, len(news)*2)
		for _, newsItem := range news {
			filesIDs = append(filesIDs, newsItem.ContentID)
			if newsItem.PreviewID != "" {
				filesIDs = append(filesIDs, newsItem.PreviewID)
			}
		}
		var files []*models.File
		err = s.db.SelectContext(ctx, &files, query2, pq.Array(filesIDs))
		if err != nil {
			return fmt.Errorf("can't get files: %w", err)
		}
		for _, newsItem := range news {
			for _, file := range files {
				if newsItem.ContentID == file.ID {
					newsItem.Content = file
				}
			}
			newsItem.AuthorInfo = &models.Profile{ID: newsItem.AuthorID}
			newsItem.UserCreatedInfo = &models.Profile{ID: newsItem.UserCreatedID}
			newsItem.UserUpdatedInfo = &models.Profile{ID: newsItem.UserUpdatedID}
		}

		return nil
	})
	if err != nil {
		return nil, errorHandle(err)
	}
	return news, nil
}

func (s *Storage) GetNewsList(ctx context.Context, limit int, offset int) ([]*models.NewsAttributes, error) {
	var news []*models.NewsAttributes
	query1 := "select id, title, author, description, content_id, preview_id, state, created_at, updated_at, user_created, user_updated from news where deleted_at is null order by created_at limit $1 offset $2"
	query2 := "select id, file_name, bucket_name, mime_type from files where id = any($1)"
	err := s.transaction(ctx, func(tx *sql.Tx) error {
		err := s.db.SelectContext(ctx, &news, query1, limit, offset)
		if err != nil {
			return fmt.Errorf("can't get news: %w", err)
		}
		filesIDs := make([]string, 0, len(news)*2)
		for _, newsItem := range news {
			filesIDs = append(filesIDs, newsItem.ContentID)
			if newsItem.PreviewID != "" {
				filesIDs = append(filesIDs, newsItem.PreviewID)
			}
		}
		var files []*models.File
		err = s.db.SelectContext(ctx, &files, query2, pq.Array(filesIDs))
		if err != nil {
			return fmt.Errorf("can't get files: %w", err)
		}
		for _, newsItem := range news {
			for _, file := range files {
				if newsItem.ContentID == file.ID {
					newsItem.Content = file
				}
			}
			newsItem.AuthorInfo = &models.Profile{ID: newsItem.AuthorID}
			newsItem.UserCreatedInfo = &models.Profile{ID: newsItem.UserCreatedID}
			newsItem.UserUpdatedInfo = &models.Profile{ID: newsItem.UserUpdatedID}
		}

		return nil
	})
	if err != nil {
		return nil, errorHandle(err)
	}
	return news, nil
}

func (s *Storage) DeleteNewsByID(ctx context.Context, newsID uuid.UUID) error {
	query := "update news set deleted_at = now() where id = $1"

	_, err := s.db.ExecContext(ctx, query, newsID)
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) SetHardDeletedFilesByIDs(ctx context.Context, files []uuid.UUID) error {
	query := "update files set hard_deleted = true where id = any($1)"
	_, err := s.db.ExecContext(ctx, query, pq.Array(files))
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) GetDeletingFiles(ctx context.Context) ([]*models.File, error) {
	var filesIDs []*models.File
	query := "select id, mime_type, bucket_name, file_name from files where deleted_at is not null and hard_deleted = false"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errorHandle(err)
	}
	for rows.Next() {
		var file models.File
		err = rows.Scan(&file.ID, &file.MimeType, &file.BucketName, &file.FileName)
		if err != nil {
			return nil, errorHandle(err)
		}
		filesIDs = append(filesIDs, &file)
	}
	return filesIDs, nil
}

func (s *Storage) GetNewsFileByID(ctx context.Context, fileID uuid.UUID) (*models.File, error) {
	var file models.File
	query := "select id, mime_type, bucket_name, file_name from files where id = $1"
	err := s.db.QueryRowContext(ctx, query, fileID).Scan(&file.ID, &file.MimeType, &file.BucketName, &file.FileName)
	if err != nil {
		return nil, errorHandle(err)
	}
	return &file, nil
}
