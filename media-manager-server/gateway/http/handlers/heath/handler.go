package handler

import (
	"gateway/http/handlers/baseHandler"
	"gateway/http/registry"
	"gateway/types"
	"net/http"
)

// HealthHandler represents a health check handler
type HealthHandler struct {
	*baseHandler.HttpHandler
}

func init() {
	registry.Register(NewHealthHandler())
}

// NewHealthHandler creates a new HealthHandler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		HttpHandler: baseHandler.NewHttpHandler(),
	}
}

// RegisterRoutes registers the health check route
func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
}

// Health is the health check endpoint
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := types.Response{
		Status:  http.StatusOK,
		Message: "Service is up and running",
	}
	h.SendResponse(w, response)
}
