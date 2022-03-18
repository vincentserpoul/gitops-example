package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sentimenter/pkg/analyze"
	"sentimenter/pkg/configuration"
	"sentimenter/pkg/observability"

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
		"orchestration-sentimenter",
		fmt.Sprintf(
			"%s:%d",
			cfg.Observability.Collector.Host,
			cfg.Observability.Collector.Port,
		),
	)
	if err != nil {
		sugar.Warnf("init provider: %v", err)

		return
	}

	defer func() {
		if err := shutdown(); err != nil {
			sugar.Warnf("shutdown: %v", err)

			return
		}
	}()

	// initialize tracer
	// trcr := otel.Tracer("sentimenter")

	// router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	analyzer, err := analyze.Analyze()
	if err != nil {
		sugar.Warnf("init sentiment analyzer: %v", err)

		return
	}

	version := "v1"

	r.Post(fmt.Sprintf("/%s", version), analyze.Handler(analyzer))

	// serve router
	fmt.Printf("Listening on port %d\n", cfg.Application.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), r); err != nil {
		sugar.Warnf("err %v", err)
	}
}
