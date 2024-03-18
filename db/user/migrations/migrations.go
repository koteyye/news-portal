package usermigration

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
	"log/slog"
	"os"
)

//go:embed *.sql
var fsys embed.FS

func lazyInit() error {
	logger := slog.NewLogLogger(slog.NewTextHandler(os.Stdout, nil), slog.LevelInfo)
	goose.SetLogger(logger)
	goose.SetBaseFS(fsys)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatalf("migrations: set dialect: %s", err)
		return err
	}
	return nil
}

// Up запускает миграцию в БД.
func Up(ctx context.Context, db *sql.DB) error {
	err := lazyInit()
	if err != nil {
		return fmt.Errorf("initializing the migrator: %w", err)
	}
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("migrations: up migrations: %w", err)
	}
	return nil
}

// Down откатывает миграцию в БД.
func Down(ctx context.Context, db *sql.DB) error {
	err := lazyInit()
	if err != nil {
		return fmt.Errorf("initializing the migrator: %w", err)
	}
	if err := goose.DownContext(ctx, db, "."); err != nil {
		return fmt.Errorf("migrations: down migrations: %w", err)
	}
	return nil
}
