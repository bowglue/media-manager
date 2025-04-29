package server

import (
	"context"
	"gateway/http/middleware"
	"gateway/http/registry"
	"net/http"
	"shared/logger"
)

type ServerConfig struct {
	GatewayPort string
	Log         logger.Logger
}

type RouteHandler interface {
	RegisterRoutes(mux *http.ServeMux)
}

type Server struct {
	server *http.Server
	router *http.ServeMux
	config ServerConfig
}

// NewServer initializes a new server instance
func NewServer(config ServerConfig) *Server {
	return &Server{
		router: http.NewServeMux(),
		config: config,
	}
}

// setupHandlers registers routes for the HTTP server
func (s *Server) setupHandlers() {
	for _, handler := range registry.GetHandlers() {
		handler.RegisterRoutes(s.router)
	}
	s.config.Log.Info("HTTP routes registered")
}

// setupMiddleware sets up the middleware stack
func (s *Server) setupMiddleware() http.Handler {
	return middleware.CreateStack(
		middleware.Log(s.config.Log),
	)(s.router)
}

// Start initializes and starts the HTTP server
func (s *Server) Start() error {
	// Start the server
	s.config.Log.Info("Starting HTTP Gateway server on port " + s.config.GatewayPort)

	// Set up route
	s.setupHandlers()

	// Create HTTP server with middleware stack
	s.server = &http.Server{
		Addr:    s.config.GatewayPort,
		Handler: s.setupMiddleware(),
	}

	return s.server.ListenAndServe() // <--- Start and return any error

	// if err := server.ListenAndServe(); err != nil {
	// 	s.config.Log.Error("Failed to start HTTP server", err)
	// }
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.config.Log.Info("Shutting down server gracefully...")

	if err := s.server.Shutdown(ctx); err != nil {
		s.config.Log.Error("Server forced to shutdown:", err)
		return err
	}

	s.config.Log.Info("Server exited properly")
	return nil
}
