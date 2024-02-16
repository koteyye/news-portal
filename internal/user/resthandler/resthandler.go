package resthandler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/koteyye/news-portal/internal/user/service"
)

// RESTHandler HTTP обработчик сервиса.
type RESTHandler struct {
	service *service.Service
	logger *slog.Logger
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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/healthCheck", h.healthCheck)
		r.Route("/user", func(r chi.Router) {
			r.Post("/signup", h.signUp)
			r.Post("/signin", h.signIn)
			r.Patch("/pass", h.changePassword)
		})
	})
	return r
}

func (h *RESTHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *RESTHandler) signUp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) signIn(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *RESTHandler) changePassword (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}