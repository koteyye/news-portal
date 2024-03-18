package grpchandler

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
)

const IPHeader = "X-Real-IP" // IPheader заголовок запроса, содержащий IP адрес

func (g *GRPCHandler) SubnetInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var ip string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		val := md.Get(IPHeader)
		if len(val) > 0 {
			ip = val[0]
		}
	}
	if ip == "" || !g.subnet.Contains(net.ParseIP(ip)) {
		return nil, status.Errorf(codes.Unavailable, "not available on this subnet")
	}
	return handler(ctx, req)
}