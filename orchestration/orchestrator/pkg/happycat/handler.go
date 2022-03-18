package happycat

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

func Handler(
	sugar *zap.SugaredLogger,
	timeoutCatFacts, timeoutSentimenter, timeoutArchiver time.Duration,
	catFactsURL, sentimenterURL, archiverURL string,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cf, errs := get1HappyCatFact(
			r.Context(),
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

			w.WriteHeader(http.StatusInternalServerError)
			sugar.Errorf("catfact issue: %v", strings.Join(errsS, " - "))

			return
		}

		if err := saveHappycatFact(r.Context(), timeoutArchiver, archiverURL, cf); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			sugar.Errorf("archiver issue: %v", err)

			return
		}

		if _, err := w.Write(([]byte)(cf)); err != nil {
			sugar.Errorf("error not written: %v", err)
		}
	})
}
