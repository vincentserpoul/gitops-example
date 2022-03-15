package happycat

import (
	"archiver/pkg/internal/db"
	"fmt"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func AddRoutes(version string, router chi.Router, q db.Querier, sugar *zap.SugaredLogger) {
	router.Route(fmt.Sprintf("/%s/happycatfact", version), func(r chi.Router) {
		r.Get("/{happyCatFactID}", GetHandler(q, sugar))
		r.Get("/", ListHandler(q, sugar))
		r.Post("/", CreateHandler(q, sugar))
	})
}
