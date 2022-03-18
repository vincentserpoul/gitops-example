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
	url string,
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

	ctxArchiver, cancel := context.WithDeadline(ctx, time.Now().Add(timeoutArchiver))
	defer cancel()

	res, err := otelhttp.Post(ctxArchiver, url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error requesting: %w", err)
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
