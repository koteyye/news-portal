package adminresthandler

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

// NewAdminRESTHandler возвращает новый экземпляр AdminRESTHandler
func NewAdminRESTHandler(service *service.Service, logger *slog.Logger, subnet *net.IPNet) *AdminRESTHandler {
	return &AdminRESTHandler{service: service, logger: logger, subnet: subnet}
}

func (h AdminRESTHandler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(h.checkSubnet)
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"userID": userID.String()})
}

func (h *AdminRESTHandler) editUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	input, err := models.ParseProfile(r.Body)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return 
	}

	err = h.service.EditUser(ctx, input)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return 
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AdminRESTHandler) deleteUserByIDs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
	defer cancel()

	var userIds []string
	err := json.NewDecoder(r.Body).Decode(&userIds)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return 
	}

	err = h.service.DeleteUsersByIDs(ctx, userIds)
	if err != nil {
		resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusBadRequest, Err: err, ContentType: resp.CtJSON})
		return 
	}
	w.WriteHeader(http.StatusOK)
}