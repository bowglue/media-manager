package server

import (
	"mms/common/database/repository"
	usermodule "mms/modules/user-module"
)

func (s *Server) registerModules() {
	queries := repository.New(s.db)

	// Register User Module
	s.config.Log.Info("Register User Module")
	userModule := usermodule.NewUserModule(queries)
	s.router.Handle("/user", userModule.Handler())

}
