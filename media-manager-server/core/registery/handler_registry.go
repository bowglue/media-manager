package registery

import (
	userHandlers "mms/modules/user-module/handlers"
	"net/http"
)

type RouteHandler interface {
	RegisterRoutes(mux *http.ServeMux)
}

var registeredHttpHandlers []RouteHandler

func RegisterAllHandlers() {
	userHandler := userHandlers.NewUserHandler()

	RegisterHttpHandlers(userHandler)
}

func RegisterHttpHandlers(h RouteHandler) {
	registeredHttpHandlers = append(registeredHttpHandlers, h)
}

func GetHttpHandlers() []RouteHandler {
	return registeredHttpHandlers
}
