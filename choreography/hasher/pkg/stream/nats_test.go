package stream

import (
	"context"
	"encoding/json"
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

func TestLaunchProcessWordSubmissions(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	words := []string{"1", "2", "3", "5"}

	if err := pushSomeMsgs(js, words); err != nil {
		t.Errorf("pushSomeMsgs: %v", err)

		return
	}

	if err := LaunchProcessWordSubmissions(ctx, js, trcr); err != nil {
		t.Errorf("LaunchProcessWordSubmissions() error = %v", err)
	}

	// wait for ctx cancellation thx to the timeout
	<-ctx.Done()

	streamInfo, err := js.StreamInfo("wordstream")
	if err != nil {
		t.Errorf("ConsumerInfo() error = %v", err)
	}

	if int(streamInfo.State.LastSeq) != len(words)*2 {
		t.Errorf("LaunchProcessWordSubmissions() wrong number of messages = %d", int(streamInfo.State.LastSeq))
	}
}

func pushSomeMsgs(js nats.JetStreamContext, words []string) error {
	for _, w := range words {
		wb, _ := json.Marshal(WordMsg{Word: w})
		if _, err := js.PublishAsync("WORDS.submitted", wb); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
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
