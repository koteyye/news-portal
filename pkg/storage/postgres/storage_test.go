package postgres

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/internal/news/config"
	"github.com/koteyye/news-portal/pkg/random"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testDSN = "postgresql://postgres:postgres@localhost:5433/news?sslmode=disable"

const (
	editedTitle = "edited title"
)

func initTestDB(t *testing.T) (*Storage, func()) {
	ctx := context.Background()

	storage, err := NewStorage(&config.Config{DBDSN: testDSN})
	assert.NoError(t, err)
	t.Cleanup(func() {
		storage.Close()
	})
	assert.NoError(t, storage.Up(ctx))

	return storage, func() {
		assert.NoError(t, storage.Down(ctx))
	}
}

func TestStorage_News(t *testing.T) {
	storage, teardown := initTestDB(t)
	defer teardown()

	ctx := context.Background()

	testNews := random.InitTestNewsAttributes(t)

	// Создание News
	newsID, err := storage.CreateNews(ctx, testNews)
	assert.NoError(t, err)
	assert.NotNil(t, newsID)

	// Редактирование News
	testNews.Title = editedTitle
	testNews.Preview = random.InitTestFile(t)
	testNews.Content = random.InitTestFile(t)
	userUUID, err := uuid.FromString(testNews.AuthorInfo.ID)
	assert.NoError(t, err)
	err = storage.EditNewsByID(ctx, newsID, userUUID, testNews)
	assert.NoError(t, err)

	// Проверка новости
	news, err := storage.GetNewsByIDs(ctx, []uuid.UUID{newsID})
	assert.NoError(t, err)
	assert.Equal(t, news[0].Title, editedTitle)

	// Получение всех новостей
	news, err = storage.GetNewsList(ctx, 5, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(news))

	// Получение файла по ID
	fileUUID, err := uuid.FromString(testNews.Content.ID)
	assert.NoError(t, err)
	file, err := storage.GetNewsFileByID(ctx, fileUUID)
	assert.NoError(t, err)
	assert.Equal(t, file, testNews.Content)

	// Удаление новости по ID
	err = storage.DeleteNewsByID(ctx, newsID)
	assert.NoError(t, err)
	// Проверка удаления новости
	news, err = storage.GetNewsByIDs(ctx, []uuid.UUID{newsID})
	assert.NoError(t, err)
	assert.Nil(t, news)

}
