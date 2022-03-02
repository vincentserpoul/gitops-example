package stream

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"hasher/pkg/hash"
	"strings"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/trace"
)

const PublishAsyncMaxPending = 256

// 	don't forget to defer nc.Close()
func JetstreamConnect(url string, insecureSkipVerify bool) (*nats.Conn, nats.JetStreamContext, error) {
	var opts []nats.Option
	if strings.Contains(url, "tls://") {
		opts = append(opts, nats.Secure(&tls.Config{InsecureSkipVerify: insecureSkipVerify}))
	}

	// Connect to NATS
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to nats: %w", err)
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(PublishAsyncMaxPending))
	if err != nil {
		return nc, nil, fmt.Errorf("error connecting to jetstream: %w", err)
	}

	return nc, js, nil
}

func LaunchProcessWordSubmissions(ctx context.Context, js nats.JetStreamContext, trcr trace.Tracer) error {
	errC := make(chan error)

	// receive loop
	receiveSenderC, err := handleSubmittedWord(ctx, errC, js)
	if err != nil {
		return fmt.Errorf("receive loop error: %w", err)
	}

	// send loop
	sendReceiverC := handleHashedWord(ctx, errC, js)

	// process loop
	go hashProcessor(ctx, receiveSenderC, sendReceiverC)

	return nil
}

func handleSubmittedWord(ctx context.Context, errC chan error, js nats.JetStreamContext) (chan *WordMsg, error) {
	// Simple Pull Consumer
	sub, err := js.PullSubscribe("WORDS.submitted", "wordstream-hasher")
	if err != nil {
		return nil, fmt.Errorf("PullSubscribe: %w", err)
	}

	receiveSenderC := make(chan *WordMsg)

	go func(sub *nats.Subscription) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := processWordSubmission(sub, receiveSenderC)
				errC <- fmt.Errorf("unmarshaller %w", err)
			}
		}
	}(sub)

	return receiveSenderC, nil
}

func processWordSubmission(sub *nats.Subscription, receiveSenderC chan *WordMsg) error {
	msgs, err := sub.Fetch(10)
	if err != nil {
		return fmt.Errorf("fetch %w", err)
	}

	for _, msg := range msgs {
		if err := msg.Ack(); err != nil {
			return fmt.Errorf("ack %w", err)
		}

		var wh WordMsg
		if err := json.Unmarshal(msg.Data, &wh); err != nil {
			return fmt.Errorf("unmarshaller %w", err)
		}

		receiveSenderC <- &wh
	}

	return nil
}

func handleHashedWord(ctx context.Context, errC chan error, js nats.JetStreamContext) chan *HashedWordMsg {
	sendReceiverC := make(chan *HashedWordMsg)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case whm := <-sendReceiverC:
				m, err := json.Marshal(whm)
				if err != nil {
					errC <- fmt.Errorf("marshal word hash msg %w", err)

					continue
				}

				if _, err := js.PublishAsync("WORDS.hashed", m); err != nil {
					errC <- fmt.Errorf("marshal word hash msg %w", err)

					continue
				}
			}
		}
	}()

	return sendReceiverC
}

func hashProcessor(
	ctx context.Context,
	receiveSenderC chan *WordMsg,
	sendReceiverC chan *HashedWordMsg,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case wm := <-receiveSenderC:
			hw := hash.Hash(wm.Word)
			sendReceiverC <- &HashedWordMsg{
				HashedWord:  hw,
				SpanContext: wm.SpanContext,
			}
		}
	}
}
