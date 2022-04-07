package main

import (
	"context"
	"errors"
	"fmt"
	"gohttp/docs"
	"gohttp/pkg/configuration"
	handlerhttp "gohttp/pkg/handler/http"
	"gohttp/pkg/observability"
	"gohttp/pkg/user"
	userhttp "gohttp/pkg/user/http"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Swagger gohttp API
// @description This is a sample server db save service.

// @contact.name Vince
// @contact.url https://vincent.serpoul.com
// @contact.email v@po.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	// context
	ctx := context.Background()

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
		ctx,
		"simple-gohttp",
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

	// storage
	storage := user.NewStorage()

	// router
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	version := "v1"

	// swagger
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = version
	docs.SwaggerInfo.Host = cfg.Application.URL.Host
	docs.SwaggerInfo.Schemes = cfg.Application.URL.Schemes

	r.Get("/"+version+"swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("%s/%s/%s/swagger/doc.json",
				cfg.Application.URL.Schemes[0], cfg.Application.URL.Host, version,
			),
		),
	))

	r.Route("/"+version, func(r chi.Router) {
		r.Use(otelchi.Middleware("chi", otelchi.WithChiRoutes(r)))
		userhttp.AddRoutes(r, log.Logger, storage, chiNamedURLParamsGetter)
	})

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
