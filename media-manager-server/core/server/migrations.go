package server

import (
	"fmt"

	"mms/common/database"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func (s *Server) runMigration() error {

	sourceDriver, err := iofs.New(database.Migrations, "migrations")
	if err != nil {
		return err
	}

	databaseURL := fmt.Sprintf("sqlite3://%s?_foreign_keys=on", s.config.DBPath)
	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
