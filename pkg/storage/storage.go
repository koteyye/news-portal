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
	Activities
}

// Authorizarion регистрация и авторизация пользователя.
//
//go:generate mockgen -source=storage.go -destination=mock/mock.go
type Authorizarion interface {
	// CreateLogin создание пользователя
	CreateLogin(ctx context.Context, login string, hashPassword string) (uuid.UUID, error)

	// GetHashedPasswordByLogin получение захешированного пароля пользователя.
	GetHashedPasswordByLogin(ctx context.Context, login string) (uuid.UUID, string, error)
}

// Users CRUD операции над пользователем.
type Users interface {
	// CreateProfileByUserID создание профиль пользователя по UserID.
	CreateProfileByUserID(ctx context.Context, userID uuid.UUID, profile *models.Profile) error

	// GetUserListByIDs получить список пользователей по ID.
	GetUserListByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Profile, error)

	// EditUserByID редактировать профиль пользователя по ID.
	EditUserByID(ctx context.Context, profile *models.Profile) error

	// DeleteUserByIDs удалить пользователя по ID.
	DeleteUserByIDs(ctx context.Context, userIDs []uuid.UUID) error

	// SetUserRoles устанавливает роли пользователю.
	SetUserRoles(ctx context.Context, userID uuid.UUID, roles []string) error

	// EditRoles удаляет текущие роли пользователя и добавляет новые.
	EditRoles(ctx context.Context, userID uuid.UUID, roles []string) error
}

// Avatar CRUD операции над аватаром пользователя.
type Avatar interface {
	// TOBE
}

// News CRUD с новостями.
type News interface {
	// CreateNews создание новости.
	CreateNews(ctx context.Context, newsAttr *models.NewsAttributes) (uuid.UUID, error)

	// GetNewsByIDs получить список новостей по ID.
	GetNewsByIDs(ctx context.Context, newsIDs []uuid.UUID) ([]*models.NewsAttributes, error)

	// GetNewsList получить список новостей.
	GetNewsList(ctx context.Context, limit int, offset int) ([]*models.NewsAttributes, error)

	// DeleteNewsByID удалить новость по ID.
	DeleteNewsByID(ctx context.Context, newsID uuid.UUID) error

	// EditNewsByID редактирование статьи по ID.
	EditNewsByID(ctx context.Context, newsID uuid.UUID, userUpdated uuid.UUID, newsAttr *models.NewsAttributes) error

	// SetHardDeletedFilesByIDs проставляет отметку о hard-delete файлов.
	SetHardDeletedFilesByIDs(ctx context.Context, files []uuid.UUID) error

	// GetDeletingFiles получить ID помеченных на удаление файлов.
	GetDeletingFiles(ctx context.Context) ([]*models.File, error)

	// GetNewsFileByID получить файл по ID
	GetNewsFileByID(ctx context.Context, fileID uuid.UUID) (*models.File, error)
}

// Activities CRUD над новостными активностями.
type Activities interface {
	// GetLikesByNewsID получить все лайки по ID новости.
	GetLikesByNewsID(ctx context.Context, newsID uuid.UUID) (map[string]models.Like, error)

	// CreateLike создать лайк для новости по ID.
	CreateLike(ctx context.Context, newsID uuid.UUID, likerID uuid.UUID) error

	// DeleteLike удалить лайк по ID новости и ID лайкера.
	DeleteLike(ctx context.Context, newsID uuid.UUID, likerID uuid.UUID) error

	// CreateComment создать комментарий к новости
	CreateComment(ctx context.Context, newsID uuid.UUID, comment *models.Comment) (uuid.UUID, error)

	// EditComment редактировать комментарий к новости
	EditComment(ctx context.Context, comment *models.Comment) error

	// DeleteComment удалить комментарий к новости
	DeleteComment(ctx context.Context, commentID uuid.UUID) error

	// GetCommentsByNewsID получить комментарии к новости
	GetCommentsByNewsID(ctx context.Context, newsID uuid.UUID) ([]*models.Comment, error)
}

// ошибки storage
var (
	ErrDuplicate = errors.New("duplicate value")
	ErrNotFound  = errors.New("value not found")
	ErrOther     = errors.New("other storage error")
)
