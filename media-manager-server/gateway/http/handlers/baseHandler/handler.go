package baseHandler

import (
	"encoding/json"
	"gateway/types"
	"net/http"
)

// HttpHandler represents a generic HTTP handler
type HttpHandler struct{}

// NewHttpHandler creates a new handler instance
func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

// RegisterRoutes is a placeholder method for any handler
func (h *HttpHandler) RegisterRoutes(mux *http.ServeMux) {
	// This method can be overridden by specific handlers
}

// SendResponse is a common method for sending JSON responses
func (h *HttpHandler) SendResponse(w http.ResponseWriter, response types.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If JSON encoding fails, write a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to encode response",
		})
		return
	}
}
