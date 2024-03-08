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
	News
}

// Authorization регистрация и авторизация пользователя.
type Authorizarion interface {
	// CreateLogin создание пользователя
	CreateLogin(ctx context.Context, login string, hashPassword string) (uuid.UUID, error)

	// GetHashedPasswordByLogin получение захешированного пароля пользователя
	GetHashedPasswordByLogin(ctx context.Context, login string) (uuid.UUID, string, error)
}

// Users CRUD операции над пользователем.
type Users interface {
	// CreateProfileByUserID создание профиль пользователя по UserID
	CreateProfileByUserID(ctx context.Context, userID uuid.UUID, profile *models.Profile) error

	// GetUserListByIDs получить список пользователей по ID.
	GetUserListByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error)

	// EditUserByID редактировать профиль пользователя по ID
	EditUserByID(ctx context.Context, profile *models.Profile) error

	// DeleteUserByID удалить пользователя по ID
	DeleteUserByIDs(ctx context.Context, userIDs []uuid.UUID) error

	// SetUserRoles устанавливает роли пользователю
	SetUserRoles(ctx context.Context, userID uuid.UUID, roles []string) error

	// EditRole удаляет текущие роли пользователя и добавляет новые
	EditRoles(ctx context.Context, userID uuid.UUID, roles []string) error
}

// Avatar CRUD операции над аватаром пользователя
type Avatar interface {
	// TOBE
}

type News interface {
	//CreateNews создание новости
	CreateNews(ctx context.Context, newsAttr *models.NewsAttributes) (uuid.UUID, error)
}

// ошибки storage
var (
	ErrDuplicate = errors.New("duplicate value")
	ErrNotFound  = errors.New("value not found")
	ErrOther     = errors.New("other storage error")
)
