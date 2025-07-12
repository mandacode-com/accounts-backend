package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	grpcserver "mandacode.com/accounts/token/cmd/server/grpc"
	"mandacode.com/accounts/token/config"
	handlerv1 "mandacode.com/accounts/token/internal/handler/v1"
	tokengen "mandacode.com/accounts/token/internal/infra/token"
	token "mandacode.com/accounts/token/internal/usecase/token"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	// tokenGenerator, err := tokengen.NewTokenGenerator()
	accesTokenGen, err := tokengen.NewTokenGeneratorByStr(
		cfg.AccessPrivateKey,
		cfg.AccessTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create access token generator", zap.Error(err))
	}
	refreshTokenGen, err := tokengen.NewTokenGeneratorByStr(
		cfg.RefreshPrivateKey,
		cfg.RefreshTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create refresh token generator", zap.Error(err))
	}
	emailVerificationTokenGen, err := tokengen.NewTokenGeneratorByStr(
		cfg.EmailVerificationPrivateKey,
		cfg.EmailVerificationTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create email verification token generator", zap.Error(err))
	}

	tokenUsecase := token.NewTokenUsecase(
		accesTokenGen,
		refreshTokenGen,
		emailVerificationTokenGen,
	)

	tokenHandler, err := handlerv1.NewTokenHandler(tokenUsecase, logger)
	if err != nil {
		logger.Fatal("failed to create token handler", zap.Error(err))
	}

	// Create the gRPC server
	servingStatus := []string{
		"token.v1.TokenService",
	}
	grpcServer, err := grpcserver.NewGRPCServer(
		cfg.Port,
		logger,
		tokenHandler,
		servingStatus,
	)

	manager := server.NewServerManager([]server.Server{grpcServer})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	go func() {
		sig := <-signalChan
		logger.Info("received signal, shutting down", zap.String("signal", sig.String()))
		cancel() // Cancel the context to stop the server
	}()

	if err := manager.Run(ctx); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
