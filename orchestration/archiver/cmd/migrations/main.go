package main

import (
	"archiver/pkg/configuration"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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
			log.Printf("getConfig: %v", err)

			return
		}

		log.Printf("getConfig: %v", err)
	}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Databases.Write.Username, cfg.Databases.Write.Password,
			cfg.Databases.Write.Host, cfg.Databases.Write.Port, cfg.Databases.Write.DatabaseName, cfg.Databases.Write.SSLMode,
		),
	)
	if err != nil {
		log.Printf("db connection: %v", err)

		return
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("db driver: %v", err)

		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./sql/migrations",
		"postgres", driver)
	if err != nil {
		log.Printf("migration: %v", err)

		return
	}

	if err := m.Up(); err != nil {
		log.Printf("up: %v", err)

		return
	}
}
