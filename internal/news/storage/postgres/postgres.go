package postgres

import (
	"context"
	"database/sql"
	"fmt"

	migration "github.com/koteyye/news-portal/db/news/migrations"
	"github.com/koteyye/news-portal/internal/news/config"
	"github.com/koteyye/news-portal/internal/news/storage"
)

var _ storage.Storage = (*Storage)(nil)

// Storage определяет структуру хранилища.
type Storage struct {
	db *sql.DB
}

// NewStorage возвращает новй экземпляр Storage.
func NewStorage(c *config.Config) (*Storage, error) {
	db, err := connect(c.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к бд: %w", err)
	}
	return &Storage{db: db}, nil
}


func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать подключение к бд: %w", err)
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("не удалось пингануть бд: %w", err)
	}

	return db, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Up(ctx context.Context) error {
	return migration.Up(ctx, s.db)
}

func (s *Storage) Down(ctx context.Context) error {
	return migration.Down(ctx, s.db)
}

func (s *Storage) transaction(
    ctx context.Context,
    fn func(*sql.Tx) error,
) error {
    tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback()
    err = fn(tx)
    if err != nil {
        return fmt.Errorf("transaction: %w", err)
    }
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }
    return nil
}

