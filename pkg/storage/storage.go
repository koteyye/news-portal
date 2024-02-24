package storage

import (
	"context"
	"errors"
	"io"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
)

// Storage хранилище.
type Storage interface {
	io.Closer

	Authorizarion
	Users
	Avatar
} 

// Authorization регистрация и авторизация пользователя.
type Authorizarion interface {
	// SignUp регистрация пользователя.
	SignUp(ctx context.Context, login string, hashPassword string, profile *models.Profile) (*models.Profile, error)

	// SignIn авторизация пользователя
	SignIn(ctx context.Context, login string, hashPassword string) (*models.Profile, error)
}

// Users CRUD операции над пользователем.
type Users interface {
	// GetUserListByIDs получить список пользователей по ID.
	GetUserListByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error)

	// EditUserByID редактировать профиль пользователя по ID
	EditUserByID(ctx context.Context, userID uuid.UUID) error

	// DeleteUserByID удалить пользователя по ID
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error
}

// Avatar CRUD операции над аватаром пользователя
type Avatar interface {
	// TOBE
}

// ошибки storage
var (
	ErrDuplicate = errors.New("duplicate value")
	ErrNotFound  = errors.New("value not found")
	ErrOther     = errors.New("other storage error")
)