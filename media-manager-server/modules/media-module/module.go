package mediamodule

import (
	"mms/common/database/repository"
	"mms/modules/media-module/services"
)

type Module struct {
	Service *services.MediaService
}

func New(queries *repository.Queries) *Module {
	service := services.NewMediaService(queries)
	return &Module{Service: service}
}
