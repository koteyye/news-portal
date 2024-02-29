package resthandler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/koteyye/news-portal/internal/news/service"
)

// RESTHandler HTTP обработчик сервиса.
type RESTHandler struct {
	service     *service.Service
	logger      *slog.Logger
	corsAllowed []string
}

// NewRESTHandler получить новый экземпляр RESTHandler.
func NewRESTHandler(service *service.Service, logger *slog.Logger, corsAllowed []string) *RESTHandler {
	return &RESTHandler{service: service, logger: logger, corsAllowed: corsAllowed}
}

// InitRoutes инициализация mux.
func (h RESTHandler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: h.corsAllowed,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETED"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/news", func(r chi.Router) {
			r.Get("/newslist", h.getNewsList)
			r.Post("/create", h.createNews)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.getNewsByID)
				r.Patch("/", h.editNews)
				r.Delete("/", h.deleteNews)
				r.Route("/content", func(r chi.Router) {
					r.Post("/", h.uploadContent)
					r.Get("/", h.downloadContent)
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
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) getNewsByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) createNews(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) editNews(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) deleteNews(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) uploadContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) downloadContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
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
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) editProfile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
