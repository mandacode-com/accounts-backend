package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	grpcserver "mandacode.com/accounts/user/cmd/server/grpc"
	httpserver "mandacode.com/accounts/user/cmd/server/http"
	"mandacode.com/accounts/user/config"
	grpchandlerv1 "mandacode.com/accounts/user/internal/handler/v1/grpc"
	httphandlerv1 "mandacode.com/accounts/user/internal/handler/v1/http"
	dbinfra "mandacode.com/accounts/user/internal/infra/database"
	dbrepo "mandacode.com/accounts/user/internal/repository/database"
	usereventrepo "mandacode.com/accounts/user/internal/repository/userevent"
	"mandacode.com/accounts/user/internal/usecase/admin"
	"mandacode.com/accounts/user/internal/usecase/user"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("failed to initialize logger", zap.Error(err))
	}
	defer logger.Sync()

	validator := validator.New()

	cfg, err := config.LoadConfig(validator)
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	// Initialize Client
	dbClient, err := dbinfra.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create database client", zap.Error(err))
	}
	userEventWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.UserEventWriter.Address),
		Topic:                  cfg.UserEventWriter.Topic,
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}

	// Initialize repository
	userRepo := dbrepo.NewUserRepository(dbClient)
	userEventRepo := usereventrepo.NewUserEventEmitter(userEventWriter)

	// Initialize use cases
	userUsecase := user.NewUserUsecase(userRepo, userEventRepo)
	adminUsecase := admin.NewAdminUsecase()

	httpUserHandler := httphandlerv1.NewUserHandler(userUsecase)
	httpAdminHandler := httphandlerv1.NewAdminHandler(adminUsecase, userUsecase)

	grpcUserHandler := grpchandlerv1.NewUserHandler(userUsecase)

	httpServer := httpserver.NewServer(cfg.HTTPServer.Port, logger, httpAdminHandler, httpUserHandler)
	grpcServer, err := grpcserver.NewGRPCServer(cfg.GRPCServer.Port, logger, grpcUserHandler, []string{
		"user.v1.UserService",
	})
	if err != nil {
		logger.Fatal("failed to create gRPC server", zap.Error(err))
	}

	serverManager := server.NewServerManager([]server.Server{
		httpServer,
		grpcServer,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	go func() {
		sig := <-signalChan
		logger.Info("received signal, shutting down", zap.String("signal", sig.String()))
		cancel() // Cancel the context to stop the server
	}()

	if err := serverManager.Run(ctx); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
