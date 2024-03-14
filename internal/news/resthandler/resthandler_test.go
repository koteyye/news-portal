package resthandler

import (
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/koteyye/news-portal/internal/news/service"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/koteyye/news-portal/pkg/random"
	"github.com/koteyye/news-portal/pkg/s3"
	mock_s3 "github.com/koteyye/news-portal/pkg/s3/mock"
	"github.com/koteyye/news-portal/pkg/signer"
	mock_storage "github.com/koteyye/news-portal/pkg/storage/mock"
	pb "github.com/koteyye/news-portal/proto"
	mock_proto "github.com/koteyye/news-portal/proto/mocks"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var testCors = []string{"localhost:8080"}

const (
	testSecret        = "supersecretkey"
	testServerAddress = "localhost:8081"
	testBucket        = "testBucket"
	testState         = "PUBLISHED"
)

const (
	baseURL        = "http://localhost:8081"
	newsListURL    = "/api/news/newsList"
	newsCreateURL  = "/api/news/writer/create"
	newsEditURL    = "/api/news/writer/"
	newsURL        = "/api/news"
	newsFilesURL   = "/api/news/files"
	newsProfileURL = "/api/profile"
)

func initTestNewsAttributes(t *testing.T) *models.NewsAttributes {
	newsID, err := uuid.NewV4()
	assert.NoError(t, err)
	testUser := initTestProfile(t, false)

	return &models.NewsAttributes{
		ID:          newsID.String(),
		Title:       random.RandSeq(10),
		Description: random.RandSeq(10),
		Author:      testUser,
		Content:     initTestFile(t),
		Preview:     initTestFile(t),
		State:       testState,
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
		UserCreated: testUser,
		UserUpdated: testUser,
	}
}

func initTestFile(t *testing.T) *models.File {
	fileID, err := uuid.NewV4()
	assert.NoError(t, err)
	return &models.File{
		ID:         fileID.String(),
		BucketName: testBucket,
		FileName:   random.RandSeq(5) + ".txt",
	}
}

func initTestProfile(t *testing.T, writer bool) *models.Profile {
	userID, err := uuid.NewV4()
	assert.NoError(t, err)

	roles := []string{"reader"}
	if writer {
		roles = append(roles, writerRole)
	}

	return &models.Profile{
		ID:        userID.String(),
		UserName:  random.RandSeq(10),
		FirstName: random.RandSeq(10),
		LastName:  random.RandSeq(10),
		SurName:   random.RandSeq(10),
		Roles:     roles,
	}
}

//r.Route("/api", func(r chi.Router) {
//	r.Route("/news", func(r chi.Router) {
//		r.Route("/writer", func(r chi.Router) {
//			r.Use(h.checkWriter)
//			r.Post("/create", h.createNews)
//			r.Patch("/{id}", h.editNews)
//			r.Delete("/{id}", h.deleteNews)
//		})
//		r.Get("/newsList", h.getNewsList)
//		r.Route("/{id}", func(r chi.Router) {
//			r.Get("/", h.getNewsByID)
//			r.Route("/likes", func(r chi.Router) {
//				r.Get("/{id}", h.getLikesByNewsID)
//				r.Patch("/like", h.incrementLike)
//				r.Patch("/dislike", h.decrementLike)
//			})
//			r.Route("/comment", func(r chi.Router) {
//				r.Post("/{newsID}", h.createComment)
//				r.Patch("/", h.editComment)
//				r.Delete("/{id}", h.deleteComment)
//				r.Get("/{newsID}", h.getComments)
//			})
//		})
//		r.Route("/files", func(r chi.Router) {
//			r.Get("/{id}", h.downloadContent)
//		})
//	})
//	r.Route("/profile", func(r chi.Router) {
//		r.Get("/me", h.me)
//	})
//})

func TestNewRESTHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testService := service.Service{}
		testSiger := signer.New([]byte(testSecret))
		handler := NewRESTHandler(&testService, &slog.Logger{}, testCors, testSiger)

		assert.Equal(t, &RESTHandler{
			service:     &testService,
			logger:      &slog.Logger{},
			signer:      testSiger,
			corsAllowed: testCors,
		}, handler)
	})
}

func initTestRESTHandler(t *testing.T) (*RESTHandler, *mock_storage.MockStorage, *mock_s3.MockS3, *mock_proto.MockUserClient) {
	c := gomock.NewController(t)
	defer c.Finish()

	db := mock_storage.NewMockStorage(c)
	s3repo := mock_s3.NewMockS3(c)
	signer := signer.New([]byte(testSecret))

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	grpcClient := mock_proto.NewMockUserClient(c)

	service := service.NewService(db, &s3.Handler{S3: s3repo}, logger, grpcClient, testServerAddress)

	restHandler := NewRESTHandler(service, logger, testCors, signer)
	return restHandler, db, s3repo, grpcClient
}

func TestGetNewsList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, db, _, grpcC := initTestRESTHandler(t)
		r := httptest.NewRequest(http.MethodGet, newsListURL, nil)
		w := httptest.NewRecorder()

		grpcUsersIDs := make([]string, 0, 5)
		grpcUsers := make([]*pb.Users, 0, 5)
		newsList := make([]*models.NewsAttributes, 0, 5)
		for _, news := range newsList {
			news = initTestNewsAttributes(t)
			grpcUsersIDs = append(grpcUsersIDs, news.Author.ID)
			grpcUsers = append(grpcUsers, &pb.Users{
				UserID:    news.Author.ID,
				Username:  news.Author.UserName,
				Firstname: news.Author.UserName,
				Lastname:  news.Author.LastName,
				Surname:   news.Author.SurName,
				Roles:     news.Author.Roles,
			})
		}

		db.EXPECT().GetNewsList(gomock.Any(), gomock.Any(), gomock.Any()).Return(newsList, error(nil))
		grpcC.EXPECT().GetUserByIDs(gomock.Any(), &pb.UserByIDsRequest{Userids: grpcUsersIDs}).Return(&pb.UserByIDsResponse{Users: grpcUsers}, error(nil))

		q := r.URL.Query()
		q.Add("limit", "5")
		q.Add("page", "5")
		r.URL.RawQuery = q.Encode()

		h.getNewsList(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
