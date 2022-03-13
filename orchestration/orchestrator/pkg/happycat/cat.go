package happycat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type CatFacts []string

var ErrAPICatFact = errors.New("api cat facts error")

func getCatFacts(
	ctx context.Context,
	trcr trace.Tracer,
	timeout time.Duration,
	url string,
) (CatFacts, error) {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   timeout,
	}

	ctx, span := trcr.Start(ctx, "req cat fact")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d: %w", res.StatusCode, ErrAPICatFact)
	}

	cfsResponse := []struct {
		Text string
	}{}
	if err := json.NewDecoder(res.Body).Decode(&cfsResponse); err != nil {
		return nil, fmt.Errorf("decode cat facts: %w", err)
	}

	cfs := make([]string, len(cfsResponse))
	for i, cf := range cfsResponse {
		cfs[i] = cf.Text
	}

	return cfs, nil
}
