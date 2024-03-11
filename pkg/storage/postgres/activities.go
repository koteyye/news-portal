package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
)

func (s *Storage) GetLikesByNewsID(ctx context.Context, newsID uuid.UUID) ([]*models.Like, error) {
	var likes []*models.Like
	query := "select id, liker, created_at, updated_at from likes where is_active = true and news_id = $1"

	rows, err := s.db.QueryContext(ctx, query, newsID)
	if err != nil {
		return nil, errorHandle(err)
	}
	for rows.Next() {
		var like models.Like
		var liker models.Profile
		err = rows.Scan(&like.ID, &liker.ID, &like.CreatedAt, &like.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("can't scan like: %w", err)
		}
		like.Liker = &liker
		likes = append(likes, &like)
	}
	return likes, nil
}

func (s *Storage) CreateLike(ctx context.Context, newsID uuid.UUID, likerID uuid.UUID) error {
	query1 := "select id from likes where liker = $1 and news_id = $2 and is_active = true"
	query2 := "insert into likes (liker, news_id) values ($1, $2)"

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		var likeID uuid.UUID
		err := s.db.QueryRowContext(ctx, query1, likerID, newsID).Scan(&likeID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("can't get like: %w", err)
			}
			err = nil
		}
		if likeID != uuid.Nil {
			return errors.New("duplicate like")
		}

		_, err = s.db.ExecContext(ctx, query2, likerID, newsID)

		return nil
	})

	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) DeleteLike(ctx context.Context, newsID uuid.UUID, likerID uuid.UUID) error {
	query := "update likes set is_active = false where liker = $1 and news_id = $2"
	_, err := s.db.ExecContext(ctx, query, likerID, newsID)
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) CreateComment(ctx context.Context, newsID uuid.UUID, comment *models.Comment) (uuid.UUID, error) {
	var commentID uuid.UUID
	query := "insert into comments (news_id, author, content) values ($1, $2, $3) returning id"
	err := s.db.QueryRowContext(ctx, query, newsID, comment.Author.ID, comment.TextComment).Scan(&commentID)
	if err != nil {
		return uuid.Nil, errorHandle(err)
	}
	return commentID, nil
}

func (s *Storage) EditComment(ctx context.Context, comment *models.Comment) error {
	query := "update comments set content = $1, updated_at = now() where id = $2"
	_, err := s.db.ExecContext(ctx, query, comment.TextComment, comment.ID)
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	query := "update comments set deleted_at = now() where id = $1"
	_, err := s.db.ExecContext(ctx, query, commentID)
	if err != nil {
		return errorHandle(err)
	}
	return nil
}

func (s *Storage) GetCommentsByNewsID(ctx context.Context, newsID uuid.UUID) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := "select id, author, created_at, updated_at, content from comments where deleted_at is null"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errorHandle(err)
	}
	for rows.Next() {
		var author *models.Profile
		var comment *models.Comment
		err = rows.Scan(&comment.ID, &author.ID, &comment.CreatedAt, &comment.UpdatedAt, &comment.TextComment)
		if err != nil {
			return nil, errorHandle(err)
		}
		comment.Author = author
		comments = append(comments, comment)
	}
	return comments, nil
}
