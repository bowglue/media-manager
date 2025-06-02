package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"mms/common/graphql/resolvers"
	"mms/common/logger"

	_ "github.com/mattn/go-sqlite3"
)

type ServerConfig struct {
	GatewayPort string
	Log         logger.Logger
	DBPath      string
}

type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
	config     ServerConfig
	db         *sql.DB
	resolver   *resolvers.Resolver
}

func NewServer(config ServerConfig) *Server {
	db, err := sql.Open("sqlite3", config.DBPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	return &Server{
		router:   http.NewServeMux(),
		config:   config,
		db:       db,
		resolver: &resolvers.Resolver{},
	}
}

func (s *Server) Start() error {
	s.config.Log.Info("Running database migrations...")
	if err := s.runMigration(); err != nil {
		return err
	}
	s.config.Log.Info("Database migration completed successfully")

	s.RegisterModules()
	s.setupHandlers()

	s.config.Log.Info("Starting HTTP Gateway server on port " + s.config.GatewayPort)

	s.httpServer = &http.Server{
		Addr:    s.config.GatewayPort,
		Handler: s.setupMiddleware(),
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.config.Log.Info("Shutting down server gracefully...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.config.Log.Error("Server forced to shutdown:", err)
		return err
	}
	s.config.Log.Info("Server exited properly")
	return nil
}
