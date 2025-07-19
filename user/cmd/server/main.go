package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	httpserver "mandacode.com/accounts/user/cmd/server/http"
	"mandacode.com/accounts/user/config"

	httphandlerv1 "mandacode.com/accounts/user/internal/handler/v1/http"
	authinfra "mandacode.com/accounts/user/internal/infra/auth"
	dbinfra "mandacode.com/accounts/user/internal/infra/database"
	profileinfra "mandacode.com/accounts/user/internal/infra/profile"
	tokeninfra "mandacode.com/accounts/user/internal/infra/token"
	authrepo "mandacode.com/accounts/user/internal/repository/auth"
	coderepo "mandacode.com/accounts/user/internal/repository/code"
	dbrepo "mandacode.com/accounts/user/internal/repository/database"
	maileventrepo "mandacode.com/accounts/user/internal/repository/mailevent"
	profilerepo "mandacode.com/accounts/user/internal/repository/profile"
	tokenrepo "mandacode.com/accounts/user/internal/repository/token"
	usereventrepo "mandacode.com/accounts/user/internal/repository/userevent"
	"mandacode.com/accounts/user/internal/usecase/admin"
	manage "mandacode.com/accounts/user/internal/usecase/management"
	"mandacode.com/accounts/user/internal/usecase/signup"
	"mandacode.com/accounts/user/internal/util"
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

	// Util
	mailCodeGenerator := util.NewRandomStringGenerator(32)

	// Initialize Redis clients
	emailCodeStore := redis.NewClient(&redis.Options{
		Addr:     cfg.EmailCodeStore.Address,
		Password: cfg.EmailCodeStore.Password,
		DB:       cfg.EmailCodeStore.DB,
	})

	// Initialize Kafka writers
	userEventWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.UserEventWriter.Address...),
		Topic:                  cfg.UserEventWriter.Topic,
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}
	mailEventWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.EmailEventWriter.Address...),
		Topic:                  cfg.EmailEventWriter.Topic,
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}

	// Initialize Client
	dbClient, err := dbinfra.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create database client", zap.Error(err))
	}
	localUserClient, _, err := authinfra.NewLocalUserClient(cfg.AuthClient.Address)
	if err != nil {
		logger.Fatal("failed to create local user client", zap.Error(err))
	}
	oauthUserClient, _, err := authinfra.NewOAuthUserClient(cfg.AuthClient.Address)
	if err != nil {
		logger.Fatal("failed to create OAuth user client", zap.Error(err))
	}
	profileClient, _, err := profileinfra.NewProfileClient(cfg.ProfileClient.Address)
	if err != nil {
		logger.Fatal("failed to create profile client", zap.Error(err))
	}
	tokenClient, _, err := tokeninfra.NewTokenClient(cfg.TokenClient.Address)
	if err != nil {
		logger.Fatal("failed to create token client", zap.Error(err))
	}

	syncCodeGenerator := util.NewRandomStringGenerator(16)

	// Initialize repository
	userRepo := dbrepo.NewUserRepository(dbClient, syncCodeGenerator)
	sentEmailRepo := dbrepo.NewSentEmailRepository(dbClient)
	userEventRepo := usereventrepo.NewUserEventEmitter(userEventWriter)
	authRepo := authrepo.NewAuthRepository(localUserClient, oauthUserClient)
	profileRepo := profilerepo.NewProfileRepository(profileClient)
	mailTokenRepo := tokenrepo.NewTokenRepository(tokenClient)
	mailEventRepo := maileventrepo.NewMailEventEmitter(mailEventWriter)
	mailCodeManager := coderepo.NewCodeManager(mailCodeGenerator, cfg.EmailCodeStore.Timeout, emailCodeStore, cfg.EmailCodeStore.Prefix)

	// Initialize use cases
	adminUsecase := admin.NewAdminUsecase()
	adminManageUsecase := manage.NewAdminManageUsecase(userRepo, userEventRepo)
	selfManageUsecase := manage.NewSelfManageUsecase(userRepo, userEventRepo)
	signupUsecase := signup.NewSignupUsecase(authRepo, profileRepo, userRepo, userEventRepo)
	verifyEmailUsecase := signup.NewVerifyEmailUsecase(sentEmailRepo, authRepo, mailTokenRepo, mailEventRepo, mailCodeManager, cfg.EmailVerificationLink, cfg.MaxSentEmails, cfg.MaxSentEmailsDuration)

	// Initialize HTTP handlers
	httpUserHandler := httphandlerv1.NewUserHandler(selfManageUsecase, cfg.UserIDHeaderKey, logger)
	httpAdminHandler := httphandlerv1.NewAdminHandler(adminUsecase, adminManageUsecase)
	httpSignupHandler := httphandlerv1.NewSignupHandler(signupUsecase, verifyEmailUsecase, validator, logger)

	// Initialize HTTP server
	httpServer := httpserver.NewServer(cfg.HTTPServer.Port, logger, httpAdminHandler, httpUserHandler, httpSignupHandler)

	serverManager := server.NewServerManager([]server.Server{
		httpServer,
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
