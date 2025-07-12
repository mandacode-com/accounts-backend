package httpserver

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"

	httphandlerv1 "mandacode.com/accounts/profile/internal/handler/v1/http"
	httpmiddleware "mandacode.com/accounts/profile/internal/middleware/http"
)

type Server struct {
	http         *http.Server
	engine       *gin.Engine
	logger       *zap.Logger
	userHandler  *httphandlerv1.UserProfileHandler
	adminHandler *httphandlerv1.AdminProfileHandler
	port         int
}

// Start implements server.Server.
func (s *Server) Start(ctx context.Context) error {
	s.engine.Use(gin.Recovery())
	s.engine.Use(httpmiddleware.ErrorHandler(s.logger))

	userGroup := s.engine.Group("/v1/user")
	s.userHandler.RegisterRoutes(userGroup)

	adminGroup := s.engine.Group("/v1/admin")
	s.adminHandler.RegisterRoutes(adminGroup)

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

func NewServer(port int, logger *zap.Logger, userHandler *httphandlerv1.UserProfileHandler, adminHandler *httphandlerv1.AdminProfileHandler) (server.Server, error) {
	engine := gin.Default()
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: engine,
	}

	return &Server{
		http:         httpServer,
		engine:       engine,
		logger:       logger,
		userHandler:  userHandler,
		adminHandler: adminHandler,
		port:         port,
	}, nil
}
