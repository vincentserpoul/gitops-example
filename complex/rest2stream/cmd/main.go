package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
)

type Config struct {
	Application struct {
		Port int
	}
	Observability struct {
		Collector struct {
			Host string
			Port int
		}
	}
	Stream struct {
		Protocol           string
		Host               string
		Port               int
		InsecureSkipVerify bool
	}
}

func main() {
	// configuration
	currEnv := "local"
	if e := os.Getenv("APP_ENVIRONMENT"); e != "" {
		currEnv = e
	}

	cfg, err := getConfig(currEnv)
	if err != nil {
		if errors.Is(err, missingBaseConfigError{}) {
			log.Printf("getConfig: %v", err)

			return
		}

		log.Printf("getConfig: %v", err)
	}

	// observability
	shutdown, err := InitProvider(
		"complex-rest2stream",
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
	trcr := otel.Tracer("gochi")

	natsURL := fmt.Sprintf(
		"%s://%s:%d",
		cfg.Stream.Protocol,
		cfg.Stream.Host,
		cfg.Stream.Port,
	)

	nc, js, err := JetstreamConnect(natsURL, cfg.Stream.InsecureSkipVerify)
	if err != nil {
		log.Printf("nats: %v", err)

		return
	}
	defer nc.Close()

	// router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))
	r.Post("/word", WordSubmissionHandler(js, trcr))
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// serve router
	fmt.Printf("Listening on port %d\n", cfg.Application.Port)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), r)
}
