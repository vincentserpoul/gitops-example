package happycat

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrNoPositiveCatFact = errors.New("could not find any positive cat fact")

const MinPositiveScore = 1

func get1HappyCatFact(
	ctx context.Context,
	timeoutCatFacts, timeoutSentimenter time.Duration,
	catFactsURL, sentiumURL string,
) (string, []error) {
	ctxCat, cancelCatFact := context.WithDeadline(ctx, time.Now().Add(timeoutCatFacts))
	defer cancelCatFact()

	cfs, err := getCatFacts(ctxCat, catFactsURL)
	if err != nil {
		dl, _ := ctxCat.Deadline()
		fmt.Printf("%s %s", time.Now(), dl)
		fmt.Println(err)
		return "", []error{fmt.Errorf("%w", err)}
	}

	firstPositiveCatFact := make(chan string, 1)
	errorsC := make(chan error, len(cfs))

	ctxSentiment, cancel := context.WithTimeout(ctx, timeoutSentimenter)

	waitC := make(chan struct{}, 1)
	wg := sync.WaitGroup{}

	go func() {
		for i, cf := range cfs {
			cf := cf
			i := i

			wg.Add(1)

			go func() {
				defer wg.Done()

				sentiment, err := getSentiment(ctxSentiment, sentiumURL, cf)
				if err != nil {
					errorsC <- fmt.Errorf("wg %d: %w", i, err)

					return
				}

				if sentiment < MinPositiveScore {
					return
				}

				select {
				case <-ctxSentiment.Done():
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
			fmt.Println("get a positive")
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
