package usermodule

import (
	"mms/common/database/repository"
	"mms/modules/user-module/services"
)

type Module struct {
	Service *services.UserService
}

func New(queries *repository.Queries) *Module {
	service := services.NewUserService(queries)
	return &Module{Service: service}
}
