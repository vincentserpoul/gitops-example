package stream

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/trace"
)

func TestJetstreamConnect(t *testing.T) {
	t.Parallel()

	nc, js, err := JetstreamConnect(ns.ClientURL(), true)
	if err != nil {
		t.Errorf("JetstreamConnect() error = %v", err)

		return
	}

	defer nc.Close()

	if nc == nil {
		t.Errorf("JetstreamConnect()  missing connection")
	}

	if js == nil {
		t.Errorf("JetstreamConnect missing jetstream")
	}
}

var (
	ns   *server.Server
	nc   *nats.Conn
	js   nats.JetStreamContext
	sub  *nats.Subscription
	trcr trace.Tracer
)

func TestMain(m *testing.M) {
	var err error

	ns, nc, js, sub, err = runNATSTestServerWithStreams()
	if err != nil {
		panic(err)
	}

	defer func() {
		sub.Unsubscribe()
		sub.Drain()
		nc.Close()
		ns.Shutdown()
	}()

	trcr = trace.NewNoopTracerProvider().Tracer("testing")

	code := m.Run()

	os.Exit(code)
}

func runNATSTestServer() *server.Server {
	rand.Seed(time.Now().UnixNano())

	min := 53000
	max := 60000
	port := rand.Intn(max-min+1) + min

	opts := &natsserver.DefaultTestOptions
	opts.JetStream = true
	opts.Port = port

	ns := natsserver.RunServer(opts)

	return ns
}

func runNATSTestServerWithStreams() (
	*server.Server,
	*nats.Conn,
	nats.JetStreamContext,
	*nats.Subscription,
	error,
) {
	ns := runNATSTestServer()

	nc, js, err := JetstreamConnect(ns.ClientURL(), true)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("JetstreamConnect()werror = %w", err)
	}

	if _, err := js.AddStream(&nats.StreamConfig{
		Name:     "wordstream",
		Subjects: []string{"WORDS.*"},
		Storage:  nats.MemoryStorage,
	}); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("AddStream error = %w", err)
	}

	if _, err := js.AddConsumer("wordstream", &nats.ConsumerConfig{
		Durable:       "wordstream-hasher",
		FilterSubject: "WORDS.submitted",
		AckPolicy:     nats.AckExplicitPolicy,
	}); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("AddConsumer error = %w", err)
	}

	// Simple Pull Consumer
	sub, err := js.PullSubscribe("WORDS.submitted", "wordstream-hasher")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("PullSubscribe: %w", err)
	}

	return ns, nc, js, sub, nil
}
