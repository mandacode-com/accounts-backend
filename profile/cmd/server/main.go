package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	grpcserver "mandacode.com/accounts/profile/cmd/server/grpc"
	httpserver "mandacode.com/accounts/profile/cmd/server/http"
	kafkaserver "mandacode.com/accounts/profile/cmd/server/kafka"
	"mandacode.com/accounts/profile/config"
	grpchandlerv1 "mandacode.com/accounts/profile/internal/handler/v1/grpc"
	httphandlerv1 "mandacode.com/accounts/profile/internal/handler/v1/http"
	kafkahandlerv1 "mandacode.com/accounts/profile/internal/handler/v1/kafka"
	dbinfra "mandacode.com/accounts/profile/internal/infra/database"
	dbrepo "mandacode.com/accounts/profile/internal/repository/database"
	"mandacode.com/accounts/profile/internal/usecase/admin"
	"mandacode.com/accounts/profile/internal/usecase/system"
	"mandacode.com/accounts/profile/internal/usecase/user"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logger.Sync()

	validator := validator.New()

	cfg, err := config.LoadConfig(validator)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Initialize the database connection
	dbClient, err := dbinfra.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create database client", zap.Error(err))
	}
	// Initialize Kafka reader for user events
	userEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.UserEventReader.Brokers,
		Topic:   cfg.UserEventReader.Topic,
		GroupID: cfg.UserEventReader.GroupID,
	})

	// Initialize repositories
	profileRepo := dbrepo.NewProfileRepository(dbClient)

	// Initialize use cases
	adminProfileUsecase := admin.NewProfileUsecase(profileRepo)
	userProfileUsecase := user.NewProfileUsecase(profileRepo)
	systemProfileUsecase := system.NewProfileUsecase(profileRepo)

	// Initialize HTTP handlers
	userHandler, err := httphandlerv1.NewUserProfileHandler(userProfileUsecase, cfg.HTTPServer.UIDHeader)
	if err != nil {
		logger.Fatal("failed to create user profile handler", zap.Error(err))
	}
	adminHandler, err := httphandlerv1.NewAdminProfileHandler(adminProfileUsecase, cfg.HTTPServer.UIDHeader)
	if err != nil {
		logger.Fatal("failed to create admin profile handler", zap.Error(err))
	}

	// Initialize GRPC handlers
	grpcHandler := grpchandlerv1.NewProfileHandler(systemProfileUsecase, logger)

	// Initialize Kafka handlers
	userEventHandler := kafkahandlerv1.NewUserEventHandler(systemProfileUsecase)

	// Server initialization
	httpServer, err := httpserver.NewServer(cfg.HTTPServer.Port, logger, userHandler, adminHandler)
	if err != nil {
		logger.Fatal("failed to create HTTP server", zap.Error(err))
	}
	grpcServer, err := grpcserver.NewGRPCServer(cfg.GRPCServer.Port, logger, grpcHandler, []string{"profile.v1.ProfileService"})
	if err != nil {
		logger.Fatal("failed to create gRPC server", zap.Error(err))
	}
	kafkaServer := kafkaserver.NewKafkaServer(logger, []*kafkaserver.ReaderHandler{
		{
			Reader:  userEventReader,
			Handler: userEventHandler,
		},
	})

	serverManager := server.NewServerManager([]server.Server{
		httpServer,
		grpcServer,
		kafkaServer,
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
