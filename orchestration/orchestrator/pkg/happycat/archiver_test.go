package happycat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_saveHappycatFact(t *testing.T) {
	t.Parallel()

	var noTimeout time.Duration

	tests := []struct {
		name             string
		serverStatusCode int
		timeout          time.Duration
		wantErr          bool
	}{
		{
			name:             "happy path",
			serverStatusCode: http.StatusCreated,
			wantErr:          false,
		},
		{
			name:             "error api",
			serverStatusCode: http.StatusInternalServerError,
			wantErr:          true,
		},
		{
			name:             "timeout",
			serverStatusCode: http.StatusCreated,
			timeout:          10 * time.Millisecond,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			s := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if tt.timeout != noTimeout {
						time.Sleep(2 * tt.timeout)
					}

					w.WriteHeader(tt.serverStatusCode)
				}),
			)
			defer s.Close()

			err := saveHappycatFact(ctx, tt.timeout, s.URL, "great cats")
			if (err != nil) == tt.wantErr {
				t.Errorf("saveHappycatFact() error = %v, wantErr %t", err, tt.wantErr)

				return
			}
		})
	}
}
