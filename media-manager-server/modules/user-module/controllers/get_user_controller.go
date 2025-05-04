package controllers

import (
	"mms/common/types"
	"net/http"
)

func (h *UserController) getUserById(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL parameters
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userData := h.UserService.GetUserById(userID)

	response := types.Response{
		Status: http.StatusOK,
		Data:   userData,
	}

	h.SendResponse(w, response)
}
