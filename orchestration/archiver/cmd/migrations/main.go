package main

import (
	"archiver/pkg/configuration"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// configuration
	currEnv := "local"
	if e := os.Getenv("APP_ENVIRONMENT"); e != "" {
		currEnv = e
	}

	cfg, err := configuration.GetConfig(currEnv)
	if err != nil {
		if errors.Is(err, configuration.MissingBaseConfigError{}) {
			log.Fatalf("getConfig: %v", err)

			return
		}

		log.Printf("getConfig: %v", err)
	}

	m, err := migrate.New(
		"file://./sql/migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Database.Username, cfg.Database.Password,
			cfg.Database.Host, cfg.Database.Port, cfg.Database.DatabaseName, cfg.Database.SSLMode,
		),
	)
	if err != nil {
		log.Fatalf("migration: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("up: %v", err)
	}
}
