package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"orchestrator/pkg/configuration"
	"orchestrator/pkg/happycat"
	"orchestrator/pkg/observability"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"go.uber.org/zap"
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
		"orchestration-orchestrator",
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
			sugar.Errorf("shutdown: %v", err)

			return
		}
	}()

	// router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	version := "v1"

	r.Get(fmt.Sprintf("/%s/happycat", version), happycat.Handler(
		sugar,
		cfg.CatFact.Timeout, cfg.Sentimenter.Timeout, cfg.Archiver.Timeout,
		cfg.CatFact.URL, cfg.Sentimenter.URL, cfg.Archiver.URL,
	))

	// serve router
	sugar.Infof("Listening on port %d", cfg.Application.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), r); err != nil {
		sugar.Warnf("err %v", err)
	}
}
