package happycat

import (
	"context"
	"testing"
	"time"
)

func Test_get1HappyCatFact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                        string
		catFactServerResponse       string
		sentimenterServerResponse   string
		timeoutCatFactFrequency     int
		timeoutSentimenterFrequency int
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
			name: "all timeout cat facts",
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
			timeoutCatFactFrequency:   1,
			want1:                     false,
			errsLen:                   1,
		},
		{
			name: "all timeout sentimenter",
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
			sentimenterServerResponse:   `6`,
			timeoutSentimenterFrequency: 1,
			want1:                       false,
			errsLen:                     4,
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
			sentimenterServerResponse:   `0`,
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
			sentimenterServerResponse:   `1`,
			timeoutSentimenterFrequency: 4,
			want1:                       true,
			errsLen:                     0,
		},
	}

	timeoutCatFacts := 10 * time.Millisecond
	timeoutSentimenter := 10 * time.Millisecond

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			sCF := testServer(t, tt.catFactServerResponse, timeoutCatFacts, tt.timeoutCatFactFrequency)
			defer sCF.Close()
			sS := testServer(t, tt.sentimenterServerResponse, timeoutSentimenter, tt.timeoutSentimenterFrequency)
			defer sS.Close()

			got, errs := get1HappyCatFact(ctx, timeoutCatFacts, timeoutSentimenter, sCF.URL, sS.URL)
			if got != "" && tt.errsLen != len(errs) {
				t.Errorf("Get1HappyCatFact() error = %v, want %d errors", errs, tt.errsLen)
			}

			if tt.want1 && got == "" {
				t.Errorf("Get1HappyCatFact() expected one fact, got none")
			}
		})
	}
}
