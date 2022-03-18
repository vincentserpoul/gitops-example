package happycat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var ErrAPIArchiver = errors.New("api archiver error")

func saveHappycatFact(
	ctx context.Context,
	timeoutArchiver time.Duration,
	baseURL string,
	fact string,
) error {
	body, err := json.Marshal(struct {
		ID   uuid.UUID
		Fact string
	}{
		ID:   uuid.New(),
		Fact: fact,
	})
	if err != nil {
		return fmt.Errorf("error mashalling text: %s - %w", fact, err)
	}

	ctxArchiver, cancel := context.WithTimeout(ctx, timeoutArchiver)
	defer cancel()

	url := baseURL + "/happycatfact"

	res, err := otelhttp.Post(ctxArchiver, url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error requesting %s: %w", url, err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			"http err, status %d instead of %d: %w",
			res.StatusCode, http.StatusCreated, ErrAPIArchiver,
		)
	}

	return nil
}
