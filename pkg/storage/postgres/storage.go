package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	newsConfig "github.com/koteyye/news-portal/internal/news/config"
	userConfig "github.com/koteyye/news-portal/internal/user/config"
	newsMigration "github.com/koteyye/news-portal/db/news/migrations"
	userMigration "github.com/koteyye/news-portal/db/user/migrations"
	"github.com/koteyye/news-portal/pkg/storage"
)

var _ storage.Storage = (*Storage)(nil)

// Storage определяет структуру хранилища.
type Storage struct {
	db *sql.DB
	cfg string
}

const (
	cfgNews = "news"
	cfgUser = "user"
)

func NewNewsStorage(c *newsConfig.Config) (*Storage, error) {
	db, err := connect(c.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к бд: %w", err)
	}
	return &Storage{db: db}, nil
}

// NewStorage возвращает новый экземпляр Storage.
func NewStorage(c any) (*Storage, error) {
	var storage Storage
	var dbdsn string
	switch conf := c.(type) {
	case *newsConfig.Config:
		dbdsn = conf.DBDSN
		storage.cfg = cfgNews
	case *userConfig.Config:
		dbdsn = conf.DBDSN
		storage.cfg = cfgUser
	default:
		return nil, errors.New("не удалось получить dbdsn")
	}
	db, err := connect(dbdsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к бд: %w", err)
	}
	storage.db = db
	return &storage, nil
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
	if s.cfg == cfgNews {
		return newsMigration.Up(ctx, s.db)
	}
	return userMigration.Up(ctx, s.db)
}

func (s *Storage) Down(ctx context.Context) error {
	if s.cfg == cfgNews {
		return newsMigration.Down(ctx, s.db)
	}
	return userMigration.Down(ctx, s.db)
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
