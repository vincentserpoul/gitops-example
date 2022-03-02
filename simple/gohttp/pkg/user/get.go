package user

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func getUser(ctx context.Context, trcr trace.Tracer, id string) string {
	_, span := trcr.Start(ctx, "getUser", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	if id == "123" {
		return "otelchi tester"
	}

	return "unknown"
}
