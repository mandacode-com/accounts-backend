package httpserver

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	"mandacode.com/accounts/auth/internal/handler/v1/http"
	httpmiddleware "mandacode.com/accounts/auth/internal/middleware/http"
)

type Server struct {
	http             *http.Server
	engine           *gin.Engine
	logger           *zap.Logger
	localAuthHandler *httphandlerv1.LocalAuthHandler
	oauthHandler     *httphandlerv1.OAuthHandler
	port             int
	sessionName      string
	sessionStore     sessions.Store
}

// Start implements server.Server.
func (s *Server) Start(ctx context.Context) error {
	s.engine.Use(gin.Recovery())
	s.engine.Use(sessions.Sessions(s.sessionName, s.sessionStore))
	s.engine.Use(httpmiddleware.ErrorHandler(s.logger))

	localAuthGroup := s.engine.Group("/v1/auth/local")
	s.localAuthHandler.RegisterRoutes(localAuthGroup)

	oauthGroup := s.engine.Group("/v1/auth/oauth")
	s.oauthHandler.RegisterRoutes(oauthGroup)

	s.logger.Info("starting HTTP server", zap.Int("port", s.port))
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("failed to start HTTP server", zap.Error(err))
		return err
	}
	s.logger.Info("HTTP server is running", zap.Int("port", s.port))
	return nil
}

// Stop implements server.Server.
func (s *Server) Stop(ctx context.Context) error {
	if err := s.http.Shutdown(ctx); err != nil {
		s.logger.Error("failed to gracefully shutdown HTTP server", zap.Error(err))
		return err
	}
	s.logger.Info("HTTP server stopped gracefully")
	return nil
}

func NewServer(
	port int,
	logger *zap.Logger,
	localAuthHandler *httphandlerv1.LocalAuthHandler,
	oauthHandler *httphandlerv1.OAuthHandler,
	sessionName string,
	sessionStore sessions.Store,
) server.Server {
	engine := gin.Default()
	return &Server{
		http:             &http.Server{Addr: ":" + strconv.Itoa(port), Handler: engine},
		engine:           engine,
		logger:           logger,
		port:             port,
		localAuthHandler: localAuthHandler,
		oauthHandler:     oauthHandler,
		sessionName:      sessionName,
		sessionStore:     sessionStore,
	}
}
