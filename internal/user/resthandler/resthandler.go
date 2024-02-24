package resthandler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/koteyye/news-portal/internal/user/service"
	"github.com/koteyye/news-portal/pkg/signer"
)

// RESTHandler HTTP обработчик сервиса.
type RESTHandler struct {
	service *service.Service
	logger *slog.Logger
	corsAllowed []string
	signer signer.Signer
}

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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("health", func(r chi.Router) {
			r.Use(h.auth)
			r.Get("/check", h.healthCheck)
		})
		r.Route("/user", func(r chi.Router) {
			r.Post("/signup", h.signUp)
			r.Post("/signin", h.signIn)
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