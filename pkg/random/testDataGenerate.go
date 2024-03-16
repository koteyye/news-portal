package random

import (
	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	testSecret        = "supersecretkey"
	testServerAddress = "localhost:8081"
	testBucket        = "testBucket"
	testState         = "PUBLISHED"
	testFileName      = "test_file.jpeg"
	writerRole        = "writer"
)

func InitTestNewsAttributes(t *testing.T) *models.NewsAttributes {
	newsID, err := uuid.NewV4()
	assert.NoError(t, err)
	testUser := InitTestProfile(t, false)

	return &models.NewsAttributes{
		ID:              newsID.String(),
		Title:           RandSeq(10),
		Description:     RandSeq(10),
		AuthorInfo:      testUser,
		Content:         InitTestFile(t),
		Preview:         InitTestFile(t),
		State:           testState,
		CreatedAt:       time.Now().String(),
		UpdatedAt:       time.Now().String(),
		UserCreatedInfo: testUser,
		UserUpdatedInfo: testUser,
	}
}

func InitTestLike(t *testing.T) models.Like {
	likeID, err := uuid.NewV4()
	assert.NoError(t, err)
	testUser := InitTestProfile(t, false)
	return models.Like{
		ID:        likeID.String(),
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
		Liker:     testUser,
	}
}

func InitTestFile(t *testing.T) *models.File {
	fileID, err := uuid.NewV4()
	assert.NoError(t, err)
	return &models.File{
		ID:         fileID.String(),
		BucketName: testBucket,
		FileName:   RandSeq(5) + ".txt",
		MimeType:   "test",
	}
}

func InitTestProfile(t *testing.T, writer bool) *models.Profile {
	userID, err := uuid.NewV4()
	assert.NoError(t, err)

	roles := []string{"reader"}
	if writer {
		roles = append(roles, writerRole)
	}

	return &models.Profile{
		ID:        userID.String(),
		UserName:  RandSeq(10),
		FirstName: RandSeq(10),
		LastName:  RandSeq(10),
		SurName:   RandSeq(10),
		Roles:     roles,
	}
}
