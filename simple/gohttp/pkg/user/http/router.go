package http

import (
	handlerhttp "gohttp/pkg/handler/http"
	"gohttp/pkg/user"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func AddRoutes(
	router chi.Router, log zerolog.Logger,
	storage *user.Storage, paramsGetter handlerhttp.NamedURLParamsGetter,
) {
	router.Route("/user", func(r chi.Router) {
		r.Get("/{userID}", handlerhttp.Wrapper(log, getHandler(storage, paramsGetter)))
		r.Post("/", handlerhttp.Wrapper(log, createHandler(storage)))
	})
}
