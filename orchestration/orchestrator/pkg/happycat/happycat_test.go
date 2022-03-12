package happycat

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func Test_get1HappyCatFact(t *testing.T) {
	t.Parallel()

	type args struct {
		timeoutCatFacts    time.Duration
		timeoutSentimenter time.Duration
	}

	tests := []struct {
		name                        string
		catFactServerResponse       string
		sentimenterServerResponse   string
		timeoutSentimenterFrequency int
		args                        args
		want1                       bool
		errsLen                     int
	}{
		{
			name: "happy path",
			catFactServerResponse: `
			[
				{
					"text": "Cats make about 100 different sounds. Dogs make only about 10."
				},
				{
					"text": "Cats are always cute."
				},
				{
					"text": "Cats are serious."
				},
				{
					"text": "Domestic cats spend about 70 percent of the day sleeping and 15 percent of the day grooming."
				}
			]`,
			sentimenterServerResponse: `6`,
			want1:                     true,
			errsLen:                   0,
		},
		{
			name: "timeout",
			catFactServerResponse: `
			[
				{
					"text": "Cats make about 100 different sounds. Dogs make only about 10."
				},
				{
					"text": "Cats are always cute."
				},
				{
					"text": "Cats are serious."
				},
				{
					"text": "Domestic cats spend about 70 percent of the day sleeping and 15 percent of the day grooming."
				}
			]`,
			sentimenterServerResponse: `6`,
			args: args{
				timeoutCatFacts: 10 * time.Millisecond,
			},
			want1:   false,
			errsLen: 1,
		},
		{
			name: "timeout",
			catFactServerResponse: `
			[
				{
					"text": "Cats make about 100 different sounds. Dogs make only about 10."
				},
				{
					"text": "Cats are always cute."
				},
				{
					"text": "Cats are serious."
				},
				{
					"text": "Domestic cats spend about 70 percent of the day sleeping and 15 percent of the day grooming."
				}
			]`,
			sentimenterServerResponse: `6`,
			args: args{
				timeoutSentimenter: 10 * time.Millisecond,
			},
			want1:   false,
			errsLen: 4,
		},
		{
			name: "1 timeout in neutral answers",
			catFactServerResponse: `
			[
				{
					"text": "Cats make about 100 different sounds. Dogs make only about 10."
				},
				{
					"text": "Cats are always cute."
				},
				{
					"text": "Cats are serious."
				},
				{
					"text": "Domestic cats spend about 70 percent of the day sleeping and 15 percent of the day grooming."
				}
			]`,
			sentimenterServerResponse: `0`,
			args: args{
				timeoutSentimenter: 10 * time.Millisecond,
			},
			timeoutSentimenterFrequency: 4,
			want1:                       false,
			errsLen:                     4,
		},
		{
			name: "1 timeout in positive answers",
			catFactServerResponse: `
			[
				{
					"text": "Cats make about 100 different sounds. Dogs make only about 10."
				},
				{
					"text": "Cats are always cute."
				},
				{
					"text": "Cats are serious."
				},
				{
					"text": "Domestic cats spend about 70 percent of the day sleeping and 15 percent of the day grooming."
				}
			]`,
			sentimenterServerResponse: `1`,
			args: args{
				timeoutSentimenter: 10 * time.Millisecond,
			},
			timeoutSentimenterFrequency: 4,
			want1:                       true,
			errsLen:                     0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			trcr := trace.NewNoopTracerProvider().Tracer("test")

			sCF := testServer(tt.catFactServerResponse, tt.args.timeoutCatFacts, 1)
			defer sCF.Close()
			sS := testServer(tt.sentimenterServerResponse, tt.args.timeoutSentimenter, tt.timeoutSentimenterFrequency)
			defer sS.Close()

			got, errs := get1HappyCatFact(ctx, trcr, tt.args.timeoutCatFacts, tt.args.timeoutSentimenter, sCF.URL, sS.URL)
			if got != "" && tt.errsLen != len(errs) {
				t.Errorf("Get1HappyCatFact() error = %v, want %d errors", errs, tt.errsLen)

				return
			}

			if tt.want1 && got == "" {
				t.Errorf("Get1HappyCatFact() expected one fact, got none")
			}
		})
	}
}
