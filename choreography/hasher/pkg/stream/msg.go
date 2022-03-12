package stream

import "go.opentelemetry.io/otel/trace"

type HashedWordMsg struct {
	SpanContext trace.SpanContext `json:"span_context"`
	HashedWord  []byte            `json:"hashed_word"`
}

type WordMsg struct {
	SpanContext trace.SpanContext `json:"span_context"`
	Word        string            `json:"word"`
}
