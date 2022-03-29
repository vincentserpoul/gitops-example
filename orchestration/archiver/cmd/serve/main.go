package main

import (
	"archiver/docs"
	"archiver/pkg/configuration"
	handlerhttp "archiver/pkg/handler/http"
	happycathttp "archiver/pkg/happycat/http"
	"archiver/pkg/observability"
	"archiver/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Swagger archiver API
// @description This is a sample server db save service.

// @contact.name Vince
// @contact.url https://vincent.serpoul.com
// @contact.email v@po.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host archiver.orchestration.dev
// @BasePath /v1
func main() {
	bi := buildInfo()
	if bi == nil {
		log.Printf("failed to read build info")

		return
	}

	log.Printf("version: %v", bi)

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

	dbConn, q, err := postgres.New(ctx, &cfg.Database)
	if err != nil {
		log.Warn().Err(err).Msg("postgres")

		return
	}
	defer dbConn.Close()

	// router
	r := chi.NewRouter()
	r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	version := "v1"

	// swagger
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = cfg.Application.URL.Host
	docs.SwaggerInfo.Schemes = cfg.Application.URL.Schemes

	r.Get("/"+version+"/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("%s/%s/%s/swagger/doc.json",
				cfg.Application.URL.Schemes[0], cfg.Application.URL.Host, version,
			),
		),
	))

	happycathttp.AddRoutes(version, r, log.Logger, q, chiNamedURLParamsGetter)

	// serve router
	log.Info().Int("port", cfg.Application.Port).Msg("listening")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), r); err != nil {
		log.Warn().Err(err).Msg("")
	}
}

func chiNamedURLParamsGetter(ctx context.Context, key string) (string, *handlerhttp.ErrorResponse) {
	v := chi.URLParamFromCtx(ctx, key)
	if v == "" {
		return "", handlerhttp.MissingParamError{Name: key}.ToErrorResponse()
	}

	return v, nil
}

type BuildInfo struct {
	Revision   string
	LastCommit time.Time
	DirtyBuild bool
}

func buildInfo() *BuildInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	var bi BuildInfo

	fmt.Printf("Build Info: %v", info)

	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			bi.Revision = kv.Value
		case "vcs.time":
			bi.LastCommit, _ = time.Parse(time.RFC3339, kv.Value)
		case "vcs.modified":
			bi.DirtyBuild = kv.Value == "true"
		}
	}

	return &bi
}
