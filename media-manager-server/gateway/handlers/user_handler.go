package handlers

import (
	"encoding/json"
	"net/http"
	// "log"
)

// HandleUserRequest handles requests for user data (REST API).
func HandleUserRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Example: Fetch user details from user-service via gRPC
		// Here we're just sending a dummy response for simplicity
		userData := map[string]interface{}{
			"id":   1,
			"name": "John Doe",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}
