package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	sessionredis "github.com/gin-contrib/sessions/redis"
	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	grpcserver "mandacode.com/accounts/auth/cmd/server/grpc"
	httpserver "mandacode.com/accounts/auth/cmd/server/http"
	kafkaserver "mandacode.com/accounts/auth/cmd/server/kafka"
	"mandacode.com/accounts/auth/config"
	"mandacode.com/accounts/auth/ent/authaccount"

	_ "mandacode.com/accounts/auth/ent/runtime"
	grpchandlerv1 "mandacode.com/accounts/auth/internal/handler/v1/grpc"
	"mandacode.com/accounts/auth/internal/handler/v1/http"
	kafkahandlerv1 "mandacode.com/accounts/auth/internal/handler/v1/kafka"
	dbinfra "mandacode.com/accounts/auth/internal/infra/database"
	"mandacode.com/accounts/auth/internal/infra/oauthapi"
	signupinfra "mandacode.com/accounts/auth/internal/infra/signup"
	tokeninfra "mandacode.com/accounts/auth/internal/infra/token"
	coderepo "mandacode.com/accounts/auth/internal/repository/code"
	dbrepository "mandacode.com/accounts/auth/internal/repository/database"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
	"mandacode.com/accounts/auth/internal/usecase/authuser"
	"mandacode.com/accounts/auth/internal/usecase/login"
	"mandacode.com/accounts/auth/internal/usecase/userevent"
	"mandacode.com/accounts/auth/internal/util"
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

	// Initialize Redis clients and session store
	loginCodeStore := redis.NewClient(&redis.Options{
		Addr:     cfg.LoginCodeStore.Address,
		Password: cfg.LoginCodeStore.Password,
		DB:       cfg.LoginCodeStore.DB,
	})
	sessionStore, err := sessionredis.NewStore(
		cfg.SessionStore.DB,
		"tcp",
		cfg.SessionStore.Address,
		"",
		cfg.SessionStore.Password,
		[]byte(cfg.SessionStore.HashKey),
	)
	if err != nil {
		logger.Fatal("failed to create session store", zap.Error(err))
	}

	// Initialize database and token clients
	dbClient, err := dbinfra.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create database client", zap.Error(err))
	}
	tokenClient, _, err := tokeninfra.NewTokenClient(cfg.TokenClient.Address)
	if err != nil {
		logger.Fatal("failed to create token client", zap.Error(err))
	}

	userEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.UserEventReader.Brokers,
		Topic:   cfg.UserEventReader.Topic,
		GroupID: cfg.UserEventReader.GroupID,
	})

	// Initialize OAuth APIs
	googleApi, err := oauthapi.NewGoogleAPI(cfg.GoogleOAuth.ClientID, cfg.GoogleOAuth.ClientSecret, cfg.GoogleOAuth.RedirectURL, validator)
	if err != nil {
		logger.Fatal("failed to create Google OAuth API", zap.Error(err))
	}
	naverApi, err := oauthapi.NewNaverAPI(cfg.NaverOAuth.ClientID, cfg.NaverOAuth.ClientSecret, cfg.NaverOAuth.RedirectURL, validator)
	if err != nil {
		logger.Fatal("failed to create Naver OAuth API", zap.Error(err))
	}
	kakaoApi, err := oauthapi.NewKakaoAPI(cfg.KakaoOAuth.ClientID, cfg.KakaoOAuth.ClientSecret, cfg.KakaoOAuth.RedirectURL, validator)
	if err != nil {
		logger.Fatal("failed to create Kakao OAuth API", zap.Error(err))
	}
	oauthApis := map[authaccount.Provider]oauthapi.OAuthAPI{
		authaccount.ProviderGoogle: googleApi,
		authaccount.ProviderNaver:  naverApi,
		authaccount.ProviderKakao:  kakaoApi,
	}
	singupApi, err := signupinfra.NewSignupApi(
		cfg.SignupAPI.Endpoint,
		&http.Client{
			Timeout: cfg.SignupAPI.Timeout,
		},
		validator,
	)

	// Initialize random code generators
	loginCodeGenerator := util.NewRandomGenerator(32)

	// Initialize repositories
	authAccountRepo := dbrepository.NewAuthAccountRepository(dbClient)
	tokenRepo := tokenrepo.NewTokenRepository(tokenClient)

	// Initialize code managers
	loginCodeManager := coderepo.NewCodeManager(loginCodeGenerator, cfg.LoginCodeStore.Timeout, loginCodeStore, cfg.LoginCodeStore.Prefix)

	// Initialize use cases
	localUserUsecase := authuser.NewLocalUserUsecase(authAccountRepo)
	oauthUserUsecase := authuser.NewOAuthUserUsecase(authAccountRepo, oauthApis)
	localLoginUsecase := login.NewLocalLoginUsecase(authAccountRepo, tokenRepo, loginCodeManager)
	oauthLoginUsecase := login.NewOAuthLoginUsecase(authAccountRepo, tokenRepo, loginCodeManager, singupApi, oauthApis)
	userEventUsecase := userevent.NewUserEventUsecase(authAccountRepo)

	// Initialize handlers
	localUserHandler := grpchandlerv1.NewLocalUserHandler(localUserUsecase, logger)
	oauthUserHandler := grpchandlerv1.NewOAuthUserHandler(oauthUserUsecase, logger)

	localAuthHandler, err := httphandlerv1.NewLocalAuthHandler(localLoginUsecase, logger, validator)
	if err != nil {
		logger.Fatal("failed to create local auth handler", zap.Error(err))
	}
	oauthHandler, err := httphandlerv1.NewOAuthHandler(oauthLoginUsecase, logger, validator)
	if err != nil {
		logger.Fatal("failed to create OAuth handler", zap.Error(err))
	}
	userEventHandler := kafkahandlerv1.NewUserEventHandler(userEventUsecase)

	// Initialize servers
	httpServer := httpserver.NewServer(
		cfg.HTTPServer.Port,
		logger,
		localAuthHandler,
		oauthHandler,
		cfg.SessionStore.SessionName,
		sessionStore,
	)
	kafkaServer := kafkaserver.NewKafkaServer(logger, []*kafkaserver.ReaderHandler{
		{
			Reader:  userEventReader,
			Handler: userEventHandler,
		},
	})
	grpcServer, err := grpcserver.NewGRPCServer(
		cfg.GRPCServer.Port,
		logger,
		localUserHandler,
		oauthUserHandler,
		[]string{
			"accounts.auth.v1.LocalUserService",
			"accounts.auth.v1.OAuthUserService",
		},
	)
	if err != nil {
		logger.Fatal("failed to create gRPC server", zap.Error(err))
	}

	serverManager := server.NewServerManager(
		[]server.Server{
			httpServer,
			kafkaServer,
			grpcServer,
		},
	)
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
