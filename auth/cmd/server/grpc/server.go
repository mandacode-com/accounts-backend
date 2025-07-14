package grpcserver

import (
	"context"
	"net"
	"strconv"

	authv1 "github.com/mandacode-com/accounts-proto/go/auth/v1"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server           *grpc.Server
	localUserHandler authv1.LocalUserServiceServer
	oauthUserHandler authv1.OAuthUserServiceServer
	logger           *zap.Logger
	port             int
}

func NewGRPCServer(
	port int,
	logger *zap.Logger,
	localUserHandler authv1.LocalUserServiceServer,
	oauthUserHandler authv1.OAuthUserServiceServer,
	servingServices []string,
) (server.Server, error) {
	server := grpc.NewServer()

	// Register health check service
	healthServer := health.NewServer()
	for _, service := range servingServices {
		healthServer.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Register reflection service on gRPC server
	reflection.Register(server)

	// Register the token handler
	authv1.RegisterLocalUserServiceServer(server, localUserHandler)
	authv1.RegisterOAuthUserServiceServer(server, oauthUserHandler)

	return &GRPCServer{
		server:           server,
		localUserHandler: localUserHandler,
		oauthUserHandler: oauthUserHandler,
		logger:           logger,
		port:             port,
	}, nil
}

func (g *GRPCServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(g.port))
	if err != nil {
		g.logger.Error("failed to listen on port", zap.Int("port", g.port), zap.Error(err))
		return err
	}

	g.logger.Info("gRPC server is running", zap.String("address", ":"+strconv.Itoa(g.port)))
	return g.server.Serve(lis)
}

func (g *GRPCServer) Stop(ctx context.Context) error {
	g.logger.Info("stopping gRPC server")
	g.server.GracefulStop()
	g.logger.Info("gRPC server stopped gracefully")
	return nil
}
