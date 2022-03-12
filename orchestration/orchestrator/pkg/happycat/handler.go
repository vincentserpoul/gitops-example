package happycat

import (
	"log"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func Handler(
	trcr trace.Tracer,
	timeoutCatFacts, timeoutSentimenter time.Duration,
	catFactsURL, sentimenterURL string,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cf, errs := get1HappyCatFact(
			r.Context(),
			trcr,
			timeoutCatFacts,
			timeoutSentimenter,
			catFactsURL,
			sentimenterURL,
		)
		if cf == "" && len(errs) > 0 {
			errsS := make([]string, len(errs))
			for i, e := range errs {
				errsS[i] = e.Error()
			}

			log.Printf("catfact issue: %v", strings.Join(errsS, " - "))

			return
		}

		if _, err := w.Write(([]byte)(cf)); err != nil {
			log.Printf("error not written: %v", err)
		}
	})
}
