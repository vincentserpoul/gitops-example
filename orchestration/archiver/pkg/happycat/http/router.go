package http

import (
	handlerhttp "archiver/pkg/handler/http"
	"archiver/pkg/internal/db"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func AddRoutes(
	version string, router chi.Router, log zerolog.Logger,
	q db.Querier, paramsGetter handlerhttp.NamedURLParamsGetter,
) {
	router.Route(fmt.Sprintf("/%s/happycatfact", version), func(r chi.Router) {
		r.Get("/{happyCatFactID}", handlerhttp.Wrapper(log, getHandler(q, paramsGetter)))
		r.Get("/", handlerhttp.Wrapper(log, listHandler(q)))
		r.Post("/", handlerhttp.Wrapper(log, createHandler(q)))
	})
}
