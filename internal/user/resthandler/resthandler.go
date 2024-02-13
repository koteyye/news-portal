package resthandler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/koteyye/news-portal/internal/user/service"
)

// RESTHandler HTTP обработчик сервиса
type RESTHandler struct {
	service *service.Service
	logger *slog.Logger
	corsAllowed []string
}

func NewRESTHandler(service *service.Service, logger *slog.Logger, corsAllowed []string) *RESTHandler {
	return &RESTHandler{service: service, logger: logger, corsAllowed: corsAllowed}
}

func (h RESTHandler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: h.corsAllowed,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETED"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/healthCheck", h.healthCheck)
	})
	return r
}

func (h *RESTHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}