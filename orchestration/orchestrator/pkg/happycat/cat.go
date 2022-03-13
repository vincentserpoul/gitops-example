package happycat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type CatFacts []string

var ErrAPICatFact = errors.New("api cat facts error")

func getCatFacts(
	ctx context.Context,
	url string,
) (CatFacts, error) {
	res, err := otelhttp.Get(ctx, url)
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
