package server

import (
	"mms/common/database/repository"
	"mms/common/graphql/resolvers"
	mediamodule "mms/modules/media-module"
	usermodule "mms/modules/user-module"
)

func (s *Server) RegisterModules() {
	queries := repository.New(s.db)
	userModule := usermodule.New(queries)
	mediaModule := mediamodule.New(queries)

	s.resolver = &resolvers.Resolver{
		UserService:  userModule.Service,
		MediaService: mediaModule.Service,
	}
}
