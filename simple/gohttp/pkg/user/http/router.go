package http

import (
	"gohttp/pkg/user"

	"github.com/go-chi/chi/v5"
	"github.com/induzo/httpwrapper"
	"github.com/rs/zerolog"
)

func AddRoutes(
	router chi.Router, log zerolog.Logger,
	storage *user.Storage, paramsGetter httpwrapper.NamedURLParamsGetter,
) {
	router.Route("/user", func(r chi.Router) {
		r.Get("/{userID}", httpwrapper.Wrapper(log, getHandler(storage, paramsGetter)))
		r.Post("/", httpwrapper.Wrapper(log, createHandler(storage)))
	})
}
