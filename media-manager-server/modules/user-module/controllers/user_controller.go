package controllers

import (
	httpUtil "mms/common/http-util"
	"mms/modules/user-module/services"
	"net/http"
)

type UserController struct {
	*httpUtil.HttpHandler
	*services.UserService
}

func NewUserHandler() *UserController {
	return &UserController{
		HttpHandler: httpUtil.NewHttpHandler(),
		UserService: services.NewUserService(),
	}
}

func (h *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /user/{id}", h.getUserById)
}
