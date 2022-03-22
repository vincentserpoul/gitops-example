package main

import (
	"archiver/pkg/configuration"
	happycathttp "archiver/pkg/happycat/http"
	"archiver/pkg/observability"
	"archiver/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

// nolint: cyclop
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
	if cfg.Application.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

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
		log.Err(err).Msg("init provider")

		return
	}

	defer func() {
		if err := shutdown(); err != nil {
			log.Warn().Err(err).Msg("shutdown")
		}
	}()

	// querier
	ctx := context.Background()

	dbConn, q, err := postgres.New(ctx, &cfg.Databases.Write)
	if err != nil {
		log.Warn().Err(err).Msg("postgres")

		return
	}

	defer func() {
		if err := q.Close(); err != nil {
			log.Warn().Err(err).Msg("querier close")
		}

		if err := dbConn.Close(); err != nil {
			log.Warn().Err(err).Msg("db close")
		}
	}()

	// router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	version := "v1"

	happycathttp.AddRoutes(version, r, log.Logger, q, chiNamedURLParamsGetter)

	// serve router
	log.Info().Int("port", cfg.Application.Port).Msg("listening")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), r); err != nil {
		log.Warn().Err(err).Msg("")
	}
}

func chiNamedURLParamsGetter(ctx context.Context, key string) string {
	return chi.URLParamFromCtx(ctx, key)
}
