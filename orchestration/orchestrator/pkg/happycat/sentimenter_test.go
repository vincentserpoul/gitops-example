package happycat

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func Test_getSentiment(t *testing.T) {
	t.Parallel()

	type args struct {
		timeout       time.Duration
		textToAnalyze string
	}

	tests := []struct {
		name           string
		serverResponse string
		args           args
		want           int
		errIs          error
	}{
		{
			name:           "happy path",
			serverResponse: `1`,
			args: args{
				textToAnalyze: "super",
				timeout:       0,
			},
			want:  1,
			errIs: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trcr := trace.NewNoopTracerProvider().Tracer("test")
			ctx := context.Background()
			s := testServer(t, tt.serverResponse, tt.args.timeout, 1)
			defer s.Close()

			got, err := getSentiment(ctx, trcr, tt.args.timeout, s.URL, tt.args.textToAnalyze)
			if !errors.Is(tt.errIs, err) {
				t.Errorf("getSentiment() error = %v, wantErr %v", err, tt.errIs)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSentiment() = %v, want %v", got, tt.want)
			}
		})
	}
}
