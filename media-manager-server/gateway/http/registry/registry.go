package registry

import "net/http"

type RouteHandler interface {
	RegisterRoutes(mux *http.ServeMux)
}

var registeredHandlers []RouteHandler

func Register(h RouteHandler) {
	registeredHandlers = append(registeredHandlers, h)
}

func GetHandlers() []RouteHandler {
	return registeredHandlers
}
