package database

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // necessary
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

type migrator struct {
	mig *migrate.Migrate
}

func NewMigrator(dbURL string) (*migrator, error) {
	migFS, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("iofs: %w", err)
	}

	mig, err := migrate.NewWithSourceInstance("iofs", migFS, dbURL)
	if err != nil {
		return nil, fmt.Errorf("new: %w", err)
	}

	return &migrator{mig: mig}, nil
}

func (m *migrator) Up() error {
	err := m.mig.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("up: %w", err)
	}

	return nil
}

func (m *migrator) Down() error {
	err := m.mig.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("down: %w", err)
	}

	return nil
}
