package resthandler

import (
	"bytes"
	"context"
	"github.com/go-chi/chi"
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
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"mime/multipart"
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
	testFileName      = "test_file.jpeg"
)

const (
	baseURL        = "http://localhost:8081/"
	newsListURL    = "/api/news/newsList"
	newsCreateURL  = "/api/news/writer/create"
	newsEditURL    = "/api/news/writer"
	newsURL        = "/api/news"
	newsFilesURL   = "/api/news/files/"
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

func initTestLike(t *testing.T) models.Like {
	likeID, err := uuid.NewV4()
	assert.NoError(t, err)
	testUser := initTestProfile(t, false)
	return models.Like{
		ID:        likeID.String(),
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
		Liker:     testUser,
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
		grpcC.EXPECT().GetUserByIDs(gomock.Any(), gomock.Any()).Return(&pb.UserByIDsResponse{Users: grpcUsers}, error(nil))

		q := r.URL.Query()
		q.Add("limit", "5")
		q.Add("page", "5")
		r.URL.RawQuery = q.Encode()

		h.getNewsList(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCreateNews(t *testing.T) {
	testRequestAttr := `{
	"title": "test",
	"description": "test",
	"state": "PUBLISHED"
}`
	t.Run("success", func(t *testing.T) {
		h, db, s3, _ := initTestRESTHandler(t)

		file, err := os.Open("./" + testFileName)
		assert.NoError(t, err)
		defer file.Close()

		buf := new(bytes.Buffer)
		bw := multipart.NewWriter(buf)

		pw, err := bw.CreateFormField(newsKeyAttr)
		pw.Write([]byte(testRequestAttr))

		fw1, err := bw.CreateFormFile(newsKeyFile, testFileName)
		fw1, err = bw.CreateFormFile(previewKeyFile, testFileName)
		io.Copy(fw1, file)
		bw.Close()

		r := httptest.NewRequest(http.MethodPost, newsCreateURL, buf)
		w := httptest.NewRecorder()

		r.Header.Add("Content-Type", bw.FormDataContentType())

		profile := initTestProfile(t, true)
		ctx := context.WithValue(r.Context(), profileIDKey, profile)

		newsUUID, err := uuid.NewV4()
		assert.NoError(t, err)

		s3.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(minio.UploadInfo{}, "image", error(nil))
		s3.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(minio.UploadInfo{}, "image", error(nil))
		db.EXPECT().CreateNews(gomock.Any(), gomock.Any()).Return(newsUUID, error(nil))

		h.createNews(w, r.WithContext(ctx))
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestEditNews(t *testing.T) {
	testRequestAttr := `{
	"title": "test",
	"description": "test",
	"state": "PUBLISHED"
}`
	t.Run("success", func(t *testing.T) {
		h, db, s3, _ := initTestRESTHandler(t)

		newsUUID, err := uuid.NewV4()
		assert.NoError(t, err)

		file, err := os.Open("./" + testFileName)
		assert.NoError(t, err)
		defer file.Close()

		buf := new(bytes.Buffer)
		bw := multipart.NewWriter(buf)

		pw, err := bw.CreateFormField(newsKeyAttr)
		pw.Write([]byte(testRequestAttr))

		fw1, err := bw.CreateFormFile(newsKeyFile, testFileName)
		fw1, err = bw.CreateFormFile(previewKeyFile, testFileName)
		io.Copy(fw1, file)
		bw.Close()
		r := httptest.NewRequest(http.MethodPatch, newsEditURL+"{id}", buf)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", newsUUID.String())

		r.Header.Add("Content-Type", bw.FormDataContentType())

		profile := initTestProfile(t, true)
		ctx := context.WithValue(r.Context(), profileIDKey, profile)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		s3.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(minio.UploadInfo{}, "image", error(nil))
		s3.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(minio.UploadInfo{}, "image", error(nil))
		db.EXPECT().EditNewsByID(gomock.Any(), newsUUID, gomock.Any(), gomock.Any()).Return(error(nil))

		h.editNews(w, r.WithContext(ctx))
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteNews(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, db, _, _ := initTestRESTHandler(t)

		newsUUID, err := uuid.NewV4()
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodDelete, newsFilesURL+"/{id}", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", newsUUID.String())

		profile := initTestProfile(t, true)
		ctx := context.WithValue(r.Context(), profileIDKey, profile)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		db.EXPECT().DeleteNewsByID(gomock.Any(), newsUUID).Return(error(nil))

		h.deleteNews(w, r.WithContext(ctx))
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetLikesByNews(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, db, _, grpc := initTestRESTHandler(t)

		newsUUID, err := uuid.NewV4()
		assert.NoError(t, err)

		likes := make(map[string]models.Like)
		likers := make([]*pb.Users, 0, 5)
		for i := 0; i < 5; i++ {
			like := initTestLike(t)
			likes[like.Liker.ID] = like
			likers = append(likers, &pb.Users{
				UserID:    like.Liker.ID,
				Username:  like.Liker.UserName,
				Firstname: like.Liker.FirstName,
				Lastname:  like.Liker.LastName,
				Surname:   like.Liker.SurName,
			})
		}

		r := httptest.NewRequest(http.MethodGet, newsURL+"/{id}"+"/likes", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", newsUUID.String())

		profile := initTestProfile(t, true)
		ctx := context.WithValue(r.Context(), profileIDKey, profile)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		db.EXPECT().GetLikesByNewsID(gomock.Any(), newsUUID).Return(likes, error(nil))
		grpc.EXPECT().GetUserByIDs(gomock.Any(), gomock.Any()).Return(&pb.UserByIDsResponse{Users: likers}, error(nil))

		h.getLikesByNewsID(w, r.WithContext(ctx))

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
