package happycat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func getSentiment(
	ctx context.Context,
	url string,
	textToAnalyze string,
) (int, error) {
	body, err := json.Marshal(struct{ Text string }{Text: textToAnalyze})
	if err != nil {
		return 0, fmt.Errorf("error mashalling text: %s - %w", textToAnalyze, err)
	}

	res, err := otelhttp.Post(ctx, url, "application/json", bytes.NewReader(body))
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
