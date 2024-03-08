package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"strconv"
)

func (s *Storage) CreateNews(ctx context.Context, newsAttr *models.NewsAttributes) (uuid.UUID, error) {
	var newsID uuid.UUID

	query1 := "insert into files (id, file_name, bucket_name, mime_type) values"
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
		var values []interface{}
		for i, file := range files {
			values = append(values, file.ID, file.FileName, file.BucketName, file.MimeType)
			numFields := 4
			n := i * numFields

			query1 = `(`
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

		err = s.db.QueryRowContext(
			ctx,
			query2,
			newsAttr.Title,
			newsAttr.Author,
			newsAttr.Description,
			newsAttr.Content.ID,
			newsAttr.Preview.ID,
			newsAttr.State,
			newsAttr.Author,
			newsAttr.Author).Scan(&newsID)
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
