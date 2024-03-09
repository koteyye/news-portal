package grpchandler

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/gofrs/uuid"
	"github.com/koteyye/news-portal/internal/user/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/koteyye/news-portal/proto"
)

// GRPCHandler структура GRPC обработчика
type GRPCHandler struct {
	service *service.Service
	subnet  *net.IPNet
	pb.UserServer
}

// InitGRPCHandlers возвращает новый экземпляр GRPCHandler
func InitGRPCHandlers(service *service.Service, subnet *net.IPNet) *GRPCHandler {
	return &GRPCHandler{service: service, subnet: subnet}
}

func (g *GRPCHandler) GetUserByIDs(ctx context.Context, in *pb.UserByIDsRequest) (*pb.UserByIDsResponse, error) {
	userUUIDs := make([]uuid.UUID, 0, len(in.Userids))
	for _, userID := range in.Userids {
		userUUID, err := uuid.FromString(userID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Errorf("%s %w", userID, errors.New("can't parse userID")).Error())
		}
		userUUIDs = append(userUUIDs, userUUID)
	}
	userProfiles, err := g.service.GetUsersByIDs(ctx, userUUIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, errors.New("can't get users").Error())
	}
	userResponse := make([]*pb.Users, 0, len(userProfiles))
	for _, profile := range userProfiles {
		userResponse = append(userResponse, &pb.Users{
			UserID:    profile.ID,
			Username:  profile.UserName,
			Firstname: profile.FirstName,
			Lastname:  profile.LastName,
			Surname:   profile.SurName,
			Avatar:    profile.AvatarID.String(),
			Roles:     profile.Roles,
		})
	}
	return &pb.UserByIDsResponse{Users: userResponse}, nil
}
