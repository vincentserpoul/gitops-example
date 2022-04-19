package db

import (
	"archiver/pkg/configuration"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var q *Queries

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	port, err := strconv.Atoi(resource.GetPort("5432/tcp"))
	if err != nil {
		log.Fatalf("Could not get port: %s", err)
	}

	log.Println("Connecting to database")

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	dbConf := &configuration.Database{
		Host:         resource.GetBoundIP("5432/tcp"),
		Port:         port,
		Username:     "user_name",
		Password:     "secret",
		DatabaseName: "dbname",
		SSLMode:      "disable",
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		var err error
		_, q, err = NewQ(context.Background(), dbConf)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// migrate
	mig, errM := migrate.New(
		"file://../../../sql/migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			dbConf.Username, dbConf.Password,
			dbConf.Host, dbConf.Port, dbConf.DatabaseName, dbConf.SSLMode,
		),
	)
	if errM != nil {
		log.Fatalf("migration: %v", errM)
	}

	if err := mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("up: %v", err)
	}

	// Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func ConnectStringFromConfig(dbConf *configuration.Database) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.DatabaseName, dbConf.SSLMode,
	)
}

func NewQ(ctx context.Context, dbConf *configuration.Database) (*pgxpool.Pool, *Queries, error) {
	psqlInfo := ConnectStringFromConfig(dbConf)

	dbConn, err := pgxpool.Connect(ctx, psqlInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("opening conn: %w", err)
	}

	if err := dbConn.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	return dbConn, New(dbConn), nil
}
