package happycat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var ErrAPISentiment = errors.New("error API sentiment")

func getSentiment(
	ctx context.Context,
	baseURL string,
	textToAnalyze string,
) (int, error) {
	body, err := json.Marshal(struct{ Text string }{Text: textToAnalyze})
	if err != nil {
		return 0, fmt.Errorf("error mashalling text: %s - %w", textToAnalyze, err)
	}

	url := baseURL

	res, err := otelhttp.Post(ctx, url, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("error requesting: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(
			"url %s status %d: %w",
			url, res.StatusCode, ErrAPISentiment,
		)
	}

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
