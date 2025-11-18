package repository

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigrations(dsn string, migrationsPath string) error {
	log.Println("Running migrations...")

	sourceURL := "file://" + migrationsPath

	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to apply.")
			return nil
		}
		return err
	}

	log.Println("Migrations applied successfully.")
	return nil
}
