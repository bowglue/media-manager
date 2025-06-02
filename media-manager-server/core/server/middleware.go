package server

import (
	"mms/core/middleware"
	"net/http"
)

func (s *Server) setupMiddleware() http.Handler {
	return middleware.CreateStack(middleware.Log(s.config.Log))(s.router)
}
