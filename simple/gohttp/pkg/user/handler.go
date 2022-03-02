package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/trace"
)

func Handler(trcr trace.Tracer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		name := getUser(r.Context(), trcr, id)
		reply := fmt.Sprintf("user %s (id %s)\n", name, id)

		if _, err := w.Write(([]byte)(reply)); err != nil {
			log.Fatalf("error not written: %v", err)
		}
	})
}
