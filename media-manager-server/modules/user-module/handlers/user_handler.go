package handlers

import (
	httpUtil "mms/common/http-util"
	"mms/modules/user-module/services"
	"net/http"
)

type UserHandler struct {
	*httpUtil.HttpHandler
	*services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		HttpHandler: httpUtil.NewHttpHandler(),
		UserService: services.NewUserService(),
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /user/{id}", h.getUserById)
}
