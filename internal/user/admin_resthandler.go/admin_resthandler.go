package adminresthandlergo

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/koteyye/news-portal/internal/user/service"
	"github.com/koteyye/news-portal/pkg/models"

	resp "github.com/koteyye/news-portal/pkg/restresponser"
)

const defaultTimeout = 10 * time.Second

// AdminRESTHandler HTTP обработчик админки сервиса.
type AdminRESTHandler struct {
	service *service.Service
	logger *slog.Logger
	subnet *net.IPNet
}

func NewAdminRESTHandler(service *service.Service, logger *slog.Logger, subnet *net.IPNet) *AdminRESTHandler {
	return &AdminRESTHandler{service: service, logger: logger}
}

func (h AdminRESTHandler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/", h.createUser)
			r.Patch("/", h.editUserByID)
			r.Delete("/", h.deleteUserByIDs)
		})
	})
	return r
}

func (h *AdminRESTHandler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()
	
	input, err := models.ParseUserData(r.Body)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return 
	}
	userID, err := h.service.CreateUser(ctx, input)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return 
	}
	json.NewEncoder(w).Encode(userID)
}

func (h *AdminRESTHandler) editUserByID(w http.ResponseWriter, r *http.Request) {

}

func (h *AdminRESTHandler) deleteUserByIDs(w http.ResponseWriter, r *http.Request) {

}