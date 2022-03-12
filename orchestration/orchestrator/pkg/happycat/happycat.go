package happycat

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
)

var ErrNoPositiveCatFact = errors.New("could not find any positive cat fact")

const MinPositiveScore = 1

func get1HappyCatFact(
	ctx context.Context,
	trcr trace.Tracer,
	timeoutCatFacts, timeoutSentium time.Duration,
	catFactsURL, sentiumURL string,
) (string, []error) {
	_, span := trcr.Start(ctx, "get1HappyCatFact")
	defer span.End()

	cfs, err := getCatFacts(ctx, trcr, timeoutCatFacts, catFactsURL)
	if err != nil {
		return "", []error{fmt.Errorf("%w", err)}
	}

	firstPositiveCatFact := make(chan string, 1)
	errorsC := make(chan error, len(cfs))

	ctxWithCancel, cancel := context.WithCancel(ctx)

	waitC := make(chan struct{}, 1)
	wg := sync.WaitGroup{}

	go func() {
		for i, cf := range cfs {
			cf := cf
			i := i

			wg.Add(1)

			go func() {
				defer wg.Done()

				sentiment, err := getSentiment(ctxWithCancel, trcr, timeoutSentium, sentiumURL, cf)
				if err != nil {
					errorsC <- fmt.Errorf("wg %d: %w", i, err)

					return
				}

				if sentiment < MinPositiveScore {
					return
				}

				select {
				case <-ctxWithCancel.Done():
				case firstPositiveCatFact <- cf:
				default:
				}
			}()
		}

		wg.Wait()

		waitC <- struct{}{}
	}()

	var errs []error

	var s string

receiveResults:
	for {
		select {
		case s = <-firstPositiveCatFact:
			break receiveResults
		case e := <-errorsC:
			errs = append(errs, e)
		case <-waitC:
			break receiveResults
		}
	}

	cancel()
	wg.Wait()

	return s, errs
}
