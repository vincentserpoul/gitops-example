package word

import "go.opentelemetry.io/otel/trace"

type WordMsg struct {
	SpanContext trace.SpanContext `json:"span_context"`
	Word        string            `json:"word"`
}
