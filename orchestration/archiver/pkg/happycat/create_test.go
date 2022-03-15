package happycat

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestCreateHandler(t *testing.T) {
	t.Parallel()

	sugar := zap.NewNop().Sugar()

	tests := []struct {
		name           string
		payload        string
		wantStatusCode int
		wantErr        bool
	}{
		{
			name:           "happy path",
			payload:        fmt.Sprintf("{\"id\": \"%s\",\"fact\": \"fact\"}", uuid.New().String()),
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "wrong formatted json",
			payload:        "{",
			wantStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := newTestQuerier(nil)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/test/happycatfact", bytes.NewReader([]byte(tt.payload)))

			r := chi.NewRouter()
			AddRoutes("test", r, q, sugar)
			r.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatusCode {
				t.Errorf(
					"handler returned unexpected status: got %d want %d",
					rr.Code, tt.wantStatusCode,
				)
			}

			if !tt.wantErr && len(q.list) != 1 {
				t.Errorf("expected a newly created cat fact, got %d", len(q.list))
			}
		})
	}
}

func BenchmarkCreateHandler(b *testing.B) {
	q := newTestQuerier(nil)
	r := chi.NewRouter()

	sugar := zap.NewNop().Sugar()

	AddRoutes("test", r, q, sugar)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()

		payload := fmt.Sprintf("{\"id\": \"%s\",\"fact\": \"fact\"}", uuid.New().String())

		req := httptest.NewRequest("POST", "/test/happycatfact", bytes.NewReader([]byte(payload)))

		b.StartTimer()

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		b.StopTimer()

		res := rr.Result()
		res.Body.Close()

		b.StartTimer()
	}
}
