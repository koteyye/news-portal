package resthandler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/pkg/models"
	"github.com/koteyye/news-portal/pkg/signer"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/koteyye/news-portal/internal/news/service"

	resp "github.com/koteyye/news-portal/pkg/restresponser"
)

// RESTHandler HTTP обработчик сервиса.
type RESTHandler struct {
	service     *service.Service
	logger      *slog.Logger
	signer      signer.Signer
	corsAllowed []string
}

const (
	newsKeyFile    = "newsFile"
	newsKeyAttr    = "newsAttr"
	previewKeyFile = "previewFile"
)

// NewRESTHandler получить новый экземпляр RESTHandler.
func NewRESTHandler(service *service.Service, logger *slog.Logger, corsAllowed []string, signer signer.Signer) *RESTHandler {
	return &RESTHandler{service: service, logger: logger, corsAllowed: corsAllowed, signer: signer}
}

// InitRoutes инициализация mux.
func (h RESTHandler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: h.corsAllowed,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETED"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))
	r.Use(h.auth)

	r.Route("/api", func(r chi.Router) {
		r.Route("/news", func(r chi.Router) {
			r.Get("/newsList", h.getNewsList)
			r.Post("/create", h.createNews)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.getNewsByID)
				r.Patch("/", h.editNews)
				r.Delete("/", h.deleteNews)
				r.Route("/content", func(r chi.Router) {
					r.Get("/", h.downloadContent)
				})
				r.Route("/likes", func(r chi.Router) {
					r.Get("/{id}", h.getLikesByNewsID)
					r.Patch("/like", h.incrementLike)
					r.Patch("/dislike", h.decrementLike)
				})
				r.Patch("/like", h.incrementLike)
				r.Patch("/dislike", h.decrementLike)
				r.Route("/comment", func(r chi.Router) {
					r.Post("/", h.createComment)
					r.Patch("/{id}", h.editComment)
					r.Delete("/id", h.deleteComment)
				})
			})
		})
		r.Route("/profile", func(r chi.Router) {
			r.Get("/me", h.me)
			r.Patch("/{id}", h.editProfile)
		})
	})

	return r
}

func (h *RESTHandler) getNewsList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("page")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: errors.New("invalid limit"), ContentType: resp.CtJSON})
		return
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: errors.New("invalid page"), ContentType: resp.CtJSON})
		return
	}
	newsList, err := h.service.GetNewsList(ctx, limit, offset)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(newsList)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}
}

func (h *RESTHandler) getNewsByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	newsID := chi.URLParam(r, "id")
	newsUUID, err := uuid.FromString(newsID)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}
	news, err := h.service.GetNewsByIDs(ctx, []uuid.UUID{newsUUID})
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&news[0])
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}
}

func (h *RESTHandler) createNews(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	profile := ctx.Value(profileIDKey).(*models.Profile)

	newsFile, newsFileHeader, err := getFileFromMultipartform(w, r, newsKeyFile)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}

	previewFile, previewFileHeader, err := getFileFromMultipartform(w, r, previewKeyFile)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: err, ContentType: resp.CtJSON})
		return
	}

	attr := r.FormValue(newsKeyAttr)
	var newsAttribues models.NewsAttributes

	err = json.Unmarshal([]byte(attr), &newsAttribues)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: err, ContentType: resp.CtJSON})
		return
	}

	newsID, err := h.service.CreateNews(ctx, &newsAttribues, newsFile, newsFileHeader, previewFile, previewFileHeader, profile.ID)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: err, ContentType: resp.CtJSON})
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]uuid.UUID{"news_id": newsID})
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: err, ContentType: resp.CtJSON})
		return
	}
}

func (h *RESTHandler) editNews(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	profile := ctx.Value(profileIDKey).(*models.Profile)

	newsID := chi.URLParam(r, "id")
	newsUUID, err := uuid.FromString(newsID)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}

	newsFile, newsFileHeader, err := getFileFromMultipartform(w, r, newsKeyFile)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}

	previewFile, previewFileHeader, err := getFileFromMultipartform(w, r, previewKeyFile)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: err, ContentType: resp.CtJSON})
		return
	}

	attr := r.FormValue(newsKeyAttr)
	var newsAttribues models.NewsAttributes

	err = json.Unmarshal([]byte(attr), &newsAttribues)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: err, ContentType: resp.CtJSON})
		return
	}

	userID, err := uuid.FromString(profile.ID)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}

	err = h.service.EditNews(ctx, newsUUID, &newsAttribues, newsFile, newsFileHeader, previewFile, previewFileHeader, userID)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *RESTHandler) deleteNews(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	newsID := chi.URLParam(r, "id")
	newsUUID, err := uuid.FromString(newsID)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}

	err = h.service.DeleteNewsByID(ctx, newsUUID)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *RESTHandler) downloadContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) getLikesByNewsID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	newsID := chi.URLParam(r, "id")
	newsUUID, err := uuid.FromString(newsID)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return
	}

}

func (h *RESTHandler) incrementLike(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) decrementLike(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) createComment(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) editComment(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) deleteComment(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) me(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	profile := ctx.Value(profileIDKey).(*models.Profile)
	if profile == nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: errors.New("profile is empty"), ContentType: resp.CtJSON})
		return
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(profile)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: err, ContentType: resp.CtJSON})
		return
	}
}

func (h *RESTHandler) editProfile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
