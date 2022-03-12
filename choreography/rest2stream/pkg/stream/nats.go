package stream

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
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
