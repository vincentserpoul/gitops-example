package word

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/trace"
)

func SubmissionHandler(js nats.JetStreamContext, trcr trace.Tracer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, spanW := trcr.Start(ctx, "read word from req body")

		word, err := io.ReadAll(r.Body)
		if err != nil {
			spanW.End()
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		spanW.End()

		_, spanN := trcr.Start(ctx, "submit word in nats stream")

		wm := &WordMsg{
			SpanContext: spanN.SpanContext(),
			Word:        string(word),
		}

		m, err := json.Marshal(wm)
		if err != nil {
			log.Printf("%v", err)

			spanN.End()

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if _, err := js.PublishAsync("WORDS.submitted", m); err != nil {
			log.Printf("%v", err)

			spanN.End()

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		spanN.End()

		w.WriteHeader(http.StatusCreated)
	})
}
