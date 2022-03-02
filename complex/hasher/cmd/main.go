package main

import (
	"context"
	"errors"
	"fmt"
	"hasher/pkg/configuration"
	"hasher/pkg/observability"
	"hasher/pkg/stream"
	"log"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel"
)

func main() {
	// configuration
	currEnv := "local"
	if e := os.Getenv("APP_ENVIRONMENT"); e != "" {
		currEnv = e
	}

	cfg, err := configuration.GetConfig(currEnv)
	if err != nil {
		if errors.Is(err, configuration.MissingBaseConfigError{}) {
			log.Printf("getConfig: %v", err)

			return
		}

		log.Printf("getConfig: %v", err)
	}

	// observability
	shutdown, err := observability.InitProvider(
		"complex-hasher",
		fmt.Sprintf(
			"%s:%d",
			cfg.Observability.Collector.Host,
			cfg.Observability.Collector.Port,
		),
	)
	if err != nil {
		log.Printf("init provider: %v", err)

		return
	}

	defer func() {
		if err := shutdown(); err != nil {
			log.Printf("shutdown: %v", err)

			return
		}
	}()

	// initialize tracer
	trcr := otel.Tracer("hasher")

	natsURL := fmt.Sprintf(
		"%s://%s:%d",
		cfg.Stream.Protocol,
		cfg.Stream.Host,
		cfg.Stream.Port,
	)

	nc, js, err := stream.JetstreamConnect(natsURL, cfg.Stream.InsecureSkipVerify)
	if err != nil {
		log.Printf("nats: %v", err)

		return
	}
	defer nc.Close()

	ctx := context.Background()

	if err := stream.LaunchProcessWordSubmissions(ctx, js, trcr); err != nil {
		log.Printf("nats: %v", err)

		return
	}

	log.Printf("waiting for ctrl-c")

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx.Done()

	log.Printf("Exiting")
}
