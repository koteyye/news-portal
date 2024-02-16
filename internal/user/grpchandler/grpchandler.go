package grpchandler

import (
	"github.com/koteyye/news-portal/internal/user/service"

	pb "github.com/koteyye/news-portal/proto"
)

// GRPCHandler структура GRPC обработчика
type GRPCHandler struct {
	service *service.Service
	pb.UserServer
}

// InitGRPCHandlers возвращает новый экземпляр GRPCHandler
func InitGRPCHandlers(service *service.Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}