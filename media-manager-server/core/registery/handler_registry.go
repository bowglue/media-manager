package registery

import (
	userController "mms/modules/user-module/controllers"
	"net/http"
)

type RouteHandler interface {
	RegisterRoutes(mux *http.ServeMux)
}

var registeredHttpHandlers []RouteHandler

func RegisterAllHandlers() {
	userHandler := userController.NewUserHandler()

	RegisterHttpHandlers(userHandler)
}

func RegisterHttpHandlers(h RouteHandler) {
	registeredHttpHandlers = append(registeredHttpHandlers, h)
}

func GetHttpHandlers() []RouteHandler {
	return registeredHttpHandlers
}
