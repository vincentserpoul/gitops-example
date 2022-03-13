package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"orchestrator/pkg/configuration"
	"orchestrator/pkg/happycat"
	"orchestrator/pkg/observability"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
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
		log.Printf("init provider: %v", err)

		return
	}

	defer func() {
		if err := shutdown(); err != nil {
			log.Printf("shutdown: %v", err)

			return
		}
	}()

	// initialize tracer
	trcr := otel.Tracer("orchestrator")

	// router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	version := "v1"

	r.Get(fmt.Sprintf("/%s/happycat", version), happycat.Handler(
		trcr,
		cfg.CatFact.Timeout, cfg.Sentimenter.Timeout,
		cfg.CatFact.URL, cfg.Sentimenter.URL,
	))

	// serve router
	fmt.Printf("Listening on port %d\n", cfg.Application.Port)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), r)
}