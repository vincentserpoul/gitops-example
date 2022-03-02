package main

import (
	"errors"
	"fmt"
	"gohttp/pkg/configuration"
	"gohttp/pkg/observability"
	"gohttp/pkg/user"
	"log"
	"net/http"
	"os"

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
		"simple-gohttp",
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
	trcr := otel.Tracer("gohttp")

	// router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))
	r.HandleFunc("/user/{id:[0-9]+}", user.Handler(trcr))
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	fmt.Println("www")
	// serve router
	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), r)

	log.Printf("Exiting")
}
