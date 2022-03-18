package main

import (
	"archiver/pkg/configuration"
	"archiver/pkg/happycat"
	"archiver/pkg/observability"
	"archiver/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func main() { // nolint:cyclop
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

	// logging
	var logger *zap.Logger
	if cfg.Application.LogPresetDev {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("logger sync: %v", err)
		}
	}()

	sugar := logger.Sugar()

	// observability
	shutdown, err := observability.InitProvider(
		"orchestration-archiver",
		fmt.Sprintf(
			"%s:%d",
			cfg.Observability.Collector.Host,
			cfg.Observability.Collector.Port,
		),
	)
	if err != nil {
		sugar.Errorf("init provider: %v", err)

		return
	}

	defer func() {
		if err := shutdown(); err != nil {
			sugar.Warnf("shutdown: %v", err)
		}
	}()

	// querier
	ctx := context.Background()

	dbConn, q, err := postgres.New(ctx, &cfg.Databases.Write)
	if err != nil {
		sugar.Warnf("shutdown: %v", err)

		return
	}

	defer func() {
		if err := q.Close(); err != nil {
			sugar.Warnf("querier close: %v", err)
		}

		if err := dbConn.Close(); err != nil {
			sugar.Warnf("db close: %v", err)
		}
	}()

	// router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	version := "v1"

	happycat.AddRoutes(version, r, q, sugar)

	// serve router
	sugar.Infof("Listening on port %d\n", cfg.Application.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), r); err != nil {
		sugar.Warnf("err %v", err)
	}
}
