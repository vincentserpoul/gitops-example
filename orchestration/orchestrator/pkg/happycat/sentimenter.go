package happycat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func getSentiment(
	ctx context.Context,
	trcr trace.Tracer,
	timeout time.Duration,
	url string,
	textToAnalyze string,
) (int, error) {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   timeout,
	}

	ctx, span := trcr.Start(ctx, "req sentiment")
	defer span.End()

	b, err := json.Marshal(struct{ Text string }{Text: textToAnalyze})
	if err != nil {
		return 0, fmt.Errorf("error mashalling text: %s - %w", textToAnalyze, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error requesting: %w", err)
	}
	defer res.Body.Close()

	scoreS, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("read body: %w", err)
	}

	score, err := strconv.Atoi(string(scoreS))
	if err != nil {
		return 0, fmt.Errorf("convert %s: %w", string(scoreS), err)
	}

	return score, nil
}
